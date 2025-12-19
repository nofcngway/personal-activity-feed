# Personal Activity Feed

Система персональной ленты активности, построенная на микросервисной архитектуре с использованием Go, Kafka, Redis и PostgreSQL.

## Архитектура

Проект состоит из двух микросервисов:

1.  **`auth-action-service`**
    *   **Роль**: Входная точка (Gateway) для клиентов.
    *   **Функции**: Регистрация, аутентификация, публикация действий (лайки, посты, подписки).
    *   **Технологии**: HTTP/gRPC API, Redis (сессии), Kafka Producer.

2.  **`feed-service`**
    *   **Роль**: Сервис формирования и отдачи ленты.
    *   **Функции**: Чтение событий, сохранение в БД, выдача ленты пользователю.
    *   **Технологии**: Kafka Consumer, PostgreSQL, gRPC API.

## Быстрый старт

### 1. Подготовка инфраструктуры

В корне проекта выполните команду для запуска Docker-контейнеров (Postgres, Redis, Kafka, Zookeeper):

```bash
make up
```

> **Важно**: Контейнерный Postgres доступен на порту **`55432`**, чтобы не конфликтовать с локальным экземпляром на 5432.

### 2. Запуск сервисов

Сервисы необходимо запускать в отдельных терминалах.

**Терминал 1: auth-action-service**

```bash
cd auth-action-service
make run
```

**Терминал 2: feed-service**

```bash
cd feed-service
make run
```

## API и Документация

У обоих сервисов есть встроенный Swagger UI для просмотра и тестирования API.

### Auth & Actions
*   **URL**: [http://localhost:8082/docs/](http://localhost:8082/docs/)
*   **Эндпоинты**:
    *   `POST /register`, `POST /login` — Получение токена.
    *   `POST /posts`, `POST /like/{id}`, `POST /follow/{id}` — Действия (требуют `Authorization`).

### Feed (Лента)
*   **URL**: [http://localhost:8083/docs/](http://localhost:8083/docs/)
*   **Эндпоинты**:
    *   `GET /feed` — Получение ленты.