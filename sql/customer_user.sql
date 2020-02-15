DROP TABLE IF EXISTS customer_user;
CREATE TABLE customer_user (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    open_id VARCHAR(64) NOT NULL,
    union_id VARCHAR(64) DEFAULT "",
    token VARCHAR(64) DEFAULT "",
    expired_at DATETIME DEFAULT NOW()
) DEFAULT CHARSET=utf8;

ALTER TABLE `customer_user` ADD UNIQUE(open_id), ADD UNIQUE(token);
