// @title API 文档
// @version 1.0
// @description 大脑训练网页开发
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000

package main

import (
	"braintraining/backend/config"
	"braintraining/backend/dao"
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

func DAORoutes(router *gin.Engine, db *gorm.DB) {
	//舒尔特矩阵
	schlteDAO := dao.NewSchulteDAO(db)

	//多色文字
	ColorWordDAO := dao.NewColorWordDAO(db)
	ColorWordHandler := routers.NewColorWordHandler(ColorWordDAO)

	//瞬间记忆
	MemoryDAO := dao.NewMemoryDAO(db)
	MemoryHandler := routers.NewMemoryHandler(MemoryDAO)

	//公交车人数
	BusDAO := dao.NewBusDAO(db)
	BusHandler := routers.NewBusHandler(BusDAO)

	//用户个人界面
	userDAO := dao.NewUserDAO(db)
	userHandler := routers.NewUserHandler(userDAO)

	//添加路由
	routers.SchulteRoutes(router, schlteDAO)
	ColorWordHandler.ColorWordsRouter(router)
	MemoryHandler.MemoryRoutes(router)
	BusHandler.BusRoutes(router)
	userHandler.UserRoutes(router)
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
