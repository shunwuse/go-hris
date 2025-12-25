package controllers

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewHealthController,
	NewUserController,
	NewApprovalController,
	NewAuthController,
)
