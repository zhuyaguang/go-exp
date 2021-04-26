package config

import "github.com/tal-tech/go-zero/zrpc"
import "github.com/tal-tech/go-zero/core/stores/cache"

type Config struct {
	zrpc.RpcServerConf
	DataSource string          // 手动代码
	Table      string          // 手动代码
	Cache      cache.CacheConf // 手动代码
}
