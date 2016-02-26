package memoize

// DB stores cached file stats and command memoization information.
type DB interface {
	GetStats(path string) ([]Stat, error)
	WriteStats(path string, s []Stat)
}
