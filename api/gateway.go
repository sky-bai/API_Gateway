package main

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/config/cert_file"
	"API_Gateway/api/internal/handler"
	"API_Gateway/api/internal/http_proxy_router"
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
	fmt.Println(*configFile)

	// 配置数据库
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)

	// 2.启动http代理服务
	go func() {
		new(http_proxy_router.ServiceManager).LoadOnce()

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
