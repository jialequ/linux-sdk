package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdBuilderScheme(t *testing.T) {
	assert.Equal(t, EtcdScheme, new(etcdBuilder).Scheme())
}
