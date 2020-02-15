DROP TABLE IF EXISTS customer_address;
CREATE TABLE customer_address (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    cuid INT UNSIGNED NOT NULL,
    label VARCHAR(128),
    reciever VARCHAR(32),
    address TEXT,
    phone VARCHAR(32),
    update_at DATETIME DEFAULT NOW(),
    create_at DATETIME DEFAULT NOW()
) DEFAULT CHARSET=utf8;

ALTER TABLE customer_address ADD INDEX(cuid);
