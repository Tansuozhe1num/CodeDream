package model

import (
	"time"
)

type Problem struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Link      string `gorm:"size:500;not null"`
	Level     string `gorm:"type:enum('easy','medium','hard');not null"`
	Score     int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DailyProblem struct {
	ID         uint      `gorm:"primaryKey"`
	ProblemID  uint      `gorm:"not null;index:idx_level_date,unique"`
	Level      string    `gorm:"type:enum('easy','medium','hard');not null;index:idx_level_date,unique"`
	PickedDate time.Time `gorm:"type:date;not null;index:idx_level_date,unique"`
	Problem    Problem   `gorm:"foreignKey:ProblemID"`
}
