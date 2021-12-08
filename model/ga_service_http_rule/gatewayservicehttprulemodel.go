package ga_service_http_rule

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	gatewayServiceHttpRuleFieldNames          = builderx.RawFieldNames(&GatewayServiceHttpRule{})
	gatewayServiceHttpRuleRows                = strings.Join(gatewayServiceHttpRuleFieldNames, ",")
	gatewayServiceHttpRuleRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceHttpRuleFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceHttpRuleRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceHttpRuleFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceHttpRuleModel interface {
		Insert(data GatewayServiceHttpRule) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceHttpRule, error)
		Update(data GatewayServiceHttpRule) error
		Delete(id int64) error

		// FindOneByServiceId 根据服务ID去查出一条http数据
		FindOneByServiceId(serviceID int) (*GatewayServiceHttpRule, error)

		// FindOneByRule 根据ruleType(匹配类型) 和 rule规则查找数据
		FindOneByRule(ruleType int, rule string) (int64, error)
	}

	defaultGatewayServiceHttpRuleModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceHttpRule struct {
		Id             int64  `db:"id"`              // 自增主键
		ServiceId      int64  `db:"service_id"`      // 服务id
		RuleType       int64  `db:"rule_type"`       // 匹配类型 0=url前缀url_prefix 1=域名domain
		Rule           string `db:"rule"`            // type=domain表示域名，type=url_prefix时表示url前缀
		NeedHttps      int64  `db:"need_https"`      // 支持https 1=支持
		NeedStripUri   int64  `db:"need_strip_uri"`  // 启用strip_uri 1=启用
		NeedWebsocket  int64  `db:"need_websocket"`  // 是否支持websocket 1=支持
		UrlRewrite     string `db:"url_rewrite"`     // url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
		HeaderTransfor string `db:"header_transfor"` // header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔
	}
)

func NewGatewayServiceHttpRuleModel(conn sqlx.SqlConn) GatewayServiceHttpRuleModel {
	return &defaultGatewayServiceHttpRuleModel{
		conn:  conn,
		table: "`gateway_service_http_rule`",
	}
}

func (m *defaultGatewayServiceHttpRuleModel) Insert(data GatewayServiceHttpRule) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, gatewayServiceHttpRuleRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ServiceId, data.RuleType, data.Rule, data.NeedHttps, data.NeedStripUri, data.NeedWebsocket, data.UrlRewrite, data.HeaderTransfor)
	return ret, err
}

func (m *defaultGatewayServiceHttpRuleModel) FindOne(id int64) (*GatewayServiceHttpRule, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceHttpRuleRows, m.table)
	var resp GatewayServiceHttpRule
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneByRule 根据ruleType(匹配类型) 和 rule规则查找数据
func (m *defaultGatewayServiceHttpRuleModel) FindOneByRule(ruleType int, rule string) (int64, error) {
	query := fmt.Sprintf("select service_id from %s where `rule_type` = ? and `rule` = ? limit 1", m.table)
	var resp int64
	err := m.conn.QueryRow(&resp, query, ruleType, rule)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

// FindOneByServiceId 根据服务ID去查出一条http数据
func (m *defaultGatewayServiceHttpRuleModel) FindOneByServiceId(serviceId int) (*GatewayServiceHttpRule, error) {
	query := fmt.Sprintf("select %s from %s where `service_id` = ? limit 1", gatewayServiceHttpRuleRows, m.table)
	var resp GatewayServiceHttpRule
	err := m.conn.QueryRow(&resp, query, serviceId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGatewayServiceHttpRuleModel) Update(data GatewayServiceHttpRule) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceHttpRuleRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ServiceId, data.RuleType, data.Rule, data.NeedHttps, data.NeedStripUri, data.NeedWebsocket, data.UrlRewrite, data.HeaderTransfor, data.Id)
	return err
}

func (m *defaultGatewayServiceHttpRuleModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
