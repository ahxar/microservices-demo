'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Package, ShoppingCart, Users, DollarSign } from 'lucide-react';
import { useOrders } from '@/hooks/use-orders';
import { useProducts } from '@/hooks/use-products';
import { adminApi } from '@/lib/api';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function AdminDashboardPage() {
  const { data: productsData } = useProducts({ page: 1, page_size: 50 });
  const { data: ordersData } = useOrders({ page: 1, page_size: 20 });
  const { data: usersData } = useQuery({
    queryKey: ['admin', 'users'],
    queryFn: () => adminApi.listUsers({ page: 1, page_size: 100 }),
  });

  const totalRevenue =
    ordersData?.orders.reduce((sum, order) => sum + (order.total?.amount_cents ?? 0), 0) ?? 0;

  const stats = [
    {
      title: 'Revenue (visible)',
      value: formatMoney(totalRevenue, 'USD'),
      icon: DollarSign,
    },
    {
      title: 'Orders',
      value: String(ordersData?.pagination?.total_count ?? 0),
      icon: ShoppingCart,
    },
    {
      title: 'Products',
      value: String(productsData?.pagination?.total_count ?? 0),
      icon: Package,
    },
    {
      title: 'Users',
      value: String(usersData?.pagination?.total_count ?? 0),
      icon: Users,
    },
  ];

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Dashboard Overview</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Recent Orders</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {(ordersData?.orders ?? []).slice(0, 5).map((order) => (
                <div key={order.id} className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Order #{order.id.slice(0, 8)}</p>
                    <p className="text-sm text-muted-foreground">
                      {new Date(order.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <span className="font-semibold">
                    {formatMoney(order.total?.amount_cents, order.total?.currency)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top Products</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {(productsData?.products ?? []).slice(0, 5).map((product) => (
                <div key={product.id} className="flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 bg-muted rounded-md overflow-hidden shrink-0">
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
                    <Link className="font-medium hover:underline" href={`/products/${product.id}`}>
                      {product.name}
                    </Link>
                    <p className="text-sm text-muted-foreground">
                      Stock: {product.stock_quantity}
                    </p>
                    </div>
                  </div>
                  <span className="font-semibold">
                    {formatMoney(product.price?.amount_cents, product.price?.currency)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
