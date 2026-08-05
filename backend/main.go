// @title Brain Training API
// @version 1.0
// @description 大脑训练应用 API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"braintraining/backend/config"
	"braintraining/backend/dao"
	_ "braintraining/backend/docs"
	"braintraining/backend/middleware"
	"braintraining/backend/models"
	"braintraining/backend/routers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())
	router.MaxMultipartMemory = 2 << 20
	_ = router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := os.MkdirAll("avatars", 0o755); err == nil {
		router.Static("/uploads/avatars", "avatars")
	}

	registerAPIRoutes(router, db)
	serveFrontend(router)
	return router
}

func registerAPIRoutes(router *gin.Engine, db *gorm.DB) {
	userDAO := dao.NewUserDAO(db)
	authHandler := routers.NewAuthHandler(userDAO)
	router.POST("/api/v1/login", authHandler.Login)
	router.POST("/api/v1/register", authHandler.Register)

	router.GET("/api/v1/schulte/matrix", routers.SchulteMatrix)
	colorWordHandler := routers.NewColorWordHandler(dao.NewColorWordDAO(db))
	router.GET("/api/v1/color_words/matrix", colorWordHandler.ColorWordsMatrix)
	router.GET("/api/v1/color_words/color", colorWordHandler.ColorStream)
	memoryHandler := routers.NewMemoryHandler(dao.NewMemoryDAO(db))
	router.GET("/api/v1/memory/matrix", memoryHandler.MemoryMatrix)

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(userDAO))
	protected.POST("/schulte/scores", func(c *gin.Context) {
		routers.SubmitScore(c, dao.NewSchulteDAO(db))
	})
	protected.POST("/color_words/scores", colorWordHandler.SubmitScore)
	protected.POST("/memory/scores", memoryHandler.SubmitScore)
	protected.POST("/bus/scores", routers.NewBusHandler(dao.NewBusDAO(db)).SubmitScore)
	protected.POST("/sudoku/scores", routers.NewSudokuHandler(dao.NewSudokuDAO(db)).SubmitScore)

	userHandler := routers.NewUserHandler(userDAO)
	userGroup := protected.Group("/user")
	userGroup.GET("/scores", userHandler.GetUserScores)
	userGroup.GET("/me", userHandler.GetUserInfo)
	userGroup.GET("/users", userHandler.GetUserInfo)
	userGroup.PATCH("/me", userHandler.UpdateUserInfo)
	userGroup.POST("/change", userHandler.UpdateUserInfo)
	userGroup.POST("/photo", userHandler.UploadAvatar)
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.SchulteScore{},
		&models.ColorWordRecord{},
		&models.ColorWordConfig{},
		&models.MemoryRecord{},
		&models.BusRecord{},
		&models.SudokuRecord{},
	)
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigin := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGIN"))
	return func(c *gin.Context) {
		if allowedOrigin != "" && c.GetHeader("Origin") == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func serveFrontend(router *gin.Engine) {
	frontendDir := strings.TrimSpace(os.Getenv("FRONTEND_DIR"))
	if frontendDir == "" {
		frontendDir = "./static"
	}
	indexPath := filepath.Join(frontendDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"app": "Brain Training", "version": "1.0"})
		})
		return
	}

	router.Static("/assets", filepath.Join(frontendDir, "assets"))
	router.GET("/", func(c *gin.Context) { c.File(indexPath) })
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.File(indexPath)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "NotFound", "message": "接口不存在"})
	})
}

func main() {
	_ = godotenv.Load()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8000"
	}
	if err := setupRouter(db).Run("0.0.0.0:" + port); err != nil {
		log.Fatal(err)
	}
}
