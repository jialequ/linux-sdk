package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sleepInterval = time.Millisecond * 10

func TestElapsedTimerDuration(t *testing.T) {
	timer := NewElapsedTimer()
	time.Sleep(sleepInterval)
	assert.True(t, timer.Duration() >= sleepInterval)
}

func TestElapsedTimerElapsed(t *testing.T) {
	timer := NewElapsedTimer()
	time.Sleep(sleepInterval)
	duration, err := time.ParseDuration(timer.Elapsed())
	assert.Nil(t, err)
	assert.True(t, duration >= sleepInterval)
}

func TestElapsedTimerElapsedMs(t *testing.T) {
	timer := NewElapsedTimer()
	time.Sleep(sleepInterval)
	duration, err := time.ParseDuration(timer.ElapsedMs())
	assert.Nil(t, err)
	assert.True(t, duration >= sleepInterval)
}

func TestCurrent(t *testing.T) {
	currentMillis := CurrentMillis()
	currentMicros := CurrentMicros()
	assert.True(t, currentMillis > 0)
	assert.True(t, currentMicros > 0)
	assert.True(t, currentMillis*1000 <= currentMicros)
}
