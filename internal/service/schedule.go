package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Tansuozhe1num/codedream/internal/model"
	"gorm.io/gorm"
)

var levels = []string{"easy", "medium", "hard"}

func UpdateDaily(db *gorm.DB) {
	// 取今天零点
	today := time.Now().Truncate(24 * time.Hour)

	// 用事务包裹删除和插入
	err := db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// 1. 删除今天已有的数据
		if err := tx.Where("picked_date = ?", today).
			Delete(&model.DailyProblem{}).Error; err != nil {
			return err
		}

		// 2. 为每个难度随机选题并插入
		for _, lvl := range levels {
			var total int64
			if err := tx.Model(&model.Problem{}).
				Where("level = ?", lvl).
				Count(&total).Error; err != nil {
				log.Printf("[UpdateDaily] count problems(level=%s) failed: %v", lvl, err)
				// 跳过该难度，继续下一个
				continue
			}
			if total == 0 {
				log.Printf("[UpdateDaily] no problems found for level=%s", lvl)
				continue
			}

			// 随机偏移
			offset := rand.Int63n(total)
			var prob model.Problem
			if err := tx.
				Where("level = ?", lvl).
				Offset(int(offset)).
				Limit(1).
				First(&prob).Error; err != nil {
				log.Printf("[UpdateDaily] fetch random problem(level=%s) failed: %v", lvl, err)
				continue
			}

			// 插入 daily_problems
			dp := model.DailyProblem{
				ProblemID:  prob.ID,
				Level:      lvl,
				PickedDate: today,
			}
			if err := tx.Create(&dp).Error; err != nil {
				log.Printf("[UpdateDaily] insert daily problem(level=%s) failed: %v", lvl, err)
				continue
			}
			log.Printf("[UpdateDaily] picked problem ID=%d for level=%s", prob.ID, lvl)
		}

		return nil
	})

	if err != nil {
		log.Printf("[UpdateDaily] transaction failed: %v", err)
	}
}
