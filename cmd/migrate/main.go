package main

import (
	"log/slog"
	"os"

	"github.com/PritamDasTharu/conduit/internal/config"
	"github.com/PritamDasTharu/conduit/internal/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	db, err := database.NewPostgres(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	logger.Info("running database migrations...")
	if err := database.AutoMigrate(db.DB); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}

	logger.Info("database migration completed successfully")
}
