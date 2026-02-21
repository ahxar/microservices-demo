import apiClient from './client';

export interface Money {
  amount_cents: number;
  currency: string;
}

export interface PaginationResponse {
  page: number;
  page_size: number;
  total_pages: number;
  total_count: number;
}

export interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: Money;
  category_id: string;
  image_urls: string[];
  stock_quantity: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductsResponse {
  products: Product[];
  pagination: PaginationResponse;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description: string;
  parent_id?: string;
  created_at: string;
}

export interface CategoriesResponse {
  categories: Category[];
}

export const productsApi = {
  getProducts: async (params?: {
    page?: number;
    page_size?: number;
    category_id?: string;
  }): Promise<ProductsResponse> => {
    const response = await apiClient.get('/api/v1/products', { params });
    return response.data;
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await apiClient.get(`/api/v1/products/${id}`);
    return response.data;
  },

  searchProducts: async (query: string): Promise<ProductsResponse> => {
    const response = await apiClient.get('/api/v1/products/search', {
      params: { q: query },
    });
    return response.data;
  },

  getCategories: async (): Promise<Category[]> => {
    const response = await apiClient.get('/api/v1/categories');
    const data: CategoriesResponse = response.data;
    return data.categories ?? [];
  },
};
