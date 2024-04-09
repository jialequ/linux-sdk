package executors

import (
	"fmt"
	"testing"
	"time"

	"github.com/jialequ/linux-sdk/core/timex"
	"github.com/stretchr/testify/assert"
)

func TestLessExecutorDoOrDiscard(t *testing.T) {
	executor := NewLessExecutor(time.Minute)
	assert.True(t, executor.DoOrDiscard(func() { fmt.Print("123") }))
	assert.False(t, executor.DoOrDiscard(func() { fmt.Print("123") }))
	executor.lastTime.Set(timex.Now() - time.Minute - time.Second*30)
	assert.True(t, executor.DoOrDiscard(func() { fmt.Print("123") }))
	assert.False(t, executor.DoOrDiscard(func() { fmt.Print("123") }))
}

func BenchmarkLessExecutor(b *testing.B) {
	exec := NewLessExecutor(time.Millisecond)
	for i := 0; i < b.N; i++ {
		exec.DoOrDiscard(func() {
			fmt.Print("123")
		})
	}
}
