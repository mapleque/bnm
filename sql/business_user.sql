DROP TABLE IF EXISTS business_user;
CREATE TABLE business_user (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    open_id VARCHAR(64) NOT NULL,
    union_id VARCHAR(64) DEFAULT "",
    token VARCHAR(64) DEFAULT "",
    expired_at DATETIME DEFAULT NOW()
) DEFAULT CHARSET=utf8;

ALTER TABLE `business_user` ADD UNIQUE(open_id), ADD UNIQUE(token);
