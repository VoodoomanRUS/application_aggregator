package main

import (
	"application_aggregator/migrations"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//var (
	//	command = flag.String("command", "up", "Migration command: up, down, status")
	//)
	//flag.Parse()

	//db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("Failed to establish connection to database: %v", err)
	//}
	//
	//migrations.Migrator{} := migrations.NewMigrator(db)
	//m.AddMigration(migrations.CreateOrganizationTable())
	//m.AddMigration(migrations.CreateLoanApplicationTable())
	//
	//ctx := context.Background()
	//
	//switch *command {
	//case "up":
	//	if err := m.Up(); err != nil {
	//		log.Fatalf("Failed to run migration: %v", err)
	//	}
	//	log.Printf("Migration completed successfully!")
	//case "down":
	//	if err := m.Down(); err != nil {
	//		log.Fatalf("Failed to rollback migration: %v", err)
	//	}
	//	log.Printf("Migration rollback completed successfully!")
	//case "status":
	//	if err := m.Status(ctx); err != nil {
	//		log.Fatalf("Failed to get migration status: %v", err)
	//	}
	//default:
	//	log.Fatalf("Unknown command: %s Please use : up, down, status", *command)
	//}
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Failed to load .env.local.\n")
	}
	dbPath := os.Getenv("DB_DSN")
	if dbPath == "" {
		log.Fatal("DB_DSN environment variable not set!")
	}
	fmt.Println(" Connecting with DSN:", dbPath)
	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ failed to connect database: ", err)
	}
	m := migrations.NewMigrator(db)
	m.AddMigration(migrations.CreateOrganizationTable())
	m.AddMigration(migrations.CreateLoanApplicationTable())

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [up|down|status]")
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal("❌ Migration UP failed: ", err)
		}
		fmt.Println("✅ All migrations applied successfully!")

	case "down":
		if err := m.Down(); err != nil {
			log.Fatal("❌ Migration DOWN failed: ", err)
		}
		fmt.Println("✅ Last migration rolled back successfully!")

	case "status":
		if err := m.Status(context.Background()); err != nil {
			log.Fatal("❌ Status check failed: ", err)
		}

	default:
		log.Fatal("Unknown command: ", cmd)
	}
}
