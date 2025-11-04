import { useState, useEffect } from 'react';
import { ShoppingCart, Plus, Trash2, Loader2 } from 'lucide-react';
import type { Item } from '../types';
import { api } from '../services/api.ts';

interface OrderItemInput {
  id: number;
  quantity: number;
}

interface CreateOrderFormProps {
  onOrderCreated: () => void;
}

export function CreateOrderForm({ onOrderCreated }: CreateOrderFormProps) {
  const [items, setItems] = useState<Item[]>([]);
  const [orderItems, setOrderItems] = useState<OrderItemInput[]>([{ id: 0, quantity: 1 }]);
  const [loading, setLoading] = useState(false);
  const [loadingItems, setLoadingItems] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadItems();
  }, []);

  const loadItems = async () => {
    try {
      setLoadingItems(true);
      const data = await api.getItems();
      setItems(data.filter((item) => item.quantity > 0));
    } catch (err) {
      setError('Ошибка при загрузке доступных товаров');
    } finally {
      setLoadingItems(false);
    }
  };

  const addOrderItem = () => {
    setOrderItems([...orderItems, { id: 0, quantity: 1 }]);
  };

  const removeOrderItem = (index: number) => {
    setOrderItems(orderItems.filter((_, i) => i !== index));
  };

  const updateOrderItem = (index: number, field: keyof OrderItemInput, value: number) => {
    const updated = [...orderItems];
    updated[index] = { ...updated[index], [field]: value };
    setOrderItems(updated);
  };

  const calculateTotal = () => {
    return orderItems.reduce((total, orderItem) => {
      const item = items.find((i) => i.id === orderItem.id);
      if (item && orderItem.quantity > 0) {
        return total + item.price * orderItem.quantity;
      }
      return total;
    }, 0);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validItems = orderItems.filter((item) => item.id > 0 && item.quantity > 0);

    if (validItems.length === 0) {
      setError('Добавьте хотя бы один товар в заказ');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await api.createOrder({ items: validItems });
      setOrderItems([{ id: 0, quantity: 1 }]);
      onOrderCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось создать заказ');
    } finally {
      setLoading(false);
    }
  };

  if (loadingItems) {
    return (
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <div className="flex items-center justify-center py-8">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
        </div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <div className="text-center py-8">
          <ShoppingCart className="w-12 h-12 mx-auto text-gray-400 mb-3" />
          <p className="text-gray-600">Нет доступных товаров для заказа</p>
          <p className="text-sm text-gray-500 mt-1">Сначала добавьте товары в инвентарь</p>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-lg p-6">
      <div className="flex items-center gap-2 mb-4">
        <ShoppingCart className="w-5 h-5 text-blue-600" />
        <h3 className="text-lg font-semibold text-gray-900">Создать заказ</h3>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 rounded-lg p-3">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      <div className="space-y-3 mb-4">
        {orderItems.map((orderItem, index) => (
          <div key={index} className="flex gap-2 items-center">
            <div className="flex-1">
              <select
                value={orderItem.id}
                onChange={(e) => updateOrderItem(index, 'id', parseInt(e.target.value, 10))}
                className="w-full h-12 px-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition appearance-none"
                disabled={loading}
              >
                <option value={0}>Select item...</option>
                {items.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.name} - ${item.price.toFixed(2)} ({item.quantity} available)
                  </option>
                ))}
              </select>
            </div>

            <div className="w-24">
              <input
                type="number"
                value={orderItem.quantity}
                onChange={(e) => updateOrderItem(index, 'quantity', parseInt(e.target.value, 10) || 0)}
                min={1}
                max={items.find((i) => i.id === orderItem.id)?.quantity || 1}
                className="w-full h-12 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition text-center"
                disabled={loading}
              />
            </div>

            <div>
              <button
                type="button"
                onClick={() => removeOrderItem(index)}
                disabled={orderItems.length === 1 || loading}
                className="w-10 h-10 flex items-center justify-center text-red-600 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors"
                aria-label="Remove item"
              >
                <Trash2 className="w-5 h-5" />
              </button>
            </div>
          </div>
        ))}
      </div>

      <button
        type="button"
        onClick={addOrderItem}
        disabled={loading}
        className="w-full mb-4 px-4 py-2 border-2 border-dashed border-gray-300 hover:border-blue-500 hover:bg-blue-50 disabled:opacity-50 disabled:cursor-not-allowed text-gray-600 hover:text-blue-600 rounded-lg transition-colors flex items-center justify-center gap-2"
      >
        <Plus className="w-5 h-5" />
        Добавить ещё товар
      </button>

      <div className="border-t border-gray-200 pt-4 mb-4">
        <div className="flex justify-between items-center">
          <span className="text-lg font-semibold text-gray-900">Предполагаемая сумма:</span>
          <span className="text-2xl font-bold text-blue-600">${calculateTotal().toFixed(2)}</span>
        </div>
        <p className="text-xs text-gray-500 mt-1">Окончательная сумма будет рассчитана при подтверждении заказа</p>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium py-2.5 px-4 rounded-lg transition-colors flex items-center justify-center gap-2"
      >
        {loading ? (
          <>
            <Loader2 className="w-5 h-5 animate-spin" />
            Создание...
          </>
        ) : (
          <>
            <ShoppingCart className="w-5 h-5" />
            Создать заказ
          </>
        )}
      </button>
    </form>
  );
}
