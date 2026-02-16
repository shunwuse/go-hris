package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type source struct {
	provider koanf.Provider
	parser   koanf.Parser
}

type config struct {
	sources []source
	watch   bool
}

type Option func(*config)

func defaultOptions() *config {
	return &config{
		sources: []source{},
		watch:   false,
	}
}

// WithFile adds a file source with a parser.
func WithFile(path string, parser koanf.Parser) Option {
	return func(c *config) {
		c.sources = append(c.sources, source{
			provider: file.Provider(path),
			parser:   parser,
		})
	}
}

// WithDotEnv adds a .env file source.
func WithDotEnv(path string) Option {
	return WithFile(path, dotenv.Parser())
}

// WithEnv adds environment variables source.
func WithEnv(prefix string, delimiter string) Option {
	return func(c *config) {
		c.sources = append(c.sources, source{
			provider: env.Provider(prefix, delimiter, func(s string) string {
				return strings.ToLower(strings.TrimPrefix(s, prefix))
			}),
			parser: nil,
		})
	}
}

// WithWatch enables file watching for the last added file source.
// Note: in a real implementation, you might want to specify which files to watch.
func WithWatch(watch bool) Option {
	return func(c *config) {
		c.watch = watch
	}
}
