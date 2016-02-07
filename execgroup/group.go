package execgroup

import "sync"

// Group synchronizes dependent concurrent execution.
type Group struct {
	cond  *sync.Cond
	deps  []*Group
	err   error
	nexec int64
}

// NewGroup initializes and returns a new Group.
func NewGroup(deps []*Group) *Group {
	cond := &sync.Cond{L: &sync.Mutex{}}
	g := &Group{
		cond: cond,
		deps: deps,
	}
	return g
}

// Exec begins executing fn and prevents any waiting goroutines from resuming
// until fn returns.
func (g *Group) Exec(fn func() error) error {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()
	if g.err != nil {
		err := g.err
		g.err = nil
		return err
	}
	g.nexec++
	go g.exec(fn)
	return nil
}

func (g *Group) exec(fn func() error) {
	var err error

	defer func() {
		g.cond.L.Lock()
		if err != nil && g.err == nil {
			g.err = err
		}
		g.nexec--
		g.cond.Broadcast()
		g.cond.L.Unlock()
	}()

	for _, dep := range g.deps {
		err = dep.Wait()
		if err != nil {
			return
		}
	}

	err = fn()
}

// Wait blocks until the group has no running functions.
func (g *Group) Wait() error {
	g.cond.L.Lock()
	for g.nexec > 0 && g.err == nil {
		g.cond.Wait()
	}
	err := g.err
	g.err = nil
	g.cond.L.Unlock()
	return err
}
