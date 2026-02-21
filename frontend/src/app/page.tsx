'use client';

import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useCategories, useProducts } from '@/hooks/use-products';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function Home() {
  const { data: productsData, isLoading: productsLoading } = useProducts({ page: 1, page_size: 6 });
  const { data: categories } = useCategories();

  const products = productsData?.products ?? [];

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/30">
      <div className="container mx-auto px-4 py-12">
        <div className="text-center max-w-3xl mx-auto mb-12">
          <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4">
            Microservices E-commerce Platform
          </h1>
          <p className="text-muted-foreground text-lg mb-6">
            Browse real products, manage your cart, and complete checkout through the API gateway.
          </p>
          <div className="flex items-center justify-center gap-3">
            <Button asChild size="lg">
              <Link href="/products">Shop Products</Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/categories">View Categories</Link>
            </Button>
          </div>
        </div>

        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">Featured Products</h2>
          {productsLoading ? (
            <p className="text-muted-foreground">Loading featured products...</p>
          ) : products.length === 0 ? (
            <p className="text-muted-foreground">No products available yet. Seed data to get started.</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {products.map((product) => (
                <Card key={product.id}>
                  <CardHeader>
                    <CardTitle className="line-clamp-1">{product.name}</CardTitle>
                    <CardDescription className="line-clamp-2">{product.description}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <p className="text-2xl font-bold">
                      {formatMoney(product.price?.amount_cents, product.price?.currency)}
                    </p>
                    <Button asChild className="w-full">
                      <Link href={`/products/${product.id}`}>View Product</Link>
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </section>

        <section>
          <h2 className="text-2xl font-semibold mb-4">Popular Categories</h2>
          <div className="flex flex-wrap gap-2">
            {(categories ?? []).slice(0, 8).map((category) => (
              <Button key={category.id} variant="outline" asChild>
                <Link href={`/products?category_id=${category.id}`}>{category.name}</Link>
              </Button>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
