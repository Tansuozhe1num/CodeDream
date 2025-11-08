package model

import (
	"time"
)

// Bookmark 收藏模型
type Bookmark struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index:idx_post_user,unique"`
	UserID    uint `gorm:"not null;index:idx_post_user,unique"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 关联关系
	Post Post `gorm:"foreignKey:PostID"`
	User User `gorm:"foreignKey:UserID"`
}
