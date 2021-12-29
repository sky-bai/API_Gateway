package serviceHttp

import (
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

type AddHttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddHttpLogic {
	return AddHttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type ErrorString struct {
	errMessage string
}

func (e *ErrorString) Error() string {
	return e.errMessage
}

// AddHttp 添加http服务
func (l *AddHttpLogic) AddHttp(req types.AddHTTPResquest) (*types.CommonReponse, error) {

	errMessage := ErrorString{errMessage: ""}
	// 1. 校验参数
	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}

	// 2.需要根据rule 和 ruleType 判断是否有已存在的服务
	serviceID, err := l.svcCtx.GatewayServiceHttpRuleModel.FindOneByRule(req.RuleType, req.Rule)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	if serviceID > 0 {
		return nil, errors.New("该http服务已存在")
	}

	// 3.如果ip列表与权重列表不一样 就返回
	if len(strings.Split(req.IpList, "\n")) != len(strings.Split(req.WeightList, "\n")) {
		return nil, errors.New("ip列表和权重列表数量不一致")
	}

	// 4.数据库添加该服务
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
		logx.Error("添加HTTP服务失败", err)
		return nil, errors.New("添加HTTP服务失败")
	}

	return &types.CommonReponse{Msg: "添加http服务成功"}, nil
}
