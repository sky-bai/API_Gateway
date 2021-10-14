package ga_service_load_balance

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
	gatewayServiceLoadBalanceFieldNames          = builderx.RawFieldNames(&GatewayServiceLoadBalance{})
	gatewayServiceLoadBalanceRows                = strings.Join(gatewayServiceLoadBalanceFieldNames, ",")
	gatewayServiceLoadBalanceRowsExpectAutoSet   = strings.Join(stringx.Remove(gatewayServiceLoadBalanceFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	gatewayServiceLoadBalanceRowsWithPlaceHolder = strings.Join(stringx.Remove(gatewayServiceLoadBalanceFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	GatewayServiceLoadBalanceModel interface {
		Insert(data GatewayServiceLoadBalance) (sql.Result, error)
		FindOne(id int64) (*GatewayServiceLoadBalance, error)
		Update(data GatewayServiceLoadBalance) error
		Delete(id int64) error
	}

	defaultGatewayServiceLoadBalanceModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceLoadBalance struct {
		Id                     int64  `db:"id"`                       // 自增主键
		ServiceId              int64  `db:"service_id"`               // 服务id
		CheckMethod            int64  `db:"check_method"`             // 检查方法 0=tcpchk,检测端口是否握手成功
		CheckTimeout           int64  `db:"check_timeout"`            // check超时时间,单位s
		CheckInterval          int64  `db:"check_interval"`           // 检查间隔, 单位s
		RoundType              int64  `db:"round_type"`               // 轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash
		IpList                 string `db:"ip_list"`                  // ip列表
		WeightList             string `db:"weight_list"`              // 权重列表
		ForbidList             string `db:"forbid_list"`              // 禁用ip列表
		UpstreamConnectTimeout int64  `db:"upstream_connect_timeout"` // 建立连接超时, 单位s
		UpstreamHeaderTimeout  int64  `db:"upstream_header_timeout"`  // 获取header超时, 单位s
		UpstreamIdleTimeout    int64  `db:"upstream_idle_timeout"`    // 链接最大空闲时间, 单位s
		UpstreamMaxIdle        int64  `db:"upstream_max_idle"`        // 最大空闲链接数
	}
)

func NewGatewayServiceLoadBalanceModel(conn sqlx.SqlConn) GatewayServiceLoadBalanceModel {
	return &defaultGatewayServiceLoadBalanceModel{
		conn:  conn,
		table: "`gateway_service_load_balance`",
	}
}

func (m *defaultGatewayServiceLoadBalanceModel) Insert(data GatewayServiceLoadBalance) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, gatewayServiceLoadBalanceRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ServiceId, data.CheckMethod, data.CheckTimeout, data.CheckInterval, data.RoundType, data.IpList, data.WeightList, data.ForbidList, data.UpstreamConnectTimeout, data.UpstreamHeaderTimeout, data.UpstreamIdleTimeout, data.UpstreamMaxIdle)
	return ret, err
}

func (m *defaultGatewayServiceLoadBalanceModel) FindOne(id int64) (*GatewayServiceLoadBalance, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gatewayServiceLoadBalanceRows, m.table)
	var resp GatewayServiceLoadBalance
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

func (m *defaultGatewayServiceLoadBalanceModel) Update(data GatewayServiceLoadBalance) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceLoadBalanceRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ServiceId, data.CheckMethod, data.CheckTimeout, data.CheckInterval, data.RoundType, data.IpList, data.WeightList, data.ForbidList, data.UpstreamConnectTimeout, data.UpstreamHeaderTimeout, data.UpstreamIdleTimeout, data.UpstreamMaxIdle, data.Id)
	return err
}

func (m *defaultGatewayServiceLoadBalanceModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
