// middleware/logging.go
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}).Info("Request started")

		next.ServeHTTP(w, r)

		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
		}).Info("Request completed")
	})
}
