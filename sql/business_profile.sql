DROP TABLE IF EXISTS business_profile;
CREATE TABLE business_profile (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    buid INT UNSIGNED NOT NULL,
    name VARCHAR(512) NOT NULL,
    avatar TEXT,
    `desc` TEXT,
    qrcode VARCHAR(512) DEFAULT '',
    wxid VARCHAR(128) DEFAULT '',
    status TINYINT NOT NULL DEFAULT 0,
    create_at DATETIME DEFAULT NOW()
) DEFAULT CHARSET=utf8mb4;

ALTER TABLE `business_profile` ADD UNIQUE(buid);
