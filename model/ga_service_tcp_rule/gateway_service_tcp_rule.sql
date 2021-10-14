
CREATE TABLE `gateway_service_tcp_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口号',
  PRIMARY KEY (`id`)
)