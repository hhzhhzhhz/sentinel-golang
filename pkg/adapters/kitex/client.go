package kitex

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	sentinel "github.com/hhzhhzhhz/sentinel-golang/api"
	"github.com/hhzhhzhhz/sentinel-golang/core/base"
)

// SentinelClientMiddleware returns new client.Middleware
// Default resource name is {service's name}:{method}
// Default block fallback is returning blockError
// Define your own behavior by setting serverOptions
func SentinelClientMiddleware(opts ...Option) func(endpoint.Endpoint) endpoint.Endpoint {
	options := newOptions(opts)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			resourceName := options.ResourceExtract(ctx, req, resp)
			entry, blockErr := sentinel.Entry(
				resourceName,
				sentinel.WithResourceType(base.ResTypeRPC),
				sentinel.WithTrafficType(base.Outbound),
			)
			if blockErr != nil {
				return options.BlockFallback(ctx, req, resp, blockErr)
			}
			defer entry.Exit()
			err := next(ctx, req, resp)
			if err != nil {
				sentinel.TraceError(entry, err)
			}
			return err
		}
	}
}
