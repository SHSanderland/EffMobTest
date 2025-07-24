# EffMobTest
## Описание
REST-сервис для агрегации данных об онлайн-подписках пользователей с возможностью CRUDL-операций и расчета суммарной стоимости подписок.

Проект использует:
- **Golang** (1.24.5)
- **chi** для маршрутизации.
- **PostgreSQL** для хранения данных.


## Установка и запуск

### 1. Запуск сервиса (необходима готовая база данных):
```bash
make upServer
```
### 2. Запуск контейнера:
```bash
make upContainer
```

### 3. Документация API:
Откройте [http://localhost:8080/swagger/](http://localhost:8080/swagger/) для просмотра Swagger-документации.

## Зависимости
 - github.com/go-chi/chi/v5 v5.2.2
 - github.com/golang-migrate/migrate/v4 v4.18.3
 - github.com/google/uuid v1.6.0
 - github.com/ilyakaznacheev/cleanenv v1.5.0
 - github.com/jackc/pgx/v5 v5.7.5
 - github.com/swaggo/http-swagger/v2 v2.0.2
 - github.com/swaggo/swag v1.16.5