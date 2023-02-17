package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"tuxpa.in/a/zlog/log"
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

func (c *Cache) save(m map[string]CacheEntry) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}
	log.Trace().Str("tosave", string(payload)).Msg("saving cache")
	err = os.WriteFile(c.Root, payload, 0640)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetOrDo(key string, do func() (string, error), ttl time.Duration) (string, error) {
	conf, err := c.load()
	if err != nil {
		log.Trace().Err(err).Msg("cache failed read")
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
	log.Trace().Str("key", key).Str("val", value).Msg("saving new cache key")
	err = c.save(conf)
	if err != nil {
		log.Trace().Err(err).Msg("cache failed save")
	}
	return value, nil
}

func (c *Cache) Clear() error {
	return os.Remove(c.Root)
}
