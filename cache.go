package autostruct

import (
	"reflect"
	"sync"
)

type cache struct {
	lock sync.RWMutex
	vals map[string]reflect.Value
}

func NewCache() *cache {
	return &cache{
		vals: make(map[string]reflect.Value),
	}
}

func (c *cache) get(key string) (reflect.Value, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.vals[key]
	return val, ok
}

func (c *cache) set(key string, val reflect.Value) {
	c.lock.Lock()
	c.vals[key] = val
	c.lock.Unlock()
}
