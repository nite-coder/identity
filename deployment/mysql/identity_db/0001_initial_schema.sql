/*
 Navicat Premium Data Transfer

 Source Server         : Local
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : localhost:3306
 Source Schema         : identity_db

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 05/09/2021 09:00:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` char(36) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `namespace` varchar(256) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `type` int NOT NULL,
  `username` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `password_encrypt` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `first_name` varchar(24) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `last_name` varchar(24) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `avatar` varchar(24) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `email` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `mobile_country_code` varchar(5) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `mobile` varchar(20) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `external_id` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `failed_password_attempt` int NOT NULL,
  `otp_enable` tinyint NOT NULL,
  `otp_secret` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `otp_effective_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `otp_last_reset_at` datetime(6) NOT NULL DEFAULT '1970-01-01 00:00:00.000000',
  `client_ip` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `user_agent` varchar(512) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `notes` varchar(512) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `last_login_at` datetime NOT NULL,
  `is_admin` tinyint NOT NULL,
  `state` int NOT NULL,
  `version` int NOT NULL,
  `creator_id` int NOT NULL,
  `creator_name` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updater_id` int NOT NULL,
  `updater_name` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_uuid`(`uuid`) USING BTREE,
  UNIQUE INDEX `uniq_username`(`namespace`, `username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for accounts_roles
-- ----------------------------
DROP TABLE IF EXISTS `accounts_roles`;
CREATE TABLE `accounts_roles`  (
  `account_id` int NOT NULL,
  `role_id` int NOT NULL,
  PRIMARY KEY (`account_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` char(36) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `namespace` varchar(256) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `name` varchar(24) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `rules` json NOT NULL,
  `creator_id` int NOT NULL,
  `creator_name` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updater_id` int NOT NULL,
  `updater_name` varchar(128) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
