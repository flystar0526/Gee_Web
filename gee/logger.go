package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		// Start timer
		startTime := time.Now()
		// Process request
		ctx.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(startTime))
	}
}
