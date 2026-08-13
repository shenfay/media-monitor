package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shenfay/go-react-admin/internal/infra/config"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cli migrate <up|down|force>")
		fmt.Println("  cli migrate up              - Run all pending migrations")
		fmt.Println("  cli migrate down            - Rollback all migrations (to version 0)")
		fmt.Println("  cli migrate force <version> - Force set migration version")
		os.Exit(1)
	}

	command := os.Args[1]
	action := os.Args[2]

	if command != "migrate" {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	// 加载配置
	cfg, err := config.Load("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 构造 DSN (postgres://user:password@host:port/dbname?sslmode=disable)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("✓ Database migrations completed successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("✓ Database rollback completed successfully")

	case "force":
		if len(os.Args) < 4 {
			log.Fatal("Usage: cli migrate force <version>")
		}
		version, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		fmt.Printf("✓ Database version forced to %d\n", version)

	default:
		fmt.Printf("Unknown action: %s\n", action)
		os.Exit(1)
	}
}
