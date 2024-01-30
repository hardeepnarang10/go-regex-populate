package goregexpopulate

import (
	"crypto/sha1"
	"fmt"
	"sync"

	regen "github.com/zach-klippenstein/goregen"
)

type generatorMap struct {
	store sync.Map
}

var cache generatorMap = generatorMap{}

func (g *generatorMap) register(pattern string) (regen.Generator, error) {
	if len(pattern) == 0 {
		return nil, ErrEmptyPattern
	}

	keyHash := sha1.New()
	_, err := keyHash.Write([]byte(pattern))
	if err != nil {
		return nil, fmt.Errorf("unable to write pattern %q to hash container: %w", pattern, err)
	}
	key := string(keyHash.Sum(nil))

	generator, found := g.store.Load(key)
	if !found {
		generator, err = regen.NewGenerator(pattern,
			&regen.GeneratorArgs{
				MinUnboundedRepeatCount: minUnboundedRepeatCount,
				MaxUnboundedRepeatCount: maxUnboundedRepeatCount,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("unable to create generator for pattern %q: %w", pattern, err)
		}
		g.store.Store(key, generator)
	}
	return generator.(regen.Generator), nil
}
