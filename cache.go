package cache

import "time"

type EmptyStruct struct{}

type Cache struct {
	keyValues   map[string]string
	deadlineMap map[string]chan EmptyStruct
}

func NewCache() Cache {
	keyValues := make(map[string]string)
	deadlineMap := make(map[string]chan EmptyStruct)
	return Cache{
		keyValues,
		deadlineMap,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	value, ok := c.keyValues[key]
	return value, ok
}

func (c *Cache) Put(key, value string) {
	c.keyValues[key] = value
}

func (c *Cache) Keys() []string {
	result := []string{}
	for k := range c.keyValues {
		result = append(result, k)
	}

	return result
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	now := time.Now()
	if now.Before(deadline) {
		if deadlineChan, ok := c.deadlineMap[key]; ok {
			deadlineChan <- EmptyStruct{}
		}
		c.Put(key, value)
		go func() {
			deadlineChan := make(chan EmptyStruct)
			c.deadlineMap[key] = deadlineChan
			timer := time.NewTimer(deadline.Sub(now))
			select {
			case <-deadlineChan:
				return
			case <-timer.C:
				delete(c.keyValues, key)
				delete(c.deadlineMap, key)
			}
		}()
	}
}
