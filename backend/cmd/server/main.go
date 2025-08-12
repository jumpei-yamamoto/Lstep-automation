package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	// Logger setup
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer logger.Sync()

	// DB接続（環境変数から取得、開発時は空でもOK）
	dsn := getenv("DB_DSN", "")
	var db *sql.DB
	if dsn != "" {
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer db.Close()

		// DB接続確認
		if err := db.Ping(); err != nil {
			logger.Fatal("Failed to ping database", zap.Error(err))
		}
		logger.Info("Database connection established")
	} else {
		logger.Info("No DB_DSN provided, running without database connection")
	}

	// DI（後でリポジトリ実装時に追加）
	// userRepo := &persistence.UserRepoPG{DB: db}
	// registerUC := &usecase.RegisterUser{Repo: userRepo, Now: time.Now}
	// uh := &http.UserHandler{Register: registerUC}

	// Echo setup
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{getenv("FRONTEND_URL", "http://localhost:3000")},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
	})

	// API routes（後でルーター実装時に追加）
	// http.RegisterRoutes(e, uh)

	// Start server
	port := getenv("PORT", "8080")
	logger.Info("Server starting", zap.String("port", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}