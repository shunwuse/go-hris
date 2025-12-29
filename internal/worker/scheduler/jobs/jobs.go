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

func NewCronJobs() CronJobs {
	return []ICronJob{}
}
