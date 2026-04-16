CREATE TABLE IF NOT EXISTS activities (
    id                BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    user_id           BIGINT UNSIGNED  NOT NULL,
    category_id       BIGINT UNSIGNED  NOT NULL,
    name              VARCHAR(100)     NOT NULL,
    note              TEXT             NULL,
    started_at        DATETIME         NOT NULL,
    ended_at          DATETIME         NOT NULL,
    duration_minutes  SMALLINT UNSIGNED NOT NULL,
    deleted_at        DATETIME         NULL,
    created_at        DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_activities_user_date    (user_id, started_at),
    KEY idx_activities_user_deleted (user_id, deleted_at),
    KEY idx_activities_category     (category_id),
    CONSTRAINT fk_activities_user     FOREIGN KEY (user_id)     REFERENCES users(id)      ON DELETE CASCADE,
    CONSTRAINT fk_activities_category FOREIGN KEY (category_id) REFERENCES categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
