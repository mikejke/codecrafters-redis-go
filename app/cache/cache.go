package cache

import (
	"fmt"
	"time"
)

type Item struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}

type Cache struct {
	cache map[string]interface{}
	TTL   time.Duration
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]interface{}),
		TTL:   time.Hour,
	}
}

func (c *Cache) Set(item *Item) {
	duration := c.TTL
	if item.TTL != 0 {
		duration = item.TTL
	}

	go c.ScheduleExpirtaion(item.Key, duration)()
	c.cache[item.Key] = item.Value
}

func (c *Cache) ScheduleExpirtaion(key string, duration time.Duration) func() {
	return func() {
		time.Sleep(time.Duration(duration.Milliseconds()))
		fmt.Printf("delete expired key %s from cache \n", key)
		delete(c.cache, key)
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	value, ok := c.cache[key]
	return value, ok
}
