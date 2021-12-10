package ga_service_access_control

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
	gatewayServiceAccessControlFieldNames          = builderx.RawFieldNames(&GatewayServiceAccessControl{})
	gatewayServiceAccessControlRows                = strings.Join(gatewayServiceAccessControlFieldNames, ",")
	gatewayServiceAccessControlRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceAccessControlFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceAccessControlRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceAccessControlFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceAccessControlModel interface {
		Insert(data GatewayServiceAccessControl) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceAccessControl, error)
		Update(data GatewayServiceAccessControl) error
		Delete(id int64) error

		// 根据服务ID去查
		FindOneByServiceId(serviceId int64) (*GatewayServiceAccessControl, error)
	}

	defaultGatewayServiceAccessControlModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceAccessControl struct {
		Id                int64  `db:"id" json:"id"`                                   // 自增主键
		ServiceId         int64  `db:"service_id" json:"service_id"`                   // 服务id
		OpenAuth          int64  `db:"open_auth" json:"open_auth"`                     // 是否开启权限 1=开启
		BlackList         string `db:"black_list" json:"black_list"`                   // 黑名单ip
		WhiteList         string `db:"white_list" json:"white_list"`                   // 白名单ip
		WhiteHostName     string `db:"white_host_name" json:"white_host_name"`         // 白名单主机
		ClientipFlowLimit int64  `db:"clientip_flow_limit" json:"clientip_flow_limit"` // 客户端ip限流
		ServiceFlowLimit  int64  `db:"service_flow_limit" json:"service_flow_limit"`   // 服务端限流
	}
)

func NewGatewayServiceAccessControlModel(conn sqlx.SqlConn) GatewayServiceAccessControlModel {
	return &defaultGatewayServiceAccessControlModel{
		conn:  conn,
		table: "`gateway_service_access_control`",
	}
}

func (m *defaultGatewayServiceAccessControlModel) Insert(data GatewayServiceAccessControl) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, gatewayServiceAccessControlRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ServiceId, data.OpenAuth, data.BlackList, data.WhiteList, data.WhiteHostName, data.ClientipFlowLimit, data.ServiceFlowLimit)
	return ret, err
}
func (m *defaultGatewayServiceAccessControlModel) FindOne(id int64) (*GatewayServiceAccessControl, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceAccessControlRows, m.table)
	var resp GatewayServiceAccessControl
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
func (m *defaultGatewayServiceAccessControlModel) FindOneByServiceId(serviceId int64) (*GatewayServiceAccessControl, error) {
	query := fmt.Sprintf("select %s from %s where `service_id` = ? limit 1", gatewayServiceAccessControlRows, m.table)
	var resp GatewayServiceAccessControl
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
func (m *defaultGatewayServiceAccessControlModel) Update(data GatewayServiceAccessControl) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceAccessControlRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ServiceId, data.OpenAuth, data.BlackList, data.WhiteList, data.WhiteHostName, data.ClientipFlowLimit, data.ServiceFlowLimit, data.Id)
	return err
}
func (m *defaultGatewayServiceAccessControlModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
