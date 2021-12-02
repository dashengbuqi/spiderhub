/*
Navicat MySQL Data Transfer

Source Server         : 开发[新]
Source Server Version : 50722
Source Host           : 192.168.1.40:3306
Source Database       : spiderhub

Target Server Type    : MYSQL
Target Server Version : 50722
File Encoding         : 65001

Date: 2021-12-02 16:50:43
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for collect_app_fields
-- ----------------------------
DROP TABLE IF EXISTS `collect_app_fields`;
CREATE TABLE `collect_app_fields` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `target` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '目标类型（1爬虫2清洗）',
  `target_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '目标id',
  `content` varchar(512) NOT NULL DEFAULT '' COMMENT '字段内容',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
