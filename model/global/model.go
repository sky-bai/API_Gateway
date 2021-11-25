package global

import (
	"API_Gateway/api/interna/svc"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/pkg/errcode"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"sync"
)

type ServiceDetail struct {
	Info          *ga_service_info.GatewayServiceInfo                    `json:"info" description:"基本信息"`
	HTTPRule      *ga_service_http_rule.GatewayServiceHttpRule           `json:"http_rule" description:"http_rule"`
	TCPRule       *ga_service_tcp_rule.GatewayServiceTcpRule             `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *ga_service_grpc_rule.GatewayServiceGrpcRule           `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance     `json:"load_balance" description:"load_balance"`
	AccessControl *ga_service_access_control.GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}

var ErrNotFound = sqlx.ErrNotFound
var ServiceManagerHandler *ServiceManager

func First() {
	fmt.Println("init")
	err := ServiceManagerHandler.LoadOnce()

	if err != nil {
		fmt.Println("sdfsdfa", err)
	}
	fmt.Println("ServiceManagerHandler init", ServiceManagerHandler.ServiceSlice)
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	Once         sync.Once
	err          error
	svcCtx       *svc.ServiceContext
}

func NewServiceManager(ctx *svc.ServiceContext) *ServiceManager {
	return &ServiceManager{
		ServiceMap:   make(map[string]*ServiceDetail),
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		Once:         sync.Once{},
		err:          nil,
		svcCtx:       ctx,
	}
}

func (s *ServiceManager) LoadOnce() error {
	httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
	tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
	grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}
	loadBalance := &ga_service_load_balance.GatewayServiceLoadBalance{}
	accessControl := &ga_service_access_control.GatewayServiceAccessControl{}

	s.Once.Do(func() {
		serviceInfoList, err := s.svcCtx.GatewayServiceInfoModel.FindAllTotal()
		if err != nil && err != ErrNotFound {
			fmt.Println("err", err)
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()

		serviceInfo := serviceInfoList.([]ga_service_info.GatewayServiceInfo)
		for _, service := range serviceInfo {
			switch service.LoadType {
			case errcode.LoadTypeHTTP:
				httpRule, err = s.svcCtx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil && err != ErrNotFound {
					s.err = err
					return
				}
			case errcode.LoadTypeTCP:
				tcpRule, err = s.svcCtx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil && err != ErrNotFound {
					s.err = err
					return
				}
			default:
				grpcRule, err = s.svcCtx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil && err != ErrNotFound {
					s.err = err
					return
				}
			}
			accessControl, err = s.svcCtx.GatewayServiceAccessControlModel.FindOneByServiceId(service.Id)
			if err != nil && err != ErrNotFound {
				s.err = err
				return
			}
			loadBalance, err = s.svcCtx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(service.Id))
			if err != nil && err != ErrNotFound {
				s.err = err

				return
			}

			s1 := &ServiceDetail{
				Info:          &service,
				HTTPRule:      httpRule,
				TCPRule:       tcpRule,
				GRPCRule:      grpcRule,
				LoadBalance:   loadBalance,
				AccessControl: accessControl,
			}

			s.ServiceMap[service.ServiceName] = s1
			s.ServiceSlice = append(s.ServiceSlice, s1)
		}

	})
	return s.err
}

// 这里需要把服务信息加载到内存中
