package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

var jwtSecret = []byte("your-secret-key-32bytes-long-for-hs256") // HS256 推荐密钥长度 ≥ 32 字节

// 登录请求体结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
}

func UserRouter(r *gin.Engine, db *gorm.DB) {
	ur := r.Group("/users")
	{
		ur.GET("/", func(c *gin.Context) {
			a := LoginRequest{
				Username: "user1",
				Password: "abc123",
			}
			c.JSON(200, a)
		})

		//添加用户注册接口
		ur.POST("/users/register", func(c *gin.Context) {
			// 1. 定义输入结构体
			var input struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			// 2. 绑定 JSON
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			// 3. 检查用户名是否已存在
			var existingUser User
			result := db.Where("username = ?", input.Username).First(&existingUser)
			if result.Error == nil {
				c.JSON(400, gin.H{"error": "Username already exists"})
				return
			}
			// 4. 存入数据库
			//4.2密码hash
			hashbyte, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to hash password"})
				return
			}
			PasswordHash := string(hashbyte)

			db.Create(&User{
				Username:     input.Username,
				PasswordHash: PasswordHash,
			})
			log.Printf("Created user with hashed password: %s", PasswordHash)
			// users[input.Username] = input.Password
			c.JSON(200, gin.H{"message": "User registered successfully"})

		})

		// 添加用户登录接口
		r.POST("/users/login", func(c *gin.Context) {
			// 1. 定义输入结构体
			var input struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			// 2. 验证用户名和密码
			var user User
			result := db.Where("username = ?", input.Username).First(&user)
			if result.Error != nil {
				c.JSON(401, gin.H{"error": "Invalid user"})
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
			if err != nil {
				c.JSON(401, gin.H{"error": "Invalid username or password"})
				return
			}
			// 3. 生成 JWT
			expirationTime := time.Now().Add(2 * time.Hour)

			// 2.2 构建 Payload（自定义 claims）
			claims := jwt.MapClaims{
				"username": input.Username,        // 存储用户名
				"role":     "user",                // 存储角色（可选）
				"exp":      expirationTime.Unix(), // 过期时间（Unix 时间戳）
				"iat":      time.Now().Unix(),     // 签发时间（可选）
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			// 4. 返回 token
			c.JSON(200, LoginResponse{Token: tokenString})
		})

	}
}
