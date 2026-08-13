package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type LRUCache interface {
	Put(key string, value string)
	Get(key string) (string, error)
}

type lruCache struct {
	keyvalues map[string]lruCacheValue
	keyLocks  map[string]*sync.RWMutex
	maxSize   int
	mapLock   *sync.Mutex
}

type lruCacheValue struct {
	content  string
	lastUsed int64
}

func (lc *lruCache) Put(key string, value string) {
	lc.mapLock.Lock()
	keyLock, ok := lc.keyLocks[key]
	if !ok {
		lc.keyLocks[key] = &sync.RWMutex{}
		keyLock = lc.keyLocks[key]
	}

	if lc.maxSize == len(lc.keyLocks) {
		k := getLeastRecentlyUsed(lc.keyvalues)

		delKey := lc.keyLocks[k]
		lc.mapLock.Unlock()

		delKey.Lock()
		delete(lc.keyvalues, k)
		delete(lc.keyLocks, k)
		delKey.Unlock()
	}
	lc.mapLock.Unlock()

	keyLock.Lock()
	val, ok := lc.keyvalues[key]
	if ok {
		val.lastUsed = time.Now().Unix()
	}
	val.content = value
	lc.keyvalues[key] = val
	keyLock.Unlock()
}

func getLeastRecentlyUsed(keyvalues map[string]lruCacheValue) string {
	var minAccess int64 = -1
	minKey := ""
	for k, v := range keyvalues {
		if v.lastUsed <= minAccess {
			minKey = k
		}
	}
	return minKey
}

func (lc *lruCache) Get(key string) (string, error) {

	if lc.keyLocks[key] == nil {
		return "", errors.New("key not found")
	}

	lc.keyLocks[key].RLock()
	defer lc.keyLocks[key].RUnlock()
	val, ok := lc.keyvalues[key]
	if !ok {
		return "", errors.New("key not found")
	}
	atomic.StoreInt64(&val.lastUsed, time.Now().Unix())
	return val.content, nil
}

func NewLRUCache(maxSize int) LRUCache {
	return &lruCache{
		keyvalues: make(map[string]lruCacheValue),
		keyLocks:  make(map[string]*sync.RWMutex),
		maxSize:   maxSize,
		mapLock:   &sync.Mutex{},
	}
}
