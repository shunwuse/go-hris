package domains

import (
	"github.com/shunwuse/go-hris/internal/constants"
)

type Approval struct {
	ID           uint
	Status       constants.ApprovalStatus
	CreatorID    uint
	ApproverID   *uint
	CreatorName  string
	ApproverName *string
}

type ApprovalCreate struct {
	Status    constants.ApprovalStatus
	CreatorID uint
}

type ApprovalFilter struct {
	Status constants.ApprovalStatus
}
