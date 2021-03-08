DROP TABLE IF EXISTS `sample`;
SET character_set_client = utf8mb4 ;
CREATE TABLE `sample` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(30) NOT NULL COMMENT 'コード',
  `name` varchar(50) NOT NULL COMMENT '名前',
  `created_at` datetime DEFAULT NULL COMMENT '作成日時',
  `updated_at` datetime DEFAULT NULL COMMENT '更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='サンプル';