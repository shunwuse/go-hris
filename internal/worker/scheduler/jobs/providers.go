package jobs

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCleanupTokensJob,

	NewCronJobs,
)
