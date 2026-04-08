package main

import (
	"context"
	"sync"
	"time"
)

// Необходимо написать in-memory кэш, который будет по ключу (uuid пользователя) возвращать профиль и список его заказов.

// 1. У кэша должен быть TTL (2 сек)
// 2. Кэшем может пользоваться функция(-и), которая работает с заказами (добавляет/обновляет/удаляет). Если TTL истек, то возвращается nil. При апдейте TTL снова устанавливается 2 сек. Методы должны быть потокобезопасными
// 3. Должны быть написаны тестовые сценарии использования данного кэша
// (базовые структуры не менять)

type UserProfile struct {
	ID string
	Name string
}

type Order struct {
	ID string
	Amount int
}

type CacheItem struct {
	Profile UserProfile
	Orders []Order
	ExpiresAt time.Time
}

type Cache struct {
	mu sync.RWMutex
	data map[string]CacheItem
	ttl time.Duration
	done chan struct{}
}

func NewCache(ctx context.Context, ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]CacheItem, 10),
		ttl: ttl,
	}
	go c.startCleaner(ctx)

	return c
}

func (c *Cache) clean() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for userId, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, userId)
		}
	}
}

func (c *Cache) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.clean()
		case <-ctx.Done():
            return
		}
	}
}

func (c *Cache) Set(userId string, orders []Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[userId] = CacheItem{
		Orders: orders,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(userId string) (CacheItem, bool) {
	c.mu.RLock()
	item, ok := c.data[userId]
	c.mu.RUnlock()

	if !ok {
		return CacheItem{}, false
	}

	if time.Now().After(item.ExpiresAt) {
        c.mu.Lock()
        defer c.mu.Unlock()

        if item2, ok2 := c.data[userId]; ok2 && time.Now().After(item2.ExpiresAt) {
            delete(c.data, userId)
        }
        return CacheItem{}, false
    }

	return item, true
}
