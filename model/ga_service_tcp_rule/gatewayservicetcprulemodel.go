package ga_service_tcp_rule

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
	gatewayServiceTcpRuleFieldNames          = builderx.RawFieldNames(&GatewayServiceTcpRule{})
	gatewayServiceTcpRuleRows                = strings.Join(gatewayServiceTcpRuleFieldNames, ",")
	gatewayServiceTcpRuleRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceTcpRuleFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceTcpRuleRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceTcpRuleFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceTcpRuleModel interface {
		Insert(data GatewayServiceTcpRule) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceTcpRule, error)
		Update(data GatewayServiceTcpRule) error
		Delete(id int64) error

		// 根据服务ID查出该服务的tcp信息
		FindOneByServiceId(serviceId int) (*GatewayServiceTcpRule, error)
	}

	defaultGatewayServiceTcpRuleModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceTcpRule struct {
		Id        int64 `db:"id"`         // 自增主键
		ServiceId int64 `db:"service_id"` // 服务id
		Port      int64 `db:"port"`       // 端口号
	}
)

func NewGatewayServiceTcpRuleModel(conn sqlx.SqlConn) GatewayServiceTcpRuleModel {
	return &defaultGatewayServiceTcpRuleModel{
		conn:  conn,
		table: "`gateway_service_tcp_rule`",
	}
}

func (m *defaultGatewayServiceTcpRuleModel) Insert(data GatewayServiceTcpRule) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, gatewayServiceTcpRuleRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ServiceId, data.Port)
	return ret, err
}

func (m *defaultGatewayServiceTcpRuleModel) FindOne(id int64) (*GatewayServiceTcpRule, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceTcpRuleRows, m.table)
	var resp GatewayServiceTcpRule
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

// 根据服务ID查出该服务的tcp信息
func (m *defaultGatewayServiceTcpRuleModel) FindOneByServiceId(serviceId int) (*GatewayServiceTcpRule, error) {
	query := fmt.Sprintf("select %s from %s where `service_id` = ? limit 1", gatewayServiceTcpRuleRows, m.table)
	var resp GatewayServiceTcpRule
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

func (m *defaultGatewayServiceTcpRuleModel) Update(data GatewayServiceTcpRule) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceTcpRuleRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ServiceId, data.Port, data.Id)

	return err
}

func (m *defaultGatewayServiceTcpRuleModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
