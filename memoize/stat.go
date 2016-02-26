package memoize

import "sync"

// Stat is a statistic that can be used to determine if a file has changed.
type Stat interface {
	// Path returns the filepath described by the Stat.
	Path() string
	// Stat returns the name of the stat, for instance "mtime".
	Stat() string
}

// Codec serializes a Stat to/from []byte.  Typically a Codec only handles Stat
// values with a specific value of Stat.Stat().
type Codec interface {
	// Stat corresponds to the stat.Stat() value the codec can serialize.
	Stat() string
	// Name is a common name for the codec, like "mtime+json".
	Name() string
	// EncodeStat encodes s as a []byte.  If s.Stat() has an unexpected value
	// an error is returned.
	EncodeStat(s Stat) ([]byte, error)
	// DecodeStat decodes b as a Stat and returns it or any error encountered.
	DecodeStat(b []byte) (Stat, error)
}

type codecRegistry struct {
	access sync.RWMutex
	codecs map[string]Codec
}

func newCodecRegistry() *codecRegistry {
	return &codecRegistry{
		codecs: map[string]Codec{},
	}
}

func (r *codecRegistry) Register(c Codec) {
	r.access.Lock()
	defer r.access.Unlock()

	stat := c.Stat()
	_, ok := r.codecs[stat]
	if ok {
		panic("already registered name: " + stat)
	}
	r.codecs[stat] = c
}

func (r *codecRegistry) Get(stat string) (Codec, bool) {
	r.access.RLock()
	c, ok := r.codecs[stat]
	r.access.RUnlock()
	return c, ok
}

var codecs = newCodecRegistry()

// RegisterCodec defines a new kind of Codec than can be used to serialize
// stats.
func RegisterCodec(c Codec) {
	codecs.Register(c)
}
