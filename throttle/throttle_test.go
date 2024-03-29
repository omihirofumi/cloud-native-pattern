package throttle

import (
	"context"
	"fmt"
	"github.com/omihirofumi/cloud-native-pattern/retry"
	"testing"
	"time"
)

func callsCountFunction(callCounter *int) retry.Effector {
	return func(ctx context.Context) (string, error) {
		*callCounter++
		return fmt.Sprintf("call %d", *callCounter), nil
	}
}

func TestThrottleMax1(t *testing.T) {
	t.Parallel()
	const max uint = 1

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is bloken; got", callsCounter)
	}
	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}

}

func TestThrottleMax10(t *testing.T) {
	t.Parallel()
	const max uint = 10

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		throttle(ctx)
	}
	if callsCounter == 0 {
		t.Error("test is bloken; got", callsCounter)
	}
	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

func TestThrottleCallFrequency5Seconds(t *testing.T) {
	t.Parallel()
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 1, 1, time.Second)

	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond).C

	for range ticker {
		tickCounts++

		s, e := throttle(ctx)
		if e != nil {
			t.Log("Error", e)
		} else {
			t.Log("output:", s)
		}

		if tickCounts >= 20 {
			break
		}
	}

	if callsCounter != 5 {
		t.Error("expected 5; got", callsCounter)
	}
}

func TestThrottleVariableRefill(t *testing.T) {
	t.Parallel()
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 4, 2, 500*time.Millisecond)

	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond)
	timer := time.NewTimer(2 * time.Second)

time:
	for {
		select {
		case <-ticker.C:
			tickCounts++

			s, e := throttle(ctx)
			if e != nil {
				t.Log("Error:", e)
			} else {
				t.Log("output:", s)
			}
		case <-timer.C:
			break time
		}
	}

	if callsCounter != 8 {
		t.Error("expected 8; got", callsCounter)
	}
}

func TestThrottleContextTimeout(t *testing.T) {
	t.Parallel()
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	throttle := Throttle(effector, 1, 1, time.Second)

	s, e := throttle(ctx)
	if e != nil {
		t.Error("unexpected error:", e)
	} else {
		t.Log("output:", s)
	}

	time.Sleep(300 * time.Millisecond)

	_, e = throttle(ctx)
	if e != nil {
		t.Log("got expected error:", e)
	} else {
		t.Error("didn't get expected error")
	}
}
