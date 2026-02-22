package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimit(requestsPerSecond float64, burst int) gin.HandlerFunc {
	visitors := map[string]*visitor{}
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		if _, ok := visitors[ip]; !ok {
			visitors[ip] = &visitor{limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst)}
		}
		visitors[ip].lastSeen = time.Now()
		limiter := visitors[ip].limiter
		mu.Unlock()

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			httputil.TooManyRequests(c, "rate_limited", fmt.Sprintf("rate limit exceeded for %s", ip))
			c.Abort()
			return
		}
		c.Next()
	}
}
