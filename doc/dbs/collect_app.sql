/*
Navicat MySQL Data Transfer

Source Server         : 开发[新]
Source Server Version : 50722
Source Host           : 192.168.1.40:3306
Source Database       : spiderhub

Target Server Type    : MYSQL
Target Server Version : 50722
File Encoding         : 65001

Date: 2021-12-02 16:50:33
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for collect_app
-- ----------------------------
DROP TABLE IF EXISTS `collect_app`;
CREATE TABLE `collect_app` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '采集名称',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户',
  `crawler_token` varchar(255) NOT NULL DEFAULT '' COMMENT '爬虫令牌',
  `clean_token` varchar(255) NOT NULL DEFAULT '' COMMENT '清洗令牌',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态(0完成1执行中)',
  `schedule` varchar(255) NOT NULL DEFAULT '' COMMENT '计划任务',
  `storage` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '存储附件(1服务器2云盘)',
  `method` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '抓取方式(1重新抓取2更新3追加)',
  `error_info` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `crawler_content` text NOT NULL COMMENT '爬虫规则',
  `clean_content` text NOT NULL COMMENT '清洗规则',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
