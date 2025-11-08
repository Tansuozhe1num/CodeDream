package preload

import (
	"github.com/Tansuozhe1num/codedream/internal/service"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
)

func StartCronJob(db *gorm.DB) error {
	// 定时任务和启动服务
	c := cron.New()

	_, err := c.AddFunc("10 0 * * *", func() {
		service.UpdateDaily(db)
	})
	if err != nil {
		log.Fatalf("init DB failed")
		return err
	}

	c.Start()
	defer c.Stop()

	return nil
}
