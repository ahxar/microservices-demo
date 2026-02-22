'use client';

import { useMemo, useState } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useAddToCart } from '@/hooks/use-cart';
import { useCategories, useProductSearch, useProducts } from '@/hooks/use-products';
import { useAddToWishlist } from '@/hooks/use-user';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function ProductsPage() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get('category_id') || undefined;
  const [query, setQuery] = useState('');
  const isSearching = query.trim().length > 2;

  const { data: productsData, isLoading: productsLoading } = useProducts({
    page: 1,
    page_size: 20,
    category_id: categoryId,
  });
  const { data: searchData, isLoading: searchLoading } = useProductSearch(query.trim());
  const { data: categories = [] } = useCategories();
  const addToCart = useAddToCart();
  const addToWishlist = useAddToWishlist();

  const products = useMemo(
    () => (isSearching ? searchData?.products ?? [] : productsData?.products ?? []),
    [isSearching, productsData?.products, searchData?.products]
  );
  const isLoading = isSearching ? searchLoading : productsLoading;

  return (
    <div className="container mx-auto px-4 py-8 space-y-6">
      <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">Products</h1>
          <p className="text-muted-foreground">Browse products from the catalog service.</p>
        </div>
        <div className="w-full md:w-80">
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search products (min 3 chars)"
          />
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        <Button asChild variant={!categoryId ? 'default' : 'outline'} size="sm">
          <Link href="/products">All</Link>
        </Button>
        {categories.map((category) => (
          <Button
            asChild
            key={category.id}
            variant={categoryId === category.id ? 'default' : 'outline'}
            size="sm"
          >
            <Link href={`/products?category_id=${category.id}`}>{category.name}</Link>
          </Button>
        ))}
      </div>

      {isLoading ? (
        <p className="text-muted-foreground">Loading products...</p>
      ) : products.length === 0 ? (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            No products found.
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.map((product) => (
            <Card key={product.id}>
              <CardHeader>
                <div className="aspect-video overflow-hidden rounded-md bg-muted mb-3">
                  {product.image_urls?.[0] ? (
                    <img
                      src={product.image_urls[0]}
                      alt={product.name}
                      className="h-full w-full object-cover"
                      onError={(event) => {
                        event.currentTarget.style.display = 'none';
                      }}
                    />
                  ) : null}
                </div>
                <CardTitle className="line-clamp-1">{product.name}</CardTitle>
                <CardDescription className="line-clamp-2">{product.description}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                <p className="text-2xl font-bold">
                  {formatMoney(product.price?.amount_cents, product.price?.currency)}
                </p>
                <p className="text-sm text-muted-foreground">
                  Stock: {product.stock_quantity}
                </p>
              </CardContent>
              <CardFooter className="flex flex-col gap-2">
                <Button asChild variant="outline" className="w-full">
                  <Link href={`/products/${product.id}`}>View Details</Link>
                </Button>
                <Button
                  className="w-full"
                  onClick={() =>
                    addToCart.mutate({
                      product_id: product.id,
                      product_name: product.name,
                      quantity: 1,
                      unit_price: product.price,
                      image_url: product.image_urls?.[0],
                    })
                  }
                  disabled={addToCart.isPending || product.stock_quantity <= 0}
                >
                  Add to Cart
                </Button>
                <Button
                  variant="ghost"
                  className="w-full"
                  onClick={() => addToWishlist.mutate(product.id)}
                  disabled={addToWishlist.isPending}
                >
                  Add to Wishlist
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
