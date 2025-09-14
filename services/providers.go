package services

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleService,
	NewUserService,
	NewAuthService,
	NewApprovalService,
)
