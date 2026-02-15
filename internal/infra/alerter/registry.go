package alerter

import (
	"errors"
	"sync"

	"github.com/shunwuse/go-hris/internal/ports/infra"
)

type Provider string

const (
	ProviderTerminal Provider = "terminal"
	ProviderLog      Provider = "log"
)

var (
	providers = make(map[Provider]infra.Alerter)
	mu        sync.RWMutex
)

func Register(name Provider, p infra.Alerter) {
	mu.Lock()
	defer mu.Unlock()

	if p == nil {
		panic("alert: Register provider is nil")
	}

	providers[name] = p
}

func Get(name Provider) (infra.Alerter, error) {
	mu.RLock()
	defer mu.RUnlock()

	p, ok := providers[name]
	if !ok {
		return nil, errors.New("alert: unknown provider " + string(name))
	}

	return p, nil
}
