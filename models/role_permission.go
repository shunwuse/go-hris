package models

import (
	"time"

	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uint           `gorm:"column:id;primaryKey;autoIncrement"`
	RoleID       uint           `gorm:"column:role_id"`
	PermissionID uint           `gorm:"column:permission_id"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`

	Role       Role       `gorm:"foreignKey:ID;references:RoleID"`
	Permission Permission `gorm:"foreignKey:ID;references:PermissionID"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
