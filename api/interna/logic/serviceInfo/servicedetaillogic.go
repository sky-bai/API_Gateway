package serviceInfo

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"errors"
	"fmt"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

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

type HttpResponse struct {
	GatewayServiceInfo          *ga_service_info.GatewayServiceInfo
	GatewayServiceHttpRule      *ga_service_http_rule.GatewayServiceHttpRule
	GatewayServiceLoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance
	GatewayServiceAccessControl *ga_service_access_control.GatewayServiceAccessControl
}

// 获取该服务信息
func (l *ServiceDetailLogic) ServiceDetail(req types.ServiceDetailResquest) (interface{}, error) {
	serviceInfo, err := l.svcCtx.GatewayServiceInfoModel.FindOne(req.ID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务不存在")
	}

	httpRule, err := l.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(req.ID))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务不存在")
	}

	accessControl, err := l.svcCtx.GatewayServiceAccessControlModel.FindOneByServiceId(req.ID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务不存在")
	}

	loadBalance, err := l.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(req.ID))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该服务不存在")
	}

	HttpResponse := &HttpResponse{}
	HttpResponse.GatewayServiceInfo = serviceInfo
	HttpResponse.GatewayServiceHttpRule = httpRule
	HttpResponse.GatewayServiceAccessControl = accessControl
	HttpResponse.GatewayServiceLoadBalance = loadBalance

	return HttpResponse, nil
}
