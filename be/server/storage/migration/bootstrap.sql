SET NAMES utf8mb4;

-- ============================================================
-- NextPT development bootstrap data
-- Run after init.sql on an empty development database.
--
-- Default admin:
--   username: admin
--   password: admin123456
-- ============================================================

-- ------------------------------------------------------------
-- IAM roles and default admin user
-- ------------------------------------------------------------

INSERT INTO `iam_role`
    (`id`, `level`, `name_i18n`, `rules`, `permissions`, `is_staff`, `created_at`, `updated_at`)
VALUES
    (
        1,
        10,
        '{"zh-CN":"用户","zh-TW":"用戶","en-US":"User"}',
        '{}',
        '[
            "read:iam/invite:*",
            "create:iam/invite:*",
            "read:forum/topic:*",
            "create:forum/topic:*",
            "read:forum/reply:*",
            "create:forum/reply:*",
            "read:catalog/torrent:*",
            "create:catalog/torrent:*",
            "download:catalog/torrent:*",
            "read:catalog/subtitle:*",
            "create:catalog/subtitle:*",
            "download:catalog/subtitle:*",
            "read:catalog/comment:*",
            "create:catalog/comment:*"
        ]',
        b'0',
        NOW(),
        NOW()
    ),
    (
        2,
        100,
        '{"zh-CN":"管理员","zh-TW":"管理員","en-US":"Administrator"}',
        '{}',
        '["*"]',
        b'1',
        NOW(),
        NOW()
    )
ON DUPLICATE KEY UPDATE
    `level` = VALUES(`level`),
    `name_i18n` = VALUES(`name_i18n`),
    `rules` = VALUES(`rules`),
    `permissions` = VALUES(`permissions`),
    `is_staff` = VALUES(`is_staff`),
    `updated_at` = NOW();

INSERT INTO `iam_user`
    (`id`, `username`, `email`, `password_hash`, `passkey`, `status`, `role`, `created_at`, `updated_at`)
VALUES
    (
        1,
        'admin',
        'admin@nextpt.local',
        '$2a$10$tOc.CwWrLPz6ro2qGv5VmOEfwhD9YaxfXFBAzL3oloXfoUjlGcs6a',
        '0123456789abcdef0123456789abcdef',
        1,
        2,
        NOW(),
        NOW()
    )
ON DUPLICATE KEY UPDATE
    `username` = VALUES(`username`),
    `email` = VALUES(`email`),
    `password_hash` = VALUES(`password_hash`),
    `passkey` = VALUES(`passkey`),
    `status` = VALUES(`status`),
    `role` = VALUES(`role`),
    `updated_at` = NOW();

INSERT INTO `iam_user_profile`
    (`id`, `user_id`, `avatar`, `info`, `signature`, `created_at`, `updated_at`)
VALUES
    (1, 1, '', 'NextPT bootstrap administrator.', '', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `avatar` = VALUES(`avatar`),
    `info` = VALUES(`info`),
    `signature` = VALUES(`signature`),
    `updated_at` = NOW();

INSERT INTO `iam_user_stat`
    (`id`, `user_id`, `uploaded`, `downloaded`, `seed_time`, `leech_time`, `bonus`, `bonus_charity`, `created_at`, `updated_at`)
VALUES
    (1, 1, 0, 0, 0, 0, 0.0, 0.0, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `user_id` = VALUES(`user_id`);

INSERT INTO `iam_user_setting`
    (`id`, `user_id`, `privacy_level`, `extra`, `created_at`, `updated_at`)
VALUES
    (1, 1, 1, '{}', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `privacy_level` = VALUES(`privacy_level`),
    `extra` = VALUES(`extra`),
    `updated_at` = NOW();

-- ------------------------------------------------------------
-- Site config defaults
-- ------------------------------------------------------------

INSERT INTO `site_config`
    (`group`, `key`, `value`, `description`, `created_at`, `updated_at`)
VALUES
    ('tracker', 'announce_interval', '{"val":1800}', '客户端心跳汇报间隔（秒）', NOW(), NOW()),
    ('tracker', 'announce_min_interval', '{"val":900}', '客户端心跳最小汇报间隔（秒）', NOW(), NOW()),
    ('tracker', 'announce_url', '{"val":"http://127.0.0.1:8000/api/tracker/announce"}', 'Tracker 宣告地址', NOW(), NOW()),
    ('tracker', 'bonus_T0', '{"val":8.0}', '魔力值衰减参数 T0', NOW(), NOW()),
    ('tracker', 'bonus_N0', '{"val":7.0}', '魔力值做种人数拥挤惩罚基数 N0', NOW(), NOW()),
    ('tracker', 'bonus_B0', '{"val":100.0}', '魔力值每小时体积收益硬上限 B0', NOW(), NOW()),
    ('tracker', 'bonus_L', '{"val":300.0}', '魔力值收益收敛平滑参数 L', NOW(), NOW()),
    ('tracker', 'bonus_base', '{"val":0.4}', '每个达标种子的基础每小时奖励', NOW(), NOW()),
    ('iam', 'default_role', '{"val":"user"}', '新注册用户默认角色标识', NOW(), NOW()),
    ('iam', 'default_register_role', '{"val":1}', '新注册用户默认角色ID', NOW(), NOW()),
    ('iam', 'register_enabled', '{"val":true}', '是否开放全站注册', NOW(), NOW()),
    ('catalog', 'torrent_source', '{"val":"NextPT"}', '私有种子 source 标识', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `value` = VALUES(`value`),
    `description` = VALUES(`description`),
    `updated_at` = NOW();

-- ------------------------------------------------------------
-- Tracker client whitelist
-- ------------------------------------------------------------

INSERT INTO `tracker_agent_whitelist`
    (`id`, `family`, `peer_id_prefix`, `agent_pattern`, `min_version`, `max_version`, `allow_https`, `enabled`, `comment`, `created_at`, `updated_at`)
VALUES
    (1, 'qBittorrent', '-qB', '^qBittorrent/', '', '', b'1', b'1', 'Default qBittorrent client rule', NOW(), NOW()),
    (2, 'Transmission', '-TR', '^Transmission/', '', '', b'1', b'1', 'Default Transmission client rule', NOW(), NOW()),
    (3, 'Deluge', '-DE', '^Deluge', '', '', b'1', b'1', 'Default Deluge client rule', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `family` = VALUES(`family`),
    `peer_id_prefix` = VALUES(`peer_id_prefix`),
    `agent_pattern` = VALUES(`agent_pattern`),
    `min_version` = VALUES(`min_version`),
    `max_version` = VALUES(`max_version`),
    `allow_https` = VALUES(`allow_https`),
    `enabled` = VALUES(`enabled`),
    `comment` = VALUES(`comment`),
    `updated_at` = NOW();

-- ------------------------------------------------------------
-- Catalog categories
-- ------------------------------------------------------------

INSERT INTO `catalog_category`
    (`id`, `name_i18n`, `slug`, `sort_order`, `enabled`, `created_at`, `updated_at`)
VALUES
    (1, '{"zh-CN":"电影","zh-TW":"電影","en-US":"Movies"}', 'movies', 10, b'1', NOW(), NOW()),
    (2, '{"zh-CN":"剧集","zh-TW":"劇集","en-US":"TV"}', 'tv', 20, b'1', NOW(), NOW()),
    (3, '{"zh-CN":"音乐","zh-TW":"音樂","en-US":"Music"}', 'music', 30, b'1', NOW(), NOW()),
    (4, '{"zh-CN":"软件","zh-TW":"軟體","en-US":"Software"}', 'software', 40, b'1', NOW(), NOW()),
    (5, '{"zh-CN":"其他","zh-TW":"其他","en-US":"Other"}', 'other', 90, b'1', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `name_i18n` = VALUES(`name_i18n`),
    `slug` = VALUES(`slug`),
    `sort_order` = VALUES(`sort_order`),
    `enabled` = VALUES(`enabled`),
    `updated_at` = NOW();

-- ------------------------------------------------------------
-- Forum categories and nodes
-- ------------------------------------------------------------

INSERT INTO `forum_category`
    (`id`, `name_i18n`, `desc_i18n`, `sort_order`, `min_role_view`, `created_at`, `updated_at`)
VALUES
    (
        1,
        '{"zh-CN":"站点","zh-TW":"站點","en-US":"Site"}',
        '{"zh-CN":"站点公告和综合讨论","zh-TW":"站點公告和綜合討論","en-US":"Announcements and general discussion"}',
        10,
        0,
        NOW(),
        NOW()
    )
ON DUPLICATE KEY UPDATE
    `name_i18n` = VALUES(`name_i18n`),
    `desc_i18n` = VALUES(`desc_i18n`),
    `sort_order` = VALUES(`sort_order`),
    `min_role_view` = VALUES(`min_role_view`),
    `updated_at` = NOW();

INSERT INTO `forum_node`
    (`id`, `category_id`, `slug`, `name_i18n`, `desc_i18n`, `sort_order`, `min_role_read`, `min_role_write`, `min_role_create`, `moderators`, `created_at`, `updated_at`)
VALUES
    (
        1,
        1,
        'announcements',
        '{"zh-CN":"公告","zh-TW":"公告","en-US":"Announcements"}',
        '{"zh-CN":"站点公告和规则说明","zh-TW":"站點公告和規則說明","en-US":"Site announcements and rules"}',
        10,
        0,
        100,
        100,
        '[1]',
        NOW(),
        NOW()
    ),
    (
        2,
        1,
        'general',
        '{"zh-CN":"综合讨论","zh-TW":"綜合討論","en-US":"General"}',
        '{"zh-CN":"日常交流和问题讨论","zh-TW":"日常交流和問題討論","en-US":"General discussion and questions"}',
        20,
        0,
        0,
        0,
        '[]',
        NOW(),
        NOW()
    )
ON DUPLICATE KEY UPDATE
    `category_id` = VALUES(`category_id`),
    `slug` = VALUES(`slug`),
    `name_i18n` = VALUES(`name_i18n`),
    `desc_i18n` = VALUES(`desc_i18n`),
    `sort_order` = VALUES(`sort_order`),
    `min_role_read` = VALUES(`min_role_read`),
    `min_role_write` = VALUES(`min_role_write`),
    `min_role_create` = VALUES(`min_role_create`),
    `moderators` = VALUES(`moderators`),
    `updated_at` = NOW();
