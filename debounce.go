package chap04

import (
	"context"
	"sync"
	"time"
)

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			// 呼び出された時間+指定時間
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			// 指定時間経過前の場合、前回の結果とエラーを返却
			return result, err
		}

		// 初回または前回呼び出しから、指定時間経過後のみ呼び出される
		result, err = circuit(ctx)
		return result, err
	}
}

func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time = time.Now()
	var ticker *time.Ticker
	var result string
	var err error
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		threshold = time.Now().Add(d)

		// リセットされるまでは一回のみの実行を保証
		once.Do(func() {
			// 100msごとにチェック
			ticker = time.NewTicker(time.Millisecond * 100)

			go func() {
				defer func() {
					m.Lock()
					ticker.Stop()
					once = sync.Once{}
					m.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						m.Lock()
						if time.Now().After(threshold) {
							// 指定時間経過していた場合
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}

			}()
		})
		return result, err
	}
}
