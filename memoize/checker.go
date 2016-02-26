package memoize

import (
	"fmt"
	"sync"
)

func init() {
	RegisterChecker(&mtimeChecker{})
}

// CheckOpt contains options for the CheckFile function.
type CheckOpt struct {
	Stats []string
}

// CheckFile runs registered filecheckers and returns the Stats for path.
func CheckFile(path string, opt *CheckOpt) ([]Stat, error) {
	return checkers.CheckFile(path, opt)
}

// Changed determines if the differences between a and b suggest that the file
// they describe has changed.  If before and after contain a stat that has
// changed according to the registered FileChecker then Changed returns true,
// otherwise Change returns false.  Stats which are not present both in the
// before and after slices are ignored because they cannot be compared.
func Changed(before, after []Stat) bool {
	return checkers.Changed(before, after)
}

// FileChecker inspects filepaths and reports statics against them that can be
// compared to statics gathered at a later date to determine if the file has
// changed.
//
// A FileChecker must allow its methods to be called concurrently, potentially
// on the same filepaths or Stat objects.
type FileChecker interface {
	Stat() string
	CheckFile(path string) (Stat, error)
	Equal(a, b Stat) bool
}

type mtimeChecker struct {
}

var _ FileChecker = &mtimeChecker{}

func (c *mtimeChecker) Stat() string {
	return "mtime"
}

func (c *mtimeChecker) CheckFile(path string) (Stat, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *mtimeChecker) Equal(a, b Stat) bool {
	return false
}

type fcRegistry struct {
	access   sync.RWMutex
	checkers map[string]FileChecker
}

func newFCRegistry() *fcRegistry {
	return &fcRegistry{
		checkers: map[string]FileChecker{},
	}
}

func (r *fcRegistry) Register(c FileChecker) {
	r.access.Lock()
	defer r.access.Unlock()

	stat := c.Stat()
	_, ok := r.checkers[stat]
	if ok {
		panic("already registered name: " + stat)
	}
	r.checkers[stat] = c
}

func (r *fcRegistry) Get(stat string) (FileChecker, bool) {
	r.access.RLock()
	c, ok := r.checkers[stat]
	r.access.RUnlock()
	return c, ok
}

func (r *fcRegistry) Stats() []string {
	r.access.RLock()
	defer r.access.RUnlock()
	var s []string
	for stat := range r.checkers {
		s = append(s, stat)
	}
	return s
}

func (r *fcRegistry) CheckFile(path string, opt *CheckOpt) ([]Stat, error) {
	var run []string
	if opt == nil {
		run = opt.Stats
	}
	if len(run) == 0 {
		run = checkers.Stats()
	}
	var res []Stat
	for _, stat := range run {
		c, _ := checkers.Get(stat)
		s, err := c.CheckFile(path)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *fcRegistry) Changed(a, b []Stat) bool {
	return true
}

var checkers = newFCRegistry()

// RegisterChecker defines a new kind of FileChecker than can be used to
func RegisterChecker(c FileChecker) {
	checkers.Register(c)
}
