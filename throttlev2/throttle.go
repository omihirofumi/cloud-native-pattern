package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru/v2"
	"time"
)

const LRU_SIZE = 100

type Effector func(ctx context.Context) (string, error)

type Throttled func(context.Context, string) (bool, string, error)

type bucket struct {
	tokens uint
	time   time.Time
}

func Throttle(e Effector, max uint, refill uint, d time.Duration) Throttled {
	buckets, _ := lru.New[string, *bucket](LRU_SIZE)

	return func(ctx context.Context, uid string) (bool, string, error) {
		b, ok := buckets.Get(uid)

		if !ok {
			buckets.Add(uid, &bucket{tokens: max - 1, time: time.Now()})

			str, err := e(ctx)
			return true, str, err
		}

		refillsSince := uint(time.Since(b.time) / d)
		tokensAddedSince := refill * refillsSince
		currentTokens := b.tokens + tokensAddedSince

		if currentTokens < 1 {
			return false, "", nil
		}

		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
		} else {
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d

			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}

		str, err := e(ctx)

		return true, str, err
	}
}
