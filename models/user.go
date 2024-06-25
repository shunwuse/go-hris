package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string         `gorm:"column:username"`
	PasswordHash string         `gorm:"column:password_hash"`
	Name         string         `gorm:"column:name"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
