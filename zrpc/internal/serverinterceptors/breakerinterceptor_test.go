package serverinterceptors

import (
	"context"
	"testing"

	"github.com/jialequ/linux-sdk/core/breaker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStreamBreakerInterceptor(t *testing.T) {
	err := StreamBreakerInterceptor(nil, nil, &grpc.StreamServerInfo{
		FullMethod: "any",
	}, func(_ any, _ grpc.ServerStream) error {
		return status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}

func TestUnaryBreakerInterceptor(t *testing.T) {
	_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(_ context.Context, _ any) (any, error) {
		return nil, status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}

func TestUnaryBreakerInterceptorUnavailable(t *testing.T) {
	_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(_ context.Context, _ any) (any, error) {
		return nil, breaker.ErrServiceUnavailable
	})
	assert.NotNil(t, err)
}
