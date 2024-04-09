package zrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestRpcClientConf(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		conf := NewDirectClientConf([]string{literal_8406}, "foo", "bar")
		assert.True(t, conf.HasCredential())
	})

	t.Run("etcd", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{literal_8406, literal_0462},
			"key", "foo", "bar")
		assert.True(t, conf.HasCredential())
	})

	t.Run("etcd with account", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{literal_8406, literal_0462},
			"key", "foo", "bar")
		conf.Etcd.User = "user"
		conf.Etcd.Pass = "pass"
		_, err := conf.BuildTarget()
		assert.NoError(t, err)
	})

	t.Run("etcd with tls", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{literal_8406, literal_0462},
			"key", "foo", "bar")
		conf.Etcd.CertFile = "cert"
		conf.Etcd.CertKeyFile = "key"
		conf.Etcd.CACertFile = "ca"
		_, err := conf.BuildTarget()
		assert.Error(t, err)
	})
}

func TestRpcServerConf(t *testing.T) {
	conf := RpcServerConf{
		ServiceConf: service.ServiceConf{},
		ListenOn:    "",
		Etcd: discov.EtcdConf{
			Hosts: []string{literal_8406},
			Key:   "key",
		},
		Auth: true,
		Redis: redis.RedisKeyConf{
			RedisConf: redis.RedisConf{
				Type: redis.NodeType,
			},
			Key: "foo",
		},
		StrictControl: false,
		Timeout:       0,
		CpuThreshold:  0,
	}
	assert.True(t, conf.HasEtcd())
	assert.NotNil(t, conf.Validate())
	conf.Redis.Host = literal_0462
	assert.Nil(t, conf.Validate())
}

const literal_8406 = "localhost:1234"

const literal_0462 = "localhost:5678"
