package grpc_proxy_router

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/grpc_proxy_middleware"
	"API_Gateway/api/internal/proxy"
	"API_Gateway/api/internal/reverse_proxy"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc"
	"net"
)

var grpcServerList []*warpGrpcServer

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServer() {
	for _, serviceItem := range *global.SerInfo {
		tempDetail := serviceItem
		if tempDetail.Info.LoadType == global.LoadTypeGrpc {
			go func(serviceInfo global.ServiceDetail) {
				addr := fmt.Sprintf(":%d", serviceInfo.GRPCRule.Port)
				// 2.根据服务在数据库配置的下游服务器列表ipList和负载类型loadType 获得一个负载均衡器
				loadBalancer, err := global.LoadBalanceHandler.GetLoadBalancer(serviceInfo) // 负载均衡只维护下游服务器地址 不维护path
				//fmt.Println("loadBalancer:", Obj2Json(loadBalancer))
				if err != nil {
					logx.Error(" [INFO] GetTcpLoadBalancer %v err:%v", serviceInfo.GRPCRule.Port, err)
					return
				}

				// 启动服务器
				lis, err := net.Listen("tcp", addr)
				if err != nil {
					logx.Error(" [INFO] net.Listen err:%v", err)
					return
				}

				grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(loadBalancer)

				// 3.创建一个grpc服务器
				s := grpc.NewServer(
					grpc.ChainStreamInterceptor(
						grpc_proxy_middleware.GrpcFlowCountMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcFlowLimitMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcJwtAuthTokenMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcJwtFlowCountMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcJwtFlowLimitMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcWhiteListMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcBlackListMiddleware(&serviceInfo),
						grpc_proxy_middleware.GrpcHeaderTransferMiddleware(&serviceInfo),
					),
					grpc.CustomCodec(proxy.Codec()),
					grpc.UnknownServiceHandler(grpcHandler))

				grpcServerList = append(grpcServerList, &warpGrpcServer{
					Addr:   addr,
					Server: s,
				})
				logx.Infof(" [INFO] grpc_proxy_run %s", addr)

				// 4.启动服务器
				if err := s.Serve(lis); err != nil {
					logx.Error(" [INFO] grpc_proxy_run %v err:%v", addr, err)
				}

			}(tempDetail)
		}
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		logx.Info(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
