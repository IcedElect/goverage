package coverage

import "sync"

type cache struct {
	coverageByProfile map[string]Coverage
	mu                sync.Mutex
}

func newCache() *cache {
	return &cache{
		coverageByProfile: make(map[string]Coverage),
		mu:                sync.Mutex{},
	}
}

func (c *cache) Get(profileName string) (Coverage, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	coverage, ok := c.coverageByProfile[profileName]
	return coverage, ok
}

func (c *cache) Set(profileName string, coverage Coverage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.coverageByProfile[profileName] = coverage
}
