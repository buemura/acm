package provider

import (
	"fmt"
	"path/filepath"
)

type PathCheck struct {
	Name     string
	Patterns []string
}

type Provider struct {
	Name       string
	CachePaths []string
	Checks     []PathCheck
}

var registry = map[string]*Provider{}

func Register(p *Provider) {
	registry[p.Name] = p
}

func Get(name string) (*Provider, error) {
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

func List() []*Provider {
	providers := make([]*Provider, 0, len(registry))
	for _, p := range registry {
		providers = append(providers, p)
	}
	return providers
}

func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

func MissingPathChecks(p *Provider) []string {
	missing := make([]string, 0, len(p.Checks))
	for _, check := range p.Checks {
		found := false
		for _, pattern := range check.Patterns {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			if len(matches) > 0 {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, check.Name)
		}
	}
	return missing
}
