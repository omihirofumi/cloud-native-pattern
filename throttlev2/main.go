package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

const MAX_QUEUE = 10

var throttled = Throttle(getHostname, 1, 1, time.Second)

func getHostname(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return os.Hostname()
}

func throttledHandler(c echo.Context) error {
	ok, hostname, err := throttled(c.Request().Context(), "test")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusTooManyRequests, "Too many requests")
	}

	return c.String(http.StatusOK, hostname)
}

func main() {
	var queueCnt int64
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			atomic.AddInt64(&queueCnt, 1)
			defer atomic.AddInt64(&queueCnt, -1)

			if atomic.LoadInt64(&queueCnt) > MAX_QUEUE {
				return echo.NewHTTPError(http.StatusTooManyRequests, "server is busy.")
			}
			return next(c)
		}
	})
	e.GET("/hostname", throttledHandler)
	log.Fatal(e.Start(":8080"))
}
