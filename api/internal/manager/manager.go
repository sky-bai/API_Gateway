package manager

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/svc"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/pkg/errcode"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"net/http"
	"strings"
	"sync"
)

var S1 *ServiceManager

func init() {

	S1 = NewServiceManager()
	err := S1.LoadOnce()
	if err != nil {
		panic(err)
		return
	}

	global.SerInfo = &S1.ServiceSlice

	//for _, detail := range S1.ServiceSlice {
	//	global.SerInfo = append(global.SerInfo, detail)
	//}
	//fmt.Println("ð", global.SerInfo)
	AppHandler = NewAppManager()
}

// ServiceManager éè¦åè·åå°æ¯ä¸ªæå¡çä¿¡æ¯
type ServiceManager struct {
	ServiceMap   map[string]*global.ServiceDetail // ä¸ä¸ªæå¡åå¯¹äºä¸ä¸ªæå¡çè¯¦ç»ä¿¡æ¯
	ServiceSlice []global.ServiceDetail           // æå¡åè¡¨
	Locker       sync.RWMutex                     // å ä¸ºè¦å¯¹mapè¿è¡è¯»åæä½ï¼æä»¥éè¦å é
	Once         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   make(map[string]*global.ServiceDetail),
		ServiceSlice: make([]global.ServiceDetail, 0),
		Locker:       sync.RWMutex{},
		Once:         sync.Once{},
		err:          nil,
	}
}

// AddService æ·»å æå¡
func (s *ServiceManager) AddService(serDetail global.ServiceDetail) {
	s.ServiceSlice = append(s.ServiceSlice, serDetail)
}

// ServiceDetail å¦ä½æ°çç»æä½å»è·å¾æ°æ®åºçæé
type ServiceDetail struct {
	Info          *ga_service_info.GatewayServiceInfo                    `json:"info" description:"åºæ¬ä¿¡æ¯"`
	HTTPRule      *ga_service_http_rule.GatewayServiceHttpRule           `json:"http_rule" description:"http_rule"`
	TCPRule       *ga_service_tcp_rule.GatewayServiceTcpRule             `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *ga_service_grpc_rule.GatewayServiceGrpcRule           `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *ga_service_load_balance.GatewayServiceLoadBalance     `json:"load_balance" description:"load_balance"`
	AccessControl *ga_service_access_control.GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}

// LoadOnce è·åå°æå¡åè¡¨çè¯¦ç»ä¿¡æ¯/**/
//func (s *ServiceManager) LoadOnce() error {
//	s.Once.Do(func() {
//		//1.è¯»åéç½®æä»¶å°ç»æä½ä¸­
//		var c config.Config
//		conf.MustLoad("etc/gateway-api.yaml", &c)
//
//		// éç½®æ°æ®åº
//		ctx := svc.NewServiceContext(c)
//
//		httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
//		tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
//		grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}
//
//		// è·åserviceinfoè¡¨ææä¿¡æ¯ ç¶åéè¿serviceIDè·åè¯¦ç»ä¿¡æ¯
//		all, err := ctx.GatewayServiceInfoModel.FindAll("", 1, 99)
//		if err != nil {
//			s.err = err
//			return
//		}
//		fmt.Println("123",all)
//		pageList := all.(*util.PageList)
//		serviceList := pageList.Data.([]ga_service_info.GatewayServiceInfo)
//
//		s.Locker.Lock()
//		defer s.Locker.Unlock()
//		//serviceAddr := "unknown"
//
//		for _, serviceInfo := range serviceList {
//			// è´è½½ç±»å 0=http 1=tcp 2=grpc
//			switch serviceInfo.LoadType {
//			case errcode.LoadTypeHTTP:
//				httpRule, err = ctx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
//				if err != nil {
//					s.err = err
//					return
//				}
//			case errcode.LoadTypeTCP:
//				tcpRule, err = ctx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(serviceInfo.Id))
//				if err != nil {
//					s.err = err
//					return
//				}
//			default:
//				grpcRule, err = ctx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(serviceInfo.Id))
//				if err != nil {
//					s.err = err
//					return
//				}
//			}
//
//			//1ãhttpåç¼æ¥å¥ clusterIP+clusterPort+path
//			//2ãhttpååæ¥å¥ domain
//			//3ãtcpãgrpcæ¥å¥ clusterIP+servicePort
//			//clusterIP := ctx.Config.Cluster.ClusterIP
//			//clusterPort := ctx.Config.Cluster.ClusterPort
//			//clusterSSLPort := ctx.Config.Cluster.ClusterSslPort
//			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
//			//	httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
//			//	httpRule.NeedHttps == 1 {
//			//	serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, httpRule.Rule)
//			//}
//			//
//			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
//			//	httpRule.RuleType == errcode.HTTPRuleTypePrefixURL &&
//			//	httpRule.NeedHttps == 0 {
//			//	serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, httpRule.Rule)
//			//}
//			//
//			//if serviceInfo.LoadType == errcode.LoadTypeHTTP &&
//			//	httpRule.RuleType == errcode.HTTPRuleTypeDomain {
//			//	serviceAddr = httpRule.Rule
//			//}
//			//if serviceInfo.LoadType == errcode.LoadTypeTCP {
//			//	serviceAddr = fmt.Sprintf("%s:%d", clusterIP, tcpRule.Port)
//			//}
//			//if serviceInfo.LoadType == errcode.LoadTypeGRPC {
//			//	serviceAddr = fmt.Sprintf("%s:%d", clusterIP, grpcRule.Port)
//			//
//			//}
//			loadBalance, err := ctx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(serviceInfo.Id))
//			if err != nil {
//				s.err = err
//				return
//			}
//
//			accessControl, err := ctx.GatewayServiceAccessControlModel.FindOneByServiceId(serviceInfo.Id)
//			if err != nil {
//				s.err = err
//				return
//			}
//
//			detail := global.ServiceDetail{
//				Info:          &serviceInfo,
//				HTTPRule:      httpRule,
//				TCPRule:       tcpRule,
//				GRPCRule:      grpcRule,
//				LoadBalance:   loadBalance,
//				AccessControl: accessControl,
//			}
//			s.ServiceMap[serviceInfo.ServiceName] = &detail
//			s.ServiceSlice = append(s.ServiceSlice, detail)
//
//		}
//
//	})
//	return s.err
//}
var ErrNotFound = sqlx.ErrNotFound

func (s *ServiceManager) LoadOnce() error {
	httpRule := &ga_service_http_rule.GatewayServiceHttpRule{}
	tcpRule := &ga_service_tcp_rule.GatewayServiceTcpRule{}
	grpcRule := &ga_service_grpc_rule.GatewayServiceGrpcRule{}
	loadBalance := &ga_service_load_balance.GatewayServiceLoadBalance{}
	accessControl := &ga_service_access_control.GatewayServiceAccessControl{}

	s.Once.Do(func() {
		//1.è¯»åéç½®æä»¶å°ç»æä½ä¸­
		var c config.Config
		conf.MustLoad("etc/gateway-api.yaml", &c)
		fmt.Println("è·åæ°æ®åºéç½®æå")
		// éç½®æ°æ®åº
		ctx := svc.NewServiceContext(c)
		serviceInfoList, err := ctx.GatewayServiceInfoModel.FindAllTotal()
		if err != nil && err != ErrNotFound {
			fmt.Println("err", err)
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()

		serviceInfo := serviceInfoList.([]ga_service_info.GatewayServiceInfo)
		for _, service := range serviceInfo {
			//fmt.Println("ðservice", service.Id)
			switch service.LoadType {
			case errcode.LoadTypeHTTP:
				httpRule, err = ctx.GatewayServiceHttpRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil {
					if err == sqlc.ErrNotFound {
						s.err = fmt.Errorf("æªæ¾å°æå¡idä¸º %d httpè®°å½", service.Id)
						return
					} else {
						s.err = err
						return
					}
				}
			case errcode.LoadTypeTCP:
				tcpRule, err = ctx.GatewayServiceTcpRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil {
					if err == sqlc.ErrNotFound {
						s.err = fmt.Errorf("æªæ¾å°æå¡idä¸º %d tcpè®°å½", service.Id)
						return
					} else {
						s.err = err
						return
					}
				}
			default:
				grpcRule, err = ctx.GatewayServiceGrpcRuleModel.FindOneByServiceId(int(service.Id))
				if err != nil {
					if err == sqlc.ErrNotFound {
						s.err = fmt.Errorf("æªæ¾å°æå¡idä¸º %d grpcè®°å½", service.Id)
						grpcRule = &ga_service_grpc_rule.GatewayServiceGrpcRule{}
						return
					} else {
						s.err = err
						return
					}
				}
			}
			accessControl, err = ctx.GatewayServiceAccessControlModel.FindOneByServiceId(service.Id)
			if err != nil {
				if err == sqlc.ErrNotFound {
					s.err = fmt.Errorf("æªæ¾å°æå¡idä¸º %d accessControl", service.Id)
					accessControl = &ga_service_access_control.GatewayServiceAccessControl{}
					return
				} else {
					s.err = err
					return
				}
			}
			loadBalance, err = ctx.GatewayServiceLoadBalanceModel.FindOneByServiceId(int(service.Id))
			if err != nil {
				if err == sqlc.ErrNotFound {
					s.err = fmt.Errorf("æªæ¾å°æå¡idä¸º %d loadbalanceè®°å½", service.Id)
					loadBalance = &ga_service_load_balance.GatewayServiceLoadBalance{}
					return
				} else {
					s.err = err
					return
				}
			}
			s1 := &global.ServiceDetail{
				Info:          service,
				HTTPRule:      *httpRule,
				TCPRule:       *tcpRule,
				GRPCRule:      *grpcRule,
				LoadBalance:   *loadBalance,
				AccessControl: *accessControl,
			}
			s.ServiceMap[service.ServiceName] = s1
			s.ServiceSlice = append(s.ServiceSlice, *s1)
		}

	})
	return s.err
}

// HTTPAccessMode åç«¯è¯·æ± ä¸åç«¯httpæå¡ æ³å¯¹æ¥
func (s *ServiceManager) HTTPAccessMode(r *http.Request) (*global.ServiceDetail, error) {
	//1ãåç¼å¹é /abc ==> serviceSlice.rule
	//2ãååå¹é www.test.com ==> serviceSlice.rule
	//   æ ¹æ®è¯·æ±å¯ä»¥å¾å°è¯¥è¯·æ±ççä¸»æº
	//   åå host c.Request.Host
	//
	//   åç¼ path c.Request.URL.Path
	host := r.Host
	host = host[0:strings.Index(host, ":")]

	path := r.URL.Path

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != errcode.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return &serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return &serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}

//
//var LoadBalancerHandler *LoadBalancer
//
//// LoadBalancer æ¯ä¸ªæå¡éè¦å¯¹åºèªå·±çè´è½½åè¡¡ç®æ³
//type LoadBalancer struct {
//	LoadBalanceMap   map[string]load_balance.LoadBalance
//	LoadBalanceSlice []load_balance.LoadBalance
//	Locker           sync.RWMutex
//}
//
//func NewLoadBalancer() *LoadBalancer {
//	return &LoadBalancer{
//		LoadBalanceSlice: []load_balance.LoadBalance{},
//		LoadBalanceMap:   map[string]load_balance.LoadBalance{},
//		Locker:           sync.RWMutex{}}
//}
//
//func (lbr *LoadBalancer) GetLoadBalancer(service global.ServiceDetail) (load_balance.LoadBalance, error) {
//	schema := "http"
//	if service.HTTPRule.NeedHttps == 1 {
//		schema = "https"
//	}
//	prefix := ""
//	if service.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
//		prefix = service.HTTPRule.Rule
//	}
//
//	ipList := service.LoadBalance.IpList
//	weightList := service.LoadBalance.WeightList
//
//	ipConf:=map[string]string{}
//	for ipIndex, ipItem := range ipList {
//		ipConf[string(ipItem)] = string(weightList[ipIndex])
//	}
//	checkConf, err := load_balance.NewLoadBalanceCheckConf(
//		fmt.Sprintf("%s:%s%s", schema, prefix),ipConf)
//	if err != nil {
//		return nil, err
//	}
//	load_balance.LoadBalanceZkConfInterface(service.LoadBalance.RoundType,checkConf)
//
//}
