'use client';

import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useOrders } from '@/hooks/use-orders';

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

export default function OrdersPage() {
  const { data, isLoading } = useOrders({ page: 1, page_size: 20 });
  const orders = data?.orders ?? [];

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">My Orders</h1>

      {isLoading ? (
        <p className="text-muted-foreground">Loading orders...</p>
      ) : orders.length === 0 ? (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            No orders yet.
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <Link key={order.id} href={`/orders/${order.id}`}>
              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>Order #{order.id.slice(0, 8)}</CardTitle>
                      <CardDescription>
                        Placed on {new Date(order.created_at).toLocaleDateString()}
                      </CardDescription>
                    </div>
                    <Badge className={statusClass(order.status)}>
                      {statusLabel(order.status)}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex justify-between text-sm">
                    <span>{order.items?.length ?? 0} item(s)</span>
                    <span className="font-semibold">
                      {formatMoney(order.total?.amount_cents, order.total?.currency)}
                    </span>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
