package repository

import (
	"context"
)

type MonitorRepository interface {
	CheckDatabase(ctx context.Context) bool
	CheckRedis(ctx context.Context) bool
}
