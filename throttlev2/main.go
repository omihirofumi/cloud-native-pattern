package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"time"
)

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
	e := echo.New()
	e.GET("/hostname", throttledHandler)
	log.Fatal(e.Start(":8080"))
}
