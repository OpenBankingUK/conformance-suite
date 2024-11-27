# Local Database Setup
./scripts/database
export DATABASE_URL="postgresql://postgres:passw0rd@localhost:7312/postgres?sslmode=disable"

# Run migrations
go run cmd/migrate/main.go -command up -migrations ./migrations

# Load fixtures
go run cmd/migrate/main.go -command fixture -fixtures ./fixtures

# Rollback
go run cmd/migrate/main.go -command down -migrations ./migrations