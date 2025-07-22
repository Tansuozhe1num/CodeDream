package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListProblems 渲染刷题单页面
func ListProblems(c *gin.Context) {
	c.HTML(http.StatusOK, "problems.html", gin.H{"title": "刷题单"})
}

// APIListProblems 返回刷题单 JSON
func APIListProblems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{}})
}
