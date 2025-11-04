import type {
  Item,
  Order,
  CreateItemRequest,
  AdjustQuantityRequest,
  CreateOrderRequest,
  ApiError,
} from '../types';

const INVENTORY_API_URL = import.meta.env.VITE_INVENTORY_API_URL || 'http://localhost:8001';
const ORDERS_API_URL = import.meta.env.VITE_ORDERS_API_URL || 'http://localhost:8002';

class ApiService {
  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }
    return response.json();
  }

  async getItems(): Promise<Item[]> {
    const response = await fetch(`${INVENTORY_API_URL}/items`);
    return this.handleResponse<Item[]>(response);
  }

  async getItem(id: number): Promise<Item> {
    const response = await fetch(`${INVENTORY_API_URL}/items/${id}`);
    return this.handleResponse<Item>(response);
  }

  async createItem(data: CreateItemRequest): Promise<Item> {
    const response = await fetch(`${INVENTORY_API_URL}/items`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return this.handleResponse<Item>(response);
  }

  async adjustItemQuantity(id: number, delta: number): Promise<Item> {
    const response = await fetch(`${INVENTORY_API_URL}/items/${id}/adjust`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ delta } as AdjustQuantityRequest),
    });
    return this.handleResponse<Item>(response);
  }

  async getOrders(): Promise<Order[]> {
    const response = await fetch(`${ORDERS_API_URL}/orders`);
    return this.handleResponse<Order[]>(response);
  }

  async getOrder(id: number): Promise<Order> {
    const response = await fetch(`${ORDERS_API_URL}/orders/${id}`);
    return this.handleResponse<Order>(response);
  }

  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const response = await fetch(`${ORDERS_API_URL}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return this.handleResponse<Order>(response);
  }
}

export const api = new ApiService();
