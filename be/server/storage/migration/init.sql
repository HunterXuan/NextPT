SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- 模块: 用户系统
-- ============================================================

-- 动态用户角色表
CREATE TABLE `iam_role` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `level`           SMALLINT        NOT NULL DEFAULT 0 COMMENT '等级权重(用于权限比对，值越大权限越高，如普通用户10，管理员100)',
    `name_i18n`       JSON            NOT NULL COMMENT '角色名称多语言映射字典',
    `rules`           JSON            NULL     COMMENT '角色规则(JSON: promotion 表示升级到该角色的条件, demotion 表示当前角色触发降级的条件)',
    `permissions`     JSON            NULL     COMMENT '角色关联的权限标识符列表',
    `is_staff`        BIT(1)      NOT NULL DEFAULT 0 COMMENT '是否为管理组成员',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色';

-- 用户专属细粒度权限 (ACL)
CREATE TABLE `iam_user_permission` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `perm_key`        VARCHAR(100)    NOT NULL COMMENT '权限标识符，如 update:catalog/torrent:123 或 update:catalog/torrent:*',
    `source_type`     TINYINT         NOT NULL DEFAULT 1 COMMENT '来源类型: 1=manual 2=user_mod',
    `source_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源记录ID，manual=0，user_mod=mod_user_log.id',
    `expire_at`       DATETIME        NULL     COMMENT '权限过期时间，NULL=永久',
    `is_active`       BIT(1)          NOT NULL DEFAULT 1 COMMENT '是否生效',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_perm_source` (`user_id`, `perm_key`, `source_type`, `source_id`),
    KEY `idx_user_active_expire` (`user_id`, `is_active`, `expire_at`),
    KEY `idx_source` (`source_type`, `source_id`),
    KEY `idx_perm_key` (`perm_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户独立权限(ACL)';

-- 用户资料表（只保留认证和身份相关的最小字段集）
CREATE TABLE `iam_user` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username`        VARCHAR(40)     NOT NULL DEFAULT '',
    `email`           VARCHAR(120)    NOT NULL DEFAULT '',
    `password_hash`   VARCHAR(255)    NOT NULL DEFAULT '' COMMENT 'bcrypt hash',
    `passkey`         CHAR(32)        NOT NULL DEFAULT '' COMMENT 'Tracker passkey',
    `status`          TINYINT         NOT NULL DEFAULT 0  COMMENT '0=pending 1=confirmed 2=disabled',
    `role`            INT UNSIGNED    NOT NULL DEFAULT 1  COMMENT '当前角色ID (关联 user_role 表)',
    `vip_until`       DATETIME        NULL     COMMENT 'VIP 过期时间',
    `vip_remark`      VARCHAR(100)    NOT NULL DEFAULT '' COMMENT 'VIP 身份获取备注/来源',
    `two_step_type`   TINYINT         NOT NULL DEFAULT 0  COMMENT '两步验证方式: 0=关闭 1=TOTP(Authenticator) 2=邮件验证码',
    `two_step_secret` VARCHAR(64)     NOT NULL DEFAULT '' COMMENT 'TOTP 密钥 (two_step_type=1 时使用；邮件验证码走 Redis 临时存储)',
    `invited_by`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `last_login`      DATETIME        NULL     COMMENT '最后登录时间',
    `last_ip`         VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '最后登录 IP',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    `deleted_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_passkey` (`passkey`(8)),
    KEY `idx_status` (`status`),
    KEY `idx_role` (`role`),
    KEY `idx_invited_by` (`invited_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户核心表';

-- 用户档案表（精简：只保留核心展示信息）
CREATE TABLE `iam_user_profile` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `avatar`          VARCHAR(500)    NOT NULL DEFAULT '',
    `info`            TEXT            NULL     COMMENT '个人简介 (Markdown)',
    `signature`       VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '论坛签名',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户档案表';

-- 用户流量/做种统计表（频繁更新的计数器隔离）
CREATE TABLE `iam_user_stat` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `uploaded`        BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '总入账上传量 (bytes)',
    `downloaded`      BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '总入账下载量 (bytes)',
    `raw_uploaded`    BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '真实总上传量 (bytes)',
    `raw_downloaded`  BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '真实总下载量 (bytes)',
    `seed_time`       BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '总做种时间 (秒)',
    `leech_time`      BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '总下载时间 (秒)',
    `bonus`      DECIMAL(12,1)   NOT NULL DEFAULT 0.0 COMMENT '魔力值',
    `bonus_charity`   DECIMAL(12,1)   NOT NULL DEFAULT 0.0 COMMENT '捐赠魔力值',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    KEY `idx_uploaded` (`uploaded`),
    KEY `idx_downloaded` (`downloaded`),
    KEY `idx_seed_bonus` (`bonus`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户流量统计表';

-- 用户偏好设置表（精简：UI 偏好收进 JSON，只保留服务端必须的独立列）
CREATE TABLE `iam_user_setting` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `privacy_level`   TINYINT         NOT NULL DEFAULT 1  COMMENT '0=宽松 1=普通 2=严格',
    `extra`           JSON            NULL     COMMENT 'UI偏好/分页/语言/时区等前端设置 (JSON 扩展)',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户偏好设置表';


-- 用户警告/处罚记录
CREATE TABLE `mod_user_log` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `mod_type`        TINYINT         NOT NULL COMMENT '1=warned 2=banned 3=leech_warned 4=upload_banned 5=download_banned 6=forum_banned',
    `reason`          VARCHAR(500)    NOT NULL DEFAULT '',
    `expire_at`       DATETIME        NULL     COMMENT '过期时间，NULL=永久',
    `mod_by`          BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '操作人ID',
    `mod_comment`     TEXT            NULL,
    `is_active`       BIT(1)      NOT NULL DEFAULT 1,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_active_type_expire` (`user_id`, `is_active`, `mod_type`, `expire_at`),
    KEY `idx_expire` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户处罚记录表';


-- 用户登录日志
CREATE TABLE `iam_login_log` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `ip`              VARCHAR(64)     NOT NULL DEFAULT '',
    `user_agent`      VARCHAR(500)    NOT NULL DEFAULT '',
    `result`          TINYINT         NOT NULL DEFAULT 1  COMMENT '1=success 0=fail',
    `fail_reason`     VARCHAR(100)    NOT NULL DEFAULT '',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_ip` (`ip`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录日志';

-- 用户周期统计增量（period_type: 1=每日, 2=每月；period_key: YYYY-MM-DD 或 YYYY-MM）
CREATE TABLE `iam_user_period_stat` (
    `id`              BIGINT UNSIGNED     NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED     NOT NULL,
    `period_type`     TINYINT UNSIGNED    NOT NULL COMMENT '1=每日 2=每月',
    `period_key`      VARCHAR(10)         NOT NULL COMMENT 'YYYY-MM-DD 或 YYYY-MM',
    `uploaded`        BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增入账上传量 (bytes)',
    `downloaded`      BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增入账下载量 (bytes)',
    `raw_uploaded`    BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增真实上传量 (bytes)',
    `raw_downloaded`  BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增真实下载量 (bytes)',
    `seed_time`       BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增做种时间 (秒)',
    `leech_time`      BIGINT UNSIGNED     NOT NULL DEFAULT 0   COMMENT '周期新增下载时间 (秒)',
    `bonus`           DECIMAL(12,1)       NOT NULL DEFAULT 0.0 COMMENT '周期获得魔力值',
    `created_at`      DATETIME            NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_period` (`user_id`, `period_type`, `period_key`),
    KEY `idx_period` (`period_type`, `period_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户周期统计增量';

-- ============================================================
-- 模块: 邀请系统
-- ============================================================

CREATE TABLE `iam_invite` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `inviter_id`      BIGINT UNSIGNED NOT NULL COMMENT '邀请人',
    `invitee_email`   VARCHAR(120)    NOT NULL DEFAULT '' COMMENT '被邀请人邮箱',
    `invitee_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '被邀请人ID（注册后回填）',
    `hash`            CHAR(32)        NOT NULL COMMENT '邀请码 (发放名额时就预生成唯一码)',
    `status`          TINYINT         NOT NULL DEFAULT 0  COMMENT '0=未分配/待发送 1=已发送 2=已注册 3=已过期 4=已回收',
    `is_temporary`    BIT(1)      NOT NULL DEFAULT 0  COMMENT '是否限时邀请',
    `expire_at`       DATETIME        NULL,
    `used_at`         DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_inviter` (`inviter_id`),
    UNIQUE KEY `uk_hash` (`hash`),
    KEY `idx_status_expire` (`status`, `expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邀请表';

-- ============================================================
-- 模块: 种子 / Tracker
-- ============================================================

-- 种子分类
CREATE TABLE `catalog_category` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `name_i18n`       JSON            NOT NULL COMMENT '多语言名称映射',
    `slug`            VARCHAR(60)     NOT NULL DEFAULT '' COMMENT 'URL-friendly',
    `sort_order`      SMALLINT        NOT NULL DEFAULT 0,
    `enabled`         BIT(1)      NOT NULL DEFAULT 1,
    `upload_config`   JSON            NULL     COMMENT '分类发布表单配置',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`),
    KEY `idx_sort` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子分类表';

-- 标签分组 (例如: 分辨率、视频编码)
CREATE TABLE `catalog_tag_group` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `name_i18n`       JSON            NOT NULL COMMENT '多语言名称映射',
    `slug`            VARCHAR(30)     NOT NULL COMMENT '英文标识如 resolution',
    `category_ids`    JSON            NULL     COMMENT '适用的分类 ID 数组 (例: [1,2])，为空则全站通用',
    `sort_order`      SMALLINT        NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签分组表';

-- 标签项
CREATE TABLE `catalog_tag` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `group_id`        INT UNSIGNED    NOT NULL COMMENT '所属分组',
    `name_i18n`       JSON            NOT NULL COMMENT '多语言名称映射',
    `value`           VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '稳定值',
    `sort_order`      SMALLINT        NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_value` (`group_id`, `value`),
    KEY `idx_group_sort` (`group_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签项表';

-- 种子核心表
CREATE TABLE `catalog_torrent` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `info_hash`       BINARY(20)      NOT NULL COMMENT 'BitTorrent info_hash',
    `name`            VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '种子标题',
    `sub_title`       VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '副标题',
    `category_id`     INT UNSIGNED    NOT NULL DEFAULT 0,
    `description`     MEDIUMTEXT      NULL     COMMENT '详情描述 (BBCode/Markdown)',
    `release_fields`  JSON            NULL     COMMENT '发布结构化字段值',
    `file_name`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '种子文件名',
    `size`            BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '总大小 (bytes)',
    `file_count`      INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '文件数量',
    `owner_id`        BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '上传者',
    `anonymous`       BIT(1)          NOT NULL DEFAULT 0  COMMENT '匿名上传',
    `status`          TINYINT         NOT NULL DEFAULT 1  COMMENT '0=待审核 1=已发布 2=已拒绝',
    `submitted_at`    DATETIME        NULL     COMMENT '最近提交审核时间',
    `published_at`    DATETIME        NULL     COMMENT '实际发布时间',
    `reviewed_by`     BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '最后审核人',
    `reviewed_at`     DATETIME        NULL     COMMENT '最后审核时间',
    `review_comment`  VARCHAR(1000)   NOT NULL DEFAULT '' COMMENT '最后审核意见',

    -- 促销
    `sp_state`        TINYINT         NOT NULL DEFAULT 0  COMMENT '0=normal 1=free 2=2x 3=2xfree 4=50%off 5=2x50% 6=30%off',
    `sp_expire_at`    DATETIME        NULL     COMMENT '促销到期时间',
    `is_featured`     BIT(1)      NOT NULL DEFAULT 0  COMMENT '是否推荐',
    `is_pinned`       BIT(1)      NOT NULL DEFAULT 0  COMMENT '是否置顶',
    `pin_weight`      SMALLINT        NOT NULL DEFAULT 0  COMMENT '置顶权重 (越大越靠前)',

    -- Tracker 统计（缓存字段）
    `seeders`         INT UNSIGNED    NOT NULL DEFAULT 0,
    `leechers`        INT UNSIGNED    NOT NULL DEFAULT 0,
    `times_completed` INT UNSIGNED    NOT NULL DEFAULT 0,
    `comments_count`  INT UNSIGNED    NOT NULL DEFAULT 0,
    `views`           INT UNSIGNED    NOT NULL DEFAULT 0,
    `like_count`      INT UNSIGNED    NOT NULL DEFAULT 0,
    `rewards_count`   INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '收到的赞赏次数',
    `rewards_amount`  DOUBLE          NOT NULL DEFAULT 0  COMMENT '收到的赞赏总金额(Bonus)',

    `banned`          BIT(1)          NOT NULL DEFAULT 0,
    `last_action`     DATETIME        NULL     COMMENT 'Tracker 最后活动时间',
    `last_reseed`     DATETIME        NULL,

    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_info_hash` (`info_hash`),
    KEY `idx_owner` (`owner_id`),
    KEY `idx_category_status` (`category_id`, `status`, `banned`, `published_at`),
    KEY `idx_published` (`status`, `banned`, `is_pinned`, `pin_weight`, `published_at`),
    KEY `idx_review_queue` (`status`, `submitted_at`),
    FULLTEXT KEY `ft_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子核心表';

-- 种子-标签关联（多对多）
CREATE TABLE `catalog_torrent_tag` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `tag_id`          INT UNSIGNED    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_torrent_tag` (`torrent_id`, `tag_id`),
    KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子标签关联表';

-- 种子外部资源身份绑定（完整元数据由外部来源按需获取并缓存）
CREATE TABLE `catalog_torrent_meta` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `imdb_id`         VARCHAR(20)     NOT NULL DEFAULT '',
    `imdb_rating`     DECIMAL(3,1)    NULL,
    `douban_id`       VARCHAR(20)     NOT NULL DEFAULT '',
    `douban_rating`   DECIMAL(3,1)    NULL,
    `bangumi_id`      VARCHAR(20)     NOT NULL DEFAULT '',
    `bangumi_rating`  DECIMAL(3,1)    NULL,
    `tmdb_id`         VARCHAR(20)     NOT NULL DEFAULT '',
    `tmdb_type`       VARCHAR(10)     NOT NULL DEFAULT '' COMMENT 'movie/tv',
    `tmdb_rating`     DECIMAL(3,1)    NULL,
    `extra`           JSON            NULL     COMMENT '其他外部标识与绑定参数',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_torrent_id` (`torrent_id`),
    KEY `idx_imdb` (`imdb_id`),
    KEY `idx_douban` (`douban_id`),
    KEY `idx_bangumi` (`bangumi_id`),
    KEY `idx_tmdb_type` (`tmdb_id`, `tmdb_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子外部资源身份绑定';

-- 种子文件列表
CREATE TABLE `catalog_torrent_file` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `file_path`       VARCHAR(1000)   NOT NULL DEFAULT '' COMMENT '文件路径',
    `size`            BIGINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_torrent` (`torrent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子文件列表';

-- 活跃 Peer 表（Tracker 核心）
CREATE TABLE `tracker_peer` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `peer_id`         BINARY(20)      NOT NULL,
    `ipv4`            VARCHAR(64)     NOT NULL DEFAULT '',
    `ipv6`            VARCHAR(64)     NOT NULL DEFAULT '',
    `port`            SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    `uploaded`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `downloaded`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `remaining`       BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '剩余大小',
    `is_seeder`       BIT(1)      NOT NULL DEFAULT 0,
    `is_connectable`  BIT(1)      NOT NULL DEFAULT 1,
    `agent`           VARCHAR(100)    NOT NULL DEFAULT '' COMMENT 'BT 客户端',
    `passkey`         CHAR(32)        NOT NULL DEFAULT '',
    `upload_offset`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `download_offset` BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `started_at`      DATETIME        NULL,
    `last_action`     DATETIME        NULL,
    `finished_at`     DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_torrent` (`torrent_id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_passkey` (`passkey`(8))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='活跃 Peer 表';
-- 注: 生产环境可考虑 ENGINE=MEMORY 或 Redis 替代

-- 下载/做种完成记录
CREATE TABLE `tracker_snatch` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `ipv4`            VARCHAR(64)     NOT NULL DEFAULT '',
    `ipv6`            VARCHAR(64)     NOT NULL DEFAULT '',
    `port`            SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    `uploaded`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `downloaded`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `remaining`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `seed_time`       INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '做种时间 (秒)',
    `leech_time`      INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '下载时间 (秒)',
    `is_finished`     BIT(1)      NOT NULL DEFAULT 0,
    `started_at`      DATETIME        NULL,
    `completed_at`    DATETIME        NULL,
    `last_action`     DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_torrent_user` (`torrent_id`, `user_id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='下载/做种完成记录表';

-- BT 客户端白名单
CREATE TABLE `tracker_agent_whitelist` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `family`          VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '客户端家族名',
    `peer_id_prefix`  VARCHAR(20)     NOT NULL DEFAULT '',
    `agent_pattern`   VARCHAR(200)    NOT NULL DEFAULT '',
    `min_version`     VARCHAR(20)     NOT NULL DEFAULT '',
    `max_version`     VARCHAR(20)     NOT NULL DEFAULT '',
    `allow_https`     BIT(1)      NOT NULL DEFAULT 0,
    `enabled`         BIT(1)      NOT NULL DEFAULT 1,
    `comment`         VARCHAR(200)    NOT NULL DEFAULT '',
    `hits`            INT UNSIGNED    NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='BT 客户端白名单';

-- 种子收藏
CREATE TABLE `catalog_torrent_bookmark` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_torrent` (`user_id`, `torrent_id`),
    KEY `idx_torrent` (`torrent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子收藏表';

-- 求种 / 续种请求
CREATE TABLE `catalog_request` (
    `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `request_type`      TINYINT UNSIGNED NOT NULL COMMENT '1=求种 2=续种',
    `requester_id`      BIGINT UNSIGNED NOT NULL,
    `category_id`       INT UNSIGNED    NOT NULL DEFAULT 0,
    `target_torrent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '续种目标种子',
    `result_torrent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '求种完成后关联的种子',
    `title`             VARCHAR(500)    NOT NULL DEFAULT '',
    `description`       TEXT            NOT NULL COMMENT '请求说明 (Markdown/BBCode)',
    `reward_amount`     DECIMAL(12,1)   NOT NULL DEFAULT 0.0 COMMENT '已托管的魔力奖励',
    `status`            TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=开放 1=已认领 2=待确认 3=已完成 4=已取消',
    `claimed_by`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `claimed_at`        DATETIME        NULL,
    `claim_expires_at`  DATETIME        NULL,
    `submitted_at`      DATETIME        NULL,
    `completed_at`      DATETIME        NULL,
    `cancelled_by`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `cancelled_at`      DATETIME        NULL,
    `cancel_reason`     VARCHAR(500)    NOT NULL DEFAULT '',
    `created_at`        DATETIME        NULL,
    `updated_at`        DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_status_created` (`status`, `created_at`),
    KEY `idx_type_status_created` (`request_type`, `status`, `created_at`),
    KEY `idx_requester_status` (`requester_id`, `status`, `created_at`),
    KEY `idx_claimed_status` (`claimed_by`, `status`, `claim_expires_at`),
    KEY `idx_category_status` (`category_id`, `status`, `created_at`),
    KEY `idx_target_status` (`target_torrent_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='求种与续种请求表';

-- 种子感谢
-- 种子赞赏
-- 作弊检测记录
CREATE TABLE `mod_cheater_log` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `uploaded`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `downloaded`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `announce_time`   INT UNSIGNED    NOT NULL DEFAULT 0,
    `seeders`         INT UNSIGNED    NOT NULL DEFAULT 0,
    `leechers`        INT UNSIGNED    NOT NULL DEFAULT 0,
    `hit_count`       TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dealt_by`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `is_dealt`        BIT(1)      NOT NULL DEFAULT 0,
    `comment`         VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '检测备注',
    `dealt_comment`   VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '处理说明',
    `dealt_at`        DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_torrent` (`torrent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作弊检测记录';

-- ============================================================
-- 模块: 通用评论系统
-- ============================================================

-- 通用评论表（支持域/资源形式的多种对象）
CREATE TABLE `catalog_comment` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `target_type`     VARCHAR(20)     NOT NULL COMMENT 'catalog_torrent/catalog_request',
    `target_id`       BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `content`         TEXT            NOT NULL COMMENT '评论内容 (Markdown/BBCode)',
    `reward_count`    INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '获得的魔力值打赏',
    `like_count`      INT UNSIGNED    NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    `deleted_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_target` (`target_type`, `target_id`, `created_at`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通用评论表';

-- ============================================================
-- 模块: 论坛系统
-- ============================================================

-- 节点分类（相当于 V2EX 的大分类，如 技术、生活）
CREATE TABLE `forum_category` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `name_i18n`       JSON            NOT NULL COMMENT '多语言名称映射',
    `desc_i18n`       JSON            NULL COMMENT '多语言描述映射',
    `sort_order`      SMALLINT        NOT NULL DEFAULT 0,
    `min_role_view`   TINYINT         NOT NULL DEFAULT 0  COMMENT '最低可见等级',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='节点分类';

-- 节点（相当于 V2EX 的 Node，如 /go/programmer）
CREATE TABLE `forum_node` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `category_id`     INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '所属节点分类',
    `slug`            VARCHAR(60)     NOT NULL COMMENT '节点英文标识符',
    `name_i18n`       JSON            NOT NULL COMMENT '多语言名称映射',
    `desc_i18n`       JSON            NULL COMMENT '多语言描述映射',
    `sort_order`      SMALLINT        NOT NULL DEFAULT 0,
    `min_role_read`   TINYINT         NOT NULL DEFAULT 0,
    `min_role_write`  TINYINT         NOT NULL DEFAULT 0,
    `min_role_create` TINYINT         NOT NULL DEFAULT 0,
    `topic_count`     INT UNSIGNED    NOT NULL DEFAULT 0,
    `reply_count`     INT UNSIGNED    NOT NULL DEFAULT 0,
    `last_topic_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `moderators`      JSON            NULL COMMENT '板块版主ID列表(前端展示用)',
    `last_reply_at`   DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`),
    KEY `idx_category_sort` (`category_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区节点';

-- 主题 (包含正文)
CREATE TABLE `forum_topic` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `node_id`         INT UNSIGNED    NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `subject`         VARCHAR(200)    NOT NULL,
    `content`         TEXT            NOT NULL COMMENT '主题正文',
    `appends`         JSON            NULL     COMMENT '追加内容 (V2EX风格) [{"content": "...", "created_at": "..."}]',
    `is_locked`       BIT(1)      NOT NULL DEFAULT 0,
    `is_sticky`       BIT(1)      NOT NULL DEFAULT 0,
    `views`           INT UNSIGNED    NOT NULL DEFAULT 0,
    `reply_count`     INT UNSIGNED    NOT NULL DEFAULT 0,
    `like_count`      INT UNSIGNED    NOT NULL DEFAULT 0,
    `last_reply_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `last_reply_at`   DATETIME        NULL,
    `last_reply_by`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    `deleted_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_node_sticky_last` (`node_id`, `is_sticky`, `last_reply_at`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区主题';

-- 论坛主题收藏表
CREATE TABLE `forum_topic_bookmark` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `topic_id`        BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_topic` (`user_id`, `topic_id`),
    KEY `idx_topic` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='论坛主题收藏表';

-- 论坛主题点赞表
CREATE TABLE `forum_topic_like` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `topic_id`        BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_topic` (`user_id`, `topic_id`),
    KEY `idx_topic` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='论坛主题点赞表';

-- 论坛回复点赞表
CREATE TABLE `forum_reply_like` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `reply_id`        BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_reply` (`user_id`, `reply_id`),
    KEY `idx_reply` (`reply_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='论坛回复点赞表';

-- 论坛回复
CREATE TABLE `forum_reply` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `topic_id`        BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `content`         TEXT            NOT NULL COMMENT '回复内容',
    `reward_count`    INT UNSIGNED    NOT NULL DEFAULT 0  COMMENT '获得的魔力值打赏',
    `like_count`      INT UNSIGNED    NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    `deleted_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_topic_created` (`topic_id`, `created_at`),
    KEY `idx_user` (`user_id`),
    FULLTEXT KEY `ft_content` (`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='论坛回复';

-- 模块: 通知与工单系统
-- ============================================================

CREATE TABLE `site_announcement` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`           VARCHAR(200)    NOT NULL DEFAULT '',
    `content`         TEXT            NOT NULL,
    `status`          TINYINT         NOT NULL DEFAULT 0 COMMENT '0=draft 1=published 2=archived',
    `created_by`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `updated_by`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `published_at`    DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_status_published` (`status`, `published_at`),
    KEY `idx_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='站点公告';

CREATE TABLE `site_announcement_read` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `announcement_id` BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `read_at`         DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_announcement_user` (`announcement_id`, `user_id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='站点公告已读记录';

-- 系统通知 (仅允许系统或管理员发送给普通用户)
CREATE TABLE `site_message` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `sender_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '0=系统通知, 或管理员ID',
    `receiver_id`     BIGINT UNSIGNED NOT NULL,
    `title`           VARCHAR(200)    NOT NULL DEFAULT '',
    `content`         TEXT            NOT NULL,
    `target_type`     VARCHAR(30)     NOT NULL DEFAULT '',
    `target_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `is_read`         BIT(1)          NOT NULL DEFAULT 0,
    `read_at`         DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_receiver_read` (`receiver_id`, `is_read`, `created_at`),
    KEY `idx_receiver_created` (`receiver_id`, `created_at`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_sender` (`sender_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统通知表';

-- 管理员工单
CREATE TABLE `mod_staff_message` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `sender_id`       BIGINT UNSIGNED NOT NULL,
    `subject`         VARCHAR(200)    NOT NULL DEFAULT '',
    `content`         TEXT            NOT NULL,
    `status`          TINYINT         NOT NULL DEFAULT 0  COMMENT '0=pending 1=answered 2=closed',
    `answered_by`     BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `answer`          TEXT            NULL,
    `answered_at`     DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_sender` (`sender_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员工单';

-- ============================================================
-- 模块: 字幕系统
-- ============================================================

-- 种子字幕表
CREATE TABLE `catalog_subtitle` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `title`           VARCHAR(255)    NOT NULL DEFAULT '',
    `file_name`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '原始文件名',
    `file_ext`        VARCHAR(10)     NOT NULL DEFAULT '',
    `file_size`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `storage_path`    VARCHAR(500)    NOT NULL DEFAULT '' COMMENT '存储路径/S3 key',
    `language`        VARCHAR(10)     NOT NULL DEFAULT '' COMMENT '语言代码 (如 zh-CN, en-US)',
    `download_count`  INT UNSIGNED    NOT NULL DEFAULT 0,
    `anonymous`       BIT(1)      NOT NULL DEFAULT 0 COMMENT '匿名上传',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_torrent` (`torrent_id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子字幕表';

CREATE TABLE `tracker_event_idempotency` (
  `msg_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`msg_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tracker事件幂等表';


-- 3.11 公告 / 活动 / 投票
-- ============================================================
-- 模块: 举报 / 审计
-- ============================================================

-- 通用举报
CREATE TABLE `mod_report` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `reporter_id`     BIGINT UNSIGNED NOT NULL,
    `target_type`     VARCHAR(20)     NOT NULL COMMENT 'catalog_torrent/catalog_comment/catalog_subtitle/forum_topic/forum_reply',
    `target_id`       BIGINT UNSIGNED NOT NULL,
    `reason`          VARCHAR(500)    NOT NULL DEFAULT '',
    `status`          TINYINT         NOT NULL DEFAULT 0  COMMENT '0=pending 1=resolved 2=rejected',
    `dealt_by`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `dealt_comment`   VARCHAR(500)    NOT NULL DEFAULT '',
    `dealt_at`        DATETIME        NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通用举报表';

-- 审计日志（记录管理员与高危系统操作）
CREATE TABLE `site_audit` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `action`          VARCHAR(100)    NOT NULL DEFAULT '',
    `target_type`     VARCHAR(30)     NOT NULL DEFAULT '',
    `target_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `detail`          TEXT            NULL,
    `ip`              VARCHAR(64)     NOT NULL DEFAULT '',
    `level`           TINYINT         NOT NULL DEFAULT 0  COMMENT '0=normal 1=important 2=critical',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_created` (`created_at`),
    KEY `idx_user` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_level_created` (`level`, `created_at`),
    KEY `idx_user_created` (`user_id`, `created_at`),
    KEY `idx_action_created` (`action`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志';

-- 用户 IP 历史记录（用于防多开/马甲查询）
CREATE TABLE `iam_user_ip_history` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `ip`              VARCHAR(64)     NOT NULL,
    `type`            VARCHAR(20)     NOT NULL DEFAULT 'login' COMMENT 'login/tracker/register',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_ip` (`ip`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户IP历史记录';

-- ============================================================
-- 模块: 系统配置 / 通用
-- ============================================================

-- 系统配置（KV）
CREATE TABLE `site_config` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `group`           VARCHAR(30)     NOT NULL DEFAULT 'general',
    `key`             VARCHAR(64)     NOT NULL,
    `value`           JSON            NOT NULL COMMENT '配置值(JSON: {"val": ...})',
    `created_at`      DATETIME        NULL,
    `updated_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_key` (`group`, `key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';


-- IP 封禁 (支持 CIDR)
CREATE TABLE `mod_ip_ban` (
    `id`              INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `network`         VARCHAR(64)     NOT NULL COMMENT 'IP或CIDR，如 192.168.1.0/24',
    `reason`          VARCHAR(255)    NOT NULL DEFAULT '',
    `banned_by`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_network` (`network`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='IP 封禁表';




-- 魔力值/奖励交易通用记录
CREATE TABLE `economy_bonus_log` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `amount`          DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '正=获得 负=消耗',
    `balance_after`   DECIMAL(14,2)   NOT NULL DEFAULT 0.00 COMMENT '变动后余额',
    `action`          VARCHAR(30)     NOT NULL COMMENT 'torrent_reward/post_reward/daily_bonus/...',
    `target_type`     VARCHAR(20)     NOT NULL DEFAULT '',
    `target_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `period`          VARCHAR(50)     NULL     DEFAULT NULL COMMENT '结算周期/幂等键',
    `remark`          VARCHAR(500)    NOT NULL DEFAULT '',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_created` (`user_id`, `created_at`),
    KEY `idx_action` (`action`),
    UNIQUE KEY `uk_user_action_period` (`user_id`, `action`, `period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='魔力值交易日志';

CREATE TABLE `economy_shop_order` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `product_key`     VARCHAR(64)     NOT NULL COMMENT '商品标识快照',
    `product_type`    VARCHAR(32)     NOT NULL COMMENT '商品类型快照',
    `product_snapshot` JSON           NOT NULL COMMENT '购买时商品配置快照',
    `price`           DECIMAL(12,2)   NOT NULL DEFAULT 0.00 COMMENT '成交价格快照',
    `status`          TINYINT         NOT NULL DEFAULT 0 COMMENT '0=待处理, 1=已完成',
    `target_type`     VARCHAR(50)     NOT NULL DEFAULT '' COMMENT '履约目标类型',
    `target_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '履约目标ID',
    `created_at`      DATETIME        NULL,
    `completed_at`    DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_created` (`user_id`, `created_at`, `id`),
    KEY `idx_product_created` (`product_key`, `created_at`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='魔力商城兑换订单表';

CREATE TABLE `economy_reward_record` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `target_type`     VARCHAR(50)     NOT NULL COMMENT 'catalog_torrent/catalog_comment/forum_topic/forum_reply',
    `target_id`       BIGINT UNSIGNED NOT NULL,
    `from_user_id`    BIGINT UNSIGNED NOT NULL COMMENT '赞赏者ID',
    `to_user_id`      BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    `amount`          DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '赞赏金额',
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    KEY `idx_target_created` (`target_type`, `target_id`, `created_at`),
    KEY `idx_target_user` (`target_type`, `target_id`, `from_user_id`),
    KEY `idx_from_created` (`from_user_id`, `created_at`),
    KEY `idx_to_created` (`to_user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通用赞赏记录表';

SET FOREIGN_KEY_CHECKS = 1;

-- ----------------------------
-- Table structure for cron_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_cron_log`;
CREATE TABLE `sys_cron_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job_name` varchar(64) NOT NULL COMMENT '任务名称',
  `node_ip` varchar(64) NOT NULL COMMENT '执行节点IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0: 执行中 1: 成功 2: 失败',
  `duration_ms` int(11) NOT NULL DEFAULT '0' COMMENT '执行耗时(毫秒)',
  `error_message` text COMMENT '错误信息',
  `created_at` datetime NOT NULL COMMENT '开始时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_job_name` (`job_name`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务执行日志';


-- ============================================================
-- 模块: 评论互动
-- ============================================================

CREATE TABLE `catalog_torrent_like` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `torrent_id`      BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_torrent` (`user_id`, `torrent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='种子点赞表';

CREATE TABLE `catalog_comment_like` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT UNSIGNED NOT NULL,
    `comment_id`      BIGINT UNSIGNED NOT NULL,
    `created_at`      DATETIME        NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_comment` (`user_id`, `comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';
