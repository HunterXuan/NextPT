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

SET @iam_permissions_limited = JSON_ARRAY(
    'read:site/announcement:*',
    'read:site/message:*',
    'read:mod/staff-message:*',
    'create:mod/staff-message:*',
    'read:economy/shop-product:*',
    'read:forum/topic:*',
    'read:forum/reply:*',
    'read:catalog/torrent:*',
    'download:catalog/torrent:*',
    'read:catalog/subtitle:*',
    'download:catalog/subtitle:*',
    'read:catalog/comment:*',
    'read:catalog/request:*'
);

SET @iam_permissions_member = JSON_ARRAY_APPEND(
    @iam_permissions_limited,
    '$',
    'read:economy/shop-order:*',
    '$',
    'create:economy/shop-order:*',
    '$',
    'read:iam/invite:*',
    '$',
    'create:iam/invite:*',
    '$',
    'create:forum/topic:*',
    '$',
    'create:forum/reply:*',
    '$',
    'create:catalog/comment:*',
    '$',
    'create:catalog/request:*'
);

SET @iam_permissions_power_user = JSON_ARRAY_APPEND(
    @iam_permissions_member,
    '$',
    'create:catalog/torrent:*',
    '$',
    'create:catalog/subtitle:*'
);

SET @iam_permissions_forum_staff = JSON_ARRAY_APPEND(
    @iam_permissions_power_user,
    '$',
    'admin:mod/report:*',
    '$',
    'admin:forum/topic:*',
    '$',
    'admin:forum/reply:*'
);

SET @iam_permissions_catalog_staff = JSON_ARRAY_APPEND(
    @iam_permissions_forum_staff,
    '$',
    'admin:catalog/torrent:*',
    '$',
    'admin:catalog/subtitle:*',
    '$',
    'admin:catalog/comment:*',
    '$',
    'admin:catalog/request:*'
);

SET @iam_permissions_admin = JSON_ARRAY_APPEND(
    @iam_permissions_catalog_staff,
    '$',
    'admin:iam/user:*',
    '$',
    'admin:iam/role:*',
    '$',
    'admin:iam/invite:*',
    '$',
    'admin:mod/cheater:*',
    '$',
    'admin:mod/user:*',
    '$',
    'admin:mod/staff-message:*',
    '$',
    'admin:site/config:*',
    '$',
    'admin:site/audit:*',
    '$',
    'admin:site/announcement:*',
    '$',
    'admin:site/message:*',
    '$',
    'admin:forum/category:*',
    '$',
    'admin:forum/node:*',
    '$',
    'admin:catalog/category:*',
    '$',
    'admin:catalog/tag:*'
);

SET @iam_permissions_sysop = JSON_ARRAY_APPEND(
    @iam_permissions_admin,
    '$',
    'admin:sys/cron:*'
);

SET @iam_permissions_all = JSON_ARRAY('*');

INSERT INTO `iam_role`
    (`id`, `level`, `name_i18n`, `rules`, `permissions`, `is_staff`, `created_at`, `updated_at`)
VALUES
    (1, 0, '{"zh-CN":"留校察看","zh-TW":"留校察看","en-US":"Peasant"}', JSON_OBJECT('demotion',JSON_ARRAY(JSON_OBJECT('downloadedGiBGt',30,'ratioLt',0.4),JSON_OBJECT('downloadedGiBGt',100,'ratioLt',0.5),JSON_OBJECT('downloadedGiBGt',200,'ratioLt',0.6),JSON_OBJECT('downloadedGiBGt',400,'ratioLt',0.7),JSON_OBJECT('downloadedGiBGt',800,'ratioLt',0.8))), @iam_permissions_limited, b'0', NOW(), NOW()),
    (2, 10, '{"zh-CN":"本科新生","zh-TW":"本科新生","en-US":"User"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('downloadedGiBLte',30),JSON_OBJECT('downloadedGiBGt',30,'downloadedGiBLte',100,'ratioGte',0.4),JSON_OBJECT('downloadedGiBGt',100,'downloadedGiBLte',200,'ratioGte',0.5),JSON_OBJECT('downloadedGiBGt',200,'downloadedGiBLte',400,'ratioGte',0.6),JSON_OBJECT('downloadedGiBGt',400,'downloadedGiBLte',800,'ratioGte',0.7),JSON_OBJECT('downloadedGiBGt',800,'ratioGte',0.8)),'demotion',JSON_ARRAY(JSON_OBJECT('downloadedGiBGt',30,'ratioLt',0.4),JSON_OBJECT('downloadedGiBGt',100,'ratioLt',0.5),JSON_OBJECT('downloadedGiBGt',200,'ratioLt',0.6),JSON_OBJECT('downloadedGiBGt',400,'ratioLt',0.7),JSON_OBJECT('downloadedGiBGt',800,'ratioLt',0.8))), @iam_permissions_member, b'0', NOW(), NOW()),
    (3, 20, '{"zh-CN":"小小学士","zh-TW":"小小學士","en-US":"PowerUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',14,'downloadedGiBGte',30,'ratioGt',1.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',1.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (4, 30, '{"zh-CN":"优秀硕士","zh-TW":"優秀碩士","en-US":"EliteUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',35,'downloadedGiBGte',50,'ratioGt',2.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',2.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (5, 40, '{"zh-CN":"初为博士","zh-TW":"初為博士","en-US":"CrazyUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',70,'downloadedGiBGte',100,'ratioGt',3.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',3.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (6, 50, '{"zh-CN":"海归博后","zh-TW":"海歸博後","en-US":"InsaneUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',105,'downloadedGiBGte',300,'ratioGt',4.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',4.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (7, 60, '{"zh-CN":"大学讲师","zh-TW":"大學講師","en-US":"VeteranUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',140,'downloadedGiBGte',500,'ratioGt',5.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',5.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (8, 70, '{"zh-CN":"晋升副教","zh-TW":"晉升副教","en-US":"ExtremeUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',210,'downloadedGiBGte',700,'ratioGt',6.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',6.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (9, 80, '{"zh-CN":"终身教授","zh-TW":"終身教授","en-US":"UltimateUser"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',280,'downloadedGiBGte',900,'ratioGt',7.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',7.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (10, 90, '{"zh-CN":"荣誉院士","zh-TW":"榮譽院士","en-US":"NexusMaster"}', JSON_OBJECT('promotion',JSON_ARRAY(JSON_OBJECT('accountAgeDaysGte',350,'downloadedGiBGte',1024,'ratioGt',8.5)),'demotion',JSON_ARRAY(JSON_OBJECT('ratioLt',8.4))), @iam_permissions_power_user, b'0', NOW(), NOW()),
    (11, 100, '{"zh-CN":"养老族","zh-TW":"養老族","en-US":"Retiree"}', '{}', @iam_permissions_power_user, b'1', NOW(), NOW()),
    (12, 110, '{"zh-CN":"保种分流员","zh-TW":"保種分流員","en-US":"Seeder"}', '{}', @iam_permissions_power_user, b'1', NOW(), NOW()),
    (13, 120, '{"zh-CN":"发布员","zh-TW":"發布員","en-US":"Uploader"}', '{}', @iam_permissions_power_user, b'1', NOW(), NOW()),
    (14, 130, '{"zh-CN":"论坛版主","zh-TW":"論壇版主","en-US":"ForumModerator"}', '{}', @iam_permissions_forum_staff, b'1', NOW(), NOW()),
    (15, 140, '{"zh-CN":"种子管理员","zh-TW":"種子管理員","en-US":"CatalogModerator"}', '{}', @iam_permissions_catalog_staff, b'1', NOW(), NOW()),
    (16, 150, '{"zh-CN":"高级管理员","zh-TW":"高級管理員","en-US":"Administrator"}', '{}', @iam_permissions_admin, b'1', NOW(), NOW()),
    (17, 160, '{"zh-CN":"维护开发员","zh-TW":"維護開發員","en-US":"Sysop"}', '{}', @iam_permissions_sysop, b'1', NOW(), NOW()),
    (18, 170, '{"zh-CN":"主管","zh-TW":"主管","en-US":"StaffLeader"}', '{}', @iam_permissions_all, b'1', NOW(), NOW())
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
        18,
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
    (`id`, `user_id`, `uploaded`, `downloaded`, `raw_uploaded`, `raw_downloaded`, `seed_time`, `leech_time`, `bonus`, `bonus_charity`, `created_at`, `updated_at`)
VALUES
    (1, 1, 0, 0, 0, 0, 0, 0, 0.0, 0.0, NOW(), NOW())
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

DELETE FROM `site_config`
WHERE `group` = 'iam' AND `key` IN ('default_role', 'invite_register_email_pattern');

INSERT INTO `site_config`
    (`group`, `key`, `value`, `created_at`, `updated_at`)
VALUES
    ('tracker', 'announce_interval', '{"val":1800}', NOW(), NOW()),
    ('tracker', 'announce_min_interval', '{"val":900}', NOW(), NOW()),
    ('tracker', 'announce_url', '{"val":"http://127.0.0.1:8000/api/tracker/announce"}', NOW(), NOW()),
    ('tracker', 'bonus_T0', '{"val":8.0}', NOW(), NOW()),
    ('tracker', 'bonus_N0', '{"val":7.0}', NOW(), NOW()),
    ('tracker', 'bonus_B0', '{"val":100.0}', NOW(), NOW()),
    ('tracker', 'bonus_L', '{"val":300.0}', NOW(), NOW()),
    ('tracker', 'bonus_base', '{"val":0.4}', NOW(), NOW()),
    ('iam', 'default_register_role', '{"val":2}', NOW(), NOW()),
    ('iam', 'register_enabled', '{"val":true}', NOW(), NOW()),
    ('iam', 'invite_bypass_email_pattern', '{"val":""}', NOW(), NOW()),
    ('site', 'maintenance_enabled', '{"val":false}', NOW(), NOW()),
    ('site', 'maintenance_message', '{"val":""}', NOW(), NOW()),
    ('site', 'tasks', JSON_OBJECT('val', JSON_ARRAY()), NOW(), NOW()),
    ('site', 'advertisements', JSON_OBJECT('val', JSON_OBJECT(
        'home', JSON_OBJECT('enabled', false, 'title', '', 'image', '', 'url', ''),
        'catalog_list', JSON_OBJECT('enabled', false, 'title', '', 'image', '', 'url', ''),
        'forum_list', JSON_OBJECT('enabled', false, 'title', '', 'image', '', 'url', '')
    )), NOW(), NOW()),
    ('economy', 'shop_products', JSON_OBJECT('val', JSON_ARRAY(
        JSON_OBJECT(
            'key', 'invite',
            'type', 'invite',
            'enabled', true,
            'price', 1000.00,
            'sortOrder', 10,
            'options', JSON_OBJECT('amount', 1)
        ),
        JSON_OBJECT(
            'key', 'upload_100_gib',
            'type', 'upload',
            'enabled', true,
            'price', 500.00,
            'sortOrder', 20,
            'options', JSON_OBJECT('amountGiB', 100)
        ),
        JSON_OBJECT(
            'key', 'download_50_gib',
            'type', 'download',
            'enabled', true,
            'price', 800.00,
            'sortOrder', 30,
            'options', JSON_OBJECT('amountGiB', 50)
        ),
        JSON_OBJECT(
            'key', 'vip_30d',
            'type', 'vip',
            'enabled', true,
            'price', 3000.00,
            'sortOrder', 40,
            'options', JSON_OBJECT('durationDays', 30)
        )
    )), NOW(), NOW()),
    ('catalog', 'torrent_source', '{"val":"NextPT"}', NOW(), NOW()),
    ('catalog', 'torrent_direct_publish_level', '{"val":20}', NOW(), NOW()),
    ('catalog', 'global_promotion', JSON_OBJECT('val', JSON_OBJECT(
        'enabled', false,
        'state', 'free',
        'expireAt', ''
    )), NOW(), NOW()),
    ('catalog', 'new_torrent_promotion', JSON_OBJECT('val', JSON_OBJECT(
        'enabled', true,
        'rules', JSON_ARRAY(
            JSON_OBJECT(
                'minGiB', 0,
                'durationHours', 72,
                'options', JSON_ARRAY(
                    JSON_OBJECT('state', 'normal', 'weight', 70),
                    JSON_OBJECT('state', 'free', 'weight', 10),
                    JSON_OBJECT('state', '2x', 'weight', 10),
                    JSON_OBJECT('state', '50_percent', 'weight', 10)
                )
            ),
            JSON_OBJECT(
                'minGiB', 10,
                'durationHours', 96,
                'options', JSON_ARRAY(
                    JSON_OBJECT('state', 'normal', 'weight', 50),
                    JSON_OBJECT('state', 'free', 'weight', 20),
                    JSON_OBJECT('state', '2x', 'weight', 20),
                    JSON_OBJECT('state', '2x_free', 'weight', 10)
                )
            ),
            JSON_OBJECT(
                'minGiB', 50,
                'durationHours', 168,
                'options', JSON_ARRAY(
                    JSON_OBJECT('state', 'free', 'weight', 40),
                    JSON_OBJECT('state', '2x_free', 'weight', 30),
                    JSON_OBJECT('state', '2x_50_percent', 'weight', 30)
                )
            )
        )
    )), NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `value` = VALUES(`value`),
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

DELETE FROM `catalog_category`;

INSERT INTO `catalog_category`
    (`id`, `name_i18n`, `slug`, `sort_order`, `enabled`, `upload_config`, `created_at`, `updated_at`)
VALUES
    (1, '{"zh-CN":"电影","zh-TW":"電影","en-US":"Movies"}', 'movies', 10, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"movie_enname","suffix":" "}]},"fields":[{"key":"movie_enname","type":"text","label":{"zh-CN":"英文名(0day)","zh-TW":"英文名(0day)","en-US":"0day title"},"description":{"zh-CN":"请填写标准英文0day名。","zh-TW":"請填寫標準英文0day名。","en-US":"Use the standard 0day release title."},"required":true,"placeholder":{"zh-CN":"例如 Blade.Runner.1982.Final.Cut.720p.HDDVD.DTS.x264-ESiR","zh-TW":"例如 Blade.Runner.1982.Final.Cut.720p.HDDVD.DTS.x264-ESiR","en-US":"e.g. Blade.Runner.1982.Final.Cut.720p.HDDVD.DTS.x264-ESiR"}}]}', NOW(), NOW()),
    (2, '{"zh-CN":"剧集","zh-TW":"劇集","en-US":"TV"}', 'tv', 20, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"series_enname","suffix":" "}]},"fields":[{"key":"series_enname","type":"text","label":{"zh-CN":"英文名","zh-TW":"英文名","en-US":"0day title"},"description":{"zh-CN":"请填写完整0day名。","zh-TW":"請填寫完整0day名。","en-US":"Use the complete 0day release title."},"required":true,"placeholder":{"zh-CN":"例如 Madam.Secretary.S01E02.720p.HDTV.X264-DIMENSION","zh-TW":"例如 Madam.Secretary.S01E02.720p.HDTV.X264-DIMENSION","en-US":"e.g. Madam.Secretary.S01E02.720p.HDTV.X264-DIMENSION"}}]}', NOW(), NOW()),
    (3, '{"zh-CN":"动漫","zh-TW":"動漫","en-US":"Anime"}', 'anime', 30, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"anime_enname","prefix":"[","suffix":"]"},{"field":"anime_serial","prefix":"[","suffix":"]"},{"field":"anime_source","prefix":"[","suffix":"]"},{"field":"anime_subtitle","prefix":"[","suffix":"]"},{"field":"anime_resolution","prefix":"[","suffix":"]"},{"field":"anime_format","prefix":"[","suffix":"]"},{"field":"anime_group","prefix":"[","suffix":"]"}]},"fields":[{"key":"anime_enname","type":"text","label":{"zh-CN":"英文名","zh-TW":"英文名","en-US":"Title"},"description":{"zh-CN":"请填写罗马音名或英文译名。","zh-TW":"請填寫羅馬音名或英文譯名。","en-US":"Use the romanized or English title."},"required":true},{"key":"anime_serial","type":"text","label":{"zh-CN":"集数/集别/卷数/张数","zh-TW":"集數/集別/卷數/張數","en-US":"Episode / volume"},"description":{"zh-CN":"注明 OVA/OAD/MOVIE；合集可写 TV 01-XX Fin +SP/OVA/MOVIE 或 Vol.1-Vol.X Fin。","zh-TW":"註明 OVA/OAD/MOVIE；合集可寫 TV 01-XX Fin +SP/OVA/MOVIE 或 Vol.1-Vol.X Fin。","en-US":"Episode, volume, OVA/OAD/MOVIE, or collection range."},"required":false},{"key":"anime_source","type":"text","label":{"zh-CN":"片源/出版社","zh-TW":"片源/出版社","en-US":"Source / publisher"},"description":{"zh-CN":"片源如 TVRip、BDRip、WEB；出版社如角川、天下；音乐可填写专辑艺人。","zh-TW":"片源如 TVRip、BDRip、WEB；出版社如角川、天下；音樂可填寫專輯藝人。","en-US":"Source such as TVRip, BDRip, WEB, publisher, or album artist."},"required":false},{"key":"anime_subtitle","type":"text","label":{"zh-CN":"字幕类型","zh-TW":"字幕類型","en-US":"Subtitle type"},"description":{"zh-CN":"例如 GB、CHS、BIG5、CHT、简繁外挂等。","zh-TW":"例如 GB、CHS、BIG5、CHT、簡繁外掛等。","en-US":"Examples: GB, CHS, BIG5, CHT, external subtitles."},"required":false},{"key":"anime_resolution","type":"text","label":{"zh-CN":"分辨率/码率/扫图者","zh-TW":"解析度/碼率/掃圖者","en-US":"Resolution / bitrate / scanner"},"description":{"zh-CN":"分辨率正常填写；码率如 320K；扫图者如 C.C、HMM。","zh-TW":"解析度正常填寫；碼率如 320K；掃圖者如 C.C、HMM。","en-US":"Resolution, bitrate such as 320K, or scanner name."},"required":false},{"key":"anime_format","type":"text","label":{"zh-CN":"动漫格式","zh-TW":"動漫格式","en-US":"Format"},"description":{"zh-CN":"动画类正常填写；漫画可写 zip、rar；音乐可写 TAK+CUE+ISO+BK 等。","zh-TW":"動畫類正常填寫；漫畫可寫 zip、rar；音樂可寫 TAK+CUE+ISO+BK 等。","en-US":"Video format, archive format, or music package format."},"required":false},{"key":"anime_group","type":"text","label":{"zh-CN":"字幕组/压制组/作者/EAC","zh-TW":"字幕組/壓制組/作者/EAC","en-US":"Group / author"},"description":{"zh-CN":"字幕或压制填写组名；漫画填写作者；音乐可填写 EAC。","zh-TW":"字幕或壓制填寫組名；漫畫填寫作者；音樂可填寫 EAC。","en-US":"Subtitle group, encode group, author, or EAC note."},"required":false}]}', NOW(), NOW()),
    (4, '{"zh-CN":"综艺","zh-TW":"綜藝","en-US":"Variety"}', 'variety', 40, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"show_enname","suffix":" "}]},"fields":[{"key":"show_enname","type":"text","label":{"zh-CN":"英文名","zh-TW":"英文名","en-US":"0day title"},"description":{"zh-CN":"请填写标准英文0day名。","zh-TW":"請填寫標準英文0day名。","en-US":"Use the standard 0day release title."},"required":true}]}', NOW(), NOW()),
    (5, '{"zh-CN":"体育","zh-TW":"體育","en-US":"Sports"}', 'sports', 50, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"sports_type","prefix":"[","suffix":"]"},{"field":"sports_time","prefix":"[","suffix":"]"},{"field":"sports_chname","prefix":"[","suffix":"]"},{"field":"sports_lang","prefix":"[","suffix":"]"},{"field":"sports_form","prefix":"[","suffix":"]"},{"field":"sports_solu","prefix":"[","suffix":"]"}]},"fields":[{"key":"sports_type","type":"select","label":{"zh-CN":"节目类型","zh-TW":"節目類型","en-US":"Program type"},"required":true,"options":{"source":"static","items":[{"value":"足球","label":{"zh-CN":"足球","zh-TW":"足球","en-US":"Football"}},{"value":"篮球","label":{"zh-CN":"篮球","zh-TW":"籃球","en-US":"Basketball"}},{"value":"网球","label":{"zh-CN":"网球","zh-TW":"網球","en-US":"Tennis"}},{"value":"台球","label":{"zh-CN":"台球","zh-TW":"撞球","en-US":"Billiards"}},{"value":"棒球","label":{"zh-CN":"棒球","zh-TW":"棒球","en-US":"Baseball"}},{"value":"羽毛球","label":{"zh-CN":"羽毛球","zh-TW":"羽毛球","en-US":"Badminton"}},{"value":"F1","label":{"zh-CN":"F1","zh-TW":"F1","en-US":"F1"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"sports_time","type":"text","label":{"zh-CN":"发行时间","zh-TW":"發行時間","en-US":"Release date"},"description":{"zh-CN":"标准格式：2010-10-01。","zh-TW":"標準格式：2010-10-01。","en-US":"Use YYYY-MM-DD, such as 2010-10-01."},"required":true,"placeholder":{"zh-CN":"2010-10-01","zh-TW":"2010-10-01","en-US":"2010-10-01"}},{"key":"sports_chname","type":"text","label":{"zh-CN":"中文名","zh-TW":"中文名","en-US":"Chinese title"},"description":{"zh-CN":"“比赛名称 参赛方A VS 参赛方B”，如 NBA常规赛 火箭 VS 湖人。","zh-TW":"「比賽名稱 參賽方A VS 參賽方B」，如 NBA常規賽 火箭 VS 湖人。","en-US":"Match name and teams, such as NBA regular season Rockets VS Lakers."},"required":true},{"key":"sports_lang","type":"text","label":{"zh-CN":"解说语言","zh-TW":"解說語言","en-US":"Commentary language"},"description":{"zh-CN":"如 CCTV/国语。","zh-TW":"如 CCTV/國語。","en-US":"Example: CCTV / Mandarin."},"required":true},{"key":"sports_form","type":"select","label":{"zh-CN":"文件格式","zh-TW":"檔案格式","en-US":"File format"},"required":true,"options":{"source":"static","items":[{"value":"AVI","label":{"zh-CN":"AVI","zh-TW":"AVI","en-US":"AVI"}},{"value":"FLV","label":{"zh-CN":"FLV","zh-TW":"FLV","en-US":"FLV"}},{"value":"MP4","label":{"zh-CN":"MP4","zh-TW":"MP4","en-US":"MP4"}},{"value":"MKV","label":{"zh-CN":"MKV","zh-TW":"MKV","en-US":"MKV"}},{"value":"WMV","label":{"zh-CN":"WMV","zh-TW":"WMV","en-US":"WMV"}},{"value":"RMVB","label":{"zh-CN":"RMVB","zh-TW":"RMVB","en-US":"RMVB"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"sports_solu","type":"select","label":{"zh-CN":"录像分辨率","zh-TW":"錄像解析度","en-US":"Video resolution"},"required":true,"options":{"source":"static","items":[{"value":"720I","label":{"zh-CN":"720I","zh-TW":"720I","en-US":"720I"}},{"value":"720P","label":{"zh-CN":"720P","zh-TW":"720P","en-US":"720P"}},{"value":"1080I","label":{"zh-CN":"1080I","zh-TW":"1080I","en-US":"1080I"}},{"value":"1080P","label":{"zh-CN":"1080P","zh-TW":"1080P","en-US":"1080P"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}}]}', NOW(), NOW()),
    (6, '{"zh-CN":"纪录","zh-TW":"紀錄","en-US":"Documentary"}', 'documentary', 60, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"doc_enname","suffix":" "}]},"fields":[{"key":"doc_enname","type":"text","label":{"zh-CN":"英文名","zh-TW":"英文名","en-US":"0day title"},"description":{"zh-CN":"请填写英文名0day名。","zh-TW":"請填寫英文名0day名。","en-US":"Use the 0day English release title."},"required":true}]}', NOW(), NOW()),
    (7, '{"zh-CN":"音乐","zh-TW":"音樂","en-US":"Music"}', 'music', 70, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"music_time","prefix":"[","suffix":"]"},{"field":"music_type","prefix":"[","suffix":"]"},{"field":"music_author","prefix":"[","suffix":"]"},{"field":"music_name","prefix":"[","suffix":"]"},{"field":"music_form","prefix":"[","suffix":"]"},{"field":"music_qual","prefix":"[","suffix":"]"}]},"fields":[{"key":"music_time","type":"text","label":{"zh-CN":"年份","zh-TW":"年份","en-US":"Date"},"description":{"zh-CN":"标准格式：2010-10-01。","zh-TW":"標準格式：2010-10-01。","en-US":"Use YYYY-MM-DD, such as 2010-10-01."},"required":true,"placeholder":{"zh-CN":"2010-10-01","zh-TW":"2010-10-01","en-US":"2010-10-01"}},{"key":"music_type","type":"select","label":{"zh-CN":"类型","zh-TW":"類型","en-US":"Type"},"required":true,"options":{"source":"static","items":[{"value":"单曲","label":{"zh-CN":"单曲","zh-TW":"單曲","en-US":"Single"}},{"value":"专辑","label":{"zh-CN":"专辑","zh-TW":"專輯","en-US":"Album"}},{"value":"合集","label":{"zh-CN":"合集","zh-TW":"合集","en-US":"Collection"}},{"value":"MV","label":{"zh-CN":"MV","zh-TW":"MV","en-US":"MV"}},{"value":"演唱会","label":{"zh-CN":"演唱会","zh-TW":"演唱會","en-US":"Live concert"}},{"value":"音乐会","label":{"zh-CN":"音乐会","zh-TW":"音樂會","en-US":"Concert"}},{"value":"戏剧","label":{"zh-CN":"戏剧","zh-TW":"戲劇","en-US":"Drama"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"music_author","type":"text","label":{"zh-CN":"作者","zh-TW":"作者","en-US":"Artist"},"description":{"zh-CN":"若为多人，请填写“群星”。","zh-TW":"若為多人，請填寫「群星」。","en-US":"Use “群星” for various artists."},"required":true},{"key":"music_name","type":"text","label":{"zh-CN":"名称","zh-TW":"名稱","en-US":"Title"},"required":true},{"key":"music_form","type":"select","label":{"zh-CN":"格式","zh-TW":"格式","en-US":"Format"},"required":true,"options":{"source":"static","items":[{"value":"720I","label":{"zh-CN":"720I","zh-TW":"720I","en-US":"720I"}},{"value":"MP3","label":{"zh-CN":"MP3","zh-TW":"MP3","en-US":"MP3"}},{"value":"AAC","label":{"zh-CN":"AAC","zh-TW":"AAC","en-US":"AAC"}},{"value":"APE","label":{"zh-CN":"APE","zh-TW":"APE","en-US":"APE"}},{"value":"FLAC","label":{"zh-CN":"FLAC","zh-TW":"FLAC","en-US":"FLAC"}},{"value":"WAV","label":{"zh-CN":"WAV","zh-TW":"WAV","en-US":"WAV"}},{"value":"VOB","label":{"zh-CN":"VOB","zh-TW":"VOB","en-US":"VOB"}},{"value":"AVI","label":{"zh-CN":"AVI","zh-TW":"AVI","en-US":"AVI"}},{"value":"MP4","label":{"zh-CN":"MP4","zh-TW":"MP4","en-US":"MP4"}},{"value":"MKV","label":{"zh-CN":"MKV","zh-TW":"MKV","en-US":"MKV"}},{"value":"MPG","label":{"zh-CN":"MPG","zh-TW":"MPG","en-US":"MPG"}},{"value":"RMVB","label":{"zh-CN":"RMVB","zh-TW":"RMVB","en-US":"RMVB"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"music_qual","type":"select","label":{"zh-CN":"质量","zh-TW":"品質","en-US":"Quality"},"required":true,"options":{"source":"static","items":[{"value":"192Kbps","label":{"zh-CN":"192Kbps","zh-TW":"192Kbps","en-US":"192Kbps"}},{"value":"256Kbps","label":{"zh-CN":"256Kbps","zh-TW":"256Kbps","en-US":"256Kbps"}},{"value":"320Kbps","label":{"zh-CN":"320Kbps","zh-TW":"320Kbps","en-US":"320Kbps"}},{"value":"VBR","label":{"zh-CN":"VBR","zh-TW":"VBR","en-US":"VBR"}},{"value":"无损","label":{"zh-CN":"无损","zh-TW":"無損","en-US":"Lossless"}},{"value":"720I","label":{"zh-CN":"720I","zh-TW":"720I","en-US":"720I"}},{"value":"720P","label":{"zh-CN":"720P","zh-TW":"720P","en-US":"720P"}},{"value":"1080I","label":{"zh-CN":"1080I","zh-TW":"1080I","en-US":"1080I"}},{"value":"1080P","label":{"zh-CN":"1080P","zh-TW":"1080P","en-US":"1080P"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}}]}', NOW(), NOW()),
    (8, '{"zh-CN":"学习","zh-TW":"學習","en-US":"Study"}', 'study', 80, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"study_type","prefix":"[","suffix":"]"},{"field":"study_name","prefix":"[","suffix":"]"},{"field":"study_format","prefix":"[","suffix":"]"}]},"fields":[{"key":"study_type","type":"select","label":{"zh-CN":"类型","zh-TW":"類型","en-US":"Type"},"required":true,"options":{"source":"static","items":[{"value":"音频","label":{"zh-CN":"音频","zh-TW":"音訊","en-US":"Audio"}},{"value":"视频","label":{"zh-CN":"视频","zh-TW":"影片","en-US":"Video"}},{"value":"文档","label":{"zh-CN":"文档","zh-TW":"文件","en-US":"Document"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"study_name","type":"text","label":{"zh-CN":"名称","zh-TW":"名稱","en-US":"Title"},"required":true},{"key":"study_format","type":"text","label":{"zh-CN":"文件格式","zh-TW":"檔案格式","en-US":"File format"},"description":{"zh-CN":"如 MP3、MP4、MKV、PDF、RAR 等。","zh-TW":"如 MP3、MP4、MKV、PDF、RAR 等。","en-US":"Examples: MP3, MP4, MKV, PDF, RAR."},"required":true}]}', NOW(), NOW()),
    (9, '{"zh-CN":"软件","zh-TW":"軟體","en-US":"Software"}', 'software', 90, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"soft_system","prefix":"[","suffix":"]"},{"field":"soft_type","prefix":"[","suffix":"]"},{"field":"soft_name","prefix":"[","suffix":"]"},{"field":"soft_ver","prefix":"[","suffix":"]"},{"field":"soft_lang","prefix":"[","suffix":"]"},{"field":"soft_other","prefix":"[","suffix":"]"}]},"fields":[{"key":"soft_system","type":"select","label":{"zh-CN":"运行环境","zh-TW":"執行環境","en-US":"Platform"},"required":true,"options":{"source":"static","items":[{"value":"Windows","label":{"zh-CN":"Windows","zh-TW":"Windows","en-US":"Windows"}},{"value":"Linux","label":{"zh-CN":"Linux","zh-TW":"Linux","en-US":"Linux"}},{"value":"Mac","label":{"zh-CN":"Mac","zh-TW":"Mac","en-US":"Mac"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"soft_type","type":"select","label":{"zh-CN":"软件类型","zh-TW":"軟體類型","en-US":"Software type"},"required":true,"options":{"source":"static","items":[{"value":"操作系统","label":{"zh-CN":"操作系统","zh-TW":"作業系統","en-US":"Operating system"}},{"value":"编程开发","label":{"zh-CN":"编程开发","zh-TW":"程式開發","en-US":"Development"}},{"value":"行业软件","label":{"zh-CN":"行业软件","zh-TW":"行業軟體","en-US":"Industry software"}},{"value":"安全软件","label":{"zh-CN":"安全软件","zh-TW":"安全軟體","en-US":"Security"}},{"value":"媒体播放","label":{"zh-CN":"媒体播放","zh-TW":"媒體播放","en-US":"Media player"}},{"value":"媒体处理","label":{"zh-CN":"媒体处理","zh-TW":"媒體處理","en-US":"Media processing"}},{"value":"网络应用","label":{"zh-CN":"网络应用","zh-TW":"網路應用","en-US":"Network app"}},{"value":"办公软件","label":{"zh-CN":"办公软件","zh-TW":"辦公軟體","en-US":"Office"}},{"value":"教育软件","label":{"zh-CN":"教育软件","zh-TW":"教育軟體","en-US":"Education"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"soft_name","type":"text","label":{"zh-CN":"软件名称","zh-TW":"軟體名稱","en-US":"Software name"},"description":{"zh-CN":"要求官方中文名或者官方英文名。","zh-TW":"要求官方中文名或者官方英文名。","en-US":"Use the official Chinese or English name."},"required":true},{"key":"soft_ver","type":"text","label":{"zh-CN":"软件版本","zh-TW":"軟體版本","en-US":"Version"},"description":{"zh-CN":"若不清楚，可通过软件的“帮助-关于”查看。","zh-TW":"若不清楚，可透過軟體的「說明-關於」查看。","en-US":"Check Help > About if unsure."},"required":true},{"key":"soft_lang","type":"select","label":{"zh-CN":"软件语言","zh-TW":"軟體語言","en-US":"Language"},"required":true,"options":{"source":"static","items":[{"value":"中文","label":{"zh-CN":"中文","zh-TW":"中文","en-US":"Chinese"}},{"value":"英文","label":{"zh-CN":"英文","zh-TW":"英文","en-US":"English"}},{"value":"多国语言","label":{"zh-CN":"多国语言","zh-TW":"多國語言","en-US":"Multilingual"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"soft_other","type":"text","label":{"zh-CN":"其它说明","zh-TW":"其他說明","en-US":"Notes"},"description":{"zh-CN":"如破解与否、32/64位等。","zh-TW":"如破解與否、32/64位等。","en-US":"Crack status, 32/64-bit, or other notes."},"required":false}]}', NOW(), NOW()),
    (10, '{"zh-CN":"游戏","zh-TW":"遊戲","en-US":"Games"}', 'games', 100, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"game_system","prefix":"[","suffix":"]"},{"field":"game_chname","prefix":"[","suffix":"]"},{"field":"game_enname","prefix":"[","suffix":"]"},{"field":"game_lang","prefix":"[","suffix":"]"},{"field":"game_form","prefix":"[","suffix":"]"}]},"fields":[{"key":"game_system","type":"select","label":{"zh-CN":"运行平台","zh-TW":"執行平台","en-US":"Platform"},"required":true,"options":{"source":"static","items":[{"value":"PC","label":{"zh-CN":"PC","zh-TW":"PC","en-US":"PC"}},{"value":"PS","label":{"zh-CN":"PS","zh-TW":"PS","en-US":"PS"}},{"value":"PS2","label":{"zh-CN":"PS2","zh-TW":"PS2","en-US":"PS2"}},{"value":"PS3","label":{"zh-CN":"PS3","zh-TW":"PS3","en-US":"PS3"}},{"value":"PSP","label":{"zh-CN":"PSP","zh-TW":"PSP","en-US":"PSP"}},{"value":"XBOX","label":{"zh-CN":"XBOX","zh-TW":"XBOX","en-US":"XBOX"}},{"value":"XBOX360","label":{"zh-CN":"XBOX360","zh-TW":"XBOX360","en-US":"XBOX360"}},{"value":"NDS","label":{"zh-CN":"NDS","zh-TW":"NDS","en-US":"NDS"}},{"value":"NGC","label":{"zh-CN":"NGC","zh-TW":"NGC","en-US":"NGC"}},{"value":"Wii","label":{"zh-CN":"Wii","zh-TW":"Wii","en-US":"Wii"}},{"value":"GBA","label":{"zh-CN":"GBA","zh-TW":"GBA","en-US":"GBA"}},{"value":"视频","label":{"zh-CN":"视频","zh-TW":"影片","en-US":"Video"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"game_chname","type":"text","label":{"zh-CN":"中文名/视频名","zh-TW":"中文名/影片名","en-US":"Chinese title / video title"},"required":true},{"key":"game_enname","type":"text","label":{"zh-CN":"英文名","zh-TW":"英文名","en-US":"English title"},"description":{"zh-CN":"如没有可不填。","zh-TW":"如沒有可不填。","en-US":"Optional if unavailable."},"required":false},{"key":"game_lang","type":"select","label":{"zh-CN":"游戏语言","zh-TW":"遊戲語言","en-US":"Language"},"required":true,"options":{"source":"static","items":[{"value":"中文","label":{"zh-CN":"中文","zh-TW":"中文","en-US":"Chinese"}},{"value":"英文","label":{"zh-CN":"英文","zh-TW":"英文","en-US":"English"}},{"value":"日文","label":{"zh-CN":"日文","zh-TW":"日文","en-US":"Japanese"}},{"value":"多国语言","label":{"zh-CN":"多国语言","zh-TW":"多國語言","en-US":"Multilingual"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"game_form","type":"select","label":{"zh-CN":"游戏格式","zh-TW":"遊戲格式","en-US":"Format"},"required":true,"options":{"source":"static","items":[{"value":"光盘镜像","label":{"zh-CN":"光盘镜像","zh-TW":"光碟映像","en-US":"Disc image"}},{"value":"压缩包","label":{"zh-CN":"压缩包","zh-TW":"壓縮包","en-US":"Archive"}},{"value":"安装包","label":{"zh-CN":"安装包","zh-TW":"安裝包","en-US":"Installer"}},{"value":"FLV","label":{"zh-CN":"FLV","zh-TW":"FLV","en-US":"FLV"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}}]}', NOW(), NOW()),
    (11, '{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}', 'other', 110, b'1', '{"title":{"mode":"generated","allowManualOverride":false,"parts":[{"field":"others_type","prefix":"[","suffix":"]"},{"field":"others_name","prefix":"[","suffix":"]"}]},"fields":[{"key":"others_type","type":"select","label":{"zh-CN":"类型","zh-TW":"類型","en-US":"Type"},"required":true,"options":{"source":"static","items":[{"value":"音频","label":{"zh-CN":"音频","zh-TW":"音訊","en-US":"Audio"}},{"value":"视频","label":{"zh-CN":"视频","zh-TW":"影片","en-US":"Video"}},{"value":"图片","label":{"zh-CN":"图片","zh-TW":"圖片","en-US":"Image"}},{"value":"文档","label":{"zh-CN":"文档","zh-TW":"文件","en-US":"Document"}},{"value":"其它","label":{"zh-CN":"其它","zh-TW":"其他","en-US":"Other"}}]}},{"key":"others_name","type":"text","label":{"zh-CN":"名称","zh-TW":"名稱","en-US":"Title"},"required":true}]}', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `name_i18n` = VALUES(`name_i18n`),
    `slug` = VALUES(`slug`),
    `sort_order` = VALUES(`sort_order`),
    `enabled` = VALUES(`enabled`),
    `upload_config` = VALUES(`upload_config`),
    `updated_at` = NOW();

INSERT INTO `catalog_tag_group`
    (`id`, `name_i18n`, `slug`, `category_ids`, `sort_order`, `created_at`, `updated_at`)
VALUES
    (1, '{"zh-CN":"分辨率","zh-TW":"解析度","en-US":"Resolution"}', 'resolution', '[3]', 10, NOW(), NOW()),
    (2, '{"zh-CN":"来源","zh-TW":"來源","en-US":"Source"}', 'source', '[3]', 20, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `name_i18n` = VALUES(`name_i18n`),
    `slug` = VALUES(`slug`),
    `category_ids` = VALUES(`category_ids`),
    `sort_order` = VALUES(`sort_order`),
    `updated_at` = NOW();

INSERT INTO `catalog_tag`
    (`id`, `group_id`, `name_i18n`, `value`, `sort_order`, `created_at`, `updated_at`)
VALUES
    (1, 1, '{"zh-CN":"720p","zh-TW":"720p","en-US":"720p"}', '720p', 10, NOW(), NOW()),
    (2, 1, '{"zh-CN":"1080p","zh-TW":"1080p","en-US":"1080p"}', '1080p', 20, NOW(), NOW()),
    (3, 1, '{"zh-CN":"2160p","zh-TW":"2160p","en-US":"2160p"}', '2160p', 30, NOW(), NOW()),
    (4, 2, '{"zh-CN":"WEB-DL","zh-TW":"WEB-DL","en-US":"WEB-DL"}', 'WEB-DL', 10, NOW(), NOW()),
    (5, 2, '{"zh-CN":"Blu-ray","zh-TW":"Blu-ray","en-US":"Blu-ray"}', 'Blu-ray', 20, NOW(), NOW()),
    (6, 2, '{"zh-CN":"HDTV","zh-TW":"HDTV","en-US":"HDTV"}', 'HDTV', 30, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    `group_id` = VALUES(`group_id`),
    `name_i18n` = VALUES(`name_i18n`),
    `value` = VALUES(`value`),
    `sort_order` = VALUES(`sort_order`),
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
