package ga_admin

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	gatewayAdminFieldNames          = builderx.RawFieldNames(&GatewayAdmin{})
	gatewayAdminRows                = strings.Join(gatewayAdminFieldNames, ",")
	gatewayAdminRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayAdminFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayAdminRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayAdminFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	gatewayAdminLogin = "`id`" + "," + "`password`"
	updatePwd         = "`password`"
)

type (
	GatewayAdminModel interface {
		Insert(data GatewayAdmin) (sql.Result, error)
		FindOne(id int64) (*GatewayAdmin, error)
		Update(data GatewayAdmin) error
		Delete(id int64) error

		// 通过用户名查找密码
		FindOneByUserName(userName string) (*GatewayAdmin, error)
		// 修改密码
		UpdatePwd(data GatewayAdmin) error
	}

	defaultGatewayAdminModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayAdmin struct {
		Id       int64     `db:"id"`        // 自增id
		UserName string    `db:"user_name"` // 用户名
		Salt     string    `db:"salt"`      // 盐
		Password string    `db:"password"`  // 密码
		CreateAt time.Time `db:"create_at"` // 新增时间
		UpdateAt time.Time `db:"update_at"` // 更新时间
		IsDelete int64     `db:"is_delete"` // 是否删除 0未删除 1已删除
	}
)

func NewGatewayAdminModel(conn sqlx.SqlConn) GatewayAdminModel {
	return &defaultGatewayAdminModel{
		conn:  conn,
		table: "`gateway_admin`",
	}
}

func (m *defaultGatewayAdminModel) Insert(data GatewayAdmin) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, gatewayAdminRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.UserName, data.Salt, data.Password, data.CreateAt, data.UpdateAt, data.IsDelete)
	return ret, err
}

func (m *defaultGatewayAdminModel) FindOne(id int64) (*GatewayAdmin, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayAdminRows, m.table)
	var resp GatewayAdmin
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

// 通过用户名查找密码
func (m *defaultGatewayAdminModel) FindOneByUserName(userName string) (*GatewayAdmin, error) {
	query := fmt.Sprintf("select * from %s where `user_name` = ? limit 1", m.table)
	var resp GatewayAdmin
	err := m.conn.QueryRow(&resp, query, userName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGatewayAdminModel) Update(data GatewayAdmin) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayAdminRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.UserName, data.Salt, data.Password, data.CreateAt, data.UpdateAt, data.IsDelete, data.Id)
	return err
}

// 修改密码
func (m *defaultGatewayAdminModel) UpdatePwd(data GatewayAdmin) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, updatePwd)
	_, err := m.conn.Exec(query, data.Id)
	return err
}

func (m *defaultGatewayAdminModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
