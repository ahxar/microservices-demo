import apiClient from './client';
import type { Address } from './user';
import type { Money, PaginationResponse } from './products';

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: Money;
  total_price: Money;
}

export interface OrderStatusHistory {
  id: string;
  status: number;
  notes: string;
  created_at: string;
}

export interface Order {
  id: string;
  user_id: string;
  status: number;
  items: OrderItem[];
  subtotal: Money;
  shipping: Money;
  tax: Money;
  total: Money;
  shipping_address: Address;
  payment_method_id: string;
  transaction_id?: string;
  tracking_number?: string;
  history?: OrderStatusHistory[];
  created_at: string;
  updated_at: string;
}

export interface OrdersResponse {
  orders: Order[];
  pagination: PaginationResponse;
}

export interface CreateOrderRequest {
  shipping_address: Address;
  payment_method_id: string;
}

export const ordersApi = {
  getOrders: async (params?: {
    page?: number;
    page_size?: number;
    status?: string;
  }): Promise<OrdersResponse> => {
    const response = await apiClient.get('/api/v1/orders', { params });
    return response.data;
  },

  getOrder: async (id: string): Promise<Order> => {
    const response = await apiClient.get(`/api/v1/orders/${id}`);
    return response.data;
  },

  createOrder: async (data: CreateOrderRequest): Promise<Order> => {
    const response = await apiClient.post('/api/v1/orders', data);
    return response.data;
  },

  cancelOrder: async (id: string, reason?: string): Promise<Order> => {
    const response = await apiClient.delete(`/api/v1/orders/${id}`, {
      params: reason ? { reason } : undefined,
    });
    return response.data;
  },
};
