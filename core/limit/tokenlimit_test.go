package limit

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jialequ/linux-sdk/core/logx"
	"github.com/jialequ/linux-sdk/core/stores/redis"
	"github.com/jialequ/linux-sdk/core/stores/redis/redistest"
	"github.com/stretchr/testify/assert"
)

func init() {
	logx.Disable()
}

func TestTokenLimitWithCtx(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.New(s.Addr()), "tokenlimit")
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ok := l.AllowCtx(ctx)
	assert.True(t, ok)

	cancel()
	for i := 0; i < total; i++ {
		ok := l.AllowCtx(ctx)
		assert.False(t, ok)
		assert.False(t, l.monitorStarted)
	}
}

func TestTokenLimitRescue(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.New(s.Addr()), "tokenlimit")
	s.Close()

	var allowed int
	for i := 0; i < total; i++ {
		time.Sleep(time.Second / time.Duration(total))
		if i == total>>1 {
			assert.Nil(t, s.Restart())
		}
		if l.Allow() {
			allowed++
		}

		// make sure start monitor more than once doesn't matter
		l.startMonitor()
	}

	assert.True(t, allowed >= burst+rate)
}

func TestTokenLimitTake(t *testing.T) {
	store := redistest.CreateRedis(t)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, store, "tokenlimit")
	var allowed int
	for i := 0; i < total; i++ {
		time.Sleep(time.Second / time.Duration(total))
		if l.Allow() {
			allowed++
		}
	}

	assert.True(t, allowed >= burst+rate)
}

func TestTokenLimitTakeBurst(t *testing.T) {
	store := redistest.CreateRedis(t)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, store, "tokenlimit")
	var allowed int
	for i := 0; i < total; i++ {
		if l.Allow() {
			allowed++
		}
	}

	assert.True(t, allowed >= burst)
}
