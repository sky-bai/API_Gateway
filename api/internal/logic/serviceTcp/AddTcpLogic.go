package serviceTcp

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"context"
	"errors"
	"fmt"
	"strings"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddTcpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTcpLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddTcpLogic {
	return AddTcpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AddTcp 增加tcp服务
func (l *AddTcpLogic) AddTcp(req types.AddTcpRequest) (*types.Response, error) {

	// 1.检查该服务名是否被占用
	serviceId, err := l.svcCtx.GatewayServiceInfoModel.FindOneByServiceName(req.ServiceName)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("检查该服务名是否被占用操作失败")
	}
	if serviceId > 0 {
		return nil, errors.New("服务名被占用，请重新输入")
	}

	// 2.检查该服务的端口号是否被占用
	isTcpUsed, err := l.svcCtx.GatewayServiceTcpRuleModel.FindOneByPort(req.Port)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("检查该服务端口是否被占用操作失败")
	}
	if isTcpUsed {
		return nil, errors.New("该服务端口被占用，请重新输入")
	}

	// 3.检查ip与权重数量是否一致
	if len(strings.Split(req.IpList, ",")) != len(strings.Split(req.WeightList, ",")) {
		return nil, errors.New("ip列表与权重设置不匹配")
	}

	// 4.服务信息表
	serviceInfo := ga_service_info.GatewayServiceInfo{}
	serviceInfo.ServiceName = req.ServiceName
	serviceInfo.ServiceDesc = req.ServiceDesc
	serviceInfo.LoadType = 1 // 负载类型 0=http 1=tcp 2=grpc

	// 5.负载均衡表
	loadBalance := ga_service_load_balance.GatewayServiceLoadBalance{}
	loadBalance.RoundType = int64(req.RoundType)
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.ForbidList = req.ForbidList

	// 6.tcp规则表
	tcpRule := ga_service_tcp_rule.GatewayServiceTcpRule{}
	tcpRule.Port = int64(req.Port)

	// 7.访问控制表
	accessControl := ga_service_access_control.GatewayServiceAccessControl{}
	accessControl.OpenAuth = int64(req.OpenAuth)
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.WhiteHostName = req.WhiteHostName
	accessControl.ClientipFlowLimit = int64(req.ClientIPFlowLimit)
	accessControl.ServiceFlowLimit = int64(req.ServiceFlowLimit)

	// 8.在表里添加tcp服务
	err = l.svcCtx.GatewayServiceInfoModel.InsertTcpService(serviceInfo, tcpRule, accessControl, loadBalance)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("添加tcp服务失败")
	}
	s1 := global.ServiceDetail{
		Info:    serviceInfo,
		TCPRule: tcpRule,

		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	tem := *global.SerInfo
	tem = append(tem, s1)
	global.SerInfo = &tem

	return &types.Response{Msg: "增加tcp服务成功"}, nil
}

// 1.增加tcp服务后后会更新全局变量里面的服务信息
// updating service info in global variable after adding tcp service
