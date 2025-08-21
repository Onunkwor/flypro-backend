DB_URL=postgres://postgres:adminpassword@localhost:5431/travel-db?sslmode=disable
MIGRATIONS_DIR=./migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) status
