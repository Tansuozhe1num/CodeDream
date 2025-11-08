package handler

import (
	"net/http"
	"time"

	"github.com/Tansuozhe1num/codedream/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DailyProblemDTO struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Level   string `json:"level"`
	Content string `json:"content"` // 如果你在 problems 表加了 content 字段
}

func HandleDailyProblems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		today := time.Now().Truncate(24 * time.Hour)
		var dps []model.DailyProblem

		// 预加载题目详情
		if err := db.
			Preload("Problem").
			Where("picked_date = ?", today).
			Find(&dps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 组装 DTO
		var out []DailyProblemDTO
		for _, dp := range dps {
			out = append(out, DailyProblemDTO{
				ID:      dp.Problem.ID,
				Title:   dp.Problem.Title,
				Link:    dp.Problem.Link,
				Level:   dp.Level,
				Content: "", // Problem 模型中没有 Content 字段，如需使用请先添加到模型中
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
	}
}
