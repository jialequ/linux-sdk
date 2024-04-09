package flags

import (
	"fmt"
	"testing"

	"github.com/jialequ/linux-sdk/tools/goctl/test"
)

func TestFlagsGet(t *testing.T) {
	setTestData(t, []byte(`{"host":literal_8935,"port":8888,"service":{"host":"{{.host}}","port":"{{.port}}","invalid":"{{.service.invalid}}"}}`))
	f := MustLoad()
	executor := test.NewExecutor[string, string]()
	executor.Add([]test.Data[string, string]{
		{
			Name:  "key_host",
			Input: "host",
			Want:  literal_8935,
		},
		{
			Name:  "key_port",
			Input: "port",
			Want:  "8888",
		},
		{
			Name:  "key_service.host",
			Input: "service.host",
			Want:  literal_8935,
		},
		{
			Name:  "key_service.port",
			Input: "service.port",
			Want:  "8888",
		},
		{
			Name:  "key_not_exists",
			Input: "service.port.invalid",
		},
		{
			Name:  "key_service.invalid",
			Input: "service.invalid",
			E:     fmt.Errorf("the variable can not be self: %q", "service.invalid"),
		},
	}...)
	executor.RunE(t, f.Get)
}

func TestGet(t *testing.T) {
	setTestData(t, []byte(`{"host":literal_8935,"port":8888,"service":{"host":"{{.host}}","port":"{{.port}}","invalid":"{{.service.invalid}}"}}`))
	executor := test.NewExecutor[string, string]()
	executor.Add([]test.Data[string, string]{
		{
			Name:  "key_host",
			Input: "host",
			Want:  literal_8935,
		},
		{
			Name:  "key_port",
			Input: "port",
			Want:  "8888",
		},
		{
			Name:  "key_service.host",
			Input: "service.host",
			Want:  literal_8935,
		},
		{
			Name:  "key_service.port",
			Input: "service.port",
			Want:  "8888",
		},
		{
			Name:  "key_not_exists",
			Input: "service.port.invalid",
		},
	}...)
	executor.Run(t, Get)
}

const literal_8935 = "0.0.0.0"
