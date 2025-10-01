
## Setup Guide
###  1. Generate Go code from SQL queries

```bash
sqlc generate
```

###  2. Generate Swagger Documentation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```
###  3. Start services with Docker Compose

```bash
docker-compose up -d

```
###  4. Access Swagger UI

```bash
http://localhost:8080/swagger/index.html

```
