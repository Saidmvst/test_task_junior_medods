# test_task_junior_medods

Решение тестового задания: Сервис управления задачами
Что реализовано:

Логика периодических задач: В слое usecase (сервис) реализован механизм генерации задач на 30 дней вперед в зависимости от выбранного интервала (daily, weekly, monthly).

Работа с БД: Настроено сохранение параметров периодичности в колонку типа JSONB в PostgreSQL.

Современный стек:

Использование стандартного логгера slog (Structured Logging).

Чистая архитектура (разбиение на transport, usecase, repository).

Конфигурация через переменные окружения.

Graceful Shutdown: Реализовано корректное завершение работы сервера при получении сигналов SIGINT/SIGTERM.

Как запустить:

Поднять окружение:

Bash
docker-compose up -d --build
Создать таблицу (если не применились миграции):

Bash
docker-compose exec postgres psql -U postgres -d taskservice -c "
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    recurrence JSONB
);"
Пример запроса:

Bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "title": "Регулярная тренировка",
    "status": "new",
    "recurrence": {
        "type": "daily",
        "interval": 1
    }
}'
