package model

import (
	"time"
)

// Post 帖子模型
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text;not null"`
	Type      string    `gorm:"type:enum('question','solution','discussion');not null;default:'question'"`
	Tags      string    `gorm:"type:json"`       // 存储JSON格式的标签数组
	VoteCount int       `gorm:"default:0;index"` // 投票总数（upvotes - downvotes）
	ViewCount int       `gorm:"default:0"`       // 浏览次数
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time

	// 关联关系
	User      User       `gorm:"foreignKey:UserID"`
	Comments  []Comment  `gorm:"foreignKey:PostID"`
	Votes     []Vote     `gorm:"foreignKey:PostID"`
	Bookmarks []Bookmark `gorm:"foreignKey:PostID"`
}
