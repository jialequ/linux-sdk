package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jialequ/linux-sdk/core/collection"
	"github.com/jialequ/linux-sdk/core/logx"
	"github.com/jialequ/linux-sdk/core/mathx"
	"github.com/jialequ/linux-sdk/core/stat"
	"github.com/jialequ/linux-sdk/core/stores/redis"
	"github.com/jialequ/linux-sdk/core/stores/redis/redistest"
	"github.com/jialequ/linux-sdk/core/syncx"
	"github.com/jialequ/linux-sdk/core/timex"
	"github.com/stretchr/testify/assert"
)

var errTestNotFound = errors.New("not found")

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCacheNodeDelCache(t *testing.T) {
	t.Run("del cache", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()
		store := redis.New(r.Addr(), redis.Cluster())

		cn := cacheNode{
			rds:            store,
			r:              rand.New(rand.NewSource(time.Now().UnixNano())),
			lock:           new(sync.Mutex),
			unstableExpiry: mathx.NewUnstable(expiryDeviation),
			stat:           NewStat("any"),
			errNotFound:    errTestNotFound,
		}
		assert.Nil(t, cn.Del())
		assert.Nil(t, cn.Del([]string{}...))
		assert.Nil(t, cn.Del(make([]string, 0)...))
		cn.Set("first", "one")
		assert.Nil(t, cn.Del("first"))
		cn.Set("first", "one")
		cn.Set("second", "two")
		assert.Nil(t, cn.Del("first", "second"))
	})

	t.Run("del cache with errors", func(t *testing.T) {
		old := timingWheel.Load()
		ticker := timex.NewFakeTicker()
		tw, err := collection.NewTimingWheelWithTicker(
			time.Millisecond, timingWheelSlots, func(key, value any) {
				clean(key, value)
			}, ticker)
		timingWheel.Store(tw)
		assert.NoError(t, err)
		t.Cleanup(func() {
			timingWheel.Store(old)
		})

		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()
		r.SetError(literal_5379)

		node := NewNode(redis.New(r.Addr(), redis.Cluster()), syncx.NewSingleFlight(),
			NewStat("any"), errTestNotFound)
		assert.NoError(t, node.Del("foo", "bar"))
		ticker.Tick()
		runtime.Gosched()
	})
}

func TestCacheNodeDelCacheWithErrors(t *testing.T) {
	store := redistest.CreateRedis(t)
	store.Type = redis.ClusterType

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}
	assert.Nil(t, cn.Del("third", "fourth"))
}

func TestCacheNodeInvalidCache(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	cn := cacheNode{
		rds:            redis.New(s.Addr()),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}
	s.Set("any", "value")
	var str string
	assert.NotNil(t, cn.Get("any", &str))
	assert.Equal(t, "", str)
	_, err = s.Get("any")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)
}

func TestCacheNodeSetWithExpire(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	assert.NotNil(t, cn.SetWithExpire("key", make(chan int), time.Second))
}

func TestCacheNodeTake(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errTestNotFound,
		WithExpiry(time.Second), WithNotFoundExpiry(time.Second))
	var str string
	err := cn.Take(&str, "any", func(v any) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.Get("any", &str))
	val, err := store.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNodeTakeBadRedis(t *testing.T) {
	r, err := miniredis.Run()
	assert.NoError(t, err)
	defer r.Close()
	r.SetError(literal_5379)

	cn := NewNode(redis.New(r.Addr()), syncx.NewSingleFlight(), NewStat("any"),
		errTestNotFound, WithExpiry(time.Second), WithNotFoundExpiry(time.Second))
	var str string
	assert.Error(t, cn.Take(&str, "any", func(v any) error {
		*v.(*string) = "value"
		return nil
	}))
}

func TestCacheNodeTakeNotFound(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		store := redistest.CreateRedis(t)

		cn := cacheNode{
			rds:            store,
			r:              rand.New(rand.NewSource(time.Now().UnixNano())),
			barrier:        syncx.NewSingleFlight(),
			lock:           new(sync.Mutex),
			unstableExpiry: mathx.NewUnstable(expiryDeviation),
			stat:           NewStat("any"),
			errNotFound:    errTestNotFound,
		}
		var str string
		err := cn.Take(&str, "any", func(v any) error {
			return errTestNotFound
		})
		assert.True(t, cn.IsNotFound(err))
		assert.True(t, cn.IsNotFound(cn.Get("any", &str)))
		val, err := store.Get("any")
		assert.Nil(t, err)
		assert.Equal(t, `*`, val)

		store.Set("any", "*")
		err = cn.Take(&str, "any", func(v any) error {
			return nil
		})
		assert.True(t, cn.IsNotFound(err))
		assert.True(t, cn.IsNotFound(cn.Get("any", &str)))

		store.Del("any")
		errDummy := errors.New("dummy")
		err = cn.Take(&str, "any", func(v any) error {
			return errDummy
		})
		assert.Equal(t, errDummy, err)
	})

	t.Run("not found with redis error", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()
		store, err := redis.NewRedis(redis.RedisConf{
			Host: r.Addr(),
			Type: redis.NodeType,
		})
		assert.NoError(t, err)

		cn := cacheNode{
			rds:            store,
			r:              rand.New(rand.NewSource(time.Now().UnixNano())),
			barrier:        syncx.NewSingleFlight(),
			lock:           new(sync.Mutex),
			unstableExpiry: mathx.NewUnstable(expiryDeviation),
			stat:           NewStat("any"),
			errNotFound:    errTestNotFound,
		}
		var str string
		err = cn.Take(&str, "any", func(v any) error {
			r.SetError(literal_5379)
			return errTestNotFound
		})
		assert.True(t, cn.IsNotFound(err))
	})
}

func TestCacheNodeTakeCtxWithRedisError(t *testing.T) {
	t.Run("not found with redis error", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()
		store, err := redis.NewRedis(redis.RedisConf{
			Host: r.Addr(),
			Type: redis.NodeType,
		})
		assert.NoError(t, err)

		cn := cacheNode{
			rds:            store,
			r:              rand.New(rand.NewSource(time.Now().UnixNano())),
			barrier:        syncx.NewSingleFlight(),
			lock:           new(sync.Mutex),
			unstableExpiry: mathx.NewUnstable(expiryDeviation),
			stat:           NewStat("any"),
			errNotFound:    errTestNotFound,
		}
		var str string
		err = cn.Take(&str, "any", func(v any) error {
			str = "foo"
			r.SetError(literal_5379)
			return nil
		})
		assert.NoError(t, err)
	})
}

func TestCacheNodeTakeNotFoundButChangedByOthers(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}

	var str string
	err := cn.Take(&str, "any", func(v any) error {
		store.Set("any", "foo")
		return errTestNotFound
	})
	assert.True(t, cn.IsNotFound(err))

	val, err := store.Get("any")
	if assert.NoError(t, err) {
		assert.Equal(t, "foo", val)
	}
	assert.True(t, cn.IsNotFound(cn.Get("any", &str)))
}

func TestCacheNodeTakeWithExpire(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	var str string
	err := cn.TakeWithExpire(&str, "any", func(v any, expire time.Duration) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.Get("any", &str))
	val, err := store.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNodeString(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	assert.Equal(t, store.Addr, cn.String())
}

func TestCacheValueWithBigInt(t *testing.T) {
	store := redistest.CreateRedis(t)

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}

	const (
		key         = "key"
		value int64 = 323427211229009810
	)

	assert.Nil(t, cn.Set(key, value))
	var val any
	assert.Nil(t, cn.Get(key, &val))
	assert.Equal(t, strconv.FormatInt(value, 10), fmt.Sprintf("%v", val))
}

const literal_5379 = "mock error"
