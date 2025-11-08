package router

import (
	"github.com/Tansuozhe1num/codedream/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 路由组 - 必须在静态文件路由之前注册
	api := r.Group("/api")
	{
		// 公开路由
		api.GET("/stats", handler.HandleStats)
		api.POST("/subscribe", handler.HandleSubscribe)
		api.GET("/features", handler.HandleFeatures)
		api.GET("/tech-stacks", handler.HandleTechStacks)
		api.GET("/problems/daily", handler.HandleDailyProblems(db))

		// 认证相关路由（不需要认证）
		api.POST("/auth/register", handler.HandleRegister(db)) // 注册
		api.POST("/auth/login", handler.HandleLogin(db))       // 登录
		api.POST("/auth/logout", handler.HandleLogout())       // 登出

		// 需要认证的路由
		auth := api.Group("/auth")
		auth.Use(handler.AuthMiddleware())
		{
			auth.GET("/me", handler.HandleGetCurrentUser(db)) // 获取当前用户信息
		}

		// 添加社区相关API端点
		community := api.Group("/community")
		{
			community.GET("/posts", handler.HandleGetPosts(db))           // 获取帖子列表（公开）
			community.GET("/hot", handler.HandleHotPosts(db))             // 热门帖子（公开）
			community.GET("/active-users", handler.HandleActiveUsers(db)) // 活跃用户（公开）

			// 需要认证的社区操作
			communityAuth := community.Group("")
			communityAuth.Use(handler.AuthMiddleware())
			{
				communityAuth.POST("/posts", handler.HandleCreatePost(db))       // 创建新帖子
				communityAuth.POST("/comments", handler.HandleAddComment(db))    // 添加评论
				communityAuth.PUT("/posts/vote", handler.HandlePostVote(db))     // 帖子投票
				communityAuth.PUT("/posts/bookmark", handler.HandleBookmark(db)) // 收藏帖子
			}
		}
	}

	// 静态文件服务 - SPA 应用（必须在API路由之后）
	r.Static("/static", "./static/static")
	r.Static("/assets", "./static/assets")
	r.StaticFile("/problems.html", "./static/problems.html")
	r.StaticFile("/community.html", "./static/community.html")
	r.StaticFile("/index.html", "./static/index.html")
	// 根路径放在最后，避免拦截API请求
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 单页面应用路由 - 所有未匹配的路由返回index.html（排除API路由）
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		path := c.Request.URL.Path
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{
				"success": false,
				"error":   "API endpoint not found",
			})
			return
		}
		c.File("./static/index.html")
	})

	return r
}
