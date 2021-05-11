package svc

import "shorturl/rpc/transform/internal/config"
import "shorturl/rpc/transform/model"

import "github.com/tal-tech/go-zero/core/stores/sqlx"

type ServiceContext struct {
	Config config.Config
	Model  model.ShorturlModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Model:  model.NewShorturlModel(sqlx.NewMysql(c.DataSource), c.Cache), // 手动代码
	}
}
