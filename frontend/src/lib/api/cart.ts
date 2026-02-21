import apiClient from './client';
import type { Money } from './products';

export interface CartItem {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: Money;
  total_price: Money;
  image_url?: string;
}

export interface Cart {
  user_id: string;
  items: CartItem[];
  total: Money;
  updated_at: string;
}

export interface AddCartItemRequest {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: Money;
  image_url?: string;
}

export const cartApi = {
  getCart: async (): Promise<Cart> => {
    const response = await apiClient.get('/api/v1/cart');
    return response.data;
  },

  addItem: async (data: AddCartItemRequest): Promise<Cart> => {
    const response = await apiClient.post('/api/v1/cart/items', data);
    return response.data;
  },

  updateItem: async (productId: string, quantity: number): Promise<Cart> => {
    const response = await apiClient.put(`/api/v1/cart/items/${productId}`, {
      quantity,
    });
    return response.data;
  },

  removeItem: async (productId: string): Promise<Cart> => {
    const response = await apiClient.delete(`/api/v1/cart/items/${productId}`);
    return response.data;
  },

  clearCart: async (): Promise<void> => {
    await apiClient.delete('/api/v1/cart');
  },
};
