package serviceGrpc

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/middleware"
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_grpc_rule"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"gopkg.in/go-playground/validator.v9"
	"strings"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

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

type ErrorString struct {
	errMessage string
}

func (e *ErrorString) Error() string {
	return e.errMessage
}

// AddGrpc 增加grpc服务
func (l *AddGrpcLogic) AddGrpc(req types.AddGrpcRequest) (*types.Reponse, error) {

	errMessage := ErrorString{errMessage: ""}

	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}

	// 1.先判断是否已存在的服务名
	serviceId, err := l.svcCtx.GatewayServiceInfoModel.FindOneByServiceName(req.ServiceName)
	if err != nil {
		if err == sqlc.ErrNotFound {
		} else {
			fmt.Println(err, "err")
			return nil, errors.New("添加该服务失败")
		}
	}
	if serviceId > 1 {
		return nil, errors.New("该服务名已存在")
	}

	// 2.判断端口是否被占用
	isGrpcUsed, err := l.svcCtx.GatewayServiceGrpcRuleModel.FindIdByPort(req.Port)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("添加该服务失败")
	}
	if isGrpcUsed {
		return nil, errors.New("端口已占用")
	}

	isTcpUsed, err := l.svcCtx.GatewayServiceTcpRuleModel.FindOneByPort(req.Port)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("添加该服务失败")
	}
	if isTcpUsed {
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
		return nil, errors.New("添加111该服务失败")
	}

	s1 := global.ServiceDetail{
		Info:     serviceInfo,
		GRPCRule: grpcRule,

		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	tem := *global.SerInfo
	tem = append(tem, s1)
	global.SerInfo = &tem

	return &types.Reponse{Msg: "新增grpc服务成功"}, nil
}
