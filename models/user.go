package models

import (
	"time"

	"github.com/shunwuse/go-hris/constants"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string         `gorm:"column:username;unique"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Password Password `gorm:"foreignkey:UserID;references:ID"`

	Roles []Role `gorm:"many2many:user_role;"`

	Permissions constants.Permissions `gorm:"-"`
}
