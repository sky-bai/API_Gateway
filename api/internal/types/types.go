// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	UserName string `json:"username"  comment:"管理员用户名"  validate:"required"` //管理员用户名
	Password string `json:"password"  validate:"required"`
}

type LoginReponse struct {
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	RefreshAfter int64  `json:"refreshAfter"`
}

type LogOutReponse struct {
	Message string `json:"message" form:"message"  validate:"required"` // 退出信息
}

type AdminInfoReponse struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	LoginTime    int      `json:"login_time"`
	Avatar       string   `json:"avatar"`
	Introduction string   `json:"introduction"`
	Roles        []string `json:"roles"`
}

type FixPwdRequest struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

type FixPwdReponse struct {
	Msg string `json:"msg"`
}

type CommonReponse struct {
	Msg string `json:"msg"`
}

type ServiceListResquest struct {
	Info     string `json:"info"`
	PageNo   int    `json:"page_no"`   //页数
	PageSize int    `json:"page_size"` //每页条数
}

type ServiceListItemReponse struct {
	ID          int64  `json:"id"`           //id
	ServiceName string `json:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc"` //服务描述
	LoadType    int    `json:"load_type"`    //类型
	ServiceAddr string `json:"service_addr"` //服务地址
	Qps         int64  `json:"qps"`          //qps
	Qpd         int64  `json:"qpd"`          //qpd
	TotalNode   int    `json:"total_node"`   //节点数
}

type DataList struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Count int         `json:"count"`
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

type PageListReponse struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Count     int         `json:"count"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"data"`
}

type ServiceResquest struct {
	ID int64 `json:"id"` // 服务id
}

type ServiceDetailResquest struct {
	ID int `json:"id"` // 服务id
}

type ServiceStatusResquest struct {
	ID int `json:"id"` // 服务id
}

type ServiceDetailResponse struct {
	ServiceName            string `json:"service_name" validate:"required,valid_service_name"` //服务名
	ServiceDesc            string `json:"service_desc" validate:"required,max=255,min=1"`      //服务描述
	RuleType               int    `json:"rule_type" validate:"max=1,min=0"`                    //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`                //域名或者前缀
	NeedHttps              int    `json:"need_https" validate:"max=1,min=0"`                   //支持https
	NeedStripUri           int    `json:"need_strip_uri"  example:"" validate:"max=1,min=0"`   //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"  example:"" validate:"max=1,min=0"`   //是否支持websocket
	UrlRewrite             string `json:"url_rewrite"  validate:"valid_url_rewrite"`           //url重写功能
	HeaderTransfor         string `json:"header_transfor"  validate:"valid_header_transfor"`   //header转换
	OpenAuth               int    `json:"open_auth"   validate:"max=1,min=0"`                  //关键词
	BlackList              string `json:"black_list"  validate:""`                             //黑名单ip
	WhiteList              string `json:"white_list"   validate:""`                            //白名单ip
	ClientipFlowLimit      int    `json:"clientip_flow_limit"  validate:"min=0"`               //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit"  validate:"min=0"`                //服务端限流
	RoundType              int    `json:"round_type"  validate:"max=3,min=0"`                  //轮询方式
	IpList                 string `json:"ip_list" validate:"required,valid_ipportlist"`        //ip列表
	WeightList             string `json:"weight_list"  validate:"required,valid_weightlist"`   //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" validate:"min=0"`           //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout"  validate:"min=0"`           //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout"  validate:"min=0"`             //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle"  validate:"min=0"`                 //最大空闲链接数
}

type ServiceStatusResponse struct {
	Today     []int `json:"today"`
	Yesterday []int `json:"yesterday"`
}

type PanelDataOutput struct {
	ServiceNum      int64 `json:"serviceNum"`
	AppNum          int64 `json:"appNum"`
	CurrentQPS      int64 `json:"currentQps"`
	TodayRequestNum int64 `json:"todayRequestNum"`
}

type ServiceFlowResponse struct {
	Today     []int `json:"today"`
	Yesterday []int `json:"yesterday"`
}

type PanelServiceStatusResponse struct {
	Legend []string              `json:"legend"`
	Data   []DashServiceStatItem `json:"data"`
}

type DashServiceStatItem struct {
	Name     string `json:"name"`
	LoadType int    `json:"load_type"`
	Value    int64  `json:"value"`
}

type ServiceDetail struct {
	Info          GatewayServiceInfo          `json:"info" description:"基本信息"`
	HTTPRule      GatewayServiceHttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       GatewayServiceTcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      GatewayServiceGrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   GatewayServiceLoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}

type GatewayServiceInfo struct {
	Id          int64  `db:"id" json:"id"`                      // 自增主键
	LoadType    int64  `db:"load_type" json:"load_type"`        // 负载类型 0=http 1=tcp 2=grpc
	ServiceName string `db:"service_name" json:"service_name"`  // 服务名称 6-128 数字字母下划线
	ServiceDesc string `db:"service_desc" json:"service_descs"` // 服务描述
	IsDelete    int64  `db:"is_delete" json:"is_delete"`        // 是否删除 1=删除 0=未删除
}

type GatewayServiceHttpRule struct {
	Id             int64  `db:"id" json:"id"`                           // 自增主键
	ServiceId      int64  `db:"service_id" json:"service_id"`           // 服务id
	RuleType       int64  `db:"rule_type" json:"rule_type"`             // 匹配类型 0=url前缀url_prefix 1=域名domain
	Rule           string `db:"rule" json:"rule"`                       // type=domain表示域名，type=url_prefix时表示url前缀
	NeedHttps      int64  `db:"need_https" json:"need_https"`           // 支持https 1=支持
	NeedStripUri   int64  `db:"need_strip_uri" json:"need_strip_uri"`   // 启用strip_uri 1=启用
	NeedWebsocket  int64  `db:"need_websocket" json:"need_websocket"`   // 是否支持websocket 1=支持
	UrlRewrite     string `db:"url_rewrite" json:"url_rewrite"`         // url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	HeaderTransfor string `db:"header_transfor" json:"header_transfor"` // header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔
}

type GatewayServiceTcpRule struct {
	Id        int64 `db:"id" json:"id"`                 // 自增主键
	ServiceId int64 `db:"service_id" json:"service_id"` // 服务id
	Port      int64 `db:"port" json:"port"`             // 端口号
}

type GatewayServiceGrpcRule struct {
	Id             int64  `db:"id" json:"id"`                           // 自增主键
	ServiceId      int64  `db:"service_id" json:"service_id"`           // 服务id
	Port           int64  `db:"port" json:"port"`                       // 端口
	HeaderTransfor string `db:"header_transfor" json:"header_transfor"` // header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔
}

type GatewayServiceLoadBalance struct {
	Id                     int64  `db:"id" json:"id"`                                             // 自增主键
	ServiceId              int64  `db:"service_id" json:"service_id"`                             // 服务id
	CheckMethod            int64  `db:"check_method" json:"check_method"`                         // 检查方法 0=tcpchk,检测端口是否握手成功
	CheckTimeout           int64  `db:"check_timeout" json:"check_timeout"`                       // check超时时间,单位s
	CheckInterval          int64  `db:"check_interval" json:"check_interval"`                     // 检查间隔, 单位s
	RoundType              int64  `db:"round_type" json:"round_type"`                             // 轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash
	IpList                 string `db:"ip_list" json:"ip_list"`                                   // ip列表
	WeightList             string `db:"weight_list" json:"weight_list"`                           // 权重列表
	ForbidList             string `db:"forbid_list" json:"forbid_list"`                           // 禁用ip列表
	UpstreamConnectTimeout int64  `db:"upstream_connect_timeout" json:"upstream_connect_timeout"` // 建立连接超时, 单位s
	UpstreamHeaderTimeout  int64  `db:"upstream_header_timeout" json:"upstream_header_timeout"`   // 获取header超时, 单位s
	UpstreamIdleTimeout    int64  `db:"upstream_idle_timeout" json:"upstream_idle_timeout"`       // 链接最大空闲时间, 单位s
	UpstreamMaxIdle        int64  `db:"upstream_max_idle" json:"upstream_max_idle"`               // 最大空闲链接数
}

type GatewayServiceAccessControl struct {
	Id                int64  `db:"id" json:"id"`                                   // 自增主键
	ServiceId         int64  `db:"service_id" json:"service_id"`                   // 服务id
	OpenAuth          int64  `db:"open_auth" json:"open_auth"`                     // 是否开启权限 1=开启
	BlackList         string `db:"black_list" json:"black_list"`                   // 黑名单ip
	WhiteList         string `db:"white_list" json:"white_list"`                   // 白名单ip
	WhiteHostName     string `db:"white_host_name" json:"white_host_name"`         // 白名单主机
	ClientipFlowLimit int64  `db:"clientip_flow_limit" json:"clientip_flow_limit"` // 客户端ip限流
	ServiceFlowLimit  int64  `db:"service_flow_limit" json:"service_flow_limit"`   // 服务端限流
}

type AddTcpRequest struct {
	ServiceName       string `json:"service_name"  validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc"  validate:"required"`
	Port              int    `json:"port"  validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor"  validate:""`
	OpenAuth          int    `json:"open_auth"  validate:""`
	BlackList         string `json:"black_list"  validate:"valid_iplist"`
	WhiteList         string `json:"white_list"  validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name"  validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit"  validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit"  validate:""`
	RoundType         int    `json:"round_type" validate:""`
	IpList            string `json:"ip_list"  validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list"  validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list"  validate:"valid_iplist"`
}

type Response struct {
	Msg string `json:"msg"`
}

type UpdateTcpRequest struct {
	ID                int64  `json:"id"  validate:"required"`                              // 服务ID
	ServiceName       string `json:"service_name"  validate:"required,valid_service_name"` // 服务名称
	ServiceDesc       string `json:"service_desc"  validate:"required"`                    // 服务描述
	Port              int    `json:"port"  validate:"required,min=8001,max=8999"`          // 端口，需要设置8001-8999范围内
	OpenAuth          int    `json:"open_auth"  validate:""`                               // 是否开启权限验证
	BlackList         string `json:"black_list" validate:"valid_iplist"`                   // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list"  validate:"valid_iplist"`                  // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name"  validate:"valid_iplist"`             // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit"   validate:""`                    // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit"  validate:""`                      // 服务端限流
	RoundType         int    `json:"round_type"  validate:""`                              // 轮询策略
	IpList            string `json:"ip_list"  validate:"required,valid_ipportlist"`        // IP列表
	WeightList        string `json:"weight_list"  validate:"required,valid_weightlist"`    // 权重列表
	ForbidList        string `json:"forbid_list"  validate:"valid_iplist"`                 // 禁用IP列表
}

type AddHTTPResquest struct {
	ServiceName            string `json:"service_name" validate:"valid_service_name"`        //服务名
	ServiceDesc            string `json:"service_desc" validate:"required,max=255,min=1"`    //服务描述
	RuleType               int    `json:"rule_type" validate:"max=1,min=0"`                  //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`              //域名或者前缀
	NeedHttps              int    `json:"need_https" validate:"max=1,min=0"`                 //支持https
	NeedStripUri           int    `json:"need_strip_uri"   validate:"max=1,min=0"`           //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"   validate:"max=1,min=0"`           //是否支持websocket
	UrlRewrite             string `json:"url_rewrite"  validate:"valid_url_rewrite"`         //url重写功能
	HeaderTransfor         string `json:"header_transfor"  validate:"valid_header_transfor"` //header转换
	OpenAuth               int    `json:"open_auth"   validate:"max=1,min=0"`                //关键词
	BlackList              string `json:"black_list"  validate:""`                           //黑名单ip
	WhiteList              string `json:"white_list"   validate:""`                          //白名单ip
	ClientipFlowLimit      int    `json:"clientip_flow_limit"  validate:"min=0"`             //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit"  validate:"min=0"`              //服务端限流
	RoundType              int    `json:"round_type"  validate:"max=3,min=0"`                //轮询方式
	IpList                 string `json:"ip_list" validate:"required,valid_ipportlist"`      //ip列表
	WeightList             string `json:"weight_list"  validate:"required,valid_weightlist"` //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" validate:"min=0"`         //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout"  validate:"min=0"`         //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout"  validate:"min=0"`           //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle"  validate:"min=0"`               //最大空闲链接数
}

type UpdateHTTPResquest struct {
	ID                     int64  `json:"id"  validate:"required,min=1"`                       //服务ID
	ServiceName            string `json:"service_name" validate:"required,valid_service_name"` //服务名
	ServiceDesc            string `json:"service_desc" validate:"required,max=255,min=1"`      //服务描述
	RuleType               int    `json:"rule_type"  validate:"max=1,min=0"`                   //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`                //域名或者前缀
	NeedHttps              int    `json:"need_https"  validate:"max=1,min=0"`                  //支持https
	NeedStripUri           int    `json:"need_strip_uri"  validate:"max=1,min=0"`              //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"  validate:"max=1,min=0"`              //是否支持websocket
	UrlRewrite             string `json:"url_rewrite"  validate:"valid_url_rewrite"`           //url重写功能
	HeaderTransfor         string `json:"header_transfor"  validate:"valid_header_transfor"`   //header转换
	OpenAuth               int    `json:"open_auth"  validate:"max=1,min=0"`                   //关键词
	BlackList              string `json:"black_list"   validate:""`                            //黑名单ip
	WhiteList              string `json:"white_list"  validate:""`                             //白名单ip
	ClientipFlowLimit      int    `json:"clientip_flow_limit"  validate:"min=0"`               //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit"  example:"" validate:"min=0"`     //服务端限流
	RoundType              int    `json:"round_type" validate:"max=3,min=0"`                   //轮询方式
	IpList                 string `json:"ip_list"  validate:"required,valid_ipportlist"`       //ip列表
	WeightList             string `json:"weight_list" validate:"required,valid_weightlist"`    //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" " validate:"min=0"`         //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout"  validate:"min=0"`           //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout"  validate:"min=0"`             //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle"  validate:"min=0"`                 //最大空闲链接数
}

type HttpReponse struct {
	Msg string `json:"msg"`
}

type HttpeDetailResquest struct {
	ID int64 `json:"id"` // 服务id
}

type HttpDetailResponse struct {
	ServiceName            string `json:"service_name" validate:"required,valid_service_name"` //服务名
	ServiceDesc            string `json:"service_desc" validate:"required,max=255,min=1"`      //服务描述
	RuleType               int    `json:"rule_type" validate:"max=1,min=0"`                    //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`                //域名或者前缀
	NeedHttps              int    `json:"need_https" validate:"max=1,min=0"`                   //支持https
	NeedStripUri           int    `json:"need_strip_uri"  example:"" validate:"max=1,min=0"`   //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"  example:"" validate:"max=1,min=0"`   //是否支持websocket
	UrlRewrite             string `json:"url_rewrite"  validate:"valid_url_rewrite"`           //url重写功能
	HeaderTransfor         string `json:"header_transfor"  validate:"valid_header_transfor"`   //header转换
	OpenAuth               int    `json:"open_auth"   validate:"max=1,min=0"`                  //关键词
	BlackList              string `json:"black_list"  validate:""`                             //黑名单ip
	WhiteList              string `json:"white_list"   validate:""`                            //白名单ip
	ClientipFlowLimit      int    `json:"clientip_flow_limit"  validate:"min=0"`               //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit"  validate:"min=0"`                //服务端限流
	RoundType              int    `json:"round_type"  validate:"max=3,min=0"`                  //轮询方式
	IpList                 string `json:"ip_list" validate:"required,valid_ipportlist"`        //ip列表
	WeightList             string `json:"weight_list"  validate:"required,valid_weightlist"`   //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" validate:"min=0"`           //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout"  validate:"min=0"`           //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout"  validate:"min=0"`             //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle"  validate:"min=0"`                 //最大空闲链接数
}

type AddGrpcRequest struct {
	ServiceName       string `json:"service_name"  validate:"required,valid_service_name"` // 服务名称
	ServiceDesc       string `json:"service_desc"  validate:"required"`                    // 服务描述
	Port              int    `json:"port"  validate:"required,min=8001,max=8999"`          // 端口，需要设置8001-8999范围内
	HeaderTransfor    string `json:"header_transfor" validate:"valid_header_transfor"`     // metadata转换
	OpenAuth          int    `json:"open_auth"  validate:""`                               // 是否开启权限验证
	BlackList         string `json:"black_list"  validate:"valid_iplist"`                  // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list"  validate:"valid_iplist"`                  // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name"  validate:"valid_iplist"`             // 白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit"  validate:""`                     // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit"  validate:""`                      // 服务端限流
	RoundType         int    `json:"round_type"  validate:""`                              // 轮询策略
	IpList            string `json:"ip_list"  validate:"required,valid_ipportlist"`        // IP列表
	WeightList        string `json:"weight_list"  validate:"required,valid_weightlist"`    // 权重列表
	ForbidList        string `json:"forbid_list"  validate:"valid_iplist"`                 // 禁用IP列表
}

type UpdateGrpcRequest struct {
	ID                int64  `json:"id"  validate:"required"`                              // 服务ID
	ServiceName       string `json:"service_name"  validate:"required,valid_service_name"` // 服务名称
	ServiceDesc       string `json:"service_desc"  validate:"required"`                    // 服务描述
	Port              int    `json:"port"  validate:"required,min=8001,max=8999"`          // 端口，需要设置8001-8999范围内
	HeaderTransfor    string `json:"header_transfor"  validate:"valid_header_transfor"`    // metadata转换
	OpenAuth          int    `json:"open_auth"  validate:""`                               // 是否开启权限验证
	BlackList         string `json:"black_list"  validate:"valid_iplist"`                  // 黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list"  validate:"valid_iplist"`                  // 白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name"  validate:"valid_iplist"`             //  白名单主机，以逗号间隔
	ClientIPFlowLimit int    `json:"clientip_flow_limit"  validate:""`                     // 客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit"  validate:""`                      // 服务端限流
	RoundType         int    `json:"round_type"  validate:""`                              // 轮询策略
	IpList            string `json:"ip_list"  validate:"required,valid_ipportlist"`        // IP列表
	WeightList        string `json:"weight_list"  validate:"required,valid_weightlist"`    // 权重列表
	ForbidList        string `json:"forbid_list"  validate:"valid_iplist"`                 // 禁用IP列表
}

type Reponse struct {
	Msg string `json:"msg"`
}

type PingReponse struct {
	Message string `json:"message"` // 返回信息
}

type HttpsReponse struct {
	Message string `json:"message"` // 返回信息
}

type AddAppRequest struct {
	AppID    string `json:"app_id"  validate:"required"` // 租户id
	Name     string `json:"name"  validate:"required"`   // 租户名称
	Secret   string `json:"secret"  validate:""`         // 密钥
	WhiteIPS string `json:"white_ips"`                   // ip白名单，支持前缀匹配
	Qpd      int64  `json:"qpd"   validate:""`           // 日请求量限制
	Qps      int64  `json:"qps"   validate:""`           // 每秒请求量限制
}

type UpdateAppRequest struct {
	ID       int64  `json:"id"  validate:"required"`
	AppID    string `json:"app_id"  validate:""`         // 租户id
	Name     string `json:"name" validate:"required"`    // 租户名称
	Secret   string `json:"secret"  validate:"required"` // 密钥
	WhiteIPS string `json:"white_ips"`                   // ip白名单，支持前缀匹配
	Qpd      int64  `json:"qpd"`                         // 日请求量限制
	Qps      int64  `json:"qps"`                         // 每秒请求量限制
}

type DeleteAppRequest struct {
	ID int `json:"id"  validate:"required"`
}

type AppDetailRequest struct {
	ID int `json:"id"  validate:"required"`
}

type AppListRequest struct {
	Info     string `json:"info"  comment:"查找信息" validate:""`
	PageSize int    `json:"page_size" comment:"页数" validate:"required,min=1,max=999"`
	PageNo   int    `json:"page_no"  comment:"页码" validate:"required,min=1,max=999"`
}

type AppStatusRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type AppResponse struct {
	Message string `json:"msg"`
}

type AppStatus struct {
	Today     []int `json:"today"`
	Yesterday []int `json:"yesterday"`
}

type GetTokenRequest struct {
	GrantType string `json:"grant_type" comment:"授权类型"  validate:"required"` //授权类型
	Scope     string `json:"scope"  comment:"权限范围"  validate:"required"`     //权限范围
}

type GetTokenResponse struct {
	AccessToken string `json:"access_token"` //access_token
	ExpiresIn   int    `json:"expires_in"`   //expires_in
	TokenType   string `json:"token_type"`   //token_type
	Scope       string `json:"scope"`        //scope
}

type APPListResponse struct {
	List  []APPListItemOutput `json:"list"  comment:"租户列表"`
	Total int64               `json:"total"  comment:"租户总数"`
}

type APPListItemOutput struct {
	ID        int64  `json:"id" gorm:"primary_key"`
	AppID     string `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配		"`
	Qpd       int64  `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64  `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	RealQpd   int64  `json:"real_qpd" description:"日请求量限制"`
	RealQps   int64  `json:"real_qps" description:"每秒请求量限制"`
	UpdatedAt string `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	CreatedAt string `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8   `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

type AppDetailResponse struct {
	ID        int64  `json:"id" gorm:"primary_key"`
	AppID     string `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64  `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64  `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt string `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	UpdatedAt string `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8   `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}
