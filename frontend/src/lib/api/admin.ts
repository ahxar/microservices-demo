import apiClient from './client';
import type { Money, PaginationResponse, Product } from './products';
import type { User } from './user';

export interface ListUsersResponse {
  users: User[];
  pagination: PaginationResponse;
}

export interface CreateProductRequest {
  name: string;
  slug: string;
  description: string;
  price: Money;
  category_id: string;
  image_urls?: string[];
  stock_quantity: number;
}

export interface UpdateProductRequest extends CreateProductRequest {
  is_active: boolean;
}

export const adminApi = {
  listUsers: async (params?: {
    page?: number;
    page_size?: number;
  }): Promise<ListUsersResponse> => {
    const response = await apiClient.get('/api/v1/admin/users', { params });
    return response.data;
  },

  createProduct: async (data: CreateProductRequest): Promise<Product> => {
    const response = await apiClient.post('/api/v1/admin/products', data);
    return response.data;
  },

  updateProduct: async (id: string, data: UpdateProductRequest): Promise<Product> => {
    const response = await apiClient.put(`/api/v1/admin/products/${id}`, data);
    return response.data;
  },

  deleteProduct: async (id: string): Promise<void> => {
    await apiClient.delete(`/api/v1/admin/products/${id}`);
  },
};
