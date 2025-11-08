package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Tansuozhe1num/codedream/internal/model"
	"github.com/Tansuozhe1num/codedream/internal/router"
	"github.com/Tansuozhe1num/codedream/internal/service"
)

// hashPassword 加密密码（临时函数，实际应该从handler导入）
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	// 1. 连接 MySQL - 从环境变量读取配置，如果没有则使用默认值
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = ""
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "codedream"
	}

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Printf("连接数据库: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// 2. 自动迁移
	// 先处理可能存在的旧表结构问题
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		// 临时禁用外键检查
		sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
		// 删除可能有问题的表（如果存在），按依赖顺序删除
		sqlDB.Exec("DROP TABLE IF EXISTS `daily_problems`")
		sqlDB.Exec("DROP TABLE IF EXISTS `problems`")
		sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")

		// 检查并添加users表的code_uid字段（如果不存在）
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'users' AND COLUMN_NAME = 'code_uid'", dbName).Scan(&count)
		if count == 0 {
			log.Println("检测到users表缺少code_uid字段，正在添加...")
			// 先添加字段（允许为空，稍后更新）
			sqlDB.Exec("ALTER TABLE `users` ADD COLUMN `code_uid` VARCHAR(50) NULL AFTER `id`")
			// 为现有用户生成code_uid
			var users []model.User
			db.Find(&users)
			for _, user := range users {
				if user.CodeUID == "" {
					codeUID := fmt.Sprintf("CODE_%s_%d", time.Now().Format("20060102150405"), user.ID)
					db.Model(&user).Update("code_uid", codeUID)
				}
			}
			// 设置为NOT NULL并添加索引
			sqlDB.Exec("ALTER TABLE `users` MODIFY COLUMN `code_uid` VARCHAR(50) NOT NULL")
			sqlDB.Exec("ALTER TABLE `users` ADD UNIQUE INDEX `idx_users_code_uid` (`code_uid`)")
			sqlDB.Exec("ALTER TABLE `users` ADD INDEX `idx_code_uid` (`code_uid`)")
			log.Println("已成功添加code_uid字段")
		}
	}

	if err := db.AutoMigrate(
		&model.Problem{},
		&model.DailyProblem{},
		&model.User{}, // 包含CodeUID字段
		&model.Post{},
		&model.Comment{},
		&model.Vote{},
		&model.Bookmark{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	// 初始化默认用户（如果不存在）
	var defaultUser model.User
	if err := db.Where("username = ?", "默认用户").First(&defaultUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 生成CodeUID
			codeUID := "CODE_DEFAULT"
			// 检查是否已存在
			var count int64
			db.Model(&model.User{}).Where("code_uid = ?", codeUID).Count(&count)
			if count > 0 {
				codeUID = "CODE_DEFAULT_" + time.Now().Format("20060102150405")
			}

			// 加密密码
			hashedPassword, _ := hashPassword("default")

			defaultUser = model.User{
				CodeUID:    codeUID,
				Username:   "默认用户",
				Email:      "default@codedream.com",
				Password:   hashedPassword,
				Reputation: 0,
			}
			if err := db.Create(&defaultUser).Error; err != nil {
				log.Printf("创建默认用户失败: %v", err)
			} else {
				log.Println("默认用户创建成功")
			}
		}
	}

	// 3. 启动定时任务
	c := cron.New()
	// 每天 00:10 执行更新（你也可以改成 00:00）
	c.AddFunc("10 0 * * *", func() { service.UpdateDaily(db) })
	c.Start()
	defer c.Stop()

	// 4. 注入 db 到 handler
	router := router.NewRouter(db)
	router.Run(":8080")
}
