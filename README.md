Проект: Система управления инвентарем и заказами (микросервисы на Go)

Структура:
- services/inventory - сервис управления товарами (порт 8001)
- services/orders - сервис управления заказами (порт 8002)
- cmd/client - простой CLI-клиент для тестирования
- docker-compose.yml - для запуска сервисов

Запуск локально:
1) Для каждого сервиса выполнить go run main.go в соответствующей папке.
   Пример:
     cd services/inventory && go run main.go
     cd services/orders && go run main.go
2) CLI-клиент:
     cd cmd/client && go run main.go list | create | order

Описание API см. в коде сервисов.

## Архитектура проекта

```mermaid
graph TD
    %% Frontend
    subgraph Frontend [Frontend]
        Client([Client Application<br/>React + TypeScript])
    end

    %% Inventory Service
    subgraph InventoryContainer [Inventory Service Container]
        direction TB
        InvRouter[Inventory Router]
        InvStore[Inventory Store]
    end

    %% Orders Service
    subgraph OrdersContainer [Orders Service Container]
        direction TB
        OrdRouter[Orders Router]
        OrdStore[Orders Store]
        InvClient[Inventory Client]
    end

    %% Data Layer
    subgraph DataLayer [Data Layer]
        DB_Inv[(Inventory DB<br/>PostgreSQL)]
        DB_Ord[(Orders DB<br/>PostgreSQL)]
    end

    %% Relationships
    Client -.->|HTTP JSON<br/>GET/POST /items| InvRouter
    Client -.->|HTTP JSON<br/>GET/POST /orders| OrdRouter

    InvRouter --> InvStore
    InvStore -->|SQL| DB_Inv

    OrdRouter --> OrdStore
    OrdRouter --> InvClient
    OrdStore -->|SQL| DB_Ord

    InvClient -.->|"HTTP GET /items/{id}"| InvRouter
```

## Диаграмма последовательности (Создание заказа)

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Orders as Orders Service
    participant Inventory as Inventory Service
    participant DB_Ord as Orders DB
    participant DB_Inv as Inventory DB

    Client->>Orders: POST /orders (items list)
    activate Orders

    loop For each item in request
        Orders->>Inventory: GET /items/{id}
        activate Inventory
        Inventory->>DB_Inv: SELECT * FROM items WHERE id=?
        DB_Inv-->>Inventory: Item Details
        Inventory-->>Orders: 200 OK (Price, Name, Stock)
        deactivate Inventory

        Orders->>Inventory: POST /items/{id}/adjust (delta: -qty)
        activate Inventory
        Inventory->>DB_Inv: UPDATE items SET quantity = quantity - qty
        DB_Inv-->>Inventory: Success/Fail
        Inventory-->>Orders: 200 OK (Stock Updated)
        deactivate Inventory
    end

    Note over Orders: If any step fails, rollback previous items

    Orders->>DB_Ord: INSERT INTO orders ...
    activate DB_Ord
    DB_Ord-->>Orders: Order Created (ID)
    deactivate DB_Ord

    Orders-->>Client: 201 Created (Order JSON)
    deactivate Orders
```

## Диаграмма последовательности (Получение списка товаров)

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Inventory as Inventory Service
    participant DB_Inv as Inventory DB

    Client->>Inventory: GET /items
    activate Inventory
    Inventory->>DB_Inv: SELECT * FROM items
    activate DB_Inv
    DB_Inv-->>Inventory: List of Items (Rows)
    deactivate DB_Inv
    Inventory-->>Client: 200 OK (JSON List)
    deactivate Inventory
```

## Сравнение архитектурных подходов

| Критерий | Монолитная архитектура | Многоуровневая архитектура | **Микросервисная архитектура (✅ Выбор проекта)** |
| :--- | :--- | :--- | :--- |
| **Масштабируемость** | **Вертикальное.** Ориентирована на увеличение ресурсов сервера. Горизонтальное расширение затруднено. | **Гибкое.** Уровни приложения могут масштабироваться независимо от уровня хранения данных. | **Идеальная.** Позволяет точечно масштабировать только те сервисы, которые испытывают нагрузку, экономя ресурсы и повышая эффективность. |
| **Производительность** | **Высокая** (для малых систем). Нет сетевых задержек внутри процесса. | **Средняя.** Возможны небольшие задержки между слоями. | **Высокая под нагрузкой.** Возможность параллельной обработки запросов и использования специализированных быстрых языков (Go) для критических узлов компенсирует сетевые расходы. |
| **Простота разработки и сопровождения** | **Высокая** (на старте). Единая кодовая база. | **Удобная.** Чёткое разделение ответственности. | **Модульная.** Система разбита на небольшие, понятные части. Легче вносить изменения в один сервис, не затрагивая остальные. Упрощает параллельную разработку командами. |
