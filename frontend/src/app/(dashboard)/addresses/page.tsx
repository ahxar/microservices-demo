'use client';

import { useState } from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAddAddress, useAddresses } from '@/hooks/use-user';

export default function AddressesPage() {
  const { data: addresses = [], isLoading } = useAddresses();
  const addAddress = useAddAddress();
  const [form, setForm] = useState({
    label: 'Home',
    street: '',
    city: '',
    state: '',
    zip_code: '',
    country: 'USA',
    is_default: false,
  });
  const [showForm, setShowForm] = useState(false);

  const submit = async () => {
    await addAddress.mutateAsync({
      label: form.label,
      is_default: form.is_default,
      address: {
        street: form.street,
        city: form.city,
        state: form.state,
        zip_code: form.zip_code,
        country: form.country,
      },
    });
    setShowForm(false);
    setForm({
      label: 'Home',
      street: '',
      city: '',
      state: '',
      zip_code: '',
      country: 'USA',
      is_default: false,
    });
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Shipping Addresses</h1>
        <Button onClick={() => setShowForm((prev) => !prev)}>
          <Plus className="h-4 w-4 mr-2" />
          {showForm ? 'Cancel' : 'Add Address'}
        </Button>
      </div>

      {showForm && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>New Address</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="label">Label</Label>
                <Input
                  id="label"
                  value={form.label}
                  onChange={(e) => setForm((prev) => ({ ...prev, label: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="street">Street</Label>
                <Input
                  id="street"
                  value={form.street}
                  onChange={(e) => setForm((prev) => ({ ...prev, street: e.target.value }))}
                />
              </div>
            </div>
            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="city">City</Label>
                <Input
                  id="city"
                  value={form.city}
                  onChange={(e) => setForm((prev) => ({ ...prev, city: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="state">State</Label>
                <Input
                  id="state"
                  value={form.state}
                  onChange={(e) => setForm((prev) => ({ ...prev, state: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="zip">ZIP</Label>
                <Input
                  id="zip"
                  value={form.zip_code}
                  onChange={(e) => setForm((prev) => ({ ...prev, zip_code: e.target.value }))}
                />
              </div>
            </div>
            <Button onClick={submit} disabled={addAddress.isPending}>
              {addAddress.isPending ? 'Saving...' : 'Save Address'}
            </Button>
          </CardContent>
        </Card>
      )}

      {isLoading ? (
        <p className="text-muted-foreground">Loading addresses...</p>
      ) : addresses.length === 0 ? (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            No addresses saved.
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {addresses.map((address) => (
            <Card key={address.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>{address.label}</CardTitle>
                  {address.is_default && (
                    <span className="text-xs bg-primary text-primary-foreground px-2 py-1 rounded">
                      Default
                    </span>
                  )}
                </div>
              </CardHeader>
              <CardContent className="text-sm">
                <p>{address.address.street}</p>
                <p>
                  {address.address.city}, {address.address.state} {address.address.zip_code}
                </p>
                <p>{address.address.country}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
