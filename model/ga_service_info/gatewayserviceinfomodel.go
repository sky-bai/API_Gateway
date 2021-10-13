package ga_service_info

import (
	"API_Gateway/util"
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
	gatewayServiceInfoFieldNames          = builderx.RawFieldNames(&GatewayServiceInfo{})
	gatewayServiceInfoRows                = strings.Join(gatewayServiceInfoFieldNames, ",")
	gatewayServiceInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceInfoModel interface {
		Insert(data GatewayServiceInfo) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceInfo, error)
		Update(data GatewayServiceInfo) error
		Delete(id int64) error

		// 模糊查询服务信息
		FindDataLike(info string, pageSize, pageNum int) (interface{}, error)
	}

	defaultGatewayServiceInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceInfo struct {
		Id          int64     `db:"id"`           // 自增主键
		LoadType    int64     `db:"load_type"`    // 负载类型 0=http 1=tcp 2=grpc
		ServiceName string    `db:"service_name"` // 服务名称 6-128 数字字母下划线
		ServiceDesc string    `db:"service_desc"` // 服务描述
		CreateAt    time.Time `db:"create_at"`    // 添加时间
		UpdateAt    time.Time `db:"update_at"`    // 更新时间
		IsDelete    int64     `db:"is_delete"`    // 是否删除 1=删除
	}
)

func NewGatewayServiceInfoModel(conn sqlx.SqlConn) GatewayServiceInfoModel {
	return &defaultGatewayServiceInfoModel{
		conn:  conn,
		table: "`gateway_service_info`",
	}
}

func (m *defaultGatewayServiceInfoModel) Insert(data GatewayServiceInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, gatewayServiceInfoRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.LoadType, data.ServiceName, data.ServiceDesc, data.CreateAt, data.UpdateAt, data.IsDelete)
	return ret, err
}

func (m *defaultGatewayServiceInfoModel) FindOne(id int64) (*GatewayServiceInfo, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceInfoRows, m.table)
	var resp GatewayServiceInfo
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

// 模糊查询服务信息
func (m *defaultGatewayServiceInfoModel) FindDataLike(info string, pageSize, pageNum int) (interface{}, error) {

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	var countNum int
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE `service_name` like ? or `service_desc` like ? AND `is_delete` = 0", m.table)
	err := m.conn.QueryRow(&countNum, countQuery, "%"+info+"%", "%"+info+"%")
	startNum := (pageNum - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where `service_name` like ? or `service_desc` like ?  AND `is_delete` = 0 ORDER BY `id` DESC LIMIT %d,%d", gatewayServiceInfoRows, m.table, startNum, pageSize)
	var resp []GatewayServiceInfo
	err = m.conn.QueryRows(&resp, query, "%"+info+"%", "%"+info+"%")
	switch err {
	case nil:
		res := util.CutPage(countNum, pageNum, pageSize, resp)
		return &res, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGatewayServiceInfoModel) Update(data GatewayServiceInfo) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceInfoRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.LoadType, data.ServiceName, data.ServiceDesc, data.CreateAt, data.UpdateAt, data.IsDelete, data.Id)
	return err
}

func (m *defaultGatewayServiceInfoModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
