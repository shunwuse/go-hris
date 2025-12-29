package jobs

import (
	"context"
)

type CronJobs []ICronJob

type ICronJob interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

func NewCronJobs(
	cleanupTokensJob *CleanupTokensJob,
) CronJobs {
	return []ICronJob{
		cleanupTokensJob,
	}
}
