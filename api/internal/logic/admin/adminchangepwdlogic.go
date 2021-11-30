package admin

import (
	"API_Gateway/model/ga_admin"
	"API_Gateway/pkg/errcode"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminChangePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminChangePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminChangePwdLogic {
	return AdminChangePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 管理员修改密码
func (l *AdminChangePwdLogic) AdminChangePwd(req types.FixPwdRequest) (*types.FixPwdReponse, error) {

	//1. 从token里面取到userId
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致
	value := l.ctx.Value("userId")
	i, _ := value.(json.Number).Int64()
	strInt64 := strconv.FormatInt(i, 10)
	userId, _ := strconv.Atoi(strInt64)

	data := ga_admin.GatewayAdmin{}
	data.Id = int64(userId)
	data.Password = req.Password

	// 从数据库查找管理员信息
	err := l.svcCtx.GatewayAdminModel.UpdatePwd(data)
	if err != nil {
		fmt.Println("查询数据库err :", err)
		return nil, errcode.NewErrCode(errcode.UserNotFound)
	}

	return &types.FixPwdReponse{Msg: "修改成功"}, nil
}
