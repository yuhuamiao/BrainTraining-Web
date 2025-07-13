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

func main() {
	_ = godotenv.Load()

	r := initGin()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err) // 打印实际错误
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.SchulteScore{}); err != nil {
		log.Fatalf("自动迁移失败：%v", err)
	} else {
		log.Println("表创建/迁移成功")
	}
	
	schlteDAO := dao.NewSchulteDAO(db)

	routers.SchulteRoutes(r, schlteDAO)

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
