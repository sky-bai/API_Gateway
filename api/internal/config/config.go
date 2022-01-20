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

	//LogConf struct {
	//	ServiceName         string `json:",optional"`
	//	Mode                string `json:",default=console,options=console|file|volume"`
	//	Path                string `json:",default=logs"`
	//	Level               string `json:",default=info,options=info|error|severe"`
	//	Compress            bool   `json:",optional"`
	//	KeepDays            int    `json:",optional"`
	//	StackCooldownMillis int    `json:",default=100"`
	//}
}
