package serviceInfo

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"errors"
	"fmt"
	"strings"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceAddHttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceAddHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceAddHttpLogic {
	return ServiceAddHttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 添加http服务
func (l *ServiceAddHttpLogic) ServiceAddHttp(req types.AddHTTPResquest) (*types.CommonReponse, error) {

	fmt.Println(req)

	// 需要根据rule 和 ruleType 判断是否有已存在的服务
	rule, err := l.svcCtx.GatewayServiceHttpRuleModel.FindOneByRule(req.RuleType, req.Rule)
	if err != nil {
		return nil, err
	}
	if rule.Id > 0 {
		return nil, errors.New("该http服务已存在")
	}

	// 如果ip列表与权重列表不一样 就返回
	if len(strings.Split(req.IpList, "/n")) != len(strings.Split(req.WeightList, "/n")) {
		return nil, errors.New("ip列表和权重列表数量不一致")
	}

	// 数据库添加该服务
	serviceInfo := ga_service_info.GatewayServiceInfo{}
	serviceInfo.ServiceDesc = req.ServiceDesc
	serviceInfo.ServiceName = req.ServiceName

	httpRule := ga_service_http_rule.GatewayServiceHttpRule{}
	httpRule.RuleType = int64(req.RuleType)
	httpRule.Rule = req.Rule
	httpRule.HeaderTransfor = req.HeaderTransfor
	httpRule.NeedWebsocket = int64(req.NeedWebsocket)
	httpRule.NeedStripUri = int64(req.NeedStripUri)
	httpRule.UrlRewrite = req.UrlRewrite
	httpRule.NeedHttps = int64(req.NeedHttps)

	accessControl := ga_service_access_control.GatewayServiceAccessControl{}
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.OpenAuth = int64(req.OpenAuth)
	accessControl.ClientipFlowLimit = int64(req.ClientipFlowLimit)
	accessControl.ServiceFlowLimit = int64(req.ServiceFlowLimit)

	loadBalance := ga_service_load_balance.GatewayServiceLoadBalance{}
	loadBalance.RoundType = int64(req.RoundType)
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.UpstreamConnectTimeout = int64(req.UpstreamConnectTimeout)
	loadBalance.UpstreamHeaderTimeout = int64(req.UpstreamHeaderTimeout)
	loadBalance.UpstreamIdleTimeout = int64(req.UpstreamIdleTimeout)
	loadBalance.UpstreamMaxIdle = int64(req.UpstreamMaxIdle)

	err = l.svcCtx.GatewayServiceInfoModel.InsertData(serviceInfo, httpRule, accessControl, loadBalance)

	if err != nil {
		return nil, err
	}

	return &types.CommonReponse{Msg: "添加http服务成功"}, nil
}
