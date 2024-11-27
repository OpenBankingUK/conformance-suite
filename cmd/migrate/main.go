// cmd/migrate/main.go
package main

import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    command := flag.String("command", "up", "migrate, rollback, or fixture")
    dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "database connection string")
    migrationsDir := flag.String("migrations", "./migrations", "migrations directory")
    fixturesDir := flag.String("fixtures", "./fixtures", "fixtures directory")
    flag.Parse()

    db, err := sql.Open("postgres", *dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    switch *command {
    case "up", "down":
        runMigrations(db, *command, *migrationsDir)
    case "fixture":
        loadFixtures(db, *fixturesDir)
    default:
        log.Fatalf("Unknown command: %s", *command)
    }
}

func runMigrations(db *sql.DB, command, migrationsDir string) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", migrationsDir),
        "postgres",
        driver,
    )
    if err != nil {
        log.Fatal(err)
    }

    var migrationErr error
    if command == "up" {
        migrationErr = m.Up()
    } else {
        migrationErr = m.Down()
    }

    if migrationErr != nil && migrationErr != migrate.ErrNoChange {
        log.Fatal(migrationErr)
    }
}

func loadFixtures(db *sql.DB, fixturesDir string) {
    files, err := os.ReadDir(fixturesDir)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".sql") {
            continue
        }

        content, err := os.ReadFile(filepath.Join(fixturesDir, file.Name()))
        if err != nil {
            log.Fatal(err)
        }

        _, err = db.Exec(string(content))
        if err != nil {
            log.Fatalf("Error executing fixture %s: %v", file.Name(), err)
        }
    }
}