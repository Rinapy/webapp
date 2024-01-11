package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	orders            map[string]Order
}

type Order struct {
	Value      interface{} // Тут надо подумать ((
	Created    time.Time
	Expiration int64
}

func InitCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	orders := make(map[string]Order)

	cache := Cache{
		orders:            orders,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.orders[key] = Order{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	order, found := c.orders[key]

	if !found {
		return nil, false
	}

	if order.Expiration > 0 {
		if time.Now().UnixNano() > order.Expiration {
			return nil, false
		}
	}

	return order.Value, true
}

func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.orders[key]; !found {
		return errors.New("Key not found in cache.")
	}

	delete(c.orders, key)

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.orders == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearOrders(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.orders {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearOrders(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.orders, k)
	}
}
