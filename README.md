# 👑 MR.KINGSMAN

> Современная платформа для управления магазинами мужской одежды с программой лояльности и интегрированным каталогом товаров

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=flat&logo=react&logoColor=%2361DAFB)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=flat&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)

</div>

## 📋 Описание

**MR.KINGSMAN** — это полнофункциональная веб-платформа для управления сетью магазинов мужской одежды премиум-класса. Система предоставляет комплексное решение для управления каталогом товаров, обработки заказов, программы лояльности клиентов и аналитики продаж.

**Целевая аудитория:**
- 🏪 Владельцы и менеджеры розничных магазинов
- 👥 Клиенты магазинов мужской одежды
- 📊 Администраторы и аналитики

## ✨ Основные возможности

### Для клиентов
- 🛒 **Каталог товаров** — удобный поиск и фильтрация по категориям
- 🛍️ **Корзина покупок** — добавление товаров и управление заказами
- 💳 **Программа лояльности** — накопление баллов и получение скидок
- 📱 **QR-код профиля** — быстрая идентификация в магазине
- 🔔 **Уведомления** — персональные предложения и новости
- 📍 **Карта магазинов** — поиск ближайших точек продаж

### Для администраторов
- 📊 **Панель управления** — контроль продаж и аналитика
- 📦 **Управление товарами** — добавление, редактирование каталога
- 👥 **Управление клиентами** — просмотр профилей и истории покупок
- 🏪 **Управление магазинами** — информация о точках продаж
- 💰 **Финансовая отчетность** — анализ доходов и прибыли

## 🏗️ Архитектура проекта

```mermaid
graph TB
    subgraph "Frontend Layer"
        React[React SPA]
        Redux[Redux Store]
        Router[React Router]
    end
    
    subgraph "Backend Layer"
        API[REST API Server]
        Auth[JWT Authentication]
        Business[Business Logic]
    end
    
    subgraph "Data Layer"
        DB[(PostgreSQL)]
        Migrations[Database Migrations]
    end
    
    subgraph "Infrastructure"
        Docker[Docker Containers]
        Nginx[Nginx Proxy]
    end
    
    React --> API
    Redux --> React
    Router --> React
    API --> Auth
    API --> Business
    Business --> DB
    Migrations --> DB
    Docker --> API
    Docker --> React
    Docker --> DB
    Nginx --> Docker
```

### Компоненты системы

| Компонент | Технология | Назначение |
|-----------|------------|------------|
| **Frontend** | React 19+ | Пользовательский интерфейс |
| **Backend API** | Go 1.23+ | REST API и бизнес-логика |
| **База данных** | PostgreSQL 15 | Хранение данных |
| **Прокси** | Nginx | Маршрутизация запросов |
| **Контейнеризация** | Docker | Развертывание и масштабирование |

## 🛠️ Технологии

<table>
<tr>
<td width="50%">

### Backend
- **Go 1.23.2** — основной язык
- **Chi Router** — HTTP маршрутизация
- **PostgreSQL** — база данных
- **JWT** — аутентификация
- **bcrypt** — хеширование паролей
- **golang-migrate** — миграции БД

</td>
<td width="50%">

### Frontend
- **React 19.1.1** — UI фреймворк
- **Redux Toolkit** — управление состоянием
- **React Router** — маршрутизация
- **Axios** — HTTP клиент
- **React Icons** — иконки
- **CSS3** — стилизация (без Tailwind)

</td>
</tr>
</table>

### DevOps & Инструменты
- 🐳 **Docker & Docker Compose** — контейнеризация
- 🔄 **Air** — hot reload для Go
- 🗃️ **pgAdmin** — управление БД
- 📝 **Golang Migrate** — миграции
- 🔧 **Nginx** — веб-сервер и прокси

## 🚀 Установка и запуск

### Предварительные требования

- **Docker** версии 20.10+
- **Docker Compose** версии 2.0+
- **Git** для клонирования репозитория

### Быстрый старт

1. **Клонирование репозитория**
   ```bash
   git clone https://github.com/username/kingsman.git
   cd kingsman
   ```

2. **Запуск в продакшн режиме**
   ```bash
   docker-compose up -d
   ```

3. **Доступ к приложению**
   - 🌐 **Frontend**: http://localhost:3000
   - 🔧 **Backend API**: http://localhost:8080
   - 🗃️ **pgAdmin**: http://localhost:5050
     - Email: `admin@kingsman.com`
     - Password: `admin123`

### Разработка

1. **Запуск в режиме разработки** (с hot reload)
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. **Просмотр логов**
   ```bash
   docker-compose logs -f
   ```

3. **Остановка сервисов**
   ```bash
   docker-compose down
   ```

### Конфигурация

Настройки находятся в файле `backend/app/.env`:

```env
# Database
DB_HOST=database
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=a6fbnmod
DB_NAME=kingsman
DB_SSLMODE=disable

# Server
PORT=8080
JWT_SECRET=your_secret_key
```

## 📖 Примеры использования

### REST API Примеры

#### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "Иван",
    "last_name": "Петров"
  }'
```

#### Авторизация
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Получение каталога товаров
```bash
curl -X GET http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Добавление товара в корзину
```bash
curl -X POST http://localhost:8080/cart/add \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Интерфейс пользователя

#### Мобильная навигация
- 📱 **Бургер-меню** — расположено в левом верхнем углу
- 🎯 **Выдвижное меню** — плавная анимация слева
- 🔄 **Адаптивный дизайн** — корректная работа на всех устройствах

#### Программа лояльности
- 💎 **Система уровней** — Bronze, Silver, Gold, Platinum
- 🎁 **Накопление баллов** — за каждую покупку
- 🏷️ **Скидочная карта** — персональная карта клиента
- 📱 **QR-код** — для быстрой идентификации

## 🖼️ Скриншоты

<div align="center">

### Главная страница
![Главная страница](docs/images/homepage.png)

### Каталог товаров
![Каталог](docs/images/catalog.png)

### Мобильное меню
![Мобильное меню](docs/images/mobile-menu.png)

### Профиль пользователя
![Профиль](docs/images/profile.png)

</div>

## 🗄️ Структура базы данных

<details>
<summary>Основные таблицы</summary>

```sql
-- Пользователи
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Товары
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Корзина
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    is_paid BOOLEAN DEFAULT FALSE
);

-- Программа лояльности
CREATE TABLE loyalty_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    tier VARCHAR(20) DEFAULT 'Bronze',
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

</details>

## 🤝 Контрибьюция

Мы приветствуем участие в развитии проекта! Вот как вы можете помочь:

### Как внести свой вклад

1. **Fork** проекта
2. Создайте **feature branch** (`git checkout -b feature/AmazingFeature`)
3. **Commit** ваши изменения (`git commit -m 'Add some AmazingFeature'`)
4. **Push** в branch (`git push origin feature/AmazingFeature`)
5. Откройте **Pull Request**

### Правила разработки

- 📝 **Код-стиль**: следуйте принятым в проекте конвенциям
- ✅ **Тестирование**: добавляйте тесты для новой функциональности
- 📖 **Документация**: обновляйте документацию при необходимости
- 🔄 **Маленькие PR**: предпочитайте небольшие, фокусированные изменения

### Отчет об ошибках

Если вы нашли ошибку, пожалуйста:

1. Проверьте, что подобная ошибка еще не была зарегистрирована
2. Создайте новый **Issue** с подробным описанием
3. Приложите шаги для воспроизведения
4. Укажите версию и окружение

## 🛡️ Безопасность

- 🔐 **JWT токены** для аутентификации
- 🔒 **bcrypt** для хеширования паролей
- 🛡️ **CORS** защита настроена
- 🔑 **Переменные окружения** для секретов
- 📝 **Валидация входных данных**

## 📄 Лицензия

Этот проект лицензирован под **MIT License** - см. файл [LICENSE](LICENSE) для деталей.

```
MIT License

Copyright (c) 2024 MR.KINGSMAN

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

<div align="center">

**[⬆ Наверх](#-mrkingsman)**

Сделано с ❤️ командой MR.KINGSMAN

[![GitHub](https://img.shields.io/badge/GitHub-Repository-black?style=flat&logo=github)](https://github.com/username/kingsman)
[![Documentation](https://img.shields.io/badge/Docs-Documentation-blue?style=flat&logo=gitbook)](docs/)
[![API Reference](https://img.shields.io/badge/API-Reference-green?style=flat&logo=swagger)](docs/api.md)

</div>