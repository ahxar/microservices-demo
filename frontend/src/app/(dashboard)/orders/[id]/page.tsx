'use client';

import { useParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useOrder } from '@/hooks/use-orders';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

function statusLabel(status: number) {
  switch (status) {
    case 1:
      return 'pending';
    case 2:
      return 'confirmed';
    case 3:
      return 'processing';
    case 4:
      return 'shipped';
    case 5:
      return 'delivered';
    case 6:
      return 'cancelled';
    case 7:
      return 'refunded';
    default:
      return 'unknown';
  }
}

function statusClass(status: number) {
  switch (status) {
    case 5:
      return 'bg-green-500';
    case 4:
      return 'bg-blue-500';
    case 3:
      return 'bg-yellow-500';
    case 6:
      return 'bg-red-500';
    default:
      return 'bg-gray-500';
  }
}

export default function OrderDetailPage() {
  const params = useParams<{ id: string }>();
  const orderId = params.id;
  const { data: order, isLoading } = useOrder(orderId);

  if (isLoading) {
    return <p className="text-muted-foreground">Loading order details...</p>;
  }

  if (!order) {
    return (
      <Card>
        <CardContent className="pt-6">Order not found.</CardContent>
      </Card>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Order #{order.id.slice(0, 8)}</h1>
        <Badge className={statusClass(order.status)}>{statusLabel(order.status)}</Badge>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Order Timeline</CardTitle>
          </CardHeader>
          <CardContent>
            {(order.history ?? []).length === 0 ? (
              <p className="text-sm text-muted-foreground">No status history available.</p>
            ) : (
              <div className="space-y-4">
                {order.history?.map((event) => (
                  <div key={event.id} className="flex items-start space-x-4">
                    <div className="w-2 h-2 mt-2 rounded-full bg-primary" />
                    <div>
                      <p className="font-semibold">{statusLabel(event.status)}</p>
                      <p className="text-sm text-muted-foreground">
                        {new Date(event.created_at).toLocaleString()}
                      </p>
                      {event.notes && (
                        <p className="text-sm text-muted-foreground">{event.notes}</p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Order Items</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {order.items?.map((item) => (
                <div key={item.id} className="flex items-center space-x-4">
                  <div className="w-20 h-20 bg-muted rounded-md" />
                  <div className="flex-1">
                    <p className="font-semibold">{item.product_name}</p>
                    <p className="text-sm text-muted-foreground">Quantity: {item.quantity}</p>
                  </div>
                  <p className="font-semibold">
                    {formatMoney(item.total_price?.amount_cents, item.total_price?.currency)}
                  </p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Shipping Address</CardTitle>
          </CardHeader>
          <CardContent className="text-sm">
            <p>{order.shipping_address?.street}</p>
            <p>
              {order.shipping_address?.city}, {order.shipping_address?.state}{' '}
              {order.shipping_address?.zip_code}
            </p>
            <p>{order.shipping_address?.country}</p>
            {order.tracking_number && (
              <p className="mt-2">Tracking: {order.tracking_number}</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Order Summary</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Subtotal</span>
              <span>{formatMoney(order.subtotal?.amount_cents, order.subtotal?.currency)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Shipping</span>
              <span>{formatMoney(order.shipping?.amount_cents, order.shipping?.currency)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Tax</span>
              <span>{formatMoney(order.tax?.amount_cents, order.tax?.currency)}</span>
            </div>
            <div className="border-t pt-2">
              <div className="flex justify-between font-bold text-lg">
                <span>Total</span>
                <span>{formatMoney(order.total?.amount_cents, order.total?.currency)}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
