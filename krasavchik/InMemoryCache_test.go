package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCache_NoDeadlock(t *testing.T) {
	mockProfile := CustomerProfile{UUID: "", Name: "", Orders: []*Order{}}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := NewCache(ctx, 100*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			cache.Set("user1", mockProfile)
			cache.Get("user1")
		}()
	}
	wg.Wait()
}

func TestCache_SetGet(t *testing.T) {
	ctx := context.Background()
	cache := NewCache(ctx, 2*time.Second)

	_, ok := cache.Get("unknown")
	if ok {
		t.Fatal("expected not exists")
	}

	userID := "123"
	now := time.Now()

	orders := []*Order{{UUID: "o1", CreatedAt: now, UpdatedAt: now}}
	profile := CustomerProfile{UUID: userID, Name: "", Orders: orders}

	cache.Set(userID, profile)

	_, ok1 := cache.Get(userID)
	if !ok1 {
		t.Fatal("expected item")
	}
}

func TestCache_TTLExpire(t *testing.T) {
	ctx := context.Background()
	cache := NewCache(ctx, 2*time.Second)

	userID := "123"

	cache.Set(userID, CustomerProfile{})

	time.Sleep(4 * time.Second)

	_, ok := cache.Get(userID)
	if ok {
		t.Fatal("expected expired")
	}
}

func TestCache_Concurrent(t *testing.T) {
	ctx := context.Background()
	cache := NewCache(ctx, 2*time.Second)

	userID := "123"
	u1 := "1"
	u2 := "2"
	u3 := "3"

	go cache.Set(userID, CustomerProfile{})
	go cache.Get(userID)
	go cache.Set(u1, CustomerProfile{})
	go cache.Set(u2, CustomerProfile{})
	go cache.Get(u3)
	go cache.Get(u1)
	go cache.Set(u3, CustomerProfile{})

	time.Sleep(100 * time.Millisecond)
}

func TestCache_HeavyLoad(t *testing.T) {
	ctx := context.Background()
	cache := NewCache(ctx, 2*time.Second)

	userID := "123"

	for i := 0; i < 1000; i++ {
		go func(i int) {
			cache.Set(userID, CustomerProfile{})
			cache.Get(userID)
		}(i)
	}

	time.Sleep(2 * time.Second)
}

func TestCache_Immutable(t *testing.T) {
	cache := NewCache(context.Background(), 2*time.Second)

	order := &Order{
		UUID: "original",
	}

	profile := CustomerProfile{
		UUID:   "user",
		Name:   "test",
		Orders: []*Order{order},
	}

	cache.Set("user", profile)

	item, _ := cache.Get("user")

	item.Orders[0].UUID = "HACK"

	item2, _ := cache.Get("user")

	if item2.Orders[0].UUID == "HACK" {
		t.Fatal("cache was mutated")
	}
}

func TestCache_SetIsolation(t *testing.T) {
	cache := NewCache(context.Background(), 2*time.Second)

	order := &Order{UUID: "original"}

	profile := CustomerProfile{
		UUID:   "user",
		Orders: []*Order{order},
	}

	cache.Set("user", profile)

	profile.Orders[0].UUID = "HACK"

	item, _ := cache.Get("user")

	if item.Orders[0].UUID == "HACK" {
		t.Fatal("cache was affected by external mutation")
	}
}
