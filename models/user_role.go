package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int            `gorm:"column:user_id"`
	RoleID    int            `gorm:"column:role_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserRole) TableName() string {
	return "user_role"
}
