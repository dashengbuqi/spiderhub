/*
Navicat MySQL Data Transfer

Source Server         : 开发[新]
Source Server Version : 50722
Source Host           : 192.168.1.40:3306
Source Database       : spiderhub

Target Server Type    : MYSQL
Target Server Version : 50722
File Encoding         : 65001

Date: 2021-12-02 16:50:57
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for system_admin
-- ----------------------------
DROP TABLE IF EXISTS `system_admin`;
CREATE TABLE `system_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) DEFAULT '',
  `mobile` varchar(32) DEFAULT '',
  `auth_key` varchar(64) DEFAULT '',
  `password_hash` varchar(255) DEFAULT '',
  `email` varchar(32) DEFAULT '',
  `status` tinyint(1) DEFAULT '1',
  `last_login_time` int(10) unsigned DEFAULT '0',
  `login_times` int(10) unsigned DEFAULT '0' COMMENT '登录次数',
  `updated_at` int(10) unsigned DEFAULT '0',
  `created_at` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of system_admin
-- ----------------------------
INSERT INTO `system_admin` VALUES ('1', 'test', '13888888888', 'RROHOU', 'f4d4ecec6473b2f94ff091fddbee6d1c', 'test@spiderhub.com', '1', '0', '0', '0', '1634108956');
