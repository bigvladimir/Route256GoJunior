package in_memory_cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"homework/internal/app/pvz/dto"
	"homework/internal/app/pvz/validation"
)

type pvzCacheUnit struct {
	pvz          dto.Pvz
	creationTime time.Time
}

type storage interface {
	GetByID(context.Context, int64) (dto.Pvz, error)
}

type cacheUpdaterOps interface {
	Subscribe(ctx context.Context, topic string, output chan<- int64) error
}

// InMemoryCache provides operations for working with database methods cache
type InMemoryCache struct {
	pvzType   map[int64]pvzCacheUnit
	pvzTypeMX sync.RWMutex

	stor         storage
	cacheUpdater cacheUpdaterOps
}

// NewInMemoryCache creates InMemoryCache
func NewInMemoryCache(ctx context.Context, stor storage, cacheUpdater cacheUpdaterOps) (*InMemoryCache, error) {
	cache := &InMemoryCache{
		pvzType:      make(map[int64]pvzCacheUnit, 1000),
		stor:         stor,
		cacheUpdater: cacheUpdater,
	}

	go clearOldCache(ctx, cache)
	cacheUpdates := make(chan int64)
	if err := cache.cacheUpdater.Subscribe(ctx, "pvz_cache_updater", cacheUpdates); err != nil {
		return nil, fmt.Errorf("CacheUpdater.Subscribe: %w", err)
	}
	go updateInvalidCache(ctx, cache, cacheUpdates)

	return cache, nil
}

// GetPvzType tries to get cached the Pvz object, returns false if not
func (c *InMemoryCache) GetPvzType(id int64) (dto.Pvz, bool) {
	c.pvzTypeMX.RLock()
	defer c.pvzTypeMX.RUnlock()
	pvzObj, ok := c.pvzType[id]
	if !ok {
		return dto.Pvz{}, false
	}

	return pvzObj.pvz, true
}

// SetPvzType writes the Pvz object in cache
func (c *InMemoryCache) SetPvzType(pvzObj dto.Pvz) {
	c.pvzTypeMX.Lock()
	defer c.pvzTypeMX.Unlock()

	c.pvzType[pvzObj.ID] = pvzCacheUnit{
		pvz:          pvzObj,
		creationTime: time.Now(),
	}
}

// clearOldCache clears out-of-date cache after certain periods of time
func clearOldCache(ctx context.Context, cache *InMemoryCache) {
	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			cache.pvzTypeMX.Lock()
			defer cache.pvzTypeMX.Unlock()
			for key, v := range cache.pvzType {
				if time.Since(v.creationTime).Hours() > 3 {
					delete(cache.pvzType, key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// updateInvalidCache updates or deletes an invalid cache using kafka messages
func updateInvalidCache(ctx context.Context, cache *InMemoryCache, updatesChan <-chan int64) {
	for {
		select {
		case id := <-updatesChan:
			if err := validation.ValidatePvzID(id); err != nil {
				log.Println("validation.ValidatePvzID:", err)
				continue
			}

			pvzObj, err := cache.stor.GetByID(ctx, id)
			if err != nil {
				cache.pvzTypeMX.Lock()
				defer cache.pvzTypeMX.Unlock()
				delete(cache.pvzType, id)
				continue
			}

			cache.SetPvzType(pvzObj)
		case <-ctx.Done():
			return
		}
	}
}
