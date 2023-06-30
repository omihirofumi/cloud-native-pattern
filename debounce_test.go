package chap04

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestDebounceFirstDataRace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	circuit := failAfter(1)
	debounce := DebounceFirst(circuit, time.Second)

	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}
	time.Sleep(time.Second * 2)
	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}
	wg.Wait()
}

func TestDebounceLastDataRace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	debounce := DebounceLast(counter(), time.Second)
	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()

	t.Log("Waiting 2 seconds")

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()
}
