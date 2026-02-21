import { useQuery } from '@tanstack/react-query';
import { productsApi } from '@/lib/api';

export function useProducts(params?: {
  page?: number;
  page_size?: number;
  category_id?: string;
}) {
  return useQuery({
    queryKey: ['products', params],
    queryFn: () => productsApi.getProducts(params),
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: ['products', id],
    queryFn: () => productsApi.getProduct(id),
    enabled: !!id,
  });
}

export function useProductSearch(query: string) {
  return useQuery({
    queryKey: ['products', 'search', query],
    queryFn: () => productsApi.searchProducts(query),
    enabled: !!query && query.length > 2,
  });
}

export function useCategories() {
  return useQuery({
    queryKey: ['categories'],
    queryFn: () => productsApi.getCategories(),
  });
}
