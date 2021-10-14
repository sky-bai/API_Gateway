package svc

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/model/ga_admin"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_tcp_rule"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                           config.Config
	GatewayAdminModel                ga_admin.GatewayAdminModel
	GatewayServiceInfoModel          ga_service_info.GatewayServiceInfoModel
	GatewayServiceGrpcRuleModel      ga_service_grpc_rule.GatewayServiceGrpcRuleModel
	GatewayServiceHttpRuleModel      ga_service_http_rule.GatewayServiceHttpRuleModel
	GatewayServiceTcpRuleModel       ga_service_tcp_rule.GatewayServiceTcpRuleModel
	GatewayServiceAccessControlModel ga_service_access_control.GatewayServiceAccessControlModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:                           c,
		GatewayAdminModel:                ga_admin.NewGatewayAdminModel(conn),
		GatewayServiceInfoModel:          ga_service_info.NewGatewayServiceInfoModel(conn),
		GatewayServiceGrpcRuleModel:      ga_service_grpc_rule.NewGatewayServiceGrpcRuleModel(conn),
		GatewayServiceHttpRuleModel:      ga_service_http_rule.NewGatewayServiceHttpRuleModel(conn),
		GatewayServiceTcpRuleModel:       ga_service_tcp_rule.NewGatewayServiceTcpRuleModel(conn),
		GatewayServiceAccessControlModel: ga_service_access_control.NewGatewayServiceAccessControlModel(conn),
	}
}
