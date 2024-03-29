package throttle

import (
	"context"
	"fmt"
	"github.com/omihirofumi/cloud-native-pattern/retry"
	"sync"
	"time"
)

func Throttle(e retry.Effector, max uint, refill uint, d time.Duration) retry.Effector {
	var tokens = max
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)
			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						m.Lock()
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
						m.Unlock()
					}
				}
			}()
		})
		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}

		tokens--

		return e(ctx)
	}
}
