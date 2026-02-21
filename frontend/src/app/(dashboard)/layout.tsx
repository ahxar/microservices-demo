import Link from 'next/link';
import { User, Package, MapPin, Heart } from 'lucide-react';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { DashboardHeader } from '@/components/dashboard/dashboard-header';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ProtectedRoute>
      <div className="min-h-screen flex flex-col">
        <DashboardHeader />

        <div className="flex-1 container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            {/* Sidebar */}
            <aside className="space-y-2">
              <Link
                href="/orders"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <Package className="h-5 w-5" />
                <span>Orders</span>
              </Link>
              <Link
                href="/profile"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <User className="h-5 w-5" />
                <span>Profile</span>
              </Link>
              <Link
                href="/addresses"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <MapPin className="h-5 w-5" />
                <span>Addresses</span>
              </Link>
              <Link
                href="/wishlist"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-muted"
              >
                <Heart className="h-5 w-5" />
                <span>Wishlist</span>
              </Link>
            </aside>

            {/* Main Content */}
            <main className="md:col-span-3">
              {children}
            </main>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
