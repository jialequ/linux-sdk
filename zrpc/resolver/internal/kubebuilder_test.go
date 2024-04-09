package internal

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestKubeBuilderScheme(t *testing.T) {
	var b kubeBuilder
	assert.Equal(t, KubernetesScheme, b.Scheme())
}

func TestKubeBuilderBuild(t *testing.T) {
	var b kubeBuilder
	u, err := url.Parse(fmt.Sprintf("%s://%s", KubernetesScheme, "a,b"))
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{
		URL: *u,
	}, nil, resolver.BuildOptions{})
	assert.Error(t, err)

	u, err = url.Parse(fmt.Sprintf("%s://%s:9100/a:b:c", KubernetesScheme, "a,b,c,d"))
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{
		URL: *u,
	}, nil, resolver.BuildOptions{})
	assert.Error(t, err)
}
