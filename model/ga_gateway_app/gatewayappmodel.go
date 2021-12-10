package ga_gateway_app

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	gatewayAppFieldNames          = builderx.RawFieldNames(&GatewayApp{})
	gatewayAppRows                = strings.Join(gatewayAppFieldNames, ",")
	gatewayAppRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayAppFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayAppRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayAppFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayAppModel interface {
		Insert(data GatewayApp) (sql.Result, error)
		FindOne(id int64) (*GatewayApp, error)
		Update(data GatewayApp) error
		Delete(id int64) error

		// FindOneByAppId 根据appId查询租户信息
		FindOneByAppId(appId string) (*GatewayApp, error)
	}

	defaultGatewayAppModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayApp struct {
		Id         int64     `db:"id"`          // 自增id
		AppId      string    `db:"app_id"`      // 租户id
		Name       string    `db:"name"`        // 租户名称
		Secret     string    `db:"secret"`      // 密钥
		WhiteIps   string    `db:"white_ips"`   // ip白名单，支持前缀匹配
		Qpd        int64     `db:"qpd"`         // 日请求量限制
		Qps        int64     `db:"qps"`         // 每秒请求量限制
		CreateTime time.Time `db:"create_time"` // 添加时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
		IsDelete   int64     `db:"is_delete"`   // 是否删除 1=删除
	}
)

func NewGatewayAppModel(conn sqlx.SqlConn) GatewayAppModel {
	return &defaultGatewayAppModel{
		conn:  conn,
		table: "`gateway_app`",
	}
}

func (m *defaultGatewayAppModel) Insert(data GatewayApp) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, gatewayAppRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.AppId, data.Name, data.Secret, data.WhiteIps, data.Qpd, data.Qps, data.IsDelete)
	return ret, err
}

func (m *defaultGatewayAppModel) FindOne(id int64) (*GatewayApp, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayAppRows, m.table)
	var resp GatewayApp
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

func (m *defaultGatewayAppModel) Update(data GatewayApp) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayAppRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.AppId, data.Name, data.Secret, data.WhiteIps, data.Qpd, data.Qps, data.IsDelete, data.Id)
	return err
}

func (m *defaultGatewayAppModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

// FindOneByAppId 根据appId查询租户信息
func (m *defaultGatewayAppModel) FindOneByAppId(appId string) (*GatewayApp, error) {
	query := fmt.Sprintf("select %s from %s where `app_id` = ? limit 1", gatewayAppRows, m.table)
	var resp GatewayApp
	err := m.conn.QueryRow(&resp, query, appId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, errors.New("根据appId查询租户信息查询失败")
	}
}
