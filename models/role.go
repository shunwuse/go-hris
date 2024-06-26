package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
