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
	r.StaticFile("/problems.html", "./static/problems.html") // 添加刷题页面
	r.Static("/static", "./static/static")
	r.Static("/assets", "./static/assets")

	// API 路由组
	api := r.Group("/api")
	{
		api.GET("/stats", handler.HandleStats)
		api.POST("/subscribe", handler.HandleSubscribe)
		api.GET("/features", handler.HandleFeatures)
		api.GET("/tech-stacks", handler.HandleTechStacks)
		api.GET("/problems/daily", handler.HandleDailyProblems) // 每日题目接口
	}

	// 单页面应用路由 - 所有未匹配的路由返回index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
