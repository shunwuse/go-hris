package main

import (
	"strings"
	"testing"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestDemoNumber(t *testing.T) {
	prefix := demoPrefix(constants.Staff)
	for _, tt := range []struct {
		username string
		want     int
	}{
		{"demo.staff.001", 1},
		{"demo.staff.1000", 1000},
		{"demo.staff.000", 0},
		{"demo.staff.1", 0},
		{"demo.staff.001.backup", 0},
		{"demo.staff.+01", 0},
		{"demo.manager.001", 0},
	} {
		t.Run(tt.username, func(t *testing.T) {
			require.Equal(t, tt.want, demoNumber(tt.username, prefix))
		})
	}
}

func TestUserRoles(t *testing.T) {
	require.Empty(t, userRoles(0))
	require.Equal(t, []constants.Role{constants.Staff}, userRoles(1))
	require.ElementsMatch(t, []constants.Role{constants.Staff, constants.Manager}, userRoles(2))
}

func TestApprovalStatuses(t *testing.T) {
	for _, tt := range []struct {
		count, approved, rejected int
	}{
		{0, 0, 0},
		{9, 0, 0},
		{10, 1, 0},
		{19, 1, 0}, // The previous rounding produced only 18 statuses.
		{20, 2, 1},
		{50, 5, 2},
		{100, 10, 5},
	} {
		statuses := approvalStatuses(tt.count)
		require.Len(t, statuses, tt.count)
		counts := make(map[constants.ApprovalStatus]int)
		for _, status := range statuses {
			counts[status]++
		}
		require.Equal(t, tt.approved, counts[constants.ApprovalStatusApproved])
		require.Equal(t, tt.rejected, counts[constants.ApprovalStatusRejected])
		require.Equal(t, tt.count-tt.approved-tt.rejected, counts[constants.ApprovalStatusPending])
	}
}

func TestExecuteRejectsInvalidOptionsBeforeConnecting(t *testing.T) {
	t.Setenv("SEED_PASSWORD", strings.Repeat("a", 73))
	require.ErrorContains(t, execute(t.Context(), -1, 0, 0), "must not be negative")
	require.ErrorContains(t, execute(t.Context(), 1, 0, 0), "must not exceed 72 bytes")
	require.NoError(t, execute(t.Context(), 0, 0, 0))

	t.Setenv("APP_ENV", constants.EnvProduction)
	require.ErrorContains(t, execute(t.Context(), 0, 1, 0), "env=development")
	t.Setenv("SEED_PASSWORD", "")
	require.ErrorContains(t, execute(t.Context(), 1, 0, 0), "env=development")
}
