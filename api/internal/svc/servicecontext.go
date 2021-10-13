package svc

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/model/ga_admin"
	"API_Gateway/model/ga_service_info"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	GatewayAdminModel       ga_admin.GatewayAdminModel
	GatewayServiceInfoModel ga_service_info.GatewayServiceInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:                  c,
		GatewayAdminModel:       ga_admin.NewGatewayAdminModel(conn),
		GatewayServiceInfoModel: ga_service_info.NewGatewayServiceInfoModel(conn),
	}
}
