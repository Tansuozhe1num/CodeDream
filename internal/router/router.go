package router

import (
	"github.com/Tansuozhe1num/codedream/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // 运行时模式：debug/release/test
	r := gin.New()

	registerWebRoutes := func(r *gin.Engine) {
		r.Static("/static", "./static")
		r.LoadHTMLGlob("templates/*")
		r.NoRoute(handler.Show404)

		web := r.Group("/")
		{
			web.GET("/", handler.ShowHome)
			web.GET("/problems", handler.ListProblems)
			web.GET("/ranking", handler.ShowRanking)
			web.GET("/login", handler.ShowLogin)
			web.POST("/login", handler.DoLogin)
			web.GET("/notes", handler.ListNotes)
			web.GET("/solutions", handler.ListSolutions)
		}
	}

	registerAPIRoutes := func(r *gin.Engine) {
		api := r.Group("/api")
		{
			api.GET("/problems", handler.APIListProblems)
			api.GET("/ranking", handler.APIGetRanking)
			api.GET("/notes", handler.APIListNotes)
			api.POST("/notes", handler.APICreateNote)
			api.GET("/solutions", handler.APIListSolutions)
		}
	}

	// 全局中registerAPIRoutes间件
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		// requestIDMiddleware(),
		// corsMiddleware(),
	)

	registerWebRoutes(r)
	registerAPIRoutes(r)
	return r
}

//package router
//
//import (
//	"github.com/Tansuozhe1num/codedream/internal/handler"
//	"github.com/gin-gonic/gin"
//)
//
//func NewRouter() *gin.Engine {
//	r := gin.Default()
//	r.Static("/static", "./static")
//	r.LoadHTMLGlob("templates/*")
//
//	r.GET("/", handler.ShowHome)
//	r.GET("/problems", handler.ListProblems)
//	r.GET("/ranking", handler.ShowRanking)
//	r.GET("/login", handler.ShowLogin)
//	r.POST("/login", handler.DoLogin)
//	r.GET("/notes", handler.ListNotes)
//	r.GET("/solutions", handler.ListSolutions)
//
//	api := r.Group("/api")
//	{
//		api.GET("/problems", handler.APIListProblems)
//		api.GET("/ranking", handler.APIGetRanking)
//		api.POST("/notes", handler.APICreateNote)
//		api.GET("/notes", handler.APIListNotes)
//		api.GET("/solutions", handler.APIListSolutions)
//	}
//
//	return r
//}
