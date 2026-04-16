CREATE TABLE IF NOT EXISTS users (
    id                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    email                   VARCHAR(255)    NOT NULL,
    password_hash           VARCHAR(255)    NOT NULL,
    display_name            VARCHAR(50)     NOT NULL,
    is_verified             BOOLEAN         NOT NULL DEFAULT FALSE,
    failed_login_attempts   TINYINT UNSIGNED NOT NULL DEFAULT 0,
    locked_until            DATETIME        NULL,
    created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
