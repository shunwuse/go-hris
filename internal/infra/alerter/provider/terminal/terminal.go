package termalerter

import (
	"context"
	"fmt"
	"strings"

	"github.com/shunwuse/go-hris/internal/infra/alerter"
)

type sender struct{}

func New() alerter.Alerter {
	return &sender{}
}

func (s *sender) Send(ctx context.Context, msg alerter.Message) error {
	// ANSI color codes.
	reset := "\033[0m"
	red := "\033[31m"
	yellow := "\033[33m"
	bold := "\033[1m"
	blue := "\033[34m"

	color := yellow
	prefix := "⚠️  ALERT"
	if msg.Level == alerter.LevelCritical {
		color = red
		prefix = "🚨 CRITICAL ALERT"
	}

	width := 60
	line := strings.Repeat("-", width)
	banner := strings.Repeat("=", width)

	fmt.Printf("\n%s%s%s\n", color, banner, reset)
	fmt.Printf("%s%s %s: %s%s\n", color, bold, prefix, msg.Title, reset)
	fmt.Printf("%s%s%s\n", color, line, reset)

	if msg.TraceID != "" {
		fmt.Printf("%s%sTrace ID:%s %s\n", blue, bold, reset, msg.TraceID)
	}

	fmt.Printf("%s%sContent:%s %s\n", color, bold, reset, msg.Content)

	if msg.StackTrace != "" {
		fmt.Printf("\n%s%sStack Trace:%s\n%s\n", red, bold, reset, msg.StackTrace)
	}

	fmt.Printf("%s%s%s\n\n", color, banner, reset)

	return nil
}

// func init() {
// 	alert.Register(alert.ProviderTerminal, &sender{})
// }
