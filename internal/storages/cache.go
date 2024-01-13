package storages

import (
	"log"
	"sync"
	"webapp/internal/models"
)

type Cache struct {
	data  map[string][]byte
	mutex *sync.RWMutex
}

func InitCache() *Cache {
	return &Cache{
		make(map[string][]byte),
		new(sync.RWMutex),
	}
}

func (c *Cache) goRlock(f func()) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	f()
}

func (c *Cache) goLock(f func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	f()
}

func (c *Cache) Add(key string, value []byte) {
	if len(key) < 1 {
		log.Printf("(Cache): error when adding value %s to cache", key)
	}
	c.goLock(func() {
		if _, ex := c.data[key]; ex {
			log.Printf("(Cache): key %s already exists", key)
			return
		}
		c.data[key] = value
	})

}

func (c *Cache) AddOrders(in []models.Order) {
	for _, ord := range in {
		c.Add(ord.ID, ord.Data)
	}
}

func (c *Cache) Get(key string) (data []byte, found bool) {
	c.goRlock(func() {
		data, found = c.data[key]
	})
	return data, found
}

func (c *Cache) GetAll() (keys []string) {
	c.goRlock(func() {
		for key := range c.data {
			keys = append(keys, key)
		}
	})
	return keys
}
