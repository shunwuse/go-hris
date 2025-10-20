package controllers

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleController,
	NewUserController,
	NewApprovalController,
)
