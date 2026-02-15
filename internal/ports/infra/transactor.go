package infra

import (
	"context"
)

type Transactor interface {
	WithTx(ctx context.Context, work func(ctx context.Context) error) error
}
