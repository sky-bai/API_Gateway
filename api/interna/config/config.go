package config

import (
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	CacheRedis cache.CacheConf
	Cluster    struct {
		ClusterIP      string
		ClusterPort    string
		ClusterSslPort string
	}

	HTTPProxy struct {
		Name string
		Host string
		Port int
	}

	HTTPSProxy struct {
		Name string
		Host string
		Port int
	}
}
