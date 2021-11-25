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

type UpdateGrpcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGrpcLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateGrpcLogic {
	return UpdateGrpcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 更新grpc服务
func (l *UpdateGrpcLogic) UpdateGrpc(req types.UpdateGrpcRequest) (*types.Reponse, error) {

	// 1.ip与权重数量一致
	if len(strings.Split(req.IpList, ",")) != len(strings.Split(req.WeightList, ",")) {
		return nil, errors.New("ip与权重数量不一致")
	}

	// 2.需要根据服务id 判断是否有已存在的服务
	serviceId, err := l.svcCtx.GatewayServiceGrpcRuleModel.FindIdByServiceId(int(req.ID))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("已存在该服务")
	}

	if serviceId < 1 {
		return nil, errors.New("该服务不存在")
	}

	// 3.数据库更新该服务
	service := ga_service_info.GatewayServiceInfo{}
	service.ServiceDesc = req.ServiceDesc
	service.ServiceName = req.ServiceName

	// 4.grpc规则表
	grpcRule := ga_service_grpc_rule.GatewayServiceGrpcRule{}
	grpcRule.Port = int64(req.Port)
	grpcRule.HeaderTransfor = req.HeaderTransfor

	// 5.访问控制表
	accessControl := ga_service_access_control.GatewayServiceAccessControl{}
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.OpenAuth = int64(req.OpenAuth)
	accessControl.ClientipFlowLimit = int64(req.ClientIPFlowLimit)
	accessControl.ServiceFlowLimit = int64(req.ServiceFlowLimit)

	// 6.负载均衡表
	loadBalance := ga_service_load_balance.GatewayServiceLoadBalance{}
	loadBalance.RoundType = int64(req.RoundType)
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList

	return &types.Reponse{}, nil
}
