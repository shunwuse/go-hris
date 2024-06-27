package dtos

type ApprovalResponse struct {
	ID           uint    `json:"id"`
	CreatorName  string  `json:"creator_name"`
	ApproverName *string `json:"approver_name"`
	Status       string  `json:"status"`
}

type ApprovalAction struct {
	ID     uint   `json:"id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=APPROVED REJECTED"`
}
