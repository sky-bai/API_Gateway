// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	UserName string `json:"username" form:"username" comment:"管理员用户名"  validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" validate:"required"`
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
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	Qps         int64  `json:"qps" form:"qps"`                   //qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   //qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
}

type PageListReponse struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Count     int         `json:"count"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"data"`
}

type ServiceResquest struct {
	ID int64 `json:"id" form:"id"` //id
}

type AddHTTPResquest struct {
	ServiceName            string `json:"service_name" validate:"required"`                  //服务名
	ServiceDesc            string `json:"service_desc" validate:"required,max=255,min=1"`    //服务描述
	RuleType               int    `json:"rule_type" validate:"max=1,min=0"`                  //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`              //域名或者前缀
	NeedHttps              int    `json:"need_https" validate:"max=1,min=0"`                 //支持https
	NeedStripUri           int    `json:"need_strip_uri"  example:"" validate:"max=1,min=0"` //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"  example:"" validate:"max=1,min=0"` //是否支持websocket
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
	ID                     int64  `json:"id"  validate:"required,min=1"`                                                                                      //服务ID
	ServiceName            string `json:"service_name" validate:"required,valid_service_name"`                                                                //服务名
	ServiceDesc            string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,max=255,min=1"` //服务描述
	RuleType               int    `json:"rule_type"  validate:"max=1,min=0"`                                                                                  //接入类型
	Rule                   string `json:"rule"  validate:"required,valid_rule"`                                                                               //域名或者前缀
	NeedHttps              int    `json:"need_https"  validate:"max=1,min=0"`                                                                                 //支持https
	NeedStripUri           int    `json:"need_strip_uri"  validate:"max=1,min=0"`                                                                             //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket"  validate:"max=1,min=0"`                                                                             //是否支持websocket
	UrlRewrite             string `json:"url_rewrite"  validate:"valid_url_rewrite"`                                                                          //url重写功能
	HeaderTransfor         string `json:"header_transfor"  validate:"valid_header_transfor"`                                                                  //header转换
	OpenAuth               int    `json:"open_auth"  validate:"max=1,min=0"`                                                                                  //关键词
	BlackList              string `json:"black_list"   validate:""`                                                                                           //黑名单ip
	WhiteList              string `json:"white_list"  validate:""`                                                                                            //白名单ip
	ClientipFlowLimit      int    `json:"clientip_flow_limit"  validate:"min=0"`                                                                              //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit"  example:"" validate:"min=0"`                                                                    //服务端限流
	RoundType              int    `json:"round_type" validate:"max=3,min=0"`                                                                                  //轮询方式
	IpList                 string `json:"ip_list"  validate:"required,valid_ipportlist"`                                                                      //ip列表
	WeightList             string `json:"weight_list" validate:"required,valid_weightlist"`                                                                   //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" " validate:"min=0"`                                                                        //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout"  validate:"min=0"`                                                                          //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout"  validate:"min=0"`                                                                            //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle"  validate:"min=0"`                                                                                //最大空闲链接数
}
