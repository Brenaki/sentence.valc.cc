CREATE TABLE IF NOT EXISTS frases (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    quote            TEXT            NOT NULL,
    author           VARCHAR(255)    NOT NULL DEFAULT '',
    work             VARCHAR(512)    NOT NULL DEFAULT '',
    categories       JSON            NOT NULL,
    like_quantity    INT             NOT NULL DEFAULT 0,
    deslike_quantity INT             NOT NULL DEFAULT 0,
    created_at       TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_frases_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
