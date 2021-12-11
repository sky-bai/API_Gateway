package serviceHttp

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"errors"
	"fmt"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpDetailLogic {
	return HttpDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type HttpResponse struct {
	ServiceInfo   *ga_service_info.GatewayServiceInfo
	HttpRule      *ga_service_http_rule.GatewayServiceHttpRule
	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance
	AccessControl *ga_service_access_control.GatewayServiceAccessControl
}

// HttpDetail 获取http服务信息
func (l *HttpDetailLogic) HttpDetail(req types.HttpeDetailResquest) (interface{}, error) {
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
	HttpResponse.ServiceInfo = serviceInfo
	HttpResponse.HttpRule = httpRule
	HttpResponse.AccessControl = accessControl
	HttpResponse.LoadBalance = loadBalance

	return HttpResponse, nil

}
