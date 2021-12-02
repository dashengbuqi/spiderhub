/*
Navicat MySQL Data Transfer

Source Server         : 开发[新]
Source Server Version : 50722
Source Host           : 192.168.1.40:3306
Source Database       : spiderhub

Target Server Type    : MYSQL
Target Server Version : 50722
File Encoding         : 65001

Date: 2021-12-02 16:51:07
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for system_menu
-- ----------------------------
DROP TABLE IF EXISTS `system_menu`;
CREATE TABLE `system_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `task_name` varchar(32) NOT NULL DEFAULT '',
  `full_name` varchar(64) NOT NULL DEFAULT '',
  `path` varchar(255) NOT NULL DEFAULT '',
  `parent_id` int(10) unsigned NOT NULL DEFAULT '0',
  `type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '类型(1:目录 2:栏目 3:菜单 4:按钮)',
  `icon` varchar(16) NOT NULL DEFAULT '',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `sort` smallint(2) unsigned NOT NULL DEFAULT '0',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=593 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of system_menu
-- ----------------------------
INSERT INTO `system_menu` VALUES ('4', '系统', '系统设置', '#', '0', '1', 'm3', '1', '1', '1634111032', '1514194107');
INSERT INTO `system_menu` VALUES ('5', '权限角色', '权限角色', '#', '4', '2', 'i7', '1', '1', '1634111053', '1514194275');
INSERT INTO `system_menu` VALUES ('6', '账号管理', '账号列表', 'user/list', '5', '3', '', '1', '1', '0', '1514194339');
INSERT INTO `system_menu` VALUES ('8', '菜单管理', '菜单列表', 'menu/list', '5', '3', '', '1', '3', '0', '1514194502');
INSERT INTO `system_menu` VALUES ('9', '账号编辑', '账号编辑', 'user/edit', '6', '4', '', '1', '1', '0', '1514292996');
INSERT INTO `system_menu` VALUES ('11', '账号删除', '后台账号删除', 'user/delete', '6', '4', '', '1', '2', '0', '1514293492');
INSERT INTO `system_menu` VALUES ('12', '账号状态', '账号启用禁用', 'user/switch', '6', '4', '', '1', '3', '0', '1514293677');
INSERT INTO `system_menu` VALUES ('18', '菜单编辑', '新增修改菜单', 'menu/edit', '8', '4', '', '1', '1', '0', '1514294748');
INSERT INTO `system_menu` VALUES ('19', '删除菜单', '删除菜单', 'menu/delete', '8', '4', '', '1', '2', '0', '1514294807');
INSERT INTO `system_menu` VALUES ('586', '数据', '数据抓取', '#', '0', '1', 'm2', '1', '2', '0', '1634110545');
INSERT INTO `system_menu` VALUES ('587', '任务管理', '任务管理', '#', '586', '2', 'i8', '1', '1', '1634111071', '1634110884');
INSERT INTO `system_menu` VALUES ('588', '我的采集', '采集列表', 'collect/list', '587', '3', '', '1', '1', '0', '1634110989');
INSERT INTO `system_menu` VALUES ('589', '清洗数据', '清洗数据', 'clean/list', '587', '3', '', '1', '2', '0', '1634111727');
INSERT INTO `system_menu` VALUES ('590', '数据导出', '数据下载', 'export/list', '587', '3', '', '1', '3', '0', '1634111960');
INSERT INTO `system_menu` VALUES ('591', '采集编辑', '采集编辑', 'collect/edit', '588', '4', '', '1', '1', '0', '1634188402');
INSERT INTO `system_menu` VALUES ('592', '删除采集', '删除采集', 'collect/remove', '588', '4', '', '1', '2', '0', '1634188426');
