package api

import (
	"fmt"
	"time"
)

type TimerMiddleware struct{}

// Calculate execution time of handler
func (tm *TimerMiddleware) MiddlewareFunc(handler HandlerFunc) HandlerFunc {
	fmt.Println("--> <TimerMiddleware:MiddlewareFunc>")
	return func(w ResponseWriter, r *Request) {
		startTime := time.Now()
		r.environment["START_TIME"] = &startTime
		handler(w, r)
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		r.environment["ELAPSED_TIME"] = &elapsedTime
	}
}
