// Package ginzap provides log handling using zap package.
// Code structure based on ginrus package.
package ginzap

import (
	"github.com/mygomod/muses/pkg/logger"
	"github.com/mygomod/muses/pkg/prom"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Ginzap(timeFormat string, utc bool, enabledMetric bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}
		// 作用域
		if len(c.Errors) > 0 {
			if enabledMetric {
				prom.HTTPServerTimer.Timing(getReqPath(c), int64(latency/time.Millisecond))
				prom.HTTPServerCounter.Incr(getReqPath(c), strconv.FormatInt(int64(c.Writer.Status()), 10))
			}

			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.DefaultLogger().Error(e)
			}
		} else {
			if enabledMetric {
				prom.HTTPServerTimer.Timing(getReqPath(c), int64(latency/time.Millisecond))
				prom.HTTPServerCounter.Incr(getReqPath(c), strconv.FormatInt(int64(c.Writer.Status()), 10))
			}
			logger.DefaultLogger().Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("reqPath", getReqPath(c)),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
			)
		}
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.DefaultLogger().Error(c.Request.URL.Path,
						zap.String("error", err.(string)),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.DefaultLogger().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.DefaultLogger().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func getReqPath(c *gin.Context) string {
	pathArr := strings.Split(c.Request.URL.Path, "/")
	for i := len(pathArr) - 1; i >= 0; i-- {
		if pathArr[i] == "" {
			pathArr = append(pathArr[:i], pathArr[i+1:]...)
		}
	}
	for i, path := range pathArr {
		if matched, err := regexp.MatchString("^[0-9]+$", path); matched && err == nil {
			pathArr[i] = "id"
		}
	}
	pathArr = append([]string{strings.ToLower(c.Request.Method)}, pathArr...)
	return strings.Join(pathArr, "_")
}
