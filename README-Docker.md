# Проект Kingsman - Настройка Docker

Этот проект содержит конфигурацию Docker для запуска стека приложений Kingsman, который включает:
- База данных PostgreSQL
- Go backend API
- React frontend
- pgAdmin (опциональный инструмент управления базой данных)

## Предварительные требования

- Docker (версия 20.10 или новее)
- Docker Compose (версия 2.0 или новее)

## Структура проекта

```
kingsman/
├── backend/
│   ├── app/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── db/migrations/
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── .env
│   │   └── .air.toml
│   ├── Dockerfile
│   └── .dockerignore
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── Dockerfile
│   ├── nginx.conf
│   └── .dockerignore
├── docker-compose.yml
├── docker-compose.dev.yml
└── README-Docker.md
```

## Быстрый старт

### Продакшн окружение

1. **Клонируйте репозиторий и перейдите в корневую папку проекта:**
   ```bash
   cd kingsman
   ```

2. **Соберите и запустите все сервисы:**
   ```bash
   docker-compose up -d
   ```

3. **Доступ к приложению:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - pgAdmin: http://localhost:5050 (admin@kingsman.com / admin123)

4. **Остановка сервисов:**
   ```bash
   docker-compose down
   ```

### Окружение разработки (с горячей перезагрузкой)

1. **Запуск среды разработки:**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. **Просмотр логов:**
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f
   ```

3. **Остановка среды разработки:**
   ```bash
   docker-compose -f docker-compose.dev.yml down
   ```

## Сервисы

### База данных (PostgreSQL)
- **Контейнер**: kingsman-db
- **Порт**: 5432
- **База данных**: kingsman
- **Пользователь**: postgres
- **Пароль**: a6fbnmod

### Backend (Go API)
- **Контейнер**: kingsman-backend
- **Порт**: 8080
- **Проверка здоровья**: GET http://localhost:8080/

### Frontend (React + Nginx)
- **Контейнер**: kingsman-frontend
- **Порт**: 3000 (продакшн), 3000 (разработка)
- **API Прокси**: запросы /api/* проксируются на backend

### pgAdmin (Управление базой данных)
- **Контейнер**: kingsman-pgadmin
- **Порт**: 5050
- **Email**: admin@kingsman.com
- **Пароль**: admin123

## Конфигурация окружения

### Переменные окружения Backend
Вы можете переопределить их в `backend/app/.env` или через Docker Compose:

```env
DB_HOST=database
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=a6fbnmod
DB_NAME=kingsman
DB_SSLMODE=disable
PORT=8080
```

### Переменные окружения Frontend
```env
REACT_APP_API_URL=http://localhost:8080
```

## Миграции базы данных

Миграции базы данных автоматически применяются при запуске контейнера PostgreSQL. Файлы миграций находятся в `backend/app/db/migrations/`.

## Полезные команды

### Просмотр логов
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f database
```

### Выполнение команд в контейнерах
```bash
# Контейнер backend
docker-compose exec backend sh

# Контейнер базы данных
docker-compose exec database psql -U postgres -d kingsman

# Контейнер frontend (разработка)
docker-compose -f docker-compose.dev.yml exec frontend sh
```

### Пересборка сервисов
```bash
# Пересборка всех сервисов
docker-compose build

# Пересборка конкретного сервиса
docker-compose build backend
docker-compose build frontend

# Пересборка и перезапуск
docker-compose up -d --build
```

### Очистка
```bash
# Остановка и удаление контейнеров, сетей
docker-compose down

# Удаление томов также (⚠️ Это удалит данные базы данных)
docker-compose down -v

# Удаление неиспользуемых образов, контейнеров и сетей
docker system prune -f
```

## Рабочий процесс разработки

### Для разработки Backend:
1. Используйте `docker-compose.dev.yml` для горячей перезагрузки
2. Изменяйте файлы в `backend/app/`
3. Изменения автоматически обнаруживаются и сервер перезапускается

### Для разработки Frontend:
1. Используйте `docker-compose.dev.yml` для горячей перезагрузки
2. Изменяйте файлы в `frontend/src/`
3. Изменения автоматически отражаются в браузере

## Устранение неполадок

### Частые проблемы

1. **Конфликты портов**: Убедитесь, что порты 3000, 5432, 8080 и 5050 не используются
2. **Проблемы с подключением к базе данных**: Дождитесь прохождения проверки здоровья базы данных
3. **Ошибки сборки**: Проверьте логи Docker и убедитесь, что все зависимости доступны

### Сброс базы данных
```bash
docker-compose down -v
docker-compose up -d database
# Дождитесь готовности базы данных, затем запустите другие сервисы
docker-compose up -d
```

### Проверка здоровья сервисов
```bash
docker-compose ps
```

## Развертывание в продакшн

Для развертывания в продакшн рассмотрите:

1. **Использование конфигурационных файлов для конкретной среды**
2. **Настройка правильного управления секретами**
3. **Настройка обратного прокси (например, Nginx) для терминации SSL**
4. **Настройка мониторинга и логирования**
5. **Использование Docker Swarm или Kubernetes для оркестрации**

## Заметки о безопасности

- Измените пароли по умолчанию в продакшн
- Используйте переменные окружения для чувствительных данных
- Рассмотрите использование Docker secrets для продакшн развертываний
- Включите SSL/TLS для продакшн окружений