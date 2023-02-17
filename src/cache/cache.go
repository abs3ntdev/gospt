package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	Root string
}

type CacheEntry struct {
	Expire time.Time `json:"e"`
	Value  string    `json:"v"`
}

func DefaultCache() *Cache {
	return &Cache{
		Root: filepath.Join(os.TempDir(), "gospt.cache"),
	}
}

func (c *Cache) load() (map[string]CacheEntry, error) {
	out := map[string]CacheEntry{}
	cache, err := os.Open(c.Root)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(cache).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Cache) save(map[string]CacheEntry) error {
	out := map[string]CacheEntry{}
	payload, err := json.Marshal(out)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.Root, payload, 0640)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetOrDo(key string, do func() (string, error), ttl time.Duration) (string, error) {
	conf, err := c.load()
	if err != nil {
		return c.Do(key, do, ttl)
	}
	val, ok := conf[key]
	if !ok {
		return c.Do(key, do, ttl)
	}
	if time.Now().After(val.Expire) {
		return c.Do(key, do, ttl)
	}
	return val.Value, nil
}

func (c *Cache) Do(key string, do func() (string, error), ttl time.Duration) (string, error) {
	if do == nil {
		return "", nil
	}
	res, err := do()
	if err != nil {
		return "", err
	}
	return c.Put(key, res, ttl)
}
func (c *Cache) Put(key string, value string, ttl time.Duration) (string, error) {
	conf, err := c.load()
	if err != nil {
		conf = map[string]CacheEntry{}
	}
	conf[key] = CacheEntry{
		Expire: time.Now().Add(ttl),
		Value:  value,
	}
	return value, nil
}
