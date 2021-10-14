
CREATE TABLE `gateway_service_load_balance` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `check_method` tinyint NOT NULL DEFAULT '0' COMMENT '检查方法 0=tcpchk,检测端口是否握手成功',
  `check_timeout` int NOT NULL DEFAULT '0' COMMENT 'check超时时间,单位s',
  `check_interval` int NOT NULL DEFAULT '0' COMMENT '检查间隔, 单位s',
  `round_type` tinyint NOT NULL DEFAULT '2' COMMENT '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash',
  `ip_list` varchar(2000) NOT NULL DEFAULT '' COMMENT 'ip列表',
  `weight_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '权重列表',
  `forbid_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '禁用ip列表',
  `upstream_connect_timeout` int NOT NULL DEFAULT '0' COMMENT '建立连接超时, 单位s',
  `upstream_header_timeout` int NOT NULL DEFAULT '0' COMMENT '获取header超时, 单位s',
  `upstream_idle_timeout` int NOT NULL DEFAULT '0' COMMENT '链接最大空闲时间, 单位s',
  `upstream_max_idle` int NOT NULL DEFAULT '0' COMMENT '最大空闲链接数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=192 DEFAULT CHARSET=utf8 COMMENT='网关负载表\n\n\n每个服务的负载均衡表';
