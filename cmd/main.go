package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nikitaenmi/AvitoTest/internal/config"
	"github.com/nikitaenmi/AvitoTest/internal/database"
	"github.com/nikitaenmi/AvitoTest/internal/handlers"
	"github.com/nikitaenmi/AvitoTest/internal/repository"
	"github.com/nikitaenmi/AvitoTest/internal/service"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := connectDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Database init failed:", err)
	}

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db, userRepo)
	prRepo := repository.NewPRRepository(db)

	svc := service.NewService(userRepo, teamRepo, prRepo)
	h := handlers.NewHandlers(svc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/team/add", h.CreateTeam)
	e.GET("/team/get", h.GetTeam)
	e.POST("/users/setIsActive", h.SetUserActive)
	e.GET("/users/getReview", h.GetUserReviewPRs)
	e.POST("/pullRequest/create", h.CreatePR)
	e.POST("/pullRequest/merge", h.MergePR)
	e.POST("/pullRequest/reassign", h.ReassignReviewer)
	e.GET("/health", h.HealthCheck)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Printf("Server is running on %s", srv.Addr)

	if err := e.StartServer(srv); err != nil {
		log.Fatal("Server not running:", err)
	}
}

func connectDatabase(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

	var db *gorm.DB
	var err error

	for i := 0; i < 15; i++ {
		db, err = database.NewPostgresDB(dsn)
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
