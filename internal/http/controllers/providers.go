package controllers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthController,
	NewUserController,
	NewApprovalController,
	NewAuthController,
)
