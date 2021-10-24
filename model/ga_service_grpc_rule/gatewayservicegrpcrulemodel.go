package ga_service_grpc_rule

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_info"
	"API_Gateway/model/ga_service_load_balance"
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

		// 根据服务ID查找是否有服务ID
		FindIdByServiceId(serviceId int) (id int, err error)

		// 根据端口号查询是否存在该服务
		FindIdByPort(port int) (int, error)

		// 添加grpc服务
		InsertGrpcData(req ga_service_info.GatewayServiceInfo, data GatewayServiceGrpcRule, ac ga_service_access_control.GatewayServiceAccessControl, lb ga_service_load_balance.GatewayServiceLoadBalance) error
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

// 根据服务ID查找是否有服务ID
func (m *defaultGatewayServiceGrpcRuleModel) FindIdByServiceId(serviceId int) (id int, err error) {
	query := fmt.Sprintf("select id from %s where `service_id` = ? limit 1", m.table)
	var resp int
	err = m.conn.QueryRow(&resp, query, serviceId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

// 根据端口号查询是否存在该服务
func (m *defaultGatewayServiceGrpcRuleModel) FindIdByPort(port int) (int, error) {
	query := fmt.Sprintf("select id from %s where `port` = ? limit 1", m.table)
	var resp int
	err := m.conn.QueryRow(&resp, query, port)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
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

// 添加grpc服务
func (m *defaultGatewayServiceGrpcRuleModel) InsertGrpcData(req ga_service_info.GatewayServiceInfo, data GatewayServiceGrpcRule, ac ga_service_access_control.GatewayServiceAccessControl, lb ga_service_load_balance.GatewayServiceLoadBalance) error {
	insertServiceSql := fmt.Sprintf("insert into %s (load_type,service_name,service_desc,is_delete) values (%d,'%s','%s',%d)", m.table, 2, req.ServiceName, req.ServiceDesc, 0)
	fmt.Println("insertServiceSql", insertServiceSql)

	err := m.conn.Transact(func(session sqlx.Session) error {

		// 1.在serviceInfo表里面添加一条数据
		stmt, err := session.Prepare(insertServiceSql)
		if err != nil {
			fmt.Println("insertServiceSql prepare", insertServiceSql)
			fmt.Println(err)
			return err
		}
		defer stmt.Close()
		sqlResult, err := stmt.Exec()
		if err != nil {
			fmt.Println("insertServiceSql exec", err)
			return err
		}

		// 2.拿到刚刚插入服务表的那条记录ID 写入grpc规则表
		InsertId, _ := sqlResult.LastInsertId()
		insertRuleSql := fmt.Sprintf("insert into gateway_service_grpc_rule (service_id,port,header_transfor,need_https,need_strip_uri,need_websocket,url_rewrite,header_transfor) values (%d,%d,'%s')", InsertId, data.Port, data.HeaderTransfor)
		fmt.Println("insertRuleSql :", insertRuleSql)

		stmt1, err := session.Prepare(insertRuleSql)
		if err != nil {
			fmt.Println("InsertFrozendSql err:", err)
			return err
		}
		defer stmt1.Close()
		if _, err := stmt1.Exec(); err != nil {
			fmt.Println("insertRuleSql err:", err)
			return err
		}

		// 3.写入权限控制表
		insertAccessControlSql := fmt.Sprintf("insert into gateway_service_access_control (service_id,service_flow_limit,clientip_flow_limit,open_auth,black_list,need_websocket,url_rewrite,header_transfor) values (%d,%d,%d,%d,'%s','%s')", InsertId, ac.ServiceFlowLimit, ac.ClientipFlowLimit, ac.OpenAuth, ac.BlackList, ac.WhiteList)
		fmt.Println("insertAccessControlSql :", insertAccessControlSql)

		stmt2, err := session.Prepare(insertAccessControlSql)
		if err != nil {
			fmt.Println("insertAccessControlSql err:", err)
			return err
		}
		defer stmt2.Close()
		if _, err := stmt2.Exec(); err != nil {
			fmt.Println("insertAccessControlSql err:", err)
			return err
		}

		// 4.写入负载均衡控制表
		insertLoadBalanceSql := fmt.Sprintf("insert into gateway_service_load_balance (service_id,round_type,ip_list,weight_list,upstream_connect_timeout,upstream_header_timeout,upstream_idle_timeout,upstream_max_idle) values (%d,%d,'%s','%s',%d,%d,%d,%d)", InsertId, lb.RoundType, lb.IpList, lb.WeightList, lb.UpstreamConnectTimeout, lb.UpstreamHeaderTimeout, lb.UpstreamIdleTimeout, lb.UpstreamMaxIdle)
		fmt.Println("insertLoadBalanceSql :", insertLoadBalanceSql)

		stmt3, err := session.Prepare(insertLoadBalanceSql)
		if err != nil {
			fmt.Println("insertLoadBalanceSql err:", err)
			return err
		}
		defer stmt3.Close()
		if _, err := stmt3.Exec(); err != nil {
			fmt.Println("insertLoadBalanceSql err:", err)
			return err
		}

		return nil
	})

	return err
}
