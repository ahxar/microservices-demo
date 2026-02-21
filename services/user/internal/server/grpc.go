package server

import (
	"context"
	"math"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/user/v1"
	"github.com/safar/microservices-demo/services/user/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewGRPCServer(userService *service.UserService) *GRPCServer {
	return &GRPCServer{
		userService: userService,
	}
}

func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, profile, accessToken, refreshToken, err := s.userService.Register(
		ctx, req.Email, req.Password, req.FirstName, req.LastName,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
			Role:  user.Role,
			Profile: &pb.Profile{
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
				Phone:     profile.Phone,
				AvatarUrl: profile.AvatarURL,
			},
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		ExpiresIn: 3600,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, profile, accessToken, refreshToken, err := s.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	pbProfile := &pb.Profile{}
	if profile != nil {
		pbProfile = &pb.Profile{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Phone:     profile.Phone,
			AvatarUrl: profile.AvatarURL,
		}
	}

	return &pb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			Profile:   pbProfile,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		ExpiresIn: 3600,
	}, nil
}

func (s *GRPCServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	userID, role, err := s.userService.ValidateToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	user, profile, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Generate new tokens
	_, _, accessToken, refreshToken, err := s.userService.Register(ctx, user.Email, "", "", "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	pbProfile := &pb.Profile{}
	if profile != nil {
		pbProfile = &pb.Profile{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Phone:     profile.Phone,
			AvatarUrl: profile.AvatarURL,
		}
	}

	return &pb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Role:      role,
			Profile:   pbProfile,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		ExpiresIn: 3600,
	}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	userID, role, err := s.userService.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Role:   role,
	}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	user, profile, err := s.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	pbProfile := &pb.Profile{}
	if profile != nil {
		pbProfile = &pb.Profile{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Phone:     profile.Phone,
			AvatarUrl: profile.AvatarURL,
		}
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Profile:   pbProfile,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	firstName := ""
	lastName := ""
	phone := ""
	avatarURL := ""

	if req.Profile != nil {
		firstName = req.Profile.FirstName
		lastName = req.Profile.LastName
		phone = req.Profile.Phone
		avatarURL = req.Profile.AvatarUrl
	}

	user, profile, err := s.userService.UpdateUser(ctx, req.Id, firstName, lastName, phone, avatarURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	pbProfile := &pb.Profile{}
	if profile != nil {
		pbProfile = &pb.Profile{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Phone:     profile.Phone,
			AvatarUrl: profile.AvatarURL,
		}
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Profile:   pbProfile,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, total, err := s.userService.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Pagination: &commonv1.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalPages: totalPages,
			TotalCount: int64(total),
		},
	}, nil
}

func (s *GRPCServer) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.UserAddress, error) {
	if req.UserId == "" || req.Address == nil {
		return nil, status.Error(codes.InvalidArgument, "user ID and address are required")
	}

	address, err := s.userService.AddAddress(
		ctx,
		req.UserId,
		req.Label,
		req.Address.Street,
		req.Address.City,
		req.Address.State,
		req.Address.ZipCode,
		req.Address.Country,
		req.IsDefault,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add address: %v", err)
	}

	return &pb.UserAddress{
		Id:     address.ID,
		UserId: address.UserID,
		Label:  address.Label,
		Address: &commonv1.Address{
			Street:  address.Street,
			City:    address.City,
			State:   address.State,
			ZipCode: address.ZipCode,
			Country: address.Country,
		},
		IsDefault: address.IsDefault,
	}, nil
}

func (s *GRPCServer) ListAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	addresses, err := s.userService.ListAddresses(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list addresses: %v", err)
	}

	var pbAddresses []*pb.UserAddress
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.UserAddress{
			Id:     addr.ID,
			UserId: addr.UserID,
			Label:  addr.Label,
			Address: &commonv1.Address{
				Street:  addr.Street,
				City:    addr.City,
				State:   addr.State,
				ZipCode: addr.ZipCode,
				Country: addr.Country,
			},
			IsDefault: addr.IsDefault,
		})
	}

	return &pb.ListAddressesResponse{
		Addresses: pbAddresses,
	}, nil
}

func (s *GRPCServer) AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*commonv1.Empty, error) {
	if req.UserId == "" || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and product ID are required")
	}

	if err := s.userService.AddToWishlist(ctx, req.UserId, req.ProductId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add to wishlist: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func (s *GRPCServer) GetWishlist(ctx context.Context, req *pb.GetWishlistRequest) (*pb.WishlistResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	items, err := s.userService.GetWishlist(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wishlist: %v", err)
	}

	var pbItems []*pb.WishlistItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.WishlistItem{
			Id:        item.ID,
			ProductId: item.ProductID,
			AddedAt:   item.AddedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &pb.WishlistResponse{
		Items: pbItems,
	}, nil
}
