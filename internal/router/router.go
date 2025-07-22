package router

import (
	"github.com/Tansuozhe1num/codedream/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 静态文件服务 - SPA 应用
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/index.html", "./static/index.html")
	r.StaticFile("/problems.html", "./static/problems.html")
	r.StaticFile("/community.html", "./static/community.html") // 添加社区页面路由
	r.Static("/static", "./static/static")
	r.Static("/assets", "./static/assets")

	// API 路由组
	api := r.Group("/api")
	{
		api.GET("/stats", handler.HandleStats)
		api.POST("/subscribe", handler.HandleSubscribe)
		api.GET("/features", handler.HandleFeatures)
		api.GET("/tech-stacks", handler.HandleTechStacks)
		api.GET("/problems/daily", handler.HandleDailyProblems)

		// 添加社区相关API端点
		community := api.Group("/community")
		{
			community.GET("/posts", handler.HandleGetPosts)           // 获取帖子列表
			community.POST("/posts", handler.HandleCreatePost)        // 创建新帖子
			community.POST("/comments", handler.HandleAddComment)     // 添加评论
			community.PUT("/posts/vote", handler.HandlePostVote)      // 帖子投票
			community.PUT("/posts/bookmark", handler.HandleBookmark)  // 收藏帖子
			community.GET("/hot", handler.HandleHotPosts)             // 热门帖子
			community.GET("/active-users", handler.HandleActiveUsers) // 活跃用户
		}
	}

	// 单页面应用路由 - 所有未匹配的路由返回index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
