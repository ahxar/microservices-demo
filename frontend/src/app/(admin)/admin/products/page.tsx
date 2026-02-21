'use client';

import { useMemo, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Plus, Trash } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useCategories, useProducts } from '@/hooks/use-products';
import { adminApi } from '@/lib/api';

function formatMoney(amountCents?: number, currency = 'USD') {
  const amount = (amountCents ?? 0) / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

export default function AdminProductsPage() {
  const queryClient = useQueryClient();
  const { data: productsData, isLoading } = useProducts({ page: 1, page_size: 50 });
  const { data: categories = [] } = useCategories();
  const [showForm, setShowForm] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    category_id: '',
    price_cents: '1999',
    stock_quantity: '10',
  });

  const products = productsData?.products ?? [];
  const defaultCategoryId = useMemo(() => categories[0]?.id ?? '', [categories]);

  const createProduct = async () => {
    if (!form.name || !form.slug || !form.description || !(form.category_id || defaultCategoryId)) {
      return;
    }

    setCreating(true);
    try {
      await adminApi.createProduct({
        name: form.name,
        slug: form.slug,
        description: form.description,
        category_id: form.category_id || defaultCategoryId,
        stock_quantity: Number(form.stock_quantity),
        price: {
          amount_cents: Number(form.price_cents),
          currency: 'USD',
        },
        image_urls: [],
      });
      setShowForm(false);
      setForm({
        name: '',
        slug: '',
        description: '',
        category_id: '',
        price_cents: '1999',
        stock_quantity: '10',
      });
      queryClient.invalidateQueries({ queryKey: ['products'] });
    } finally {
      setCreating(false);
    }
  };

  const removeProduct = async (id: string) => {
    await adminApi.deleteProduct(id);
    queryClient.invalidateQueries({ queryKey: ['products'] });
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Products</h1>
        <Button onClick={() => setShowForm((prev) => !prev)}>
          <Plus className="h-4 w-4 mr-2" />
          {showForm ? 'Cancel' : 'Add Product'}
        </Button>
      </div>

      {showForm && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Create Product</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Name</Label>
                <Input value={form.name} onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))} />
              </div>
              <div className="space-y-2">
                <Label>Slug</Label>
                <Input value={form.slug} onChange={(e) => setForm((prev) => ({ ...prev, slug: e.target.value }))} />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Description</Label>
              <Input
                value={form.description}
                onChange={(e) => setForm((prev) => ({ ...prev, description: e.target.value }))}
              />
            </div>
            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label>Category</Label>
                <select
                  className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={form.category_id}
                  onChange={(e) => setForm((prev) => ({ ...prev, category_id: e.target.value }))}
                >
                  <option value="">Select category</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-2">
                <Label>Price (cents)</Label>
                <Input
                  value={form.price_cents}
                  onChange={(e) => setForm((prev) => ({ ...prev, price_cents: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label>Stock</Label>
                <Input
                  value={form.stock_quantity}
                  onChange={(e) => setForm((prev) => ({ ...prev, stock_quantity: e.target.value }))}
                />
              </div>
            </div>
            <Button onClick={createProduct} disabled={creating}>
              {creating ? 'Creating...' : 'Create Product'}
            </Button>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">Product</th>
                  <th className="text-left p-4">Price</th>
                  <th className="text-left p-4">Stock</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-right p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td className="p-4 text-muted-foreground" colSpan={5}>Loading products...</td>
                  </tr>
                ) : (
                  products.map((product) => (
                    <tr key={product.id} className="border-b last:border-b-0">
                      <td className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-muted rounded-md" />
                          <span className="font-medium">{product.name}</span>
                        </div>
                      </td>
                      <td className="p-4">
                        {formatMoney(product.price?.amount_cents, product.price?.currency)}
                      </td>
                      <td className="p-4">{product.stock_quantity}</td>
                      <td className="p-4">
                        <span
                          className={`px-2 py-1 text-xs rounded-full ${
                            product.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-700'
                          }`}
                        >
                          {product.is_active ? 'Active' : 'Inactive'}
                        </span>
                      </td>
                      <td className="p-4">
                        <div className="flex justify-end space-x-2">
                          <Button variant="destructive" size="sm" onClick={() => removeProduct(product.id)}>
                            <Trash className="h-4 w-4" />
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
