package model

import (
	"time"
)

// Vote 投票模型
type Vote struct {
	ID        uint   `gorm:"primaryKey"`
	PostID    uint   `gorm:"not null;index:idx_post_user,unique"`
	UserID    uint   `gorm:"not null;index:idx_post_user,unique"`
	Type      string `gorm:"type:enum('up','down');not null"` // up=点赞, down=点踩
	CreatedAt time.Time
	UpdatedAt time.Time

	// 关联关系
	Post Post `gorm:"foreignKey:PostID"`
	User User `gorm:"foreignKey:UserID"`
}
