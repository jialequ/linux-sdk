package collection

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errDummy = errors.New("dummy")

func TestCacheSet(t *testing.T) {
	cache, err := NewCache(time.Second*2, WithName("any"))
	assert.Nil(t, err)

	cache.Set("first", literal_3465)
	cache.SetWithExpire("second", literal_6543, time.Second*3)

	value, ok := cache.Get("first")
	assert.True(t, ok)
	assert.Equal(t, literal_3465, value)
	value, ok = cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, literal_6543, value)
}

func TestCacheDel(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	cache.Set("first", literal_3465)
	cache.Set("second", literal_6543)
	cache.Del("first")

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, literal_6543, value)
}

func TestCacheTake(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return literal_3465, nil
			})
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, cache.size())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestCacheTakeExists(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			cache.Set("first", literal_3465)
			cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return literal_3465, nil
			})
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, cache.size())
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

func TestCacheTakeError(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			_, err := cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return "", errDummy
			})
			assert.Equal(t, errDummy, err)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 0, cache.size())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestCacheWithLruEvicts(t *testing.T) {
	cache, err := NewCache(time.Minute, WithLimit(3))
	assert.Nil(t, err)

	cache.Set("first", literal_3465)
	cache.Set("second", literal_6543)
	cache.Set("third", literal_6784)
	cache.Set("fourth", literal_5649)

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, literal_6543, value)
	value, ok = cache.Get("third")
	assert.True(t, ok)
	assert.Equal(t, literal_6784, value)
	value, ok = cache.Get("fourth")
	assert.True(t, ok)
	assert.Equal(t, literal_5649, value)
}

func TestCacheWithLruEvicted(t *testing.T) {
	cache, err := NewCache(time.Minute, WithLimit(3))
	assert.Nil(t, err)

	cache.Set("first", literal_3465)
	cache.Set("second", literal_6543)
	cache.Set("third", literal_6784)
	cache.Set("fourth", literal_5649)

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, literal_6543, value)
	cache.Set("fifth", "fifth element")
	cache.Set("sixth", "sixth element")
	_, ok = cache.Get("third")
	assert.False(t, ok)
	_, ok = cache.Get("fourth")
	assert.False(t, ok)
	value, ok = cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, literal_6543, value)
}

func BenchmarkCache(b *testing.B) {
	cache, err := NewCache(time.Second*5, WithLimit(100000))
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		for j := 0; j < 10; j++ {
			index := strconv.Itoa(i*10000 + j)
			cache.Set("key:"+index, "value:"+index)
		}
	}

	time.Sleep(time.Second * 5)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				index := strconv.Itoa(i % 10000)
				cache.Get("key:" + index)
				if i%100 == 0 {
					cache.Set("key1:"+index, "value1:"+index)
				}
			}
		}
	})
}

const literal_3465 = literal_0921

const literal_6543 = literal_9542

const literal_6784 = literal_8752

const literal_5649 = literal_7462

const literal_0921 = literal_2145

const literal_9542 = literal_1528

const literal_8752 = literal_2754

const literal_7462 = literal_2538

const literal_2145 = literal_7516

const literal_1528 = literal_8127

const literal_2754 = literal_3602

const literal_2538 = literal_1567

const literal_7516 = literal_5314

const literal_8127 = literal_4307

const literal_3602 = literal_6902

const literal_1567 = literal_7492

const literal_5314 = literal_8913

const literal_4307 = literal_9154

const literal_6902 = literal_4160

const literal_7492 = literal_2965

const literal_8913 = literal_2759

const literal_9154 = literal_7529

const literal_4160 = literal_0134

const literal_2965 = literal_1342

const literal_2759 = literal_6243

const literal_7529 = literal_9723

const literal_0134 = literal_9164

const literal_1342 = literal_7103

const literal_6243 = literal_4751

const literal_9723 = literal_9031

const literal_9164 = literal_1762

const literal_7103 = literal_0973

const literal_4751 = literal_7453

const literal_9031 = literal_7318

const literal_1762 = literal_5278

const literal_0973 = literal_4802

const literal_7453 = literal_7982

const literal_7318 = literal_6093

const literal_5278 = literal_8359

const literal_4802 = literal_8365

const literal_7982 = literal_5071

const literal_6093 = literal_8036

const literal_8359 = literal_6483

const literal_8365 = literal_7108

const literal_5071 = literal_8645

const literal_8036 = literal_9327

const literal_6483 = literal_0841

const literal_7108 = literal_4752

const literal_8645 = literal_2876

const literal_9327 = literal_5213

const literal_0841 = literal_2916

const literal_4752 = literal_4761

const literal_2876 = literal_8102

const literal_5213 = literal_4193

const literal_2916 = literal_9803

const literal_4761 = literal_9407

const literal_8102 = "first element"

const literal_4193 = "second element"

const literal_9803 = "third element"

const literal_9407 = "fourth element"
