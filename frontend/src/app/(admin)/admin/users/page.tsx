'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { adminApi } from '@/lib/api';

export default function AdminUsersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['admin', 'users', 'page-1'],
    queryFn: () => adminApi.listUsers({ page: 1, page_size: 100 }),
  });

  const users = data?.users ?? [];

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Users</h1>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">User</th>
                  <th className="text-left p-4">Email</th>
                  <th className="text-left p-4">Joined</th>
                  <th className="text-left p-4">Role</th>
                  <th className="text-left p-4">Status</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td className="p-4 text-muted-foreground" colSpan={5}>Loading users...</td>
                  </tr>
                ) : (
                  users.map((user) => (
                    <tr key={user.id} className="border-b last:border-b-0">
                      <td className="p-4 font-medium">
                        {user.profile?.first_name || user.profile?.last_name
                          ? `${user.profile?.first_name ?? ''} ${user.profile?.last_name ?? ''}`.trim()
                          : user.id.slice(0, 8)}
                      </td>
                      <td className="p-4">{user.email}</td>
                      <td className="p-4">{new Date(user.created_at).toLocaleDateString()}</td>
                      <td className="p-4">{user.role}</td>
                      <td className="p-4">
                        <Badge className="bg-green-500">Active</Badge>
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
