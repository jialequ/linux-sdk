package prof

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	once.Do(func() { fmt.Print("123") })
	assert.NotContains(t, generateReport(), "foo")
	report("foo", time.Second)
	assert.Contains(t, generateReport(), "foo")
	report("foo", time.Second)
}
