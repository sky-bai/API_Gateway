package global

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
)

var SerInfo []ServiceDetail

type ServiceDetail struct {
	Info          *ga_service_info.GatewayServiceInfo                    `json:"info" description:"基本信息"`
	HTTPRule      *ga_service_http_rule.GatewayServiceHttpRule           `json:"http_rule" description:"http_rule"`
	TCPRule       *ga_service_tcp_rule.GatewayServiceTcpRule             `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *ga_service_grpc_rule.GatewayServiceGrpcRule           `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance     `json:"load_balance" description:"load_balance"`
	AccessControl *ga_service_access_control.GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}

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
//	svcCtx       svc.ServiceContext
//}
//
//func NewServiceManager(ctx svc.ServiceContext) *ServiceManager {
//	return &ServiceManager{
//		ServiceMap:   make(map[string]*ServiceDetail),
//		ServiceSlice: []*ServiceDetail{},
//		Locker:       sync.RWMutex{},
//		Once:         sync.Once{},
//		err:          nil,
//		svcCtx:       ctx,
//	}
//}
//
//// LoadOnce 获取服务列表
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
