package models

import (
	"time"

	"gorm.io/gorm"
)

type Password struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint           `gorm:"column:user_id"`
	Hash      string         `gorm:"column:hash"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
