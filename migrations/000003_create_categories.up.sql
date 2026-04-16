CREATE TABLE IF NOT EXISTS categories (
    id           BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    user_id      BIGINT UNSIGNED  NULL,
    name         VARCHAR(50)      NOT NULL,
    color        VARCHAR(7)       NOT NULL DEFAULT '#6B7280',
    icon         VARCHAR(50)      NOT NULL DEFAULT 'circle',
    is_default   BOOLEAN          NOT NULL DEFAULT FALSE,
    sort_order   TINYINT UNSIGNED NOT NULL DEFAULT 0,
    created_at   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_categories_user_id (user_id),
    CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
