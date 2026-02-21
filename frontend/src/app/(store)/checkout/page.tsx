'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useCart } from '@/hooks/use-cart';
import { useCreateOrder } from '@/hooks/use-orders';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function CheckoutPage() {
  const router = useRouter();
  const { data: cart } = useCart();
  const createOrder = useCreateOrder();
  const [error, setError] = useState('');
  const [form, setForm] = useState({
    street: '',
    city: '',
    state: '',
    zip_code: '',
    country: 'USA',
    payment_method_id: 'pm_test_visa',
  });

  const items = cart?.items ?? [];

  const placeOrder = async () => {
    setError('');

    if (!form.street || !form.city || !form.state || !form.zip_code || !form.payment_method_id) {
      setError('Please fill in all required fields.');
      return;
    }

    try {
      const order = await createOrder.mutateAsync({
        shipping_address: {
          street: form.street,
          city: form.city,
          state: form.state,
          zip_code: form.zip_code,
          country: form.country,
        },
        payment_method_id: form.payment_method_id,
      });
      router.push(`/orders/${order.id}`);
    } catch (err: any) {
      setError(err?.response?.data?.message ?? 'Failed to place order.');
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Checkout</h1>

      {items.length === 0 ? (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            Your cart is empty. Add products before checkout.
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Shipping Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="street">Street</Label>
                  <Input
                    id="street"
                    placeholder="123 Main St"
                    value={form.street}
                    onChange={(e) => setForm((prev) => ({ ...prev, street: e.target.value }))}
                  />
                </div>
                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="city">City</Label>
                    <Input
                      id="city"
                      placeholder="New York"
                      value={form.city}
                      onChange={(e) => setForm((prev) => ({ ...prev, city: e.target.value }))}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="state">State</Label>
                    <Input
                      id="state"
                      placeholder="NY"
                      value={form.state}
                      onChange={(e) => setForm((prev) => ({ ...prev, state: e.target.value }))}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="zip">ZIP</Label>
                    <Input
                      id="zip"
                      placeholder="10001"
                      value={form.zip_code}
                      onChange={(e) => setForm((prev) => ({ ...prev, zip_code: e.target.value }))}
                    />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Payment Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="paymentMethod">Payment Method ID</Label>
                  <Input
                    id="paymentMethod"
                    placeholder="pm_test_visa"
                    value={form.payment_method_id}
                    onChange={(e) =>
                      setForm((prev) => ({ ...prev, payment_method_id: e.target.value }))
                    }
                  />
                </div>
              </CardContent>
            </Card>
          </div>

          <div>
            <Card>
              <CardHeader>
                <CardTitle>Order Summary</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  {items.map((item) => (
                    <div key={item.product_id} className="flex justify-between text-sm">
                      <span>{item.product_name} x{item.quantity}</span>
                      <span>
                        {formatMoney(item.total_price?.amount_cents, item.total_price?.currency)}
                      </span>
                    </div>
                  ))}
                </div>
                <div className="border-t pt-4">
                  <div className="flex justify-between font-bold text-lg">
                    <span>Total</span>
                    <span>{formatMoney(cart?.total?.amount_cents, cart?.total?.currency)}</span>
                  </div>
                </div>
                {error && (
                  <p className="text-sm text-destructive">{error}</p>
                )}
                <Button
                  className="w-full"
                  size="lg"
                  onClick={placeOrder}
                  disabled={createOrder.isPending}
                >
                  {createOrder.isPending ? 'Placing Order...' : 'Place Order'}
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      )}
    </div>
  );
}
