package models

import (
	"time"

	"github.com/shunwuse/go-hris/constants"
	"gorm.io/gorm"
)

type Approval struct {
	ID         uint                     `gorm:"column:id;primaryKey;autoIncrement"`
	Status     constants.ApprovalStatus `gorm:"column:status"`
	CreatorID  uint                     `gorm:"column:creator_id"`
	ApproverID *uint                    `gorm:"column:approver_id"`
	CreatedAt  time.Time                `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time                `gorm:"column:updated_at;autoCreateTime:milli"`
	DeletedAt  gorm.DeletedAt           `gorm:"column:deleted_at"`

	Creator  User `gorm:"foreignkey:ID;references:CreatorID"`
	Approver User `gorm:"foreignkey:ID;references:ApproverID"`
}
