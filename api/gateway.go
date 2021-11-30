package main

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/config/cert_file"
	"API_Gateway/api/internal/handler"
	_ "API_Gateway/api/internal/http_proxy_router"
	"API_Gateway/api/internal/svc"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
	flag.Parse()
	*configFile = "etc/gateway-api.yaml"
	// 1.读取配置文件到结构体中
	var c config.Config
	conf.MustLoad(*configFile, &c)
	//fmt.Println(*configFile)

	// 配置数据库
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)

	// 2.启动http代理服务
	go func() {

		c1 := c
		c1.RestConf.Name = c.HTTPProxy.Name
		c1.RestConf.Host = c.HTTPProxy.Host
		c1.RestConf.Port = c.HTTPProxy.Port
		server1 := rest.MustNewServer(c1.RestConf)
		defer server1.Stop()

		handler.RegisterHandlers(server1, ctx)

		fmt.Printf("Starting http proxy server at %s:%d...\n", c1.Host, c1.Port)
		server1.Start()
	}()

	// 3.启动https代理服务
	go func() {
		c2 := c
		c2.RestConf.Name = c.HTTPSProxy.Name
		c2.RestConf.Host = c.HTTPSProxy.Host
		c2.RestConf.Port = c.HTTPSProxy.Port
		c2.RestConf.CertFile = cert_file.Path("server.crt")
		c2.RestConf.KeyFile = cert_file.Path("server.key")

		server2 := rest.MustNewServer(c2.RestConf)
		defer server2.Stop()
		handler.RegisterHandlers(server2, ctx)

		fmt.Printf("Starting https proxy server at %s:%d...\n", c2.Host, c2.Port)
		server2.Start()
	}()
	defer server.Stop() // 2.确定服务启动和操作数据库
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}

//
//type ServiceDetail struct {
//	Info          *ga_service_info.GatewayServiceInfo                    `json:"info" description:"基本信息"`
//	HTTPRule      *ga_service_http_rule.GatewayServiceHttpRule           `json:"http_rule" description:"http_rule"`
//	TCPRule       *ga_service_tcp_rule.GatewayServiceTcpRule             `json:"tcp_rule" description:"tcp_rule"`
//	GRPCRule      *ga_service_grpc_rule.GatewayServiceGrpcRule           `json:"grpc_rule" description:"grpc_rule"`
//	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance     `json:"load_balance" description:"load_balance"`
//	AccessControl *ga_service_access_control.GatewayServiceAccessControl `json:"access_control" description:"access_control"`
//}
//
//var ErrNotFound = sqlx.ErrNotFound
//var ServiceManagerHandler *ServiceManager
//
//func First() {
//	fmt.Println("init")
//	err := ServiceManagerHandler.LoadOnce()
//
//	if err != nil {
//		fmt.Println("sdfsdfa", err)
//	}
//	fmt.Println("ServiceManagerHandler init", ServiceManagerHandler.ServiceSlice)
//}
//
//// ServiceManager 12.2 需要有个服务管理类去管理每个服务
//type ServiceManager struct {
//	ServiceMap   map[string]*ServiceDetail
//	ServiceSlice []*ServiceDetail
//	Locker       sync.RWMutex
//	Once         sync.Once
//	err          error
//}
//
//func NewServiceManager(ctx svc.ServiceContext) *ServiceManager {
//	return &ServiceManager{
//		ServiceMap:   make(map[string]*ServiceDetail),
//		ServiceSlice: []*ServiceDetail{},
//		Locker:       sync.RWMutex{},
//		Once:         sync.Once{},
//		err:          nil,
//	}
//}
//
//func (s *ServiceManager) LoadOnce() error {
//	httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
//	tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
//	grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}
//	loadBalance := &ga_service_load_balance.GatewayServiceLoadBalance{}
//	accessControl := &ga_service_access_control.GatewayServiceAccessControl{}
//
//	s.Once.Do(func() {
//		serviceInfoList, err := s.svcCtx.GatewayServiceInfoModel.FindAllTotal()
//		if err != nil && err != ErrNotFound {
//			fmt.Println("err", err)
//			s.err = err
//			return
//		}
//		s.Locker.Lock()
//		defer s.Locker.Unlock()
//
//		serviceInfo := serviceInfoList.([]ga_service_info.GatewayServiceInfo)
//		for _, service := range serviceInfo {
//			switch service.LoadType {
//			case errcode.LoadTypeHTTP:
//				httpRule, err = s.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(service.Id))
//				if err != nil && err != ErrNotFound {
//					s.err = err
//					return
//				}
//			case errcode.LoadTypeTCP:
//				tcpRule, err = s.svcCtx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(service.Id))
//				if err != nil && err != ErrNotFound {
//					s.err = err
//					return
//				}
//			default:
//				grpcRule, err = s.svcCtx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(service.Id))
//				if err != nil && err != ErrNotFound {
//					s.err = err
//					return
//				}
//			}
//			accessControl, err = s.svcCtx.GatewayServiceAccessControlModel.FindOneByServiceId(service.Id)
//			if err != nil && err != ErrNotFound {
//				s.err = err
//				return
//			}
//			loadBalance, err = s.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(service.Id))
//			if err != nil && err != ErrNotFound {
//				s.err = err
//
//				return
//			}
//
//			s1 := &ServiceDetail{
//				Info:          &service,
//				HTTPRule:      httpRule,
//				TCPRule:       tcpRule,
//				GRPCRule:      grpcRule,
//				LoadBalance:   loadBalance,
//				AccessControl: accessControl,
//			}
//			s.ServiceMap[service.ServiceName] = s1
//			s.ServiceSlice = append(s.ServiceSlice, s1)
//		}
//
//	})
//	return s.err
//}
//
//// 这里需要把服务信息加载到内存中
//
//// HTTPAccessMode 判断是否当前请求是否有对应服务
//func (s *ServiceManager) HTTPAccessMode(w http.ResponseWriter, r *http.Request) (*ServiceDetail, error) {
//	// 1.前缀匹配 /abc --> serviceSlice.rule
//	// 2.域名匹配 www.test.com --> serviceSlice.rule
//
//	// 现在就获取1.2
//	r.Host = strings.Split(r.Host, ":")[0]
//	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
//
//	for _, serviceItem := range s.ServiceSlice {
//		if serviceItem.Info.LoadType != errcode.LoadTypeHTTP {
//			continue
//		}
//
//		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypeDomain {
//			if serviceItem.HTTPRule.Rule == r.Host {
//				return serviceItem, nil
//			}
//		}
//		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
//			if serviceItem.HTTPRule.Rule == r.URL.Path {
//				return serviceItem, nil
//			}
//		}
//	}
//
//	return nil, errors.New("not matched service ")
//}
