'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { useAddToCart } from '@/hooks/use-cart';
import { useWishlist } from '@/hooks/use-user';
import { productsApi } from '@/lib/api';

export default function WishlistPage() {
  const { data: wishlist = [], isLoading } = useWishlist();
  const addToCart = useAddToCart();
  const [loadingProduct, setLoadingProduct] = useState<string | null>(null);

  const moveToCart = async (productId: string) => {
    setLoadingProduct(productId);
    try {
      const product = await productsApi.getProduct(productId);
      await addToCart.mutateAsync({
        product_id: product.id,
        product_name: product.name,
        quantity: 1,
        unit_price: product.price,
        image_url: product.image_urls?.[0],
      });
    } finally {
      setLoadingProduct(null);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">My Wishlist</h1>

      <p className="text-sm text-muted-foreground mb-4">
        Wishlist removal is not exposed by the current backend API yet.
      </p>

      {isLoading ? (
        <p className="text-muted-foreground">Loading wishlist...</p>
      ) : wishlist.length === 0 ? (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">Your wishlist is empty.</CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {wishlist.map((item) => (
            <Card key={item.id}>
              <CardHeader>
                <div className="aspect-square bg-muted rounded-md mb-4" />
                <CardTitle className="truncate">Product {item.product_id.slice(0, 8)}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <p className="text-sm text-muted-foreground">Added {new Date(item.added_at).toLocaleDateString()}</p>
                <Button variant="outline" asChild className="w-full">
                  <Link href={`/products/${item.product_id}`}>View Product</Link>
                </Button>
              </CardContent>
              <CardFooter>
                <Button
                  className="w-full"
                  onClick={() => moveToCart(item.product_id)}
                  disabled={loadingProduct === item.product_id || addToCart.isPending}
                >
                  {loadingProduct === item.product_id ? 'Adding...' : 'Add to Cart'}
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
