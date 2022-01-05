package manager

import (
	"API_Gateway/api/internal/config"
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/svc"
	"API_Gateway/model/ga_gateway_app"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tal-tech/go-zero/core/conf"
	"sync"
)

var AppHandler *AppManager

func init() {
	AppHandler = NewAppManager()
	err := AppHandler.LoadOnce()
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("AppManagerHandler init success", AppHandler.AppSlice)
}

type AppManager struct {
	AppMap   map[string]*ga_gateway_app.GatewayApp
	AppSlice []*ga_gateway_app.GatewayApp
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

func (s *AppManager) GetAppList() []*ga_gateway_app.GatewayApp {
	return s.AppSlice
}
func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   map[string]*ga_gateway_app.GatewayApp{},
		AppSlice: []*ga_gateway_app.GatewayApp{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func (s *AppManager) LoadOnce() error {
	s.init.Do(func() {

		var c config.Config
		conf.MustLoad("etc/gateway-api.yaml", &c)
		fmt.Println("获取数据库配置成功")
		// 配置数据库
		ctx := svc.NewServiceContext(c)
		appList, err := ctx.GatewayAppModel.GetAllServiceList()
		if err != nil {
			panic(err)
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _, listItem := range appList {
			tmpItem := listItem
			s.AppMap[listItem.AppId] = &tmpItem
			s.AppSlice = append(s.AppSlice, &tmpItem)
			global.AppInfo.AppMap[listItem.AppId] = &tmpItem
			global.AppInfo.AppSlice = append(global.AppInfo.AppSlice, &tmpItem)
		}
		fmt.Println("获取应用列表成功", o2json(s.AppSlice))
	})
	return s.err
}

func o2json(o interface{}) string {
	marshal, err := jsoniter.Marshal(o)
	if err != nil {
		return "转换失败"
	}
	return string(marshal)
}
