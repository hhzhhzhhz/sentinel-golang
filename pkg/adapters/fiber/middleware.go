package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	sentinel "github.com/hhzhhzhhz/sentinel-golang/api"
	"github.com/hhzhhzhhz/sentinel-golang/core/base"
)

// SentinelMiddleware returns new gin.HandlerFunc
// Default resource name is {method}:{path}, such as "GET:/api/users/:id"
// Default block fallback is returning 429 code
// Define your own behavior by setting options
func SentinelMiddleware(opts ...Option) fiber.Handler {
	options := evaluateOptions(opts)
	return func(ctx *fiber.Ctx) error {
		resourceName := ctx.Route().Method + ":" + string(ctx.Context().Path())

		if options.resourceExtract != nil {
			resourceName = options.resourceExtract(ctx)
		}

		entry, entryErr := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)

		if entryErr != nil {
			if options.blockFallback != nil {
				return options.blockFallback(ctx)
			} else {
				return ctx.SendStatus(http.StatusTooManyRequests)
			}
		}

		defer entry.Exit()
		return ctx.Next()
	}
}
