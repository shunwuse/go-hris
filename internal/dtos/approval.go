package dtos

import (
	"slices"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
)

type ApprovalGetList struct {
	domains.CursorQuery

	Status constants.ApprovalStatus `schema:"status"`
}

func (d *ApprovalGetList) Validate() error {
	if d.Status != "" {
		if !slices.Contains(
			[]constants.ApprovalStatus{
				constants.ApprovalStatusPending,
				constants.ApprovalStatusApproved,
				constants.ApprovalStatusRejected,
			},
			d.Status,
		) {
			return errors.ErrValidationFailed.WithDetails(map[string]string{
				"status": "invalid approval status",
			})
		}
	}

	return nil
}

type ApprovalResponse struct {
	ID           uint                     `json:"id"`
	CreatorName  string                   `json:"creator_name"`
	ApproverName *string                  `json:"approver_name"`
	Status       constants.ApprovalStatus `json:"status"`
}

type ApprovalPathParams struct {
	ID uint `schema:"id"`
}

func (d *ApprovalPathParams) Validate() error {
	if d.ID == 0 {
		return errors.ErrValidationFailed.WithDetails(map[string]string{
			"id": "id must be greater than 0",
		})
	}

	return nil
}

type ApprovalAction struct {
	Action constants.ApprovalStatus `json:"action"`
}

func (d *ApprovalAction) Validate() error {
	if !slices.Contains(
		[]constants.ApprovalStatus{
			constants.ApprovalStatusApproved,
			constants.ApprovalStatusRejected,
		},
		d.Action,
	) {
		return errors.ErrValidationFailed.WithDetails(map[string]string{
			"action": "invalid action, must be approved or rejected",
		})
	}

	return nil
}
