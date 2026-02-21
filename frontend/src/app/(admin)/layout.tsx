import Link from 'next/link';
import { LayoutDashboard, Package, ShoppingCart, Users } from 'lucide-react';
import { AdminRoute } from '@/components/auth/admin-route';
import { AdminHeader } from '@/components/admin/admin-header';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AdminRoute>
      <div className="min-h-screen flex flex-col">
        <AdminHeader />

        <div className="flex-1 container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 md:grid-cols-5 gap-8">
            {/* Sidebar */}
            <aside className="space-y-2">
              <Link
                href="/admin/dashboard"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <LayoutDashboard className="h-5 w-5" />
                <span>Dashboard</span>
              </Link>
              <Link
                href="/admin/products"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <Package className="h-5 w-5" />
                <span>Products</span>
              </Link>
              <Link
                href="/admin/orders"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <ShoppingCart className="h-5 w-5" />
                <span>Orders</span>
              </Link>
              <Link
                href="/admin/users"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <Users className="h-5 w-5" />
                <span>Users</span>
              </Link>
            </aside>

            {/* Main Content */}
            <main className="md:col-span-4">
              {children}
            </main>
          </div>
        </div>
      </div>
    </AdminRoute>
  );
}
