package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID         uint   `gorm:"primaryKey"`
	CodeUID    string `gorm:"size:50;not null;uniqueIndex;index"` // 唯一标识符
	Username   string `gorm:"size:50;not null;uniqueIndex"`       // 用户名，唯一
	Email      string `gorm:"size:100;not null;uniqueIndex"`
	Password   string `gorm:"size:255;not null"` // 存储加密后的密码
	Reputation int    `gorm:"default:0"`         // 声望值
	Avatar     string `gorm:"size:500"`          // 头像URL
	Bio        string `gorm:"type:text"`         // 个人简介
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
