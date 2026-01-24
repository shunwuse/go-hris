package alerter

import (
	"errors"
	"sync"
)

type Provider string

const (
	ProviderTerminal Provider = "terminal"
	ProviderLog      Provider = "log"
)

var (
	providers = make(map[Provider]Alerter)
	mu        sync.RWMutex
)

func Register(name Provider, p Alerter) {
	mu.Lock()
	defer mu.Unlock()

	if p == nil {
		panic("alert: Register provider is nil")
	}

	providers[name] = p
}

func Get(name Provider) (Alerter, error) {
	mu.RLock()
	defer mu.RUnlock()

	p, ok := providers[name]
	if !ok {
		return nil, errors.New("alert: unknown provider " + string(name))
	}

	return p, nil
}
