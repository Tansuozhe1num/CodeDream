package service

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Tansuozhe1num/codedream/internal/model"
)

type DatabaseService struct {
	db *gorm.DB
}

func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

func (s *DatabaseService) Migrate() error {
	return s.db.AutoMigrate(
		&model.Problem{},
		&model.DailyProblem{},
		&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Vote{},
		&model.Bookmark{},
	)
}

func (s *DatabaseService) InitDefaultData() error {
	var defaultUser model.User
	if err := s.db.Where("username = ?", "默认用户").First(&defaultUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("default"), bcrypt.DefaultCost)

			defaultUser = model.User{
				CodeUID:    generateCodeUID(),
				Username:   "默认用户",
				Email:      "default@codedream.com",
				Password:   string(hashedPassword),
				Reputation: 0,
			}

			if err := s.db.Create(&defaultUser).Error; err != nil {
				return err
			}
			log.Println("默认用户创建成功")
		}
	}
	return nil
}

// 生成codeuid
func generateCodeUID() string {
	return "UU"
}
