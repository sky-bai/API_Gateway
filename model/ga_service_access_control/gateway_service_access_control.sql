
CREATE TABLE `gateway_service_access_control` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `open_auth` tinyint NOT NULL DEFAULT '0' COMMENT '是否开启权限 1=开启',
  `black_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '黑名单ip',
  `white_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单ip',
  `white_host_name` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单主机',
  `clientip_flow_limit` int NOT NULL DEFAULT '0' COMMENT '客户端ip限流',
  `service_flow_limit` int NOT NULL DEFAULT '0' COMMENT '服务端限流',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=192 DEFAULT CHARSET=utf8 COMMENT='网关权限控制表';
