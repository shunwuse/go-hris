package alerter

import (
	"context"
)

type Level string

const (
	LevelCritical Level = "critical"
	LevelError    Level = "error"
)

type Message struct {
	Level      Level
	TraceID    string
	Title      string
	Content    string
	StackTrace string
}

type Alerter interface {
	Send(ctx context.Context, msg Message) error
}
