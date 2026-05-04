ECommerce API — README

Features:

User authentication (JWT) (DOESN'T WORK)

Product CRUD

Shopping cart service (DOESN'T WORK)

Order creation (DOESN'T WORK)

PostgreSQL database

Dockerized environment

Swagger documentation

Unit tests with coverage

Requirements:

Go 1.22+

Docker & Docker Compose

Git

How to Start the Program (Docker Recommended):

Run: docker compose --profile dev up --build

Stop: docker compose down
API URL: http://localhost:8080
Swagger Docs: http://localhost:8080/swagger/index.html



Run migrations:
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable" up

Generate Swagger docs:
swag init -g cmd/api/main.go


Running Tests:

All tests: go test ./... -cover


Project Structure:
ecommerce/
cmd/api/ (entrypoint)
docs/ (swagger)
internal/
db/ (database connection)
domain/
entity/
repository/
service/
handler/http/
repository/
service/
migrations/
docker-compose.yml
Dockerfile
go.mod
README.md

Environment Variables (.env):
PORT=8080
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce
JWT_SECRET=your-secret-key

API Documentation:
http://localhost:8080/swagger/index.html
