# 🗒️ Note Manager API

Полноценное RESTful-приложение на Go для управления заметками с поддержкой авторизации (JWT), избранных, кеширования через Redis, миграций базы данных и сборки через Docker.

---

A full-fledged RSTFUL application on the control of notes with authorization support (JWT), selected, caching through Redis, database migrations and assembly through Docker.

--

## Для старта// Setup

- Docker & Docker Compose
- Go >= 1.20

## Code structure

![alt text](image.png)
## ⚙️ Установка и запуск

1. **Клонируйте репозиторий:**
bash
git clone https://github.com/your-username/note-manager-api.git
cd note-manager-api
2. Запустите всё через Docker Compose:

bash
clone https://github.com/your-username/note-manager-api.git
cd note-manager-api

2.Запусти всё через Docker Compose:
docker-compose up --build

3.Откройте в браузере:
http://localhost:8484

## ⚙️ Installation and launch

1. ** Clon the repository: **
bash
GIT CLONE https://github.com/your-username/note-manager-opi.git
CD Note-Manager-API
2. Run everything through the Docker Compose:

Bash
clone https://github.com/your-username/note-manager-api.git
CD Note-Manager-API

2. Let’s go all through Docker Compose:
Docker-Compose up-Build

3. Roll up in the browser:
http: // localhost: 8484

API Endpoints
| Метод  | URL                  | Описание                     |
| ------ | -------------------- | ---------------------------- |
| POST   | `/register`          | Регистрация пользователя / User registration    |
| POST   | `/login`             | Логин и получение JWT   /login and get JWT     |
| GET    | `/api/notes`         | Получить все заметки  / Get all note       |
| POST   | `/api/notes`         | Создать новую заметку  / Сreate New note     |
| PUT    | `/api/notes/:id`     | Обновить заметку по ID / Update by ID     |
| DELETE | `/api/notes/:id`     | Удалить заметку по ID / Delete by ID       |
| GET    | `/api/favorites`     | Получить избранные заметки / Get favorites notes  |
| POST   | `/api/favorites/:id` | Добавить заметку в избранное / Add to favorites |
| DELETE | `/api/favorites/:id` | Удалить из избранного / Delete from favorites       |

Стек технологий
Go (Gin, pgx/v5, bcrypt, JWT)

PostgreSQL (через pgx и goose)

Redis (для кеширования избранных)

Docker (приложение + PostgreSQL + Redis)

HTML интерфейс (без фреймворков, чистый JS)

Stack of technology
GO (GIN, PGX/V5, BCRYPT, JWT)

PostgreSQL (via PGX and Goose)

Redis (for caching the chosen ones)

Docker (Appendix + PostgreSQL + Redis)

HTML interface (without frameworks, pure JS)

## Примеры запуска без Docker
go mod tidy
go run cmd/main.go
Но вы должны настроить PostgreSQL и Redis вручную.

## Examples of running without Docker
go mod tidy
go run cmd/main.go
But you must configure PostgreSQL and Redis manually.

## Переменные окружения 
DATABASE_URL=postgres://user:password@note_db:5432/noteManager?sslmode=disable
REDIS_ADDR=note_manager_api_redis:6379

## Environment variables
DATABASE_URL=postgres://user:password@note_db:5432/noteManager?sslmode=disable
REDIS_ADDR=note_manager_api_redis:6379

## Миграции базы данных
Миграции лежат в migrations/ и применяются автоматически при запуске приложения через pressly/goose.

Формат имени файла:
YYYYMMDDHHMMSS_description.sql

Пример команды вручную:
goose -dir migrations postgres "$DATABASE_URL" up

## Database migrations
Migrations are stored in migrations/ and are applied automatically when the application is launched via pressly/goose.

File name format:
YYYYMMDDHHMMSS_description.sql

Example manual command:
goose -dir migrations postgres "$DATABASE_URL" up

Автор/Author
GitHub: @zifeldev

