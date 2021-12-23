package ga_service_info

import (
	"API_Gateway/model/ga_service_access_control"
	"API_Gateway/model/ga_service_http_rule"
	"API_Gateway/model/ga_service_load_balance"
	"API_Gateway/model/ga_service_tcp_rule"
	"API_Gateway/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
	"strings"
	"time"
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

		// FindDataLike 模糊查询服务信息
		FindDataLike(info string, pageSize, pageNum int) (interface{}, error)

		// FindOneByServiceName 根据服务名查找一条数据
		FindOneByServiceName(serviceName string) (int, error)

		// FindAll 查询获取所有服务信息
		FindAll(info string, pageSize, pageNum int) (interface{}, error)

		// FindAllTotal 直接获取所有服务信息不分页
		FindAllTotal() (interface{}, error)

		// InsertData 添加http服务
		InsertData(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, accessControl ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error
		// UpdateDate 更新http服务
		UpdateDate(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, ac ga_service_access_control.GatewayServiceAccessControl, lb ga_service_load_balance.GatewayServiceLoadBalance) error

		// InsertTcpService 添加tcp服务
		InsertTcpService(req GatewayServiceInfo, data ga_service_tcp_rule.GatewayServiceTcpRule, ac ga_service_access_control.GatewayServiceAccessControl, ld ga_service_load_balance.GatewayServiceLoadBalance) error
		// UpdateTcp 更新tcp服务
		UpdateTcp(req GatewayServiceInfo, data ga_service_tcp_rule.GatewayServiceTcpRule, ac ga_service_access_control.GatewayServiceAccessControl, ld ga_service_load_balance.GatewayServiceLoadBalance) error

		// GetServiceNum 获取服务数量
		GetServiceNum() (int, error)

		// GetAllNum 获取每个服务的数量
		GetAllNum() ([]ServiceNum, error)
	}

	defaultGatewayServiceInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GatewayServiceInfo struct {
		Id          int64     `db:"id" json:"id"`                      // 自增主键
		LoadType    int64     `db:"load_type" json:"load_type"`        // 负载类型 0=http 1=tcp 2=grpc
		ServiceName string    `db:"service_name" json:"service_name"`  // 服务名称 6-128 数字字母下划线
		ServiceDesc string    `db:"service_desc" json:"service_descs"` // 服务描述
		CreateTime  time.Time `db:"create_time" json:"create_time"`    // 添加时间
		UpdateTime  time.Time `db:"update_time" json:"update_time"`    // 更新时间
		IsDelete    int64     `db:"is_delete" json:"is_delete"`        // 是否删除 1=删除 0=未删除
	}

	ServiceInfo struct {
		Id          int64  `db:"id"`           // 自增主键
		LoadType    int64  `db:"load_type"`    // 负载类型 0=http 1=tcp 2=grpc
		ServiceName string `db:"service_name"` // 服务名称 6-128 数字字母下划线
	}

	ServiceNum struct {
		LoadType int   `db:"num" json:"num"`
		Value    int64 `db:"value" json:"value"`
		//Name     string `json:"name"`
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

// FindOneByServiceName 根据服务名查找一条数据
func (m *defaultGatewayServiceInfoModel) FindOneByServiceName(serviceName string) (int, error) {
	query := fmt.Sprintf("select id from %s where `service_name` = ? limit 1", m.table)
	var resp int
	err := m.conn.QueryRow(&resp, query, serviceName)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

// FindDataLike 从服务信息表中模糊查询服务信息(服务ID )
func (m *defaultGatewayServiceInfoModel) FindDataLike(info string, pageSize, pageNum int) (interface{}, error) {

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var countNum int

	if info == "" {
		countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE `is_delete` = 0", m.table)
		err := m.conn.QueryRow(&countNum, countQuery)
		startNum := (pageNum - 1) * pageSize
		query := fmt.Sprintf("select %s from %s where`is_delete` = 0 ORDER BY `id` DESC LIMIT ?,?", gatewayServiceInfoRows, m.table)
		var resp []GatewayServiceInfo
		err = m.conn.QueryRows(&resp, query, startNum, pageSize)
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

	countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE `service_name` like ? or `service_desc` like ? AND `is_delete` = 0", m.table)
	err := m.conn.QueryRow(&countNum, countQuery, "%"+info+"%", "%"+info+"%")
	startNum := (pageNum - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where `service_name` like ? or `service_desc` like ?  AND `is_delete` = 0 ORDER BY `id` DESC LIMIT ?,?", gatewayServiceInfoRows, m.table)
	var resp []GatewayServiceInfo
	err = m.conn.QueryRows(&resp, query, "%"+info+"%", "%"+info+"%", startNum, pageSize)
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
	_, err := m.conn.Exec(query, data.LoadType, data.ServiceName, data.ServiceDesc, data.IsDelete, data.Id)
	return err
}

// UpdateDate 更新http服务
func (m *defaultGatewayServiceInfoModel) UpdateDate(req GatewayServiceInfo, hr ga_service_http_rule.GatewayServiceHttpRule, ac ga_service_access_control.GatewayServiceAccessControl, lb ga_service_load_balance.GatewayServiceLoadBalance) error {

	UpdateServiceSql := fmt.Sprintf("update %s set service_name = '%s',service_desc = ? where `id` = ?", m.table, req.ServiceName)
	fmt.Println("UpdateServiceSql", UpdateServiceSql)

	err := m.conn.Transact(func(session sqlx.Session) error {

		// 1.更新serviceInfo表数据
		stmt, err := session.Prepare(UpdateServiceSql)
		if err != nil {
			fmt.Println("updateServiceInfo:err", err)
			return errors.New("更新http rule服务失败")
		}
		defer stmt.Close()
		_, err = stmt.Exec(req.ServiceDesc, req.Id)
		if err != nil {
			fmt.Println("UpdateServiceSql exec", err)
			return errors.New("更新http rule服务失败")
		}

		// 2.更新http规则表
		UpdateHttpSql := fmt.Sprintf("update gateway_service_http_rule set rule_type = ?,rule = ?,need_https = ?,need_strip_uri = ?,need_websocket = ?,url_rewrite = ?,header_transfor = ? where service_id = ? ")

		stmt1, err := session.Prepare(UpdateHttpSql)
		if err != nil {
			fmt.Println("UpdateTcpSql err:", err)
			return errors.New("更新http rule服务失败")
		}
		defer stmt1.Close()
		if _, err := stmt1.Exec(hr.RuleType, hr.Rule, hr.NeedHttps, hr.NeedStripUri, hr.NeedWebsocket, hr.UrlRewrite, hr.HeaderTransfor, req.Id); err != nil {
			fmt.Println("UpdateTcpSql err:", err)
			return errors.New("更新http rule服务失败")
		}

		// 3.写入权限控制表
		updateAccessControlSql := fmt.Sprintf("update gateway_service_access_control set open_auth = ?, black_list = ?, white_list = ?, white_host_name = ?, clientip_flow_limit = ?, service_flow_limit = ? where service_id = ?")
		fmt.Println("updateAccessControlSql :", updateAccessControlSql)

		stmt2, err := session.Prepare(updateAccessControlSql)
		if err != nil {
			fmt.Println("updateAccessControlSql err:", err)
			return errors.New("更新http rule服务失败")
		}
		defer stmt2.Close()
		if _, err := stmt2.Exec(ac.OpenAuth, ac.BlackList, ac.WhiteList, ac.WhiteHostName, ac.ClientipFlowLimit, ac.ServiceFlowLimit, req.Id); err != nil {
			fmt.Println("updateAccessControlSql err:", err)
			return errors.New("更新http rule服务失败")
		}

		// 4.写入负载均衡控制表
		updateLoadBalanceSql := fmt.Sprintf("update gateway_service_load_balance set check_method = ?,check_timeout = ?,check_interval = ?, ip_list = ?,weight_list = ?,forbid_list = ?,upstream_connect_timeout = ?,upstream_header_timeout = ?,upstream_idle_timeout = ?,upstream_max_idle = ? where service_id = ?")
		fmt.Println("updateLoadBalanceSql :", updateLoadBalanceSql)

		stmt3, err := session.Prepare(updateLoadBalanceSql)
		if err != nil {
			fmt.Println("updateLoadBalanceSql err:", err)
			return errors.New("更新http rule服务失败")
		}
		defer stmt3.Close()
		if _, err := stmt3.Exec(lb.CheckMethod, lb.CheckTimeout, lb.CheckInterval, lb.IpList, lb.WeightList, lb.ForbidList, lb.UpstreamConnectTimeout, lb.UpstreamHeaderTimeout, lb.UpstreamIdleTimeout, lb.UpstreamMaxIdle, req.Id); err != nil {
			fmt.Println("updateLoadBalanceSql err:", err)
			return errors.New("更新http rule服务失败")
		}

		return nil
	})

	return err
}

func (m *defaultGatewayServiceInfoModel) Delete(id int64) error {
	query := fmt.Sprintf("update %s set is_delete = 1 where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

// 现在我要对四张表进行数据操作 但是我可能就修改一个值
// 改动一个数据 然而对四张表进行了更新

// InsertData 添加一个http服务
func (m *defaultGatewayServiceInfoModel) InsertData(req GatewayServiceInfo, data ga_service_http_rule.GatewayServiceHttpRule, ac ga_service_access_control.GatewayServiceAccessControl, loadBalance ga_service_load_balance.GatewayServiceLoadBalance) error {

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
		insertAccessControlSql := fmt.Sprintf("insert into gateway_service_access_control (service_id,open_auth,black_list,white_list,white_host_name,clientip_flow_limit,service_flow_limit) values (%d,%d,'%s','%s','%s',%d,%d)", InsertId, ac.OpenAuth, ac.BlackList, ac.WhiteList, ac.WhiteHostName, ac.ClientipFlowLimit, ac.ServiceFlowLimit)
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

// InsertTcpService 添加一条tcp服务
func (m *defaultGatewayServiceInfoModel) InsertTcpService(req GatewayServiceInfo, data ga_service_tcp_rule.GatewayServiceTcpRule, ac ga_service_access_control.GatewayServiceAccessControl, ld ga_service_load_balance.GatewayServiceLoadBalance) error {
	insertServiceSql := fmt.Sprintf("insert into %s (%s) values (?,?,?,?)", m.table, gatewayServiceInfoRowsExpectAutoSet)
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
		sqlResult, err := stmt.Exec(1, req.ServiceName, req.ServiceDesc, 0)
		if err != nil {
			fmt.Println("insertServiceSql exec", err)
			return err
		}

		// 2.拿到刚刚插入服务表的那条记录ID 写入tcp规则表
		InsertId, _ := sqlResult.LastInsertId()
		insertRuleSql := fmt.Sprintf("insert into gateway_service_tcp_rule (service_id,port) values (?,?)")
		fmt.Println("insertTcpRuleSql :", insertRuleSql)
		stmt1, err := session.Prepare(insertRuleSql)
		if err != nil {
			fmt.Println("insertTcpRuleSql err:", err)
			return err
		}
		defer stmt1.Close()
		if _, err := stmt1.Exec(InsertId, data.Port); err != nil {
			fmt.Println("insertTcpRuleSql err:", err)
			return err
		}

		// 3.写入权限控制表
		insertAccessControlSql := fmt.Sprintf("insert into gateway_service_access_control (service_id,open_auth,black_list,white_list,white_host_name,clientip_flow_limit,service_flow_limit) values (?,?,?,?,?,?,?)")
		fmt.Println("insertAccessControlSql :", insertAccessControlSql)
		stmt2, err := session.Prepare(insertAccessControlSql)
		if err != nil {
			fmt.Println("insertAccessControlSql err:", err)
			return err
		}
		defer stmt2.Close()
		if _, err := stmt2.Exec(InsertId, ac.OpenAuth, ac.BlackList, ac.WhiteList, ac.WhiteHostName, ac.ClientipFlowLimit, ac.ServiceFlowLimit); err != nil {
			fmt.Println("insertAccessControlSql err:", err)
			return err
		}

		// 4.写入负载均衡控制表
		insertLoadBalanceSql := fmt.Sprintf("insert into gateway_service_load_balance (service_id,check_method,check_timeout,check_interval,round_type,ip_list,weight_list,forbid_list,upstream_connect_timeout,upstream_header_timeout,upstream_idle_timeout,upstream_max_idle) values (?,?,?,?,?,?,?,?,?,?,?,?)")
		fmt.Println("insertLoadBalanceSql :", insertLoadBalanceSql)
		stmt3, err := session.Prepare(insertLoadBalanceSql)
		if err != nil {
			fmt.Println("insertLoadBalanceSql err:", err)
			return err
		}
		defer stmt3.Close()
		if _, err := stmt3.Exec(InsertId, ld.CheckMethod, ld.CheckTimeout, ld.CheckInterval, ld.RoundType, ld.IpList, ld.WeightList, ld.ForbidList, ld.UpstreamConnectTimeout, ld.UpstreamHeaderTimeout, ld.UpstreamIdleTimeout, ld.UpstreamMaxIdle); err != nil {
			fmt.Println("insertLoadBalanceSql err:", err)
			return err
		}

		return nil
	})
	return err
}

// UpdateTcp 更新tcp服务
func (m *defaultGatewayServiceInfoModel) UpdateTcp(req GatewayServiceInfo, data ga_service_tcp_rule.GatewayServiceTcpRule, ac ga_service_access_control.GatewayServiceAccessControl, ld ga_service_load_balance.GatewayServiceLoadBalance) error {
	UpdateServiceSql := fmt.Sprintf("update %s set service_name = '%s',service_desc = '%s' where `id` = %d", m.table, req.ServiceName, req.ServiceDesc, req.Id)
	fmt.Println("UpdateServiceSql", UpdateServiceSql)

	err := m.conn.Transact(func(session sqlx.Session) error {

		// 1.更新serviceInfo表数据
		stmt, err := session.Prepare(UpdateServiceSql)
		if err != nil {
			fmt.Println("updetehttpsql:err", err)
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println("UpdateServiceSql exec", err)
			return err
		}

		// 2.更新tcp规则表
		UpdateTcpSql := fmt.Sprintf("update gateway_service_tcp_rule set port = %d where `service_id` = %d )", data.Port, req.Id)
		fmt.Println("UpdateTcpSql :", UpdateTcpSql)

		stmt1, err := session.Prepare(UpdateTcpSql)
		if err != nil {
			fmt.Println("UpdateTcpSql err:", err)
			return err
		}
		defer stmt1.Close()
		if _, err := stmt1.Exec(); err != nil {
			fmt.Println("UpdateTcpSql err:", err)
			return err
		}

		// 3.写入权限控制表
		updateAccessControlSql := fmt.Sprintf("update gateway_service_access_control set open_auth = %d black_list = '%s' white_list = '%s' white_host_name = '%s' clientip_flow_limit = %d service_flow_limit = %d where `service_id` = %d", ac.OpenAuth, ac.BlackList, ac.WhiteList, ac.WhiteHostName, ac.ClientipFlowLimit, ac.ServiceFlowLimit, req.Id)
		fmt.Println("updateAccessControlSql :", updateAccessControlSql)

		stmt2, err := session.Prepare(updateAccessControlSql)
		if err != nil {
			fmt.Println("updateAccessControlSql err:", err)
			return err
		}
		defer stmt2.Close()
		if _, err := stmt2.Exec(); err != nil {
			fmt.Println("updateAccessControlSql err:", err)
			return err
		}

		// 4.写入负载均衡控制表
		updateLoadBalanceSql := fmt.Sprintf("update gateway_service_load_balance set check_method = %d,check_timeout = %d,check_interval = %d,round_type = %d,ip_list = '%s',weight_list = '%s',forbid_list = '%s',upstream_connect_timeout = %d,upstream_header_timeout = %d,upstream_idle_timeout = %d,upstream_max_idle = %d where `service_id` = %d", ld.CheckMethod, ld.CheckTimeout, ld.CheckInterval, ld.WeightList, ld.IpList, ld.WeightList, ld.ForbidList, ld.UpstreamConnectTimeout, ld.UpstreamHeaderTimeout, ld.UpstreamIdleTimeout, ld.UpstreamMaxIdle, req.Id)
		fmt.Println("updateLoadBalanceSql :", updateLoadBalanceSql)

		stmt3, err := session.Prepare(updateLoadBalanceSql)
		if err != nil {
			fmt.Println("updateLoadBalanceSql err:", err)
			return err
		}
		defer stmt3.Close()
		if _, err := stmt3.Exec(); err != nil {
			fmt.Println("updateLoadBalanceSql err:", err)
			return err
		}

		return nil
	})

	return err
}

// FindAll 查询获取所有服务信息
func (m *defaultGatewayServiceInfoModel) FindAll(info string, pageSize, pageNum int) (interface{}, error) {
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var countNum int
	startNum := (pageNum - 1) * pageSize
	var resp []GatewayServiceInfo
	if info == "" {
		countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE `is_delete` = 0", m.table)
		err := m.conn.QueryRow(&countNum, countQuery)

		query := fmt.Sprintf("select %s from %s where `is_delete` = 0 ORDER BY `id` DESC LIMIT %d,%d", gatewayServiceInfoRows, m.table, startNum, pageSize)
		err = m.conn.QueryRows(&resp, query)
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

	countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE `service_name` like ? or `service_desc` like ? AND `is_delete` = 0", m.table)
	err := m.conn.QueryRow(&countNum, countQuery, "%"+info+"%", "%"+info+"%")

	query := fmt.Sprintf("select %s from %s where `service_name` like ? or `service_desc` like ?  AND `is_delete` = 0 ORDER BY `id` DESC LIMIT %d,%d", gatewayServiceInfoRows, m.table, startNum, pageSize)
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

// FindAllTotal 直接获取所有服务信息不分页
func (m *defaultGatewayServiceInfoModel) FindAllTotal() (interface{}, error) {
	query := fmt.Sprintf("select * from %s", m.table)
	var resp []GatewayServiceInfo
	err := m.conn.QueryRows(&resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetServiceNum 获取服务数量
func (m *defaultGatewayServiceInfoModel) GetServiceNum() (int, error) {
	query := fmt.Sprintf("select count(*) from %s where `is_delete` = 0", m.table)
	var countNum int
	err := m.conn.QueryRow(&countNum, query)
	switch err {
	case nil:
		return countNum, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

// GetAllNum 获取每个服务的数量
func (m *defaultGatewayServiceInfoModel) GetAllNum() ([]ServiceNum, error) {

	query := fmt.Sprintf("SELECT load_type as num,count(*) as value FROM gateway_service_info WHERE  is_delete = 0 GROUP BY load_type")

	var resp []ServiceNum
	err := m.conn.QueryRows(&resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
