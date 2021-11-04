package errcode

const (
	// OK 成功返回
	OK = 200

	/**(前3位代表业务,后三位代表具体功能)**/

	// BadReuqestError 全局错误码
	BadReuqestError  = 100001
	ReuqesParamError = 100002
	WriteDataFail    = 100003
	ReadDataFail     = 100004
	NoData           = 100005
	StartEndTimeFail = 100006
	// UserNotFound 用户模块
	UserNotFound       = 200001
	UserHasRegist      = 200002
	UserErrCode        = 200003
	UserErrPwd         = 200004
	UserErrArgs        = 200005
	UserFailFixPwd     = 200006
	UserNotAuth        = 200007
	UserNotOpenId      = 200008
	UserTokenFailSet   = 200009
	UserErrpPwd        = 200010
	UserFailBindOpenId = 200011
	UserHadOpenId      = 200012
	UserFailRegister   = 200013
	UserErrCardId      = 200014
	UserFailPay        = 200015
	UserNoBank         = 200016
	UserErrarg         = 200017

	// 服务类型
	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	// http 类型
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
)
