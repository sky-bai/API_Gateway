package serviceGrpc

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"errors"
	"fmt"
	"strings"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddGrpcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGrpcLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddGrpcLogic {
	return AddGrpcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 增加grpc服务
func (l *AddGrpcLogic) AddGrpc(req types.AddGrpcRequest) (*types.Reponse, error) {

	// 1.先判断是否已存在的服务名
	serviceId, err := l.svcCtx.GatewayServiceInfoModel.FindOneByServiceName(req.ServiceName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if serviceId < 1 {
		return nil, errors.New("该服务名已存在")
	}

	// 2.判断端口是否被占用
	grpcId, err := l.svcCtx.GatewayServiceGrpcRuleModel.FindIdByPort(req.Port)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if grpcId > 0 {
		return nil, errors.New("端口已占用")
	}
	tcpId, err := l.svcCtx.GatewayServiceTcpRuleModel.FindOneByPort(req.Port)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("端口已占用")
	}
	if tcpId > 0 {
		return nil, errors.New("端口已占用")
	}

	// 3.ip与权重数量一致
	if len(strings.Split(req.IpList, ",")) != len(strings.Split(req.WeightList, ",")) {
		return nil, errors.New("ip列表与权重设置不匹配")
	}

	// 4.服务信息表
	serviceInfo := ga_service_info.GatewayServiceInfo{}
	serviceInfo.ServiceName = req.ServiceName
	serviceInfo.ServiceDesc = req.ServiceDesc
	serviceInfo.LoadType = 2 // 负载类型 0=http 1=tcp 2=grpc

	// 5.负载均衡表
	loadBalance := ga_service_load_balance.GatewayServiceLoadBalance{}
	loadBalance.RoundType = int64(req.RoundType)
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.ForbidList = req.ForbidList

	// 7.grpc规则表
	grpcRule := ga_service_grpc_rule.GatewayServiceGrpcRule{}
	grpcRule.Port = int64(req.Port)

	// 8.访问控制表
	accessControl := ga_service_access_control.GatewayServiceAccessControl{}
	accessControl.OpenAuth = int64(req.OpenAuth)
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.WhiteHostName = req.WhiteHostName
	accessControl.ClientipFlowLimit = int64(req.ClientIPFlowLimit)
	accessControl.ServiceFlowLimit = int64(req.ServiceFlowLimit)

	// 9.写入数据
	err = l.svcCtx.GatewayServiceGrpcRuleModel.InsertGrpcData(serviceInfo, grpcRule, accessControl, loadBalance)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &types.Reponse{Msg: "新增grpc服务成功"}, nil
}
