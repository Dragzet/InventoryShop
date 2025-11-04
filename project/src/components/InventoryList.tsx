import { useState, useEffect } from 'react';
import { Package, Plus, Minus, RefreshCw, Loader2 } from 'lucide-react';
import type { Item } from '../types';
import { api } from '../services/api.ts';

interface InventoryListProps {
  onItemsChange?: () => void;
}

export function InventoryList({ onItemsChange }: InventoryListProps) {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [adjusting, setAdjusting] = useState<number | null>(null);

  const loadItems = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getItems();
      setItems(data);
      onItemsChange?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load items');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadItems();
  }, []);

  const handleAdjust = async (id: number, delta: number) => {
    try {
      setAdjusting(id);
      await api.adjustItemQuantity(id, delta);
      await loadItems();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to adjust quantity');
    } finally {
      setAdjusting(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-800">{error}</p>
        <button
          onClick={loadItems}
          className="mt-2 text-sm text-red-600 hover:text-red-700 font-medium"
        >
          Try again
        </button>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <Package className="w-16 h-16 mx-auto text-gray-400 mb-4" />
        <p className="text-gray-600 mb-4">No items in inventory yet</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Current Stock</h3>
        <button
          onClick={loadItems}
          className="text-sm text-blue-600 hover:text-blue-700 flex items-center gap-1"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {items.map((item) => (
          <div
            key={item.id}
            className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between mb-3">
              <div>
                <h4 className="font-semibold text-gray-900">{item.name}</h4>
                <p className="text-2xl font-bold text-blue-600 mt-1">
                  ${item.price.toFixed(2)}
                </p>
              </div>
              <div
                className={`px-3 py-1 rounded-full text-sm font-medium ${
                  item.quantity === 0
                    ? 'bg-red-100 text-red-800'
                    : item.quantity < 10
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-green-100 text-green-800'
                }`}
              >
                {item.quantity} in stock
              </div>
            </div>

            <div className="flex gap-2 mt-4">
              <button
                onClick={() => handleAdjust(item.id, -1)}
                disabled={adjusting === item.id || item.quantity === 0}
                className="flex-1 flex items-center justify-center gap-1 px-3 py-2 bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed text-gray-700 rounded-lg transition-colors"
              >
                {adjusting === item.id ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <Minus className="w-4 h-4" />
                    <span className="text-sm font-medium">Remove</span>
                  </>
                )}
              </button>
              <button
                onClick={() => handleAdjust(item.id, 1)}
                disabled={adjusting === item.id}
                className="flex-1 flex items-center justify-center gap-1 px-3 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {adjusting === item.id ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <Plus className="w-4 h-4" />
                    <span className="text-sm font-medium">Add</span>
                  </>
                )}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
