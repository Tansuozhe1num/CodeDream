package router

import (
	"github.com/Tansuozhe1num/codedream"
	"github.com/gin-gonic/gin/internal/handler"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", handler.ShowHome)
	r.GET("/problems", handler.ListProblems)
	r.GET("/ranking", handler.ShowRanking)
	r.GET("/login", handler.ShowLogin)
	r.POST("/login", handler.DoLogin)
	r.GET("/notes", handler.ListNotes)
	r.GET("/solutions", handler.ListSolutions)

	api := r.Group("/api")
	{
		api.GET("/problems", handler.APIListProblems)
		api.GET("/ranking", handler.APIGetRanking)
		api.POST("/notes", handler.APICreateNote)
		api.GET("/notes", handler.APIListNotes)
		api.GET("/solutions", handler.APIListSolutions)
	}

	return r
}
