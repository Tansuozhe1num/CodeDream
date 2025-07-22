package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 数据结构定义
type StatsResponse struct {
	Users        int `json:"users"`
	Courses      int `json:"courses"`
	Projects     int `json:"projects"`
	Satisfaction int `json:"satisfaction"`
}

type Feature struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type TechStack struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// HandleStats 获取平台统计数据
func HandleStats(c *gin.Context) {
	stats := StatsResponse{
		Users:        50000,
		Courses:      120,
		Projects:     1000000,
		Satisfaction: 98,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// HandleSubscribe 处理订阅请求
func HandleSubscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "邮箱格式不正确",
		})
		return
	}

	// 这里应该将邮箱保存到数据库
	// 例如: database.SaveEmail(req.Email)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅成功",
	})
}

// HandleFeatures 获取特性列表
func HandleFeatures(c *gin.Context) {
	features := []Feature{
		{
			ID:          1,
			Title:       "交互式学习环境",
			Description: "直接在浏览器中编写和运行代码，实时查看结果，无需复杂环境配置。",
			Icon:        "laptop-code",
			Color:       "blue",
		},
		{
			ID:          2,
			Title:       "专业学习路径",
			Description: "精心设计的课程体系，从基础到进阶，涵盖主流编程语言和技术栈。",
			Icon:        "graduation-cap",
			Color:       "purple",
		},
		{
			ID:          3,
			Title:       "活跃开发者社区",
			Description: "与全球开发者交流经验，解决难题，参与开源项目，共同成长。",
			Icon:        "users",
			Color:       "cyan",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    features,
	})
}

// HandleTechStacks 获取技术栈列表
func HandleTechStacks(c *gin.Context) {
	techStacks := []TechStack{
		{Name: "HTML5", Icon: "html5", Color: "orange"},
		{Name: "CSS3", Icon: "css3-alt", Color: "blue"},
		{Name: "JavaScript", Icon: "js", Color: "yellow"},
		{Name: "Python", Icon: "python", Color: "blue"},
		{Name: "React", Icon: "react", Color: "cyan"},
		{Name: "Database", Icon: "database", Color: "emerald"},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    techStacks,
	})
}
