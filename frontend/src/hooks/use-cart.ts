import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { cartApi, type AddCartItemRequest } from '@/lib/api';
import { hasAccessToken, redirectToLogin } from '@/lib/auth-redirect';

export function useCart() {
  return useQuery({
    queryKey: ['cart'],
    queryFn: () => cartApi.getCart(),
  });
}

export function useAddToCart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AddCartItemRequest) => {
      if (!hasAccessToken()) {
        redirectToLogin('action_requires_auth');
        return Promise.reject(new Error('Authentication required'));
      }
      return cartApi.addItem(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });
}

export function useUpdateCartItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ productId, quantity }: { productId: string; quantity: number }) =>
      cartApi.updateItem(productId, quantity),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });
}

export function useRemoveFromCart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (productId: string) => cartApi.removeItem(productId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });
}

export function useClearCart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => cartApi.clearCart(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });
}
