package util

var message = map[int]string{
	//成功返回
	OK: "SUCCESS",
	//全局错误码
	BadReuqestError:  "服务器繁忙,请稍后再试",
	ReuqesParamError: "参数错误",
	WriteDataFail:    "写入数据失败，请稍后重试",
	ReadDataFail:     "读取数据失败，请稍后重试",
	NoData:           "数据不存在",
	StartEndTimeFail: "开始结束时间错误",
	//用户模块
	UserNotFound:       "用户不存在",
	UserHasRegist:      "用户已注册",
	UserErrCode:        "用户验证码错误",
	UserErrPwd:         "用户密码不一致错误",
	UserErrArgs:        "用户输入参数有误",
	UserFailFixPwd:     "用户重新修改密码失败",
	UserNotAuth:        "用户没有微信绑定",
	UserNotOpenId:      "用户没有绑定openid",
	UserTokenFailSet:   "用户Token设置失败",
	UserErrpPwd:        "用户密码输入错误",
	UserFailBindOpenId: "用户绑定openId失败",
	UserHadOpenId:      "用户已绑定openId",
	UserFailRegister:   "用户注册失败",
	UserErrCardId:      "用户身份证号码位数应该为18位",
	UserFailPay:        "用户支付失败",
	UserNoBank:         "用户没有银行卡",
	UserErrarg:         "用户余额不足",
}

func MapErrMsg(errcode int) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器繁忙,请稍后再试"
	}
}
