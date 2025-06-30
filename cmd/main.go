package main

import (
	"log"
	"path/filepath"

	"L0/internal/config"
	"L0/internal/repository"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	migrationsPath, err := filepath.Abs("internal/repository/migrations")
	if err != nil {
		log.Fatalf("failed to get migrations path: %v", err)
	}
	
	err = repo.RunMigrations(migrationsPath)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

}
