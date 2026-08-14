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
		fmt.Println("Usage: cli migrate <up|down|force|status>")
		fmt.Println("  cli migrate up              - Run all pending migrations")
		fmt.Println("  cli migrate down            - Rollback all migrations (to version 0)")
		fmt.Println("  cli migrate force <version> - Force set migration version")
		fmt.Println("  cli migrate status          - Show current migration status")
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

	// 打印连接信息（隐藏密码）
	fmt.Printf("📦 Database: %s@%s:%d/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch action {
	case "up":
		// 显示当前版本
		if ver, _, err := m.Version(); err == nil {
			fmt.Printf("📍 Current migration version: %d\n", ver)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		} else if err == migrate.ErrNoChange {
			fmt.Println("ℹ️  No pending migrations, already at latest version")
		} else {
			ver, _, _ := m.Version()
			fmt.Printf("✓ Migrations applied successfully, now at version %d\n", ver)
		}

	case "down":
		if ver, _, err := m.Version(); err == nil {
			fmt.Printf("📍 Current migration version: %d\n", ver)
		}
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		} else if err == migrate.ErrNoChange {
			fmt.Println("ℹ️  Already at version 0, nothing to rollback")
		} else {
			fmt.Println("✓ Database rollback completed successfully")
		}

	case "status":
		ver, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("📭 No migrations applied yet (database is clean)\n")
		} else {
			dirtyStr := "clean"
			if dirty {
				dirtyStr = "DIRTY (needs fix)"
			}
			fmt.Printf("📍 Current version: %d (%s)\n", ver, dirtyStr)
		}

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
