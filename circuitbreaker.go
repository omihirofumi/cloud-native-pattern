package chap04

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()

		d := consecutiveFailures - int(failureThreshold)
		// 失敗カウントが失敗閾値を超えた場合
		if d >= 0 {
			// リトライ時間を設定
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			// リトライ時間より、時間が立っていれば復活の可能性
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0
		return response, nil
	}
}
