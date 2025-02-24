create database shortener;

use shortener;

CREATE TABLE user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(64) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY unique_username (username)  -- 唯一约束可以命名，便于管理和维护
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User table';

CREATE TABLE sequence (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    stub varchar(1) NOT NULL,
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_uniq_stub (stub)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;


CREATE TABLE short_url_map (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_by varchar(64) NOT NULL default 'system',
    is_del tinyint unsigned NOT NULL default 0 comment '0-no, 1-yes',

    lurl varchar(1024) default null,
    md5 char(32) default null ,
    surl varchar(11) default null ,

    PRIMARY KEY (id),
    INDEX(is_del),
    UNIQUE(md5),
    UNIQUE(surl)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment 'map long url to short url';

-- 记录每次访问
CREATE TABLE short_url_access_log (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    -- 关联 short_url_map 表的 id
    short_url_id bigint unsigned NOT NULL,
    -- 访问时间
    access_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 访问的 IP 地址
    access_ip varchar(45) NOT NULL,
    -- 访问的用户代理信息
    user_agent varchar(255) DEFAULT NULL,
    -- 访问者的国家 地区 城市
    country varchar(100) DEFAULT NULL,
    region varchar(100) DEFAULT NULL,
    city varchar(100) DEFAULT NULL,

    PRIMARY KEY (id),
    -- 关联 short_url_map 表的外键约束
    FOREIGN KEY (short_url_id) REFERENCES short_url_map(id) ON DELETE CASCADE,
    -- 为 short_url_id 创建索引，加快查询速度
    INDEX idx_short_url_id (short_url_id),
    -- 为 access_time 创建索引，方便按时间范围查询
    INDEX idx_access_time (access_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'Short URL access log';