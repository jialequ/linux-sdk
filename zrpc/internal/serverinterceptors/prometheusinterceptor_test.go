package serverinterceptors

import (
	"context"
	"testing"

	"github.com/jialequ/linux-sdk/core/prometheus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryPromMetricInterceptorDisabled(t *testing.T) {
	_, err := UnaryPrometheusInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryPromMetricInterceptorEnabled(t *testing.T) {
	prometheus.StartAgent(prometheus.Config{
		Host: "localhost",
		Path: "/",
	})
	_, err := UnaryPrometheusInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}
