package svc

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/middleware"
	"API_Gateway/model/ga_admin"
	"API_Gateway/model/ga_gateway_app"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
)

type ServiceContext struct {
	Config                                                                                                                                                                 config.Config
	GatewayAdminModel                                                                                                                                                      ga_admin.GatewayAdminModel
	GatewayServiceInfoModel                                                                                                                                                ga_service_info.GatewayServiceInfoModel
	GatewayServiceGrpcRuleModel                                                                                                                                            ga_service_grpc_rule.GatewayServiceGrpcRuleModel
	GatewayServiceHttpRuleModel                                                                                                                                            ga_service_http_rule.GatewayServiceHttpRuleModel
	GatewayServiceTcpRuleModel                                                                                                                                             ga_service_tcp_rule.GatewayServiceTcpRuleModel
	GatewayServiceAccessControlModel                                                                                                                                       ga_service_access_control.GatewayServiceAccessControlModel
	GatewayServiceLoadBalanceModel                                                                                                                                         ga_service_load_balance.GatewayServiceLoadBalanceModel
	GatewayAppModel                                                                                                                                                        ga_gateway_app.GatewayAppModel
	HTTPAccessMode, HTTPReverseProxy, ParamCheck, HeaderTransfer, UrlRewrite, StripUrl, HTTPWhiteList, HTTPBlackList, HTTPFlowCount, HTTPFlowLimit, AppToken, HTTPJwtToken rest.Middleware
	HTTPJwtFlowCount, HTTP_Jwt_Flow_Limit                                                                                                                                  rest.Middleware
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
		GatewayServiceLoadBalanceModel:   ga_service_load_balance.NewGatewayServiceLoadBalanceModel(conn),
		GatewayAppModel:                  ga_gateway_app.NewGatewayAppModel(conn),
		HTTPAccessMode:                   middleware.NewHTTPAccessModeMiddleware().Handle,
		HTTPReverseProxy:                 middleware.NewHTTPReverseProxyMiddleware().Handle,
		ParamCheck:                       middleware.NewParamCheckMiddleware().Handle,
		HeaderTransfer:                   middleware.NewHeaderTransferMiddleware().Handle,
		UrlRewrite:                       middleware.NewUrlRewriteMiddleware().Handle,
		StripUrl:                         middleware.NewStripUrlMiddleware().Handle,
		HTTPWhiteList:                    middleware.NewHTTPWhiteListMiddleware().Handle,
		HTTPBlackList:                    middleware.NewHTTPBlackListMiddleware().Handle,
		HTTPFlowCount:                    middleware.NewHTTPFlowCountMiddleware().Handle,
		HTTPFlowLimit:                    middleware.NewHTTPFlowLimitMiddleware().Handle,
		AppToken:                         middleware.NewAppTokenMiddleware().Handle,
		HTTPJwtToken:                     middleware.NewHTTPJwtTokenMiddleware().Handle,
		HTTPJwtFlowCount:                 middleware.NewHTTPJwtFlowCountMiddleware().Handle,
		HTTP_Jwt_Flow_Limit:              middleware.NewHTTP_Jwt_Flow_LimitMiddleware().Handle,
	}
}
