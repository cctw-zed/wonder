package main

import (
	"log"

	"github.com/cctw-zed/wonder/internal/infrastructure/database"
)

func main() {
	log.Println("🗑️  Resetting database...")

	// Get database connection
	conn, err := database.NewConnectionFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create migrator
	migrator := database.NewMigrator(conn.DB())

	// Drop all tables
	if err := migrator.DropAll(); err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	log.Println("✅ Dropped all tables")

	// Run migrations
	if err := migrator.MigrateAll(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	log.Println("✅ Created tables with new schema")

	// Indexes are created automatically by GORM via model tags
	log.Println("✅ Indexes created automatically via model tags")

	log.Println("🎉 Database reset completed successfully!")
}