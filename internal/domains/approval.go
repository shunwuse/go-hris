package domains

import (
	"github.com/shunwuse/go-hris/internal/constants"
)

type ApprovalCreate struct {
	Status    constants.ApprovalStatus
	CreatorID uint
}
