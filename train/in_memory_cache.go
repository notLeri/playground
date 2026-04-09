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

type CustomerProfile struct {
	UUID   string
	Name   string
	Orders []*Order
}

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CacheItem struct {
	Profile   CustomerProfile
	ExpiresAt time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
	ttl  time.Duration
}

func NewCache(ctx context.Context, ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]CacheItem),
		ttl:  ttl,
	}
	go c.startCleaner(ctx)

	return c
}

var cleanerDelay = 10 * time.Second

// утилитарные методы

func (c *Cache) clean() {
	now := time.Now()

	var expiredKeys []string

	c.mu.RLock()
	for userID, item := range c.data {
		if now.After(item.ExpiresAt) {
			expiredKeys = append(expiredKeys, userID)
		}

	}
	c.mu.RUnlock()

	if len(expiredKeys) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, userID := range expiredKeys {
		if item, ok := c.data[userID]; ok && now.After(item.ExpiresAt) {
			delete(c.data, userID)
		}
	}
}

func (c *Cache) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(cleanerDelay)
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

func cloneOrders(orders []*Order) []*Order {
	if orders == nil {
		return nil
	}

	clone := make([]*Order, len(orders))

	for i, order := range orders {
		if order == nil {
			continue
		}

		copyOrder := *order

		clone[i] = &copyOrder
	}

	return clone
}

func (p CustomerProfile) Clone() CustomerProfile {
	return CustomerProfile{
		UUID:   p.UUID,
		Name:   p.Name,
		Orders: cloneOrders(p.Orders),
	}
}

// основные методы

func (c *Cache) Set(userId string, profile CustomerProfile) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[userId] = CacheItem{
		Profile:   profile.Clone(),
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(userId string) (CustomerProfile, bool) {
	c.mu.RLock()
	cacheItem, ok := c.data[userId]

	if !ok || time.Now().After(cacheItem.ExpiresAt) {
		c.mu.RUnlock()
		return CustomerProfile{}, false
	}

	profile := cacheItem.Profile.Clone()
	c.mu.RUnlock()

	return profile, true
}
