package ga_service_grpc_rule

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
	gatewayServiceGrpcRuleFieldNames          = builderx.RawFieldNames(&GatewayServiceGrpcRule{})
	gatewayServiceGrpcRuleRows                = strings.Join(gatewayServiceGrpcRuleFieldNames, ",")
	gatewayServiceGrpcRuleRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceGrpcRuleFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceGrpcRuleRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceGrpcRuleFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceGrpcRuleModel interface {
		Insert(data GatewayServiceGrpcRule) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceGrpcRule, error)
		Update(data GatewayServiceGrpcRule) error
		Delete(id int64) error

		// 根据服务表中的服务ID进行查找该服务的grpc规则
		FindOneByServiceId(serviceId int) (*GatewayServiceGrpcRule, error)
	}

	defaultGatewayServiceGrpcRuleModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceGrpcRule struct {
		Id             int64  `db:"id"`              // 自增主键
		ServiceId      int64  `db:"service_id"`      // 服务id
		Port           int64  `db:"port"`            // 端口
		HeaderTransfor string `db:"header_transfor"` // header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔
	}
)

func NewGatewayServiceGrpcRuleModel(conn sqlx.SqlConn) GatewayServiceGrpcRuleModel {
	return &defaultGatewayServiceGrpcRuleModel{
		conn:  conn,
		table: "`gateway_service_grpc_rule`",
	}
}

func (m *defaultGatewayServiceGrpcRuleModel) Insert(data GatewayServiceGrpcRule) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, gatewayServiceGrpcRuleRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ServiceId, data.Port, data.HeaderTransfor)
	return ret, err
}

func (m *defaultGatewayServiceGrpcRuleModel) FindOne(id int64) (*GatewayServiceGrpcRule, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceGrpcRuleRows, m.table)
	var resp GatewayServiceGrpcRule
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

// 根据服务表中的服务ID进行查找该服务的grpc规则
func (m *defaultGatewayServiceGrpcRuleModel) FindOneByServiceId(serviceId int) (*GatewayServiceGrpcRule, error) {
	query := fmt.Sprintf("select %s from %s where `service_id` = ? limit 1", gatewayServiceGrpcRuleRows, m.table)
	var resp GatewayServiceGrpcRule
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

func (m *defaultGatewayServiceGrpcRuleModel) Update(data GatewayServiceGrpcRule) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceGrpcRuleRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ServiceId, data.Port, data.HeaderTransfor, data.Id)
	return err
}

func (m *defaultGatewayServiceGrpcRuleModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
