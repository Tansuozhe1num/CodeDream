package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/Tansuozhe1num/codedream/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWT密钥（实际应用中应该从环境变量读取）
var jwtSecret = []byte("codedream-secret-key-change-in-production")

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	CodeUID  string `json:"code_uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID         uint   `json:"id"`
	CodeUID    string `json:"code_uid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Reputation int    `json:"reputation"`
	Avatar     string `json:"avatar"`
}

// generateCodeUID 生成唯一的Code_uid
func generateCodeUID(db *gorm.DB) (string, error) {
	for {
		// 生成格式: CODE_ + 8位随机数字
		codeUID := "CODE_" + generateRandomString(8)

		// 检查是否已存在
		var count int64
		db.Model(&model.User{}).Where("code_uid = ?", codeUID).Count(&count)
		if count == 0 {
			return codeUID, nil
		}
	}
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	// 使用时间戳和随机数生成更随机的字符串
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*1103515245 + 12345
		b[i] = charset[seed%int64(len(charset))]
	}
	return string(b)
}

// hashPassword 加密密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword 验证密码
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken 生成JWT token
func generateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时过期

	claims := &Claims{
		UserID:   user.ID,
		CodeUID:  user.CodeUID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// parseToken 解析JWT token
func parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// 尝试从Cookie获取
			cookie, err := c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "未授权，请先登录",
				})
				c.Abort()
				return
			}
			tokenString = cookie
		} else {
			// 移除 "Bearer " 前缀
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("code_uid", claims.CodeUID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// HandleRegister 处理用户注册
func HandleRegister(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		// 检查用户名是否已存在
		var existingUser model.User
		if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "用户名已存在",
			})
			return
		}

		// 检查邮箱是否已存在
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "邮箱已被注册",
			})
			return
		}

		// 生成Code_uid
		codeUID, err := generateCodeUID(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "生成用户ID失败",
			})
			return
		}

		// 加密密码
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "密码加密失败",
			})
			return
		}

		// 创建用户
		user := model.User{
			CodeUID:    codeUID,
			Username:   req.Username,
			Email:      req.Email,
			Password:   hashedPassword,
			Reputation: 0,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建用户失败: " + err.Error(),
			})
			return
		}

		// 生成token
		token, err := generateToken(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "生成token失败",
			})
			return
		}

		// 设置Cookie
		c.SetCookie("token", token, 24*3600, "/", "", false, true)

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "注册成功",
			"token":   token,
			"user": UserResponse{
				ID:         user.ID,
				CodeUID:    user.CodeUID,
				Username:   user.Username,
				Email:      user.Email,
				Reputation: user.Reputation,
				Avatar:     user.Avatar,
			},
		})
	}
}

// HandleLogin 处理用户登录
func HandleLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		// 查找用户（支持用户名或邮箱登录）
		var user model.User
		err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "用户名或密码错误",
			})
			return
		}

		// 验证密码
		if !checkPassword(req.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "用户名或密码错误",
			})
			return
		}

		// 生成token
		token, err := generateToken(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "生成token失败",
			})
			return
		}

		// 设置Cookie
		c.SetCookie("token", token, 24*3600, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "登录成功",
			"token":   token,
			"user": UserResponse{
				ID:         user.ID,
				CodeUID:    user.CodeUID,
				Username:   user.Username,
				Email:      user.Email,
				Reputation: user.Reputation,
				Avatar:     user.Avatar,
			},
		})
	}
}

// HandleGetCurrentUser 获取当前登录用户信息
func HandleGetCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未授权",
			})
			return
		}

		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "用户不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"user": UserResponse{
				ID:         user.ID,
				CodeUID:    user.CodeUID,
				Username:   user.Username,
				Email:      user.Email,
				Reputation: user.Reputation,
				Avatar:     user.Avatar,
			},
		})
	}
}

// HandleLogout 处理用户登出
func HandleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 删除Cookie
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "登出成功",
		})
	}
}

func init() {
	// 从环境变量读取JWT密钥（如果设置了）
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}
}
