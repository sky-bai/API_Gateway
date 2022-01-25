package serviceInfo

import (
	"API_Gateway/api/internal/middleware"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_tcp_rule"
	"context"
	"fmt"
	"gopkg.in/go-playground/validator.v9"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceDetailLogic {
	return ServiceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ServiceDetail 获取服务信息
func (l *ServiceDetailLogic) ServiceDetail(req types.ServiceDetailResquest) (*types.ServiceDetail, error) {

	errMessage := ErrorString{errMessage: ""}

	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}
	serviceInfo, err := l.svcCtx.GatewayServiceInfoModel.FindOne(int64(req.ID))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务不存在")
	}
	httpRule := ga_service_http_rule.GatewayServiceHttpRule{}
	tcpRule := ga_service_tcp_rule.GatewayServiceTcpRule{}
	grpcRule := ga_service_grpc_rule.GatewayServiceGrpcRule{}
	fmt.Println(serviceInfo.LoadType)

	switch serviceInfo.LoadType {
	// http 0
	case 0:
		getHttpRule, err := l.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(req.ID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("该http服务不存在")
		}
		httpRule.Rule = getHttpRule.Rule
		httpRule.RuleType = getHttpRule.RuleType
		httpRule.NeedHttps = getHttpRule.NeedHttps
		httpRule.ServiceId = getHttpRule.ServiceId
		httpRule.UrlRewrite = getHttpRule.UrlRewrite
		httpRule.NeedStripUri = getHttpRule.NeedStripUri
		httpRule.NeedWebsocket = getHttpRule.NeedWebsocket
		httpRule.HeaderTransfor = getHttpRule.HeaderTransfor
		// tcp 1
	case 1:
		getTcpRule, err := l.svcCtx.GatewayServiceTcpRuleModel.FindOneByServiceId(req.ID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("该tcp服务不存在")
		}
		tcpRule.ServiceId = getTcpRule.ServiceId
		tcpRule.Port = getTcpRule.Port
		tcpRule.Id = getTcpRule.Id
		// grpc 2
	default:
		getGrpcRule, err := l.svcCtx.GatewayServiceGrpcRuleModel.FindOneByServiceId(req.ID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("该grpc服务不存在")
		}
		grpcRule.ServiceId = getGrpcRule.ServiceId
		grpcRule.Port = getGrpcRule.Port
		grpcRule.Id = getGrpcRule.Id
		grpcRule.HeaderTransfor = getGrpcRule.HeaderTransfor
	}

	accessControl, err := l.svcCtx.GatewayServiceAccessControlModel.FindOneByServiceId(int64(req.ID))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务访问控制不存在")
	}

	loadBalance, err := l.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(req.ID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务负载均衡不存在")
	}

	// 1.确定serviceInfo
	var newInfo types.GatewayServiceInfo
	newInfo.Id = serviceInfo.Id
	newInfo.ServiceName = serviceInfo.ServiceName
	newInfo.ServiceDesc = serviceInfo.ServiceDesc
	newInfo.LoadType = serviceInfo.LoadType

	// 2.确定httpRule
	var newHttpRule types.GatewayServiceHttpRule
	newHttpRule.Rule = httpRule.Rule
	fmt.Println("63.", newHttpRule.Rule)

	newHttpRule.RuleType = httpRule.RuleType
	newHttpRule.NeedHttps = httpRule.NeedHttps
	newHttpRule.ServiceId = httpRule.ServiceId
	newHttpRule.UrlRewrite = httpRule.UrlRewrite
	newHttpRule.NeedStripUri = httpRule.NeedStripUri
	newHttpRule.NeedWebsocket = httpRule.NeedWebsocket
	newHttpRule.HeaderTransfor = httpRule.HeaderTransfor

	// 3.确定tcpRule
	var newTcpRule types.GatewayServiceTcpRule
	newTcpRule.ServiceId = tcpRule.ServiceId
	newTcpRule.Port = tcpRule.Port

	// 4.确定grpcRule
	var newGrpcRule types.GatewayServiceGrpcRule
	newGrpcRule.Port = grpcRule.Port
	newGrpcRule.HeaderTransfor = grpcRule.HeaderTransfor

	// 5.确定accessControl
	var newAccessControl types.GatewayServiceAccessControl
	newAccessControl.ServiceId = accessControl.ServiceId
	newAccessControl.BlackList = accessControl.BlackList
	newAccessControl.WhiteList = accessControl.WhiteList
	newAccessControl.ServiceFlowLimit = accessControl.ServiceFlowLimit
	newAccessControl.ClientipFlowLimit = accessControl.ClientipFlowLimit
	newAccessControl.OpenAuth = accessControl.OpenAuth
	newAccessControl.WhiteHostName = accessControl.WhiteHostName

	// 6.确定loadBalance
	var newLoadBalance types.GatewayServiceLoadBalance
	newLoadBalance.ServiceId = loadBalance.ServiceId
	newLoadBalance.CheckMethod = loadBalance.CheckMethod
	newLoadBalance.CheckTimeout = loadBalance.CheckTimeout
	newLoadBalance.CheckInterval = loadBalance.CheckInterval
	newLoadBalance.RoundType = loadBalance.RoundType
	newLoadBalance.IpList = loadBalance.IpList
	newLoadBalance.WeightList = loadBalance.WeightList
	newLoadBalance.ForbidList = loadBalance.ForbidList
	newLoadBalance.UpstreamConnectTimeout = loadBalance.UpstreamConnectTimeout
	newLoadBalance.UpstreamIdleTimeout = loadBalance.UpstreamIdleTimeout
	newLoadBalance.UpstreamHeaderTimeout = loadBalance.UpstreamHeaderTimeout
	newLoadBalance.UpstreamMaxIdle = loadBalance.UpstreamMaxIdle

	return &types.ServiceDetail{
		Info:          newInfo,
		HTTPRule:      newHttpRule,
		TCPRule:       newTcpRule,
		GRPCRule:      newGrpcRule,
		AccessControl: newAccessControl,
		LoadBalance:   newLoadBalance,
	}, nil
}
