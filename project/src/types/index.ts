export interface Item {
  id: number;
  name: string;
  quantity: number;
  price: number;
}

export interface OrderItem {
  item_id: number;
  name: string;
  quantity: number;
  price: number;
}

export interface Order {
  id: number;
  items: OrderItem[];
  total: number;
  created_unix: number;
}

export interface CreateItemRequest {
  name: string;
  quantity: number;
  price: number;
}

export interface AdjustQuantityRequest {
  delta: number;
}

export interface CreateOrderRequest {
  items: Array<{
    id: number;
    quantity: number;
  }>;
}

export interface ApiError {
  error: string;
}
