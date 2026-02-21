'use client';

import Link from 'next/link';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Eye } from 'lucide-react';
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

export default function AdminOrdersPage() {
  const { data, isLoading } = useOrders({ page: 1, page_size: 50 });
  const orders = data?.orders ?? [];

  return (
    <div>
      <h1 className="text-3xl font-bold mb-2">Orders</h1>
      <p className="text-sm text-muted-foreground mb-6">
        This view currently lists orders visible to the signed-in user.
      </p>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">Order ID</th>
                  <th className="text-left p-4">Date</th>
                  <th className="text-left p-4">Total</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-right p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td className="p-4 text-muted-foreground" colSpan={5}>Loading orders...</td>
                  </tr>
                ) : (
                  orders.map((order) => (
                    <tr key={order.id} className="border-b last:border-b-0">
                      <td className="p-4 font-medium">#{order.id.slice(0, 8)}</td>
                      <td className="p-4">{new Date(order.created_at).toLocaleDateString()}</td>
                      <td className="p-4">
                        {formatMoney(order.total?.amount_cents, order.total?.currency)}
                      </td>
                      <td className="p-4">
                        <Badge className={statusClass(order.status)}>{statusLabel(order.status)}</Badge>
                      </td>
                      <td className="p-4">
                        <div className="flex justify-end">
                          <Button variant="outline" size="sm" asChild>
                            <Link href={`/orders/${order.id}`}>
                              <Eye className="h-4 w-4 mr-2" />
                              View
                            </Link>
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
