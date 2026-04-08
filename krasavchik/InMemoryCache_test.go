package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCache_NoDeadlock(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    cache := NewCache(ctx, 100*time.Millisecond)
    
    var wg sync.WaitGroup
    wg.Add(1000)
    
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            cache.Set("user1", []Order{})
            cache.Get("user1")
        }()
    }
    wg.Wait()
}

func TestCache_SetGet(t *testing.T) {
    ctx := context.Background()
    cache := NewCache(ctx, 2 * time.Second)

	_, ok := cache.Get("unknown")
	if ok {
		t.Fatal("expected not exists")
	}

    userID := "123"

    orders := []Order{{ID: "o1", Amount: 100}}

    cache.Set(userID, orders)

    _, ok1 := cache.Get(userID)
    if !ok1 {
        t.Fatal("expected item")
    }
}

func TestCache_TTLExpire(t *testing.T) {
    ctx := context.Background()
    cache := NewCache(ctx, 2 * time.Second)

    userID := "123"

    cache.Set(userID, nil)

    time.Sleep(4 * time.Second)

    _, ok := cache.Get(userID)
    if ok {
        t.Fatal("expected expired")
    }
}

func TestCache_Concurrent(t *testing.T) {
    ctx := context.Background()
    cache := NewCache(ctx, 2 * time.Second)

    userID := "123"
	u1 := "1"
	u2 := "2"
	u3 := "3"

    go cache.Set(userID, nil)
    go cache.Get(userID)
	go cache.Set(u1, nil)
	go cache.Set(u2, nil)
    go cache.Get(u3)
    go cache.Get(u1)
	go cache.Set(u3, nil)

    time.Sleep(100 * time.Millisecond)
}

func TestCache_HeavyLoad(t *testing.T) {
    ctx := context.Background()
    cache := NewCache(ctx, 2 * time.Second)

    userID := "123"

    for i := 0; i < 1000; i++ {
        go func(i int) {
            cache.Set(userID, nil)
            cache.Get(userID)
        }(i)
    }

    time.Sleep(2 * time.Second)
}

func BenchmarkCache_Get(b *testing.B) {
    ctx := context.Background()
    cache := NewCache(ctx, 2 * time.Second)

    userID := "123"
    cache.Set(userID, nil)

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        cache.Get(userID)
    }
}