/*
 Navicat Premium Data Transfer

 Source Server         : 抖音
 Source Server Type    : MySQL
 Source Server Version : 50650
 Source Host           : 43.138.25.60:3306
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 50650
 File Encoding         : 65001

 Date: 04/06/2022 23:43:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id，自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '评论发布用户id',
  `video_id` bigint(20) NOT NULL COMMENT '评论视频id',
  `comment_text` varchar(255) NOT NULL COMMENT '评论内容',
  `create_date` datetime NOT NULL COMMENT '评论发布时间',
  `cancel` tinyint(4) NOT NULL DEFAULT '0' COMMENT '默认评论发布为0，取消后为1',
  PRIMARY KEY (`id`),
  KEY `videoIdIdx` (`video_id`) USING BTREE COMMENT '评论列表使用视频id作为索引-方便查看视频下的评论列表'
) ENGINE=InnoDB AUTO_INCREMENT=1206 DEFAULT CHARSET=utf8 COMMENT='评论表';

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `follower_id` bigint(20) NOT NULL COMMENT '关注的用户',
  `cancel` tinyint(4) NOT NULL DEFAULT '0' COMMENT '默认关注为0，取消关注为1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `userIdToFollowerIdIdx` (`user_id`,`follower_id`) USING BTREE,
  KEY `FollowerIdIdx` (`follower_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1096 DEFAULT CHARSET=utf8 COMMENT='关注表';

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '点赞用户id',
  `video_id` bigint(20) NOT NULL COMMENT '被点赞的视频id',
  `cancel` tinyint(4) NOT NULL DEFAULT '0' COMMENT '默认点赞为0，取消赞为1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
  KEY `userIdIdx` (`user_id`) USING BTREE,
  KEY `videoIdx` (`video_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1229 DEFAULT CHARSET=utf8 COMMENT='点赞表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
  `name` varchar(255) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '用户密码',
  PRIMARY KEY (`id`),
  KEY `name_password_idx` (`name`,`password`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20044 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
  `author_id` bigint(20) NOT NULL COMMENT '视频作者id',
  `play_url` varchar(255) NOT NULL COMMENT '播放url',
  `cover_url` varchar(255) NOT NULL COMMENT '封面url',
  `publish_time` datetime NOT NULL COMMENT '发布时间戳',
  `title` varchar(255) DEFAULT NULL COMMENT '视频名称',
  PRIMARY KEY (`id`),
  KEY `time` (`publish_time`) USING BTREE,
  KEY `author` (`author_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8 COMMENT='\r\n视频表';

-- ----------------------------
-- Procedure structure for addFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `addFollowRelation`;
delimiter ;;
CREATE PROCEDURE `addFollowRelation`(IN user_id int,IN follower_id int)
BEGIN
	#Routine body goes here...
	# 声明记录个数变量。
	DECLARE cnt INT DEFAULT 0;
	# 获取记录个数变量。
	SELECT COUNT(1) FROM follows f where f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 判断是否已经存在该记录，并做出相应的插入关系、更新关系动作。
	# 插入操作。
	IF cnt = 0 THEN
		INSERT INTO follows(`user_id`,`follower_id`) VALUES(user_id,follower_id);
	END IF;
	# 更新操作
	IF cnt != 0 THEN
		UPDATE follows f SET f.cancel = 0 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for delFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `delFollowRelation`;
delimiter ;;
CREATE PROCEDURE `delFollowRelation`(IN `user_id` int,IN `follower_id` int)
BEGIN
	#Routine body goes here...
	# 定义记录个数变量，记录是否存在此关系，默认没有关系。
	DECLARE cnt INT DEFAULT 0;
	# 查看是否之前有关系。
	SELECT COUNT(1) FROM follows f WHERE f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 有关系，则需要update cancel = 1，使其关系无效。
	IF cnt = 1 THEN
		UPDATE follows f SET f.cancel = 1 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
