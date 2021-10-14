
CREATE TABLE `gateway_service_grpc_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口',
  `header_transfor` varchar(5000) NOT NULL DEFAULT '' COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
)