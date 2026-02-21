package server

import (
	"context"
	"math"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/catalog/v1"
	"github.com/safar/microservices-demo/services/catalog/internal/repository"
	"github.com/safar/microservices-demo/services/catalog/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedCatalogServiceServer
	catalogService *service.CatalogService
}

func NewGRPCServer(catalogService *service.CatalogService) *GRPCServer {
	return &GRPCServer{
		catalogService: catalogService,
	}
}

func (s *GRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	products, total, err := s.catalogService.ListProducts(ctx, page, pageSize, req.CategoryId, req.ActiveOnly)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		categoryID := ""
		if p.CategoryID.Valid {
			categoryID = p.CategoryID.String
		}

		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Price: &commonv1.Money{
				AmountCents: p.PriceCents,
				Currency:    p.Currency,
			},
			CategoryId:    categoryID,
			ImageUrls:     p.ImageURLs,
			StockQuantity: p.StockQuantity,
			IsActive:      p.IsActive,
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Pagination: &commonv1.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalPages: totalPages,
			TotalCount: int64(total),
		},
	}, nil
}

func (s *GRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	var product interface{}
	var err error

	switch id := req.Identifier.(type) {
	case *pb.GetProductRequest_Id:
		product, err = s.catalogService.GetProductByID(ctx, id.Id)
	case *pb.GetProductRequest_Slug:
		product, err = s.catalogService.GetProductBySlug(ctx, id.Slug)
	default:
		return nil, status.Error(codes.InvalidArgument, "product ID or slug is required")
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	p := product.(*repository.Product)
	categoryID := ""
	if p.CategoryID.Valid {
		categoryID = p.CategoryID.String
	}

	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price: &commonv1.Money{
			AmountCents: p.PriceCents,
			Currency:    p.Currency,
		},
		CategoryId:    categoryID,
		ImageUrls:     p.ImageURLs,
		StockQuantity: p.StockQuantity,
		IsActive:      p.IsActive,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "search query is required")
	}

	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	products, total, err := s.catalogService.SearchProducts(ctx, req.Query, page, pageSize, req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search products: %v", err)
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		categoryID := ""
		if p.CategoryID.Valid {
			categoryID = p.CategoryID.String
		}

		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Price: &commonv1.Money{
				AmountCents: p.PriceCents,
				Currency:    p.Currency,
			},
			CategoryId:    categoryID,
			ImageUrls:     p.ImageURLs,
			StockQuantity: p.StockQuantity,
			IsActive:      p.IsActive,
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Pagination: &commonv1.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalPages: totalPages,
			TotalCount: int64(total),
		},
	}, nil
}

func (s *GRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "name and slug are required")
	}

	if req.Price == nil {
		return nil, status.Error(codes.InvalidArgument, "price is required")
	}

	product, err := s.catalogService.CreateProduct(
		ctx,
		req.Name,
		req.Slug,
		req.Description,
		req.Price.AmountCents,
		req.Price.Currency,
		req.CategoryId,
		req.ImageUrls,
		req.StockQuantity,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	categoryID := ""
	if product.CategoryID.Valid {
		categoryID = product.CategoryID.String
	}

	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price: &commonv1.Money{
			AmountCents: product.PriceCents,
			Currency:    product.Currency,
		},
		CategoryId:    categoryID,
		ImageUrls:     product.ImageURLs,
		StockQuantity: product.StockQuantity,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	if req.Price == nil {
		return nil, status.Error(codes.InvalidArgument, "price is required")
	}

	product, err := s.catalogService.UpdateProduct(
		ctx,
		req.Id,
		req.Name,
		req.Slug,
		req.Description,
		req.Price.AmountCents,
		req.Price.Currency,
		req.CategoryId,
		req.ImageUrls,
		req.StockQuantity,
		req.IsActive,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	categoryID := ""
	if product.CategoryID.Valid {
		categoryID = product.CategoryID.String
	}

	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price: &commonv1.Money{
			AmountCents: product.PriceCents,
			Currency:    product.Currency,
		},
		CategoryId:    categoryID,
		ImageUrls:     product.ImageURLs,
		StockQuantity: product.StockQuantity,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*commonv1.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	if err := s.catalogService.DeleteProduct(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func (s *GRPCServer) ListCategories(ctx context.Context, req *commonv1.Empty) (*pb.ListCategoriesResponse, error) {
	categories, err := s.catalogService.ListCategories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
	}

	var pbCategories []*pb.Category
	for _, c := range categories {
		parentID := ""
		if c.ParentID.Valid {
			parentID = c.ParentID.String
		}

		pbCategories = append(pbCategories, &pb.Category{
			Id:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			ParentId:    parentID,
			CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &pb.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}

func (s *GRPCServer) CheckInventory(ctx context.Context, req *pb.CheckInventoryRequest) (*pb.CheckInventoryResponse, error) {
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items are required")
	}

	items := make(map[string]int32)
	for _, item := range req.Items {
		items[item.ProductId] = item.Quantity
	}

	available, unavailable, err := s.catalogService.CheckInventory(ctx, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check inventory: %v", err)
	}

	return &pb.CheckInventoryResponse{
		Available:             available,
		UnavailableProductIds: unavailable,
	}, nil
}

func (s *GRPCServer) ReserveInventory(ctx context.Context, req *pb.ReserveInventoryRequest) (*pb.ReserveInventoryResponse, error) {
	if req.OrderId == "" || len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID and items are required")
	}

	items := make(map[string]int32)
	for _, item := range req.Items {
		items[item.ProductId] = item.Quantity
	}

	expirationMinutes := req.ExpirationMinutes
	if expirationMinutes <= 0 {
		expirationMinutes = 15
	}

	reservationID, err := s.catalogService.ReserveInventory(ctx, req.OrderId, items, expirationMinutes)
	if err != nil {
		return &pb.ReserveInventoryResponse{
			Success: false,
		}, nil
	}

	return &pb.ReserveInventoryResponse{
		Success:       true,
		ReservationId: reservationID,
	}, nil
}
