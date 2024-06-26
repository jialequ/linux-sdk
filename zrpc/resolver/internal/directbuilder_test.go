package internal

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/jialequ/linux-sdk/core/lang"
	"github.com/jialequ/linux-sdk/core/mathx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestDirectBuilderBuild(t *testing.T) {
	tests := []int{
		0,
		1,
		2,
		subsetSize / 2,
		subsetSize,
		subsetSize * 2,
	}

	for _, test := range tests {
		test := test
		t.Run(strconv.Itoa(test), func(t *testing.T) {
			var servers []string
			for i := 0; i < test; i++ {
				servers = append(servers, fmt.Sprintf("localhost:%d", i))
			}
			var b directBuilder
			cc := new(mockedClientConn)
			target := fmt.Sprintf("%s:///%s", DirectScheme, strings.Join(servers, ","))
			uri, err := url.Parse(target)
			assert.Nil(t, err)
			cc.err = errors.New("foo")
			_, err = b.Build(resolver.Target{
				URL: *uri,
			}, cc, resolver.BuildOptions{})
			assert.NotNil(t, err)
			cc.err = nil
			_, err = b.Build(resolver.Target{
				URL: *uri,
			}, cc, resolver.BuildOptions{})
			assert.NoError(t, err)

			size := mathx.MinInt(test, subsetSize)
			assert.Equal(t, size, len(cc.state.Addresses))
			m := make(map[string]lang.PlaceholderType)
			for _, each := range cc.state.Addresses {
				m[each.Addr] = lang.Placeholder
			}
			assert.Equal(t, size, len(m))
		})
	}
}

func TestDirectBuilderScheme(t *testing.T) {
	var b directBuilder
	assert.Equal(t, DirectScheme, b.Scheme())
}
