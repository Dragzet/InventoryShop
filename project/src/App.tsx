import { useState } from 'react';
import { Package, ShoppingCart, Store } from 'lucide-react';
import { InventoryList } from './components/InventoryList.tsx';
import { AddItemForm } from './components/AddItemForm.tsx';
import { CreateOrderForm } from './components/CreateOrderForm.tsx';
import { OrdersList } from './components/OrdersList.tsx';

type Tab = 'inventory' | 'orders';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('inventory');
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleItemChange = () => {
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-gray-50 to-gray-100">
      <header className="bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center gap-3">
            <div className="bg-blue-600 p-2.5 rounded-lg">
              <Store className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Store Manager</h1>
              <p className="text-sm text-gray-600">Inventory & Order Management System</p>
            </div>
          </div>
        </div>
      </header>

      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-1">
            <button
              onClick={() => setActiveTab('inventory')}
              className={`flex items-center gap-2 px-6 py-4 font-medium transition-colors relative ${
                activeTab === 'inventory'
                  ? 'text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <Package className="w-5 h-5" />
              Inventory
            </button>
            <button
              onClick={() => setActiveTab('orders')}
              className={`flex items-center gap-2 px-6 py-4 font-medium transition-colors relative ${
                activeTab === 'orders'
                  ? 'text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <ShoppingCart className="w-5 h-5" />
              Orders
            </button>
          </div>
        </div>
      </nav>

      <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'inventory' && (
          <div className="grid lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <InventoryList onItemsChange={handleItemChange} />
            </div>
            <div>
              <AddItemForm onItemAdded={handleItemChange} />
            </div>
          </div>
        )}

        {activeTab === 'orders' && (
          <div className="grid lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <OrdersList refreshTrigger={refreshTrigger} />
            </div>
            <div>
              <CreateOrderForm onOrderCreated={handleItemChange} />
            </div>
          </div>
        )}
      </main>

      <footer className="mt-12 border-t border-gray-200 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <p className="text-center text-sm text-gray-600">
            Connected to Inventory API: {import.meta.env.VITE_INVENTORY_API_URL || 'http://localhost:8001'} |
            Orders API: {import.meta.env.VITE_ORDERS_API_URL || 'http://localhost:8002'}
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
