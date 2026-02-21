import apiClient from './client';

export interface Profile {
  first_name: string;
  last_name: string;
  phone?: string;
  avatar_url?: string;
}

export interface User {
  id: string;
  email: string;
  role: string;
  profile?: Profile;
  created_at: string;
  updated_at: string;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip_code: string;
  country: string;
}

export interface UserAddress {
  id: string;
  user_id: string;
  label: string;
  address: Address;
  is_default: boolean;
}

export interface WishlistItem {
  id: string;
  product_id: string;
  added_at: string;
}

export const userApi = {
  getMe: async (): Promise<User> => {
    const response = await apiClient.get('/api/v1/me');
    return response.data;
  },

  updateProfile: async (data: Partial<Profile>): Promise<User> => {
    const response = await apiClient.put('/api/v1/me', data);
    return response.data;
  },

  getAddresses: async (): Promise<UserAddress[]> => {
    const response = await apiClient.get('/api/v1/addresses');
    return response.data.addresses ?? [];
  },

  addAddress: async (data: { label: string; address: Address; is_default?: boolean }): Promise<UserAddress> => {
    const response = await apiClient.post('/api/v1/addresses', data);
    return response.data;
  },

  getWishlist: async (): Promise<WishlistItem[]> => {
    const response = await apiClient.get('/api/v1/wishlist');
    return response.data.items ?? [];
  },

  addToWishlist: async (productId: string): Promise<void> => {
    await apiClient.post('/api/v1/wishlist', { product_id: productId });
  },
};
