package preload

import (
	"fmt"
	"github.com/Tansuozhe1num/codedream/internal/Tool"
	"github.com/Tansuozhe1num/codedream/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitDB() (*gorm.DB, error) {
	dsn := Tool.BuildDSN()
	fmt.Println(dsn) // TODO： 修改连接提示日志

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	dbService := service.NewDatabaseService(db)

	// 迁移和初始化
	if err := dbService.Migrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	if err := dbService.InitDefaultData(); err != nil {
		log.Printf("初始化数据失败: %v", err)
	}

	return db, nil
}
