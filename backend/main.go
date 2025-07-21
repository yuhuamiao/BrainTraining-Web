// @title API 文档
// @version 1.0
// @description 大脑训练网页开发
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:8000

package main

import (
	"braintraining/backend/config"
	"braintraining/backend/dao"
	"braintraining/backend/middleware"
	"braintraining/backend/models"
	"braintraining/backend/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	//"net/http"
	_ "braintraining/backend/docs" // main 文件中导入 docs 包
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initGin() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":     "Brain Training",
			"version": "1.0",
		})
	})

	return router
}

//这是原来没有添加中间件的DAORoutes函数，所有的添加路由函数在routers包中的各个文件中
//func DAORoutes(router *gin.Engine, db *gorm.DB) {
//	//舒尔特矩阵
//	schlteDAO := dao.NewSchulteDAO(db)
//
//	//多色文字
//	ColorWordDAO := dao.NewColorWordDAO(db)
//	ColorWordHandler := routers.NewColorWordHandler(ColorWordDAO)
//
//	//瞬间记忆
//	MemoryDAO := dao.NewMemoryDAO(db)
//	MemoryHandler := routers.NewMemoryHandler(MemoryDAO)
//
//	//公交车人数
//	BusDAO := dao.NewBusDAO(db)
//	BusHandler := routers.NewBusHandler(BusDAO)
//
//	//用户个人界面
//	userDAO := dao.NewUserDAO(db)
//	userHandler := routers.NewUserHandler(userDAO)
//
//	//添加路由
//	routers.SchulteRoutes(router, schlteDAO)
//	ColorWordHandler.ColorWordsRouter(router)
//	MemoryHandler.MemoryRoutes(router)
//	BusHandler.BusRoutes(router)
//	userHandler.UserRoutes(router)
//}

func DAORoutes(router *gin.Engine, db *gorm.DB) {
	// 初始化DAO
	userDAO := dao.NewUserDAO(db)

	// 创建认证中间件
	authMiddleware := middleware.AuthMiddleware(userDAO)

	// 公开路由 (不需要认证)
	publicRoutes(router, db)

	// 受保护路由 (需要认证)
	protectedRoutes(router, db, authMiddleware)
}

// 公开路由
func publicRoutes(router *gin.Engine, db *gorm.DB) {
	// 登录注册
	authHandler := routers.NewAuthHandler(dao.NewUserDAO(db))
	router.POST("/api/v1/login", authHandler.Login)
	router.POST("/api/v1/register", authHandler.Register)

	// 其他公开API (如获取训练题目)
	//schulteDAO := dao.NewSchulteDAO(db)
	router.GET("/api/v1/schulte/matrix", routers.SchulteMatrix) // 获取舒尔特矩阵

	colorWordHandler := routers.NewColorWordHandler(dao.NewColorWordDAO(db))
	router.GET("/api/v1/color_words/matrix", colorWordHandler.ColorWordsMatrix) // 获取颜色文字
	router.GET("/api/v1/color_words/color", colorWordHandler.ColorStream)

	memoryHandler := routers.NewMemoryHandler(dao.NewMemoryDAO(db))
	router.GET("/api/v1/memory/matrix", memoryHandler.MemoryMatrix) // 获取记忆矩阵
}

// 受保护路由
func protectedRoutes(router *gin.Engine, db *gorm.DB, authMiddleware gin.HandlerFunc) {
	// 创建路由组并应用中间件
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware)

	// 舒尔特矩阵
	schulteDAO := dao.NewSchulteDAO(db)
	protected.POST("/schulte/scores", func(c *gin.Context) {
		routers.SubmitScore(c, schulteDAO)
	})

	// 多色文字
	colorWordHandler := routers.NewColorWordHandler(dao.NewColorWordDAO(db))
	protected.POST("/color_words/scores", colorWordHandler.SubmitScore)
	//protected.GET("/color_words/color", colorWordHandler.ColorStream)

	// 瞬间记忆
	memoryHandler := routers.NewMemoryHandler(dao.NewMemoryDAO(db))
	protected.POST("/memory/scores", memoryHandler.SubmitScore)

	// 公交车人数
	busHandler := routers.NewBusHandler(dao.NewBusDAO(db))
	protected.POST("/bus/scores", busHandler.SubmitScore)

	// 用户个人界面
	userHandler := routers.NewUserHandler(dao.NewUserDAO(db))
	userGroup := protected.Group("/user")
	{
		userGroup.GET("/scores", userHandler.GetUserScores)
		userGroup.GET("/users", userHandler.GetUserInfo)
		userGroup.POST("/change", userHandler.UpdateUserInfo)
		userGroup.POST("/photo", userHandler.UploadAvatar)
	}
}

func AutoMigrate(db *gorm.DB) error { //添加数据库
	//舒尔特矩阵
	if err := db.AutoMigrate(&models.SchulteScore{}); err != nil {
		log.Fatalf("SchulteScore自动迁移失败：%v", err)
		return err
	} else {
		log.Println("SchulteScore表创建/迁移成功")
	}

	//多色文字
	if err := db.AutoMigrate(&models.ColorWordRecord{}); err != nil {
		log.Fatalf("ColorWordRecord自动迁移失败：%v", err)
		return err
	} else {
		log.Println("ColorWordRecord表创建/迁移成功")
	}

	//瞬间记忆
	if err := db.AutoMigrate(&models.MemoryRecord{}); err != nil {
		log.Fatalf("MemoryRecord自动迁移失败：%v", err)
		return err
	} else {
		log.Println("MemoryRecord表创建/迁移成功")
	}

	//公交车人数
	if err := db.AutoMigrate(&models.BusRecord{}); err != nil {
		log.Fatalf("BusRecord：%v", err)
		return err
	} else {
		log.Println("BusRecord表创建/迁移成功")
	}

	//用户个人界面
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("User：%v", err)
		return err
	} else {
		log.Println("User表创建/迁移成功")
	}
	return nil
}

func main() {
	_ = godotenv.Load()

	r := initGin()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err) // 打印实际错误
	}

	// 自动迁移
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("自动迁移失败：%v", err)
	} else {
		log.Println("表创建/迁移成功")
	}

	DAORoutes(r, db)

	r.GET("/health", func(c *gin.Context) {
		if err := db.Exec("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"db": "unhealthy"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})
	if err := r.Run("0.0.0.0:8000"); err != nil {
		panic(err) // 或更优雅的错误处理
	}
}
