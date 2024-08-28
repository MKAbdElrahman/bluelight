package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
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
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					em.SendServerError(w, r, v1.ServerError{
						InternalMessage: fmt.Sprintf("%s", err),
						Code:            http.StatusInternalServerError,
					})
					return
				}
				mu.Lock()
				if _, found := clients[ip]; !found {
					clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(cfg.RPS), cfg.Burst)}
				}
				clients[ip].lastSeen = time.Now()

				if !clients[ip].limiter.Allow() {
					mu.Unlock()
					em.SendClientError(w, r, v1.TooManyRequestsError)
					return
				}
				mu.Unlock()
			}
			next.ServeHTTP(w, r)
		})
	}
}
