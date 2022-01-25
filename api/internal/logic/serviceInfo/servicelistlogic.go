package serviceInfo

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/middleware"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/pkg/errcode"
	"API_Gateway/util"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ServiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceListLogic {
	return ServiceListLogic{
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

// ServiceList 获取服务列表
func (l *ServiceListLogic) ServiceList(req types.ServiceListResquest) (*types.DataList, error) {

	errMessage := ErrorString{errMessage: ""}

	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}
	// 1.通过服务名称和服务描述模糊查询 获取到 服务ID
	dataLike, err := l.svcCtx.GatewayServiceInfoModel.FindDataLike(req.Info, req.PageSize, req.PageNo) // 其实查询出来的就是serviceInfo 但是我在上层进行了断言
	if err != nil {
		return nil, errors.New("获取服务信息失败")
	}

	// 2. 通过服务ID 获取到服务的规则 总数
	pa := dataLike.(*util.PageList)
	ServiceInfo := pa.Data.([]ga_service_info.GatewayServiceInfo)
	//for _, info := range ServiceInfo {
	//	idList:=[]int64{}
	//	idlist := append(idList, info.Id)
	//}

	httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
	tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
	grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}

	var data []types.ServiceListItemReponse //nolint:prealloc
	// 3.获取每一个服务的具体信息
	for _, serviceInfo := range ServiceInfo {
		// 负载类型 0=http 1=tcp 2=grpc 4.先确定该服务是什么类型
		switch serviceInfo.LoadType {
		case errcode.LoadTypeHTTP:
			httpRule, err = l.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				return nil, errors.New("获取http服务信息失败")
			}
		case errcode.LoadTypeTCP:
			tcpRule, err = l.svcCtx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err == sqlc.ErrNotFound {
				tcpRule = &ga_service_tcp_rule.GatewayServiceTcpRule{}
			} else if err != nil {
				return nil, errors.New("获取tcp服务信息失败")
			}
		default:
			grpcRule, err = l.svcCtx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				return nil, errors.New("获取grpc服务信息失败")
			}
		}

		// 1、http后缀接入 clusterIP+clusterPort+path
		// 2、http域名接入 domain
		// 3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknown"
		clusterIP := l.svcCtx.Config.Cluster.ClusterIP
		clusterPort := l.svcCtx.Config.Cluster.ClusterPort
		clusterSSLPort := l.svcCtx.Config.Cluster.ClusterSslPort

		// 5.如果是https服务并且是前缀接入
		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			httpRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, httpRule.Rule)
		}
		// 6.如果是http服务并且是前缀接入
		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			httpRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, httpRule.Rule)
		}
		// 7.如果是http服务并且是域名接入
		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypeDomain {
			serviceAddr = httpRule.Rule
		}
		// 8.如果是tcp服务
		if serviceInfo.LoadType == errcode.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, tcpRule.Port)
		}
		// 9.如果是grpc服务
		if serviceInfo.LoadType == errcode.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, grpcRule.Port)
		}
		// 10.获取该服务发负载配置
		loadBalance, err := l.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(serviceInfo.Id))
		if err == sqlx.ErrNotFound {
			fmt.Println("!!!!!!!", err)
			err = nil
			loadBalance = &ga_service_load_balance.GatewayServiceLoadBalance{}
		}
		if err != nil {
			return nil, errors.New("获取loadBalance服务信息失败")
		}

		ipList := strings.Split(loadBalance.IpList, ",")
		if err != nil {
			logx.Error("获取服务列表的时候，修改ipList失败")
			return nil, errors.New("获取该服务信息失败")
		}

		// 11.获取该服务的qps
		serviceCounter, err := global.FlowCounterHandler.GetCounter(global.FlowServicePrefix + serviceInfo.ServiceName)
		if err != nil {
			return nil, errors.New("获取该服务的qps失败")
		}

		ResponseData := &types.ServiceListItemReponse{
			ID:          serviceInfo.Id,
			LoadType:    int(serviceInfo.LoadType),
			ServiceName: serviceInfo.ServiceName,
			ServiceDesc: serviceInfo.ServiceDesc,
			ServiceAddr: serviceAddr, // 从每个规则表中去判断服务的地址
			Qps:         serviceCounter.QPS,
			Qpd:         serviceCounter.TotalCount,
			TotalNode:   len(ipList),
		}
		data = append(data, *ResponseData)
	}

	var serviceListInfo util.PageList
	serviceListInfo.Data = data
	serviceListInfo.Count = pa.Count
	serviceListInfo.TotalPage = pa.TotalPage
	serviceListInfo.Page = pa.Page
	serviceListInfo.Limit = pa.Limit
	return &types.DataList{Page: serviceListInfo.Page, Limit: serviceListInfo.Limit, Count: serviceListInfo.Count, Total: serviceListInfo.TotalPage, List: serviceListInfo.Data}, nil
}
