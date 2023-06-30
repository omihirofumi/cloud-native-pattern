package chap04

import (
	"context"
	"log"
	"time"
)

type Effector func(ctx context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}
			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			// ctx.Done()でキャンセルしたいので
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
