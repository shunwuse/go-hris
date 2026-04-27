package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
)

func TestMapUser(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 9, 30, 0, 0, time.UTC)
	mapped := mapUser(&entgen.User{
		ID:        7,
		Username:  "alice",
		Name:      "Alice",
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Hour),
	})

	require.NotNil(t, mapped)
	require.Equal(t, uint(7), mapped.ID)
	require.Equal(t, "alice", mapped.Username)
	require.Equal(t, "Alice", mapped.Name)
	require.Equal(t, now, mapped.CreatedAt)
	require.Equal(t, now.Add(1*time.Hour), mapped.UpdatedAt)
	require.Nil(t, mapUser(nil))
}

func TestMapUserWithRoles(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	mapped := mapUserWithRoles(&entgen.User{
		ID:        8,
		Username:  "bob",
		Name:      "Bob",
		CreatedAt: now,
		UpdatedAt: now.Add(30 * time.Minute),
		Edges: entgen.UserEdges{
			Password: &entgen.Password{Hash: "hashed-password"},
			Roles: []*entgen.Role{
				{Name: constants.Admin},
				{Name: constants.Manager},
			},
		},
	})

	require.NotNil(t, mapped)
	require.Equal(t, uint(8), mapped.ID)
	require.Equal(t, "bob", mapped.Username)
	require.Equal(t, "Bob", mapped.Name)
	require.Equal(t, "hashed-password", mapped.PasswordHash)
	require.Equal(t, []constants.Role{constants.Admin, constants.Manager}, mapped.Roles)
	require.Equal(t, now, mapped.CreatedAt)
	require.Equal(t, now.Add(30*time.Minute), mapped.UpdatedAt)
	require.Nil(t, mapUserWithRoles(nil))
}

func TestMapApproval(t *testing.T) {
	t.Parallel()

	creatorName := "Creator"
	approverName := "Approver"
	approverID := uint(12)
	mapped := mapApproval(&entgen.Approval{
		ID:         9,
		Status:     constants.ApprovalStatusPending,
		CreatorID:  3,
		ApproverID: &approverID,
		Edges: entgen.ApprovalEdges{
			Creator:  &entgen.User{Name: creatorName},
			Approver: &entgen.User{Name: approverName},
		},
	})

	require.NotNil(t, mapped)
	require.Equal(t, uint(9), mapped.ID)
	require.Equal(t, constants.ApprovalStatusPending, mapped.Status)
	require.Equal(t, uint(3), mapped.CreatorID)
	require.NotNil(t, mapped.ApproverID)
	require.Equal(t, uint(12), *mapped.ApproverID)
	require.Equal(t, creatorName, mapped.CreatorName)
	require.NotNil(t, mapped.ApproverName)
	require.Equal(t, approverName, *mapped.ApproverName)
	require.Nil(t, mapApproval(nil))
}

func TestMapRefreshToken(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 11, 0, 0, 0, time.UTC)
	mapped := mapRefreshToken(&entgen.RefreshToken{
		ID:        10,
		UserID:    4,
		TokenHash: "token-hash",
		ExpiresAt: now,
		Revoked:   true,
		RevokedAt: &now,
	})

	require.NotNil(t, mapped)
	require.Equal(t, uint(10), mapped.ID)
	require.Equal(t, uint(4), mapped.UserID)
	require.Equal(t, "token-hash", mapped.TokenHash)
	require.Equal(t, now, mapped.ExpiresAt)
	require.True(t, mapped.Revoked)
	require.Equal(t, &now, mapped.RevokedAt)
	require.Nil(t, mapRefreshToken(nil))
}

func TestMapRole(t *testing.T) {
	t.Parallel()

	mapped := mapRole(&entgen.Role{
		ID:   11,
		Name: constants.Staff,
	})

	require.NotNil(t, mapped)
	require.Equal(t, uint(11), mapped.ID)
	require.Equal(t, constants.Staff, mapped.Name)
	require.Nil(t, mapRole(nil))
}

func TestMapSlices(t *testing.T) {
	t.Parallel()

	users := mapUsers([]*entgen.User{
		{ID: 1, Username: "u1", Name: "User 1"},
		{ID: 2, Username: "u2", Name: "User 2"},
	})
	require.Len(t, users, 2)
	require.Equal(t, "u1", users[0].Username)
	require.Equal(t, "u2", users[1].Username)

	approvals := mapApprovals([]*entgen.Approval{
		{ID: 1, Status: constants.ApprovalStatusPending, CreatorID: 1},
		{ID: 2, Status: constants.ApprovalStatusApproved, CreatorID: 2},
	})
	require.Len(t, approvals, 2)
	require.Equal(t, constants.ApprovalStatusPending, approvals[0].Status)
	require.Equal(t, constants.ApprovalStatusApproved, approvals[1].Status)

	roles := mapRoles([]*entgen.Role{
		{ID: 1, Name: constants.Admin},
		{ID: 2, Name: constants.Manager},
	})
	require.Len(t, roles, 2)
	require.Equal(t, constants.Admin, roles[0].Name)
	require.Equal(t, constants.Manager, roles[1].Name)
}
