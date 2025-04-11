package cache

import (
	"sync"
	"time"
)

type Cache struct {
	keys            map[string]cacheTTL
	keepTimeSeconds int64
	checkCacheTime  int64
	mu              sync.RWMutex
}

type cacheTTL struct {
	ttl       int64
	cacheData []byte
}

func (c *Cache) NewCache(keepTimeSeconds int64, checkCacheTime int64) error {
	c.keepTimeSeconds = keepTimeSeconds
	c.checkCacheTime = checkCacheTime
	c.keys = map[string]cacheTTL{}
	go c.pruneCache()
	return nil
}

func (c *Cache) AddCache(newKey string, rawData []byte) error {
	var newTTL cacheTTL
	//10 second cache
	newTTL.ttl = time.Now().UnixMilli() + (c.keepTimeSeconds * 1000)
	//fmt.Printf("Adding key: %s, Time: %d, TTL: %d\n", newKey, time.Now().UnixMilli(), newTTL.ttl)
	newTTL.cacheData = rawData

	//Mutex the write
	c.mu.Lock()
	defer c.mu.Unlock()

	c.keys[newKey] = newTTL
	return nil
}

func (c *Cache) CacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.keys)
}

func (c *Cache) GetCache(key string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cData, ok := c.keys[key]

	if !ok {
		return []byte{}
	}

	return cData.cacheData
}

func (c *Cache) pruneCache() {
	for {
		for key, cObj := range c.keys {
			if cObj.ttl < time.Now().UnixMilli() {
				c.mu.Lock()
				//fmt.Printf("Delete key: %s\n", key)
				delete(c.keys, key)
				c.mu.Unlock()
			}
		}
		//End Mutex
		time.Sleep(time.Second * time.Duration(c.checkCacheTime))
	}

}
