package service

import "context"

type ExampleService interface {
	Ping(ctx context.Context) string
}
