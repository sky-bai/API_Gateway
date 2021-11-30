package serviceInfo

import (
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/pkg/errcode"
	"API_Gateway/util"
	"context"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
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

// 获取服务列表
func (l *ServiceListLogic) ServiceList(req types.ServiceListResquest) (interface{}, error) {

	// 通过服务名称和服务描述模糊查询 获取到 服务ID
	dataLike, err := l.svcCtx.GatewayServiceInfoModel.FindDataLike(req.Info, req.PageSize, req.PageNo) // 其实查询出来的就是serviceInfo 但是我在上层进行了断言
	if err != nil {
		return nil, err
	}
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
	// 拿到每个服务ID
	for _, serviceInfo := range ServiceInfo {
		// 负载类型 0=http 1=tcp 2=grpc
		switch serviceInfo.LoadType {
		case errcode.LoadTypeHTTP:
			httpRule, err = l.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				return nil, err
			}
		case errcode.LoadTypeTCP:
			tcpRule, err = l.svcCtx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				return nil, err
			}
		default:
			grpcRule, err = l.svcCtx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				return nil, err
			}
		}

		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknown"
		clusterIP := l.svcCtx.Config.Cluster.ClusterIP
		clusterPort := l.svcCtx.Config.Cluster.ClusterPort
		clusterSSLPort := l.svcCtx.Config.Cluster.ClusterSslPort

		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			httpRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, httpRule.Rule)
		}

		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			httpRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, httpRule.Rule)
		}

		if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			httpRule.RuleType == errcode.HTTPRuleTypeDomain {
			serviceAddr = httpRule.Rule
		}
		if serviceInfo.LoadType == errcode.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, tcpRule.Port)
		}
		if serviceInfo.LoadType == errcode.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, grpcRule.Port)

		}
		loadBalance, err := l.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(serviceInfo.Id))
		if err != nil {
			return nil, err
		}
		ipList := strings.Split(loadBalance.IpList, ",")
		if err != nil {
			return nil, err
		}

		ResponseData := &types.ServiceListItemReponse{
			ID:          serviceInfo.Id,
			LoadType:    int(serviceInfo.LoadType),
			ServiceName: serviceInfo.ServiceName,
			ServiceDesc: serviceInfo.ServiceDesc,
			ServiceAddr: serviceAddr, // 从每个规则表中去判断服务的地址
			//Qps:         serviceInfo.QPS,
			//Qpd:         serviceInfo.TotalCount,
			TotalNode: len(ipList),
		}
		data = append(data, *ResponseData)
	}

	return data, nil
}
