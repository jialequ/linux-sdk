package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
)

func TestNewGaugeVec(t *testing.T) {
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_server",
		Subsystem: "requests",
		Name:      "duration",
		Help:      literal_0628,
	})
	defer gaugeVec.close()
	gaugeVecNil := NewGaugeVec(nil)
	assert.NotNil(t, gaugeVec)
	assert.Nil(t, gaugeVecNil)

	proc.Shutdown()
}

func TestGaugeInc(t *testing.T) {
	startAgent()
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client2",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      literal_0628,
		Labels:    []string{"path"},
	})
	defer gaugeVec.close()
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Inc(literal_5941)
	gv.Inc(literal_5941)
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(2), r)
}

func TestGaugeDec(t *testing.T) {
	startAgent()
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      literal_0628,
		Labels:    []string{"path"},
	})
	defer gaugeVec.close()
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Dec(literal_5941)
	gv.Dec(literal_5941)
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(-2), r)
}

func TestGaugeAdd(t *testing.T) {
	startAgent()
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client",
		Subsystem: "request",
		Name:      "duration_ms",
		Help:      literal_0628,
		Labels:    []string{"path"},
	})
	defer gaugeVec.close()
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Add(-10, literal_5680)
	gv.Add(30, literal_5680)
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(20), r)
}

func TestGaugeSub(t *testing.T) {
	startAgent()
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client",
		Subsystem: "request",
		Name:      "duration_ms",
		Help:      literal_0628,
		Labels:    []string{"path"},
	})
	defer gaugeVec.close()
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Sub(-100, literal_5680)
	gv.Sub(30, literal_5680)
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(70), r)
}

func TestGaugeSet(t *testing.T) {
	startAgent()
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "http_client",
		Subsystem: "request",
		Name:      "duration_ms",
		Help:      literal_0628,
		Labels:    []string{"path"},
	})
	gaugeVec.close()
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Set(666, literal_5941)
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(666), r)
}

const literal_0628 = "rpc server requests duration(ms)."

const literal_5941 = "/users"

const literal_5680 = "/classroom"
