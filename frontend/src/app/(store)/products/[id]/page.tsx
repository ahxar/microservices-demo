'use client';

import Link from 'next/link';
import { useParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { useAddToCart } from '@/hooks/use-cart';
import { useProduct } from '@/hooks/use-products';
import { useAddToWishlist } from '@/hooks/use-user';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function ProductDetailPage() {
  const params = useParams<{ id: string }>();
  const productId = params.id;

  const { data: product, isLoading, isError } = useProduct(productId);
  const addToCart = useAddToCart();
  const addToWishlist = useAddToWishlist();

  if (isLoading) {
    return <div className="container mx-auto px-4 py-8 text-muted-foreground">Loading product...</div>;
  }

  if (isError || !product) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Card>
          <CardContent className="pt-6">
            <p className="mb-4">Product not found.</p>
            <Button asChild variant="outline">
              <Link href="/products">Back to Products</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="aspect-square bg-muted rounded-lg overflow-hidden">
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

        <div>
          <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
          <p className="text-2xl font-bold mb-4">
            {formatMoney(product.price?.amount_cents, product.price?.currency)}
          </p>
          <p className="text-muted-foreground mb-6">{product.description}</p>

          <div className="space-y-3">
            <Button
              size="lg"
              className="w-full md:w-auto"
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
              variant="outline"
              size="lg"
              className="w-full md:w-auto"
              onClick={() => addToWishlist.mutate(product.id)}
              disabled={addToWishlist.isPending}
            >
              Add to Wishlist
            </Button>
          </div>

          <Card className="mt-8">
            <CardContent className="pt-6">
              <h3 className="font-semibold mb-2">Product Details</h3>
              <dl className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Product ID</dt>
                  <dd className="truncate ml-4">{product.id}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Category ID</dt>
                  <dd className="truncate ml-4">{product.category_id}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Stock</dt>
                  <dd>{product.stock_quantity}</dd>
                </div>
              </dl>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
