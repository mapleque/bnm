DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    cuid INT UNSIGNED NOT NULL,
    bid INT UNSIGNED NOT NULL,
    itid INT UNSIGNED NOT NULL,
    name VARCHAR(512) NOT NULL,
    price INT UNSIGNED NOT NULL,
    counts INT UNSIGNED NOT NULL,
    reciever VARCHAR(32),
    address TEXT,
    phone VARCHAR(32),
    stage VARCHAR(32) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    additional TEXT,
    exp_no VARCHAR(64) DEFAULT "",
    update_at DATETIME DEFAULT NOW(),
    create_at DATETIME DEFAULT NOW()
) DEFAULT CHARSET=utf8mb4;

ALTER TABLE `order` ADD INDEX(cuid, status), ADD INDEX(bid, status);
