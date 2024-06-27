package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Description string         `gorm:"column:description;unique"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}
