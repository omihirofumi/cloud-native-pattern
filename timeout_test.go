package chap04

import (
	"context"
	"testing"
	"time"
)

func takeLongTime(arg string) (string, error) {
	time.Sleep(2 * time.Second)
	return arg, nil
}

func TestTimeout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	withContext := Timeout(takeLongTime)
	_, err := withContext(ctx, "hello")
	if err != nil {
		t.Log("got expected error:", err)
	} else {
		t.Error("didn't get expected error")
	}
}

func TestInTime(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	withContext := Timeout(takeLongTime)
	res, err := withContext(ctx, "hello")
	if err != nil {
		t.Error("got error:", err)
	} else {
		t.Log("get response:", res)
	}
}
