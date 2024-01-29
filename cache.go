package goregexpopulate

import (
	"crypto/sha1"
	"fmt"
	"sync"

	regen "github.com/zach-klippenstein/goregen"
)

type genCache map[string]regen.Generator

var cache genCache

var once sync.Once

func newGenCache() {
	once.Do(func() {
		cache = make(genCache)
	})
}

func (g genCache) register(pattern string) (regen.Generator, error) {
	if len(pattern) == 0 {
		return nil, ErrEmptyPattern
	}

	keyHash := sha1.New()
	_, err := keyHash.Write([]byte(pattern))
	if err != nil {
		return nil, fmt.Errorf("unable to write pattern %q to hash container: %w", pattern, err)
	}
	key := string(keyHash.Sum(nil))

	gen, found := g[key]
	if !found {
		gen, err = regen.NewGenerator(pattern,
			&regen.GeneratorArgs{
				MinUnboundedRepeatCount: minUnboundedRepeatCount,
				MaxUnboundedRepeatCount: maxUnboundedRepeatCount,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("unable to create generator for pattern %q: %w", pattern, err)
		}
		g[key] = gen
	}
	return g[key], nil
}
