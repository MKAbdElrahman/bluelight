package middleware

import (
	"net/http"
	"sync"
	"time"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"github.com/tomasen/realip" // New import
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	Burst   int
	RPS     float64
	Enabled bool
}

func RateLimiter(em *errorhandler.ErrorHandeler, cfg RateLimiterConfig) middlewareFunc {
	return func(next http.Handler) http.Handler {

		type client struct {
			limiter  *rate.Limiter
			lastSeen time.Time
		}

		var (
			mu      sync.Mutex
			clients = make(map[string]*client)
		)

		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Enabled {
				ip := realip.FromRequest(r)
				mu.Lock()
				if _, found := clients[ip]; !found {
					clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(cfg.RPS), cfg.Burst)}
				}
				clients[ip].lastSeen = time.Now()

				if !clients[ip].limiter.Allow() {
					mu.Unlock()
					em.SendClientError(w, r, apierror.TooManyRequestsError)
					return
				}
				mu.Unlock()
			}
			next.ServeHTTP(w, r)
		})
	}
}
