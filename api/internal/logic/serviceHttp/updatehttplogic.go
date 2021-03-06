package serviceHttp

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/middleware"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"strings"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateHttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateHttpLogic {
	return UpdateHttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateHttp 更新http服务
func (l *UpdateHttpLogic) UpdateHttp(req types.UpdateHTTPResquest) (*types.Reponse, error) {
	errMessage := ErrorString{errMessage: ""}

	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}
	// 需要根据id 判断是否有已存在的服务
	serviceInfo, err := l.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(req.ID))
	if err != nil {
		return nil, err
	}
	if serviceInfo.Id < 1 {
		return nil, errors.New("该http服务未存在")
	}

	// 如果ip列表与权重列表不一样 就返回
	if len(strings.Split(req.IpList, "/n")) != len(strings.Split(req.WeightList, "/n")) {
		return nil, errors.New("ip列表和权重列表数量不一致")
	}

	// 数据库更新该服务
	service := ga_service_info.GatewayServiceInfo{}
	service.Id = req.ID
	//service.ServiceDesc = req.ServiceDesc
	service.ServiceName = req.ServiceName
	service.ServiceDesc = req.ServiceDesc

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

	err = l.svcCtx.GatewayServiceInfoModel.UpdateDate(service, httpRule, accessControl, loadBalance)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("该http服务更新失败")
	}

	s1 := global.ServiceDetail{
		Info:          service,
		HTTPRule:      httpRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	tem := *global.SerInfo
	//tem = append(tem, s1)

	var nilService []global.ServiceDetail

	for _, value := range tem {
		if value.HTTPRule.Rule == req.Rule {
			nilService = append(nilService, s1)
		} else {
			nilService = append(nilService, value)
		}
	}
	global.SerInfo = &nilService

	fmt.Println(global.SerInfo)

	return &types.Reponse{Msg: "该http服务更新完成"}, nil
}
