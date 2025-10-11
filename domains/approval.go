package domains

import (
	"github.com/shunwuse/go-hris/constants"
)

type ApprovalCreate struct {
	Status    constants.ApprovalStatus
	CreatorID uint
}
