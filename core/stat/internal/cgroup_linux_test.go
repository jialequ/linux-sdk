package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningInUserNS(t *testing.T) {
	// should be false in docker
	assert.False(t, runningInUserNS())
}

func TestCgroups(t *testing.T) {
	// test cgroup legacy(v1) & hybrid
	if !isCgroup2UnifiedMode() {
		cg, err := currentCgroupV1()
		assert.NoError(t, err)
		_, err = cg.effectiveCpus()
		assert.NoError(t, err)
		_, err = cg.cpuQuota()
		assert.NoError(t, err)
		_, err = cg.cpuUsage()
		assert.NoError(t, err)
	}

	// test cgroup v2
	if isCgroup2UnifiedMode() {
		cg, err := currentCgroupV2()
		assert.NoError(t, err)
		_, err = cg.effectiveCpus()
		assert.NoError(t, err)
		_, err = cg.cpuUsage()
		assert.NoError(t, err)
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
		err   error
	}{
		{"0", 0, nil},
		{"123", 123, nil},
		{"-1", 0, nil},
		{"-18446744073709551616", 0, nil},
		{"foo", 0, fmt.Errorf("cgroup: bad int format: foo")},
	}

	for _, tt := range tests {
		got, err := parseUint(tt.input)
		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.want, got)
	}
}
