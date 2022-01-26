package main

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/grpc_proxy_router"
	"API_Gateway/api/internal/handler"
	_ "API_Gateway/api/internal/manager"
	"API_Gateway/api/internal/reverse_proxy"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/tcp_proxy_middleware"
	"API_Gateway/api/internal/tcp_server"
	"context"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/router"
	"net/http"
)

const LoadTypeTCP = 1

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

type hello struct{}

func (h hello) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("ok"))
}
func main() {
	flag.Parse()
	*configFile = "etc/gateway-api.yaml"
	// 1.读取配置文件到结构体中
	var c config.Config
	conf.MustLoad(*configFile, &c)

	var logConf logx.LogConf
	logConf.Path = "logs"
	logConf.Mode = "file"
	// 从 yaml 文件中 初始化配置
	//*logFile = "etc/config.yaml"
	//conf.MustLoad("etc/config.yaml", &logConf)
	//logx 根据配置初始化

	logx.MustSetup(logConf)
	fmt.Printf("%+v\n", logConf)
	logx.Info("logx init success")

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
		// 这里设置全局中间价
		//server1.Use(ctx.HTTPAccessMode)
		fmt.Printf("Starting http proxy server at %s:%d...\n", c1.Host, c1.Port)
		server1.Start()
	}()

	// 3.启动https代理服务
	//go func() {
	//	c2 := c
	//	c2.RestConf.Name = c.HTTPSProxy.Name
	//	c2.RestConf.Host = c.HTTPSProxy.Host
	//	c2.RestConf.Port = c.HTTPSProxy.Port
	//	c2.RestConf.CertFile = cert_file.Path("server.crt")
	//	c2.RestConf.KeyFile = cert_file.Path("server.key")
	//
	//	server2 := rest.MustNewServer(c2.RestConf)
	//	defer server2.Stop()
	//	handler.RegisterHandlers(server2, ctx)
	//
	//	fmt.Printf("Starting https proxy server at %s:%d...\n", c2.Host, c2.Port)
	//	server2.Start()
	//}()

	go func() {

		logConf := logx.LogConf{
			Mode: "console",
		}
		serviceConf := service.ServiceConf{
			Mode: "dev",
		}
		// 1.get router
		router := router.NewRouter()
		var he hello
		router.Handle("GET", "/", he)
		// 2.use rest.WithRouter
		runOption := rest.WithRouter(router)
		restConf := rest.RestConf{
			Host: "127.0.0.1",
			Port: 8907,
		}
		// 3.set conf
		restConf.ServiceConf = serviceConf
		restConf.ServiceConf.Log = logConf
		// 4.get NewServer
		newServer, err := rest.NewServer(restConf, runOption)
		if err != nil {
			fmt.Printf("new server err:%+v\n", err)
			return
		}
		// 5.use middleware
		newServer.Use(func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				logx.Info("global middleware")
				fmt.Printf("global middleware\n")
				next(w, r)
			}
		})
		defer newServer.Stop()
		// 6. start server
		newServer.Start()
	}()

	// 4.启动tcp代理服务
	for _, serverItem := range *global.SerInfo {
		tempDetail := serverItem
		if tempDetail.Info.LoadType == LoadTypeTCP {
			fmt.Println("id", tempDetail.Info.Id)
			go func(tempItem global.ServiceDetail) {

				addr := fmt.Sprintf(":%d", tempItem.TCPRule.Port)
				// 2.根据服务在数据库配置的下游服务器列表ipList和负载类型loadType 获得一个负载均衡器
				loadBalancer, err := global.LoadBalanceHandler.GetLoadBalancer(tempItem) // 负载均衡只维护下游服务器地址 不维护path
				//fmt.Println("loadBalancer:", Obj2Json(loadBalancer))
				if err != nil {
					logx.Error(" [INFO] GetTcpLoadBalancer %v err:%v\n", tempItem.TCPRule.Port, err)
					return
				}

				//构建路由及设置中间件
				router := tcp_proxy_middleware.NewTcpSliceRouter()
				router.Group("/").Use(
					tcp_proxy_middleware.TCPFlowCountMiddleware(),
					tcp_proxy_middleware.TCPFlowLimitMiddleware(),
					tcp_proxy_middleware.TCPWhiteListMiddleware(),
					tcp_proxy_middleware.TCPBlackListMiddleware(),
				)

				//构建回调handler
				routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
					func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
						return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, loadBalancer)
					}, router)
				baseCtx := context.WithValue(context.Background(), "service", tempItem)

				tcpServer := tcp_server.TcpServer{
					Addr:    addr,
					Handler: routerHandler,
					BaseCtx: baseCtx,
				}
				fmt.Println("tcp 服务已启动", tcpServer.Addr)
				if err := tcpServer.ListenAndServe(); err != nil {
					logx.Error(" [INFO] GetTcpLoadBalancer %v err:%v\n", tempItem.TCPRule.Port, err)
					fmt.Println("--")
					fmt.Println(err)
					return
				}

			}(tempDetail)
		}
	}

	// 5.启动grpc代理服务
	go grpc_proxy_router.GrpcServer()

	defer server.Stop() // 2.确定服务启动和操作数据库
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
