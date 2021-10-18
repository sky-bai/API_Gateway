package ga_service_info

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_load_balance"
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

		// 添加一个服务
		InsertData(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, accessControl ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error
		UpdateDate(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, accessControl ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error
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
		CreateTime  time.Time `db:"create_time"`  // 添加时间
		UpdateTime  time.Time `db:"update_time"`  // 更新时间
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
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, gatewayServiceInfoRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.LoadType, data.ServiceName, data.ServiceDesc, data.IsDelete)
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

func (m *defaultGatewayServiceInfoModel) Update(data GatewayServiceInfo) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceInfoRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.LoadType, data.ServiceName, data.ServiceDesc, data.IsDelete, data.Id)
	return err
}

// 更新http服务
func (m *defaultGatewayServiceInfoModel) UpdateDate(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, accessControl ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gatewayServiceInfoRowsWithPlaceHolder)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		//todo

		return err
	}
	return err
}
func (m *defaultGatewayServiceInfoModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

// 添加一个服务
func (m *defaultGatewayServiceInfoModel) InsertData(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, accessControl ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error {

	insertServiceSql := fmt.Sprintf("insert into %s (%s) values (%d,'%s','%s',%d)", m.table, gatewayServiceInfoRowsExpectAutoSet, 0, req.ServiceName, req.ServiceDesc, 0)
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

		// 2.拿到刚刚插入服务表的那条记录ID 写入http规则表
		InsertId, _ := sqlResult.LastInsertId()
		insertRuleSql := fmt.Sprintf("insert into gateway_service_http_rule (service_id,rule_type,rule,need_https,need_strip_uri,need_websocket,url_rewrite,header_transfor) values (%d,%d,'%s',%d,%d,%d,'%s','%s')", InsertId, data.RuleType, data.Rule, data.NeedHttps, data.NeedStripUri, data.NeedWebsocket, data.UrlRewrite, data.HeaderTransfor)
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
		insertAccessControlSql := fmt.Sprintf("insert into gateway_service_access_control (service_id,service_flow_limit,clientip_flow_limit,open_auth,black_list,need_websocket,url_rewrite,header_transfor) values (%d,%d,%d,%d,'%s','%s')", InsertId, accessControl.ServiceFlowLimit, accessControl.ClientipFlowLimit, accessControl.OpenAuth, accessControl.BlackList, accessControl.WhiteList)
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
		insertLoadBalanceSql := fmt.Sprintf("insert into gateway_service_load_balance (service_id,round_type,ip_list,weight_list,upstream_connect_timeout,upstream_header_timeout,upstream_idle_timeout,upstream_max_idle) values (%d,%d,'%s','%s',%d,%d,%d,%d)", InsertId, loadBalance.RoundType, loadBalance.IpList, loadBalance.WeightList, loadBalance.UpstreamConnectTimeout, loadBalance.UpstreamHeaderTimeout, loadBalance.UpstreamIdleTimeout, loadBalance.UpstreamMaxIdle)
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

// 从服务信息表中模糊查询服务信息(服务ID )
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
