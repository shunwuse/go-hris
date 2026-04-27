package repositories

import (
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

func mapUsers(users []*entgen.User) []domains.User {
	result := make([]domains.User, len(users))
	for idx, user := range users {
		result[idx] = *mapUser(user)
	}

	return result
}

func mapUser(user *entgen.User) *domains.User {
	if user == nil {
		return nil
	}

	return &domains.User{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func mapUserWithRoles(user *entgen.User) *domains.UserWithRoles {
	if user == nil {
		return nil
	}

	roles := make([]constants.Role, len(user.Edges.Roles))
	for idx, role := range user.Edges.Roles {
		roles[idx] = role.Name
	}

	passwordHash := ""
	if user.Edges.Password != nil {
		passwordHash = user.Edges.Password.Hash
	}

	return &domains.UserWithRoles{
		ID:           user.ID,
		Username:     user.Username,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		PasswordHash: passwordHash,
		Roles:        roles,
	}
}

func mapApprovals(approvals []*entgen.Approval) []domains.Approval {
	result := make([]domains.Approval, len(approvals))
	for idx, approval := range approvals {
		result[idx] = *mapApproval(approval)
	}

	return result
}

func mapApproval(approval *entgen.Approval) *domains.Approval {
	if approval == nil {
		return nil
	}

	domainApproval := &domains.Approval{
		ID:        approval.ID,
		Status:    approval.Status,
		CreatorID: approval.CreatorID,
	}

	if approval.ApproverID != nil {
		approverID := *approval.ApproverID
		domainApproval.ApproverID = &approverID
	}

	if approval.Edges.Creator != nil {
		domainApproval.CreatorName = approval.Edges.Creator.Name
	}

	if approval.Edges.Approver != nil {
		approverName := approval.Edges.Approver.Name
		domainApproval.ApproverName = &approverName
	}

	return domainApproval
}

func mapRoles(roles []*entgen.Role) []domains.Role {
	result := make([]domains.Role, len(roles))
	for idx, role := range roles {
		result[idx] = *mapRole(role)
	}

	return result
}

func mapRole(role *entgen.Role) *domains.Role {
	if role == nil {
		return nil
	}

	return &domains.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}

func mapRefreshToken(token *entgen.RefreshToken) *domains.RefreshToken {
	if token == nil {
		return nil
	}

	return &domains.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		Revoked:   token.Revoked,
		RevokedAt: token.RevokedAt,
	}
}
