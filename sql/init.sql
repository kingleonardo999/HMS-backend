-- HMS (Hotel Management System) 数据库初始化脚本
-- 基于 GORM AutoMigrate 模型定义生成
-- 数据库: hms

CREATE DATABASE IF NOT EXISTS `hms` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `hms`;

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role_name` VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_role_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 图片表
CREATE TABLE IF NOT EXISTS `imgs` (
    `id`  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `url` VARCHAR(255)    DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `login_id` VARCHAR(255)    NOT NULL,
    `password` VARCHAR(255)    NOT NULL,
    `name`     VARCHAR(255)    NOT NULL,
    `phone`    VARCHAR(255)    NOT NULL,
    `email`    VARCHAR(255)    DEFAULT NULL,
    `role_id`  BIGINT UNSIGNED NOT NULL,
    `img_id`   BIGINT UNSIGNED DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_login_id` (`login_id`),
    INDEX `idx_role_id` (`role_id`),
    INDEX `idx_img_id` (`img_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 房间类型表
CREATE TABLE IF NOT EXISTS `room_types` (
    `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `room_type_name`    VARCHAR(255)    NOT NULL,
    `room_type_price`   INT             NOT NULL,
    `type_description`  TEXT            NOT NULL,
    `bed_num`           INT             NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_room_type_name` (`room_type_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 房间状态字典表
CREATE TABLE IF NOT EXISTS `room_statuses` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `status_name` VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_status_name` (`status_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 房间表
CREATE TABLE IF NOT EXISTS `rooms` (
    `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `room_id`           VARCHAR(255)    NOT NULL,
    `room_type_id`      BIGINT UNSIGNED NOT NULL,
    `room_status_id`    BIGINT UNSIGNED NOT NULL,
    `room_description`  TEXT            DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_room_id` (`room_id`),
    INDEX `idx_room_type_id` (`room_type_id`),
    INDEX `idx_room_status_id` (`room_status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 客人表
CREATE TABLE IF NOT EXISTS `guests` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `identity_id` VARCHAR(255)    NOT NULL,
    `name`        VARCHAR(255)    NOT NULL,
    `phone`       VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_identity_id` (`identity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 入住状态字典表
CREATE TABLE IF NOT EXISTS `reside_states` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `state_name` VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_state_name` (`state_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 入住记录表
CREATE TABLE IF NOT EXISTS `resides` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `guest_id`     BIGINT UNSIGNED NOT NULL,
    `room_id`      VARCHAR(255)    NOT NULL,
    `reside_date`  VARCHAR(255)    NOT NULL,
    `leave_date`   VARCHAR(255)    DEFAULT NULL,
    `total_money`  INT             NOT NULL,
    `deposit`      INT             NOT NULL,
    `guest_num`    INT             NOT NULL,
    `reside_state` VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_guest_id` (`guest_id`),
    INDEX `idx_reside_room_id` (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 订单表
CREATE TABLE IF NOT EXISTS `orders` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `order_id`    VARCHAR(255)    NOT NULL,
    `guest_id`    BIGINT UNSIGNED NOT NULL,
    `room_id`     VARCHAR(255)    NOT NULL,
    `order_date`  VARCHAR(255)    NOT NULL,
    `leave_date`  VARCHAR(255)    NOT NULL,
    `total_money` INT             NOT NULL,
    `guest_num`   INT             NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_order_id` (`order_id`),
    INDEX `idx_order_guest_id` (`guest_id`),
    INDEX `idx_order_room_id` (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 账单表
CREATE TABLE IF NOT EXISTS `billings` (
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `amount`        INT             NOT NULL,
    `time`          VARCHAR(255)    NOT NULL,
    `guest_id`      BIGINT UNSIGNED NOT NULL,
    `room_id`       VARCHAR(255)    NOT NULL,
    `room_type_name` VARCHAR(255)   NOT NULL,
    `reside_id`     BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_billing_guest_id` (`guest_id`),
    INDEX `idx_billing_room_id` (`room_id`),
    INDEX `idx_billing_reside_id` (`reside_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜品分类表
CREATE TABLE IF NOT EXISTS `menu_types` (
    `id`    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `type`  VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_menu_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜品表
CREATE TABLE IF NOT EXISTS `menus` (
    `id`      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`    VARCHAR(255)    NOT NULL,
    `type_id` BIGINT UNSIGNED NOT NULL,
    `price`   INT             NOT NULL,
    `img_id`  BIGINT UNSIGNED NOT NULL,
    `desc`    VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_menu_type_id` (`type_id`),
    INDEX `idx_menu_img_id` (`img_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 消息表
CREATE TABLE IF NOT EXISTS `messages` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `admin_id`   BIGINT UNSIGNED NOT NULL,
    `title`      VARCHAR(255)    NOT NULL,
    `content`    TEXT            NOT NULL,
    `created_at` DATETIME(3)     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户消息关联表（软删除）
CREATE TABLE IF NOT EXISTS `user_messages` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`    BIGINT UNSIGNED NOT NULL,
    `message_id` BIGINT UNSIGNED NOT NULL,
    `deleted_at` DATETIME(3)     DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_message_id` (`message_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 初始化数据（与 config/db.go 中的 initDB 逻辑一致）
-- ============================================================

-- 插入默认角色
INSERT IGNORE INTO `roles` (`id`, `role_name`) VALUES (1, 'admin');
INSERT IGNORE INTO `roles` (`id`, `role_name`) VALUES (2, 'user');

-- 插入默认头像
INSERT IGNORE INTO `imgs` (`id`, `url`) VALUES (1, 'https://imgs.161517.xyz/2025/05/10/default-photo.jpg');

-- 插入房间状态字典
INSERT IGNORE INTO `room_statuses` (`id`, `status_name`) VALUES (1, '空闲');
INSERT IGNORE INTO `room_statuses` (`id`, `status_name`) VALUES (2, '已入住');
INSERT IGNORE INTO `room_statuses` (`id`, `status_name`) VALUES (3, '已预定');

-- 插入入住状态字典
INSERT IGNORE INTO `reside_states` (`id`, `state_name`) VALUES (1, '未结账');
INSERT IGNORE INTO `reside_states` (`id`, `state_name`) VALUES (2, '已结账');

-- 插入默认管理员 (密码: admin, 由 utils.HashPassword 生成)
-- 如需修改密码，请使用程序提供的注册接口
INSERT IGNORE INTO `users` (`id`, `login_id`, `password`, `name`, `phone`, `role_id`, `img_id`)
VALUES (1, 'admin', '$2a$10$iUSKOwCHqfVpVDi0PpKQtOtLoQT8h3CGbLptN/vu0KeT4sRqAKsGq', 'admin', '12345678901', 1, 1);