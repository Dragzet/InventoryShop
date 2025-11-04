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

