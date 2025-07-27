package coalesce

import (
	"sync"
)

type call struct {
	val  any
	err  error
	done bool
	cond *sync.Cond
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func NewGroup() *Group {
	return &Group{m: make(map[string]*call)}
}

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		for !c.done {
			// Wait() automatically unlocks the mutex while waiting, and relocks it before returning
			c.cond.Wait()
		}
		val, err := c.val, c.err
		g.mu.Unlock()
		return val, err
	}
	c := &call{cond: sync.NewCond(&g.mu)}
	g.m[key] = c
	g.mu.Unlock()

	val, err := fn()

	g.mu.Lock()
	c.val = val
	c.err = err
	c.done = true
	c.cond.Broadcast()
	g.mu.Unlock()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return val, err
}
