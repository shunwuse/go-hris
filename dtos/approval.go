package dtos

type ApprovalResponse struct {
	ID           uint    `json:"id"`
	CreatorName  string  `json:"creator_name"`
	ApproverName *string `json:"approver_name"`
	Status       string  `json:"status"`
}
