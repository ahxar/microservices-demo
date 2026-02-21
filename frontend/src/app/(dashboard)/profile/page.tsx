'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useUpdateProfile, useUser } from '@/hooks/use-user';

export default function ProfilePage() {
  const { data: user, isLoading } = useUser();
  const updateProfile = useUpdateProfile();
  const [formOverrides, setFormOverrides] = useState<Partial<{
    first_name: string;
    last_name: string;
    phone: string;
    avatar_url: string;
  }>>({});
  const [saved, setSaved] = useState(false);

  const form = {
    first_name: formOverrides.first_name ?? user?.profile?.first_name ?? '',
    last_name: formOverrides.last_name ?? user?.profile?.last_name ?? '',
    phone: formOverrides.phone ?? user?.profile?.phone ?? '',
    avatar_url: formOverrides.avatar_url ?? user?.profile?.avatar_url ?? '',
  };

  const submit = async () => {
    setSaved(false);
    await updateProfile.mutateAsync(form);
    setSaved(true);
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Profile Settings</h1>

      <Card>
        <CardHeader>
          <CardTitle>Personal Information</CardTitle>
          <CardDescription>Update your personal details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {isLoading ? (
            <p className="text-muted-foreground">Loading profile...</p>
          ) : (
            <>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="firstName">First Name</Label>
                  <Input
                    id="firstName"
                    value={form.first_name}
                    onChange={(e) =>
                      setFormOverrides((prev) => ({ ...prev, first_name: e.target.value }))
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="lastName">Last Name</Label>
                  <Input
                    id="lastName"
                    value={form.last_name}
                    onChange={(e) =>
                      setFormOverrides((prev) => ({ ...prev, last_name: e.target.value }))
                    }
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input id="email" type="email" value={user?.email ?? ''} disabled />
              </div>
              <div className="space-y-2">
                <Label htmlFor="phone">Phone</Label>
                <Input
                  id="phone"
                  type="tel"
                  value={form.phone}
                  onChange={(e) =>
                    setFormOverrides((prev) => ({ ...prev, phone: e.target.value }))
                  }
                />
              </div>
              <Button onClick={submit} disabled={updateProfile.isPending}>
                {updateProfile.isPending ? 'Saving...' : 'Save Changes'}
              </Button>
              {saved && <p className="text-sm text-green-600">Profile updated.</p>}
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
