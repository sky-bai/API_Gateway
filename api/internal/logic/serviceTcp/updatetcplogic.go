package serviceTcp

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"context"
	"errors"
	"strings"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateTcpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTcpLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateTcpLogic {
	return UpdateTcpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 更新Tcp服务
func (l *UpdateTcpLogic) UpdateTcp(req types.UpdateTcpRequest) (*types.Reponse, error) {

	// 1.如果ip列表与权重列表不一样 就返回
	if len(strings.Split(req.IpList, "/n")) != len(strings.Split(req.WeightList, "/n")) {
		return nil, errors.New("ip列表和权重列表数量不一致")
	}

	// 2.需要根据服务id 判断是否有已存在的服务
	serviceInfo, err := l.svcCtx.GatewayServiceInfoModel.FindOne(req.ID)
	if err != nil {
		return nil, err
	}
	if serviceInfo.Id < 1 {
		return nil, errors.New("该http服务未存在")
	}

	// 3.数据库更新该服务
	service := ga_service_info.GatewayServiceInfo{}
	service.ServiceDesc = req.ServiceDesc
	service.ServiceName = req.ServiceName

	// 4.tcp规则表
	tcpRule := ga_service_tcp_rule.GatewayServiceTcpRule{}
	tcpRule.Port = int64(req.Port)

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

	return &types.Reponse{Msg: "更新tcp服务成功"}, nil
}
