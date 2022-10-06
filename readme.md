# URL Shortener

## What it uses

- Redis for storing tasks;
- [Asynq](https://github.com/hibiken/asynq) for enqueuing tasks;
- [Fiber](https://github.com/gofiber/fiber) for handling HTTP requests;

## Environment Variables

You need to create `.env` file for development and put these variables there (but with real values)

```bash
POSTGRES_PASSWORD=db_password
POSTGRES_USER=db_user
POSTGRES_DB=db_name
DATABASE_HOST=localhost

REDIS_HOST=localhost
REDIS_PORT=6379
```

For development file called - `.env`
For production - `.env.production`

## Start in dev mode

```bash
make dev
```

this will start database and redis. After that, you need to manually run `go run main.go`

## Build for production

```bash
make build
```

```bash
make start
```

## Routes

### GET `/api/urls`

Returns all urls from database

### POST `/api/shorten`

Creates new url

#### Example body

```json
{
  "url": "https://github.com",
  "name": "gh",
  "auto_delete": 30
}
```

`name` and `auto_delete` are optional.

`auto_delete` should be a number in minutes.

### 
