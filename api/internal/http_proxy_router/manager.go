package http_proxy_router

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/svc"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/pkg/errcode"
	"API_Gateway/util"
	"errors"
	"github.com/tal-tech/go-zero/core/conf"
	"net/http"
	"strings"
	"sync"
)

// 我需要先获取到每个服务的信息
type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail // 一个服务名对于一个服务的详细信息
	ServiceSlice []*ServiceDetail          // 服务列表
	Locker       sync.RWMutex              // 因为要对map进行读写操作，所以需要加锁
	Once         sync.Once
	err          error
}

// 如何新的结构体去获得数据库的权限
type ServiceDetail struct {
	Info          *ga_service_info.GatewayServiceInfo                    `json:"info" description:"基本信息"`
	HTTPRule      *ga_service_http_rule.GatewayServiceHttpRule           `json:"http_rule" description:"http_rule"`
	TCPRule       *ga_service_tcp_rule.GatewayServiceTcpRule             `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *ga_service_grpc_rule.GatewayServiceGrpcRule           `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance     `json:"load_balance" description:"load_balance"`
	AccessControl *ga_service_access_control.GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}

// LoadOnce 获取到服务列表的详细信息
func (s *ServiceManager) LoadOnce() error {
	s.Once.Do(func() {
		//1.读取配置文件到结构体中
		var c config.Config
		conf.MustLoad("etc/gateway-api.yaml", &c)

		// 配置数据库
		ctx := svc.NewServiceContext(c)

		httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
		tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
		grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}

		// 获取serviceinfo表所有信息 然后通过serviceID获取详细信息
		all, err := ctx.GatewayServiceInfoModel.FindAll("", 1, 99)
		if err != nil {
			s.err = err
			return
		}
		pageList := all.(*util.PageList)
		serviceList := pageList.Data.([]ga_service_info.GatewayServiceInfo)

		s.Locker.Lock()
		defer s.Locker.Unlock()
		//serviceAddr := "unknown"

		for _, serviceInfo := range serviceList {
			// 负载类型 0=http 1=tcp 2=grpc
			switch serviceInfo.LoadType {
			case errcode.LoadTypeHTTP:
				httpRule, err = ctx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
				if err != nil {
					s.err = err
					return
				}
			case errcode.LoadTypeTCP:
				tcpRule, err = ctx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
				if err != nil {
					s.err = err
					return
				}
			default:
				grpcRule, err = ctx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(serviceInfo.Id))
				if err != nil {
					s.err = err
					return
				}
			}

			//1、http后缀接入 clusterIP+clusterPort+path
			//2、http域名接入 domain
			//3、tcp、grpc接入 clusterIP+servicePort
			//clusterIP := ctx.Config.Cluster.ClusterIP
			//clusterPort := ctx.Config.Cluster.ClusterPort
			//clusterSSLPort := ctx.Config.Cluster.ClusterSslPort

			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			//	httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			//	httpRule.NeedHttps == 1 {
			//	serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, httpRule.Rule)
			//}
			//
			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			//	httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
			//	httpRule.NeedHttps == 0 {
			//	serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, httpRule.Rule)
			//}
			//
			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
			//	httpRule.RuleType == errcode.HTTPRuleTypeDomain {
			//	serviceAddr = httpRule.Rule
			//}
			//if serviceInfo.LoadType == errcode.LoadTypeTCP {
			//	serviceAddr = fmt.Sprintf("%s:%d", clusterIP, tcpRule.Port)
			//}
			//if serviceInfo.LoadType == errcode.LoadTypeGRPC {
			//	serviceAddr = fmt.Sprintf("%s:%d", clusterIP, grpcRule.Port)
			//
			//}
			loadBalance, err := ctx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(serviceInfo.Id))
			if err != nil {
				s.err = err
				return
			}

			accessControl, err := ctx.GatewayServiceAccessControlModel.FindOneByServiceId(serviceInfo.Id)
			if err != nil {
				s.err = err
				return
			}

			detail := ServiceDetail{
				Info:          &serviceInfo,
				HTTPRule:      httpRule,
				TCPRule:       tcpRule,
				GRPCRule:      grpcRule,
				LoadBalance:   loadBalance,
				AccessControl: accessControl,
			}
			s.ServiceMap[serviceInfo.ServiceName] = &detail
			s.ServiceSlice = append(s.ServiceSlice, &detail)

		}

	})
	return s.err
}

// 前端请求 与后端http服务 想对接
func (s *ServiceManager) HTTPAccessMode(r *http.Request) (*ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//   根据请求可以得到该请求的的主机
	//   域名 host c.Request.Host
	//
	//   前缀 path c.Request.URL.Path
	host := r.Host
	host = host[0:strings.Index(host, ":")]

	path := r.URL.Path

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != errcode.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}
