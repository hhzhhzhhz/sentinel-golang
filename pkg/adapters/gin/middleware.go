package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sentinel "github.com/hhzhhzhhz/sentinel-golang/api"
	"github.com/hhzhhzhhz/sentinel-golang/core/base"
)

// SentinelMiddleware returns new gin.HandlerFunc
// Default resource name is {method}:{path}, such as "GET:/api/users/:id"
// Default block fallback is returning 429 code
// Define your own behavior by setting options
func SentinelMiddleware(opts ...Option) gin.HandlerFunc {
	options := evaluateOptions(opts)
	return func(c *gin.Context) {
		resourceName := c.Request.Method + ":" + c.FullPath()

		if options.resourceExtract != nil {
			resourceName = options.resourceExtract(c)
		}

		entry, err := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)

		if err != nil {
			if options.blockFallback != nil {
				options.blockFallback(c)
			} else {
				c.AbortWithStatus(http.StatusTooManyRequests)
			}
			return
		}

		defer entry.Exit()
		c.Next()
	}
}
