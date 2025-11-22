package main

import (
	"go_web/template/router"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	f, _ := os.Create("app.log")
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 引擎
	r := gin.Default()

	type User struct {
		gorm.Model
		Username     string `gorm:"column:username;type:varchar(100);not null;unique"`
		PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	}

	// 连接数据库
	dsn := "ecommerce_user:123456@tcp(127.0.0.1:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 测试连接语句
	db.Raw("SELECT 1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		print(err)
	}
	db.AutoMigrate(&User{})

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 配置静态文件服务
	r.Static("/static", "./frontend/build/static")
	r.LoadHTMLGlob("frontend/build/*.html")

	router.WsRouter(r)
	router.UserRouter(r, db)
	router.TestRouter(r)

	// 启动服务器，添加错误处理
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
