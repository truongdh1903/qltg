# Database Schema Document
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-16
**Database:** MySQL 8.x
**Charset:** utf8mb4 / utf8mb4_unicode_ci

---

## 1. ERD Overview

```
users (1) ──────────────── (N) user_sessions
  │
  ├── (1) ──── (N) email_verifications
  ├── (1) ──── (N) password_reset_tokens
  ├── (1) ──── (N) activities
  └── (1) ──── (N) categories [user_id nullable]
                        │
                        └── (1) ──── (N) activities

categories.user_id = NULL  →  Default/system category (dùng chung)
categories.user_id = <id>  →  Custom category của user đó
```

---

## 2. Tables

### 2.1 `users`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `email` | VARCHAR(255) | NO | — | UNIQUE, lowercase |
| `password_hash` | VARCHAR(255) | NO | — | bcrypt hash |
| `display_name` | VARCHAR(50) | NO | — | |
| `is_verified` | BOOLEAN | NO | FALSE | Email đã xác thực chưa |
| `failed_login_attempts` | TINYINT UNSIGNED | NO | 0 | Reset về 0 sau login thành công |
| `locked_until` | DATETIME | YES | NULL | NULL = không bị khóa |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |
| `updated_at` | DATETIME | NO | CURRENT_TIMESTAMP ON UPDATE | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_users_email (email)`

---

### 2.2 `user_sessions`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `user_id` | BIGINT UNSIGNED | NO | — | FK → users.id |
| `session_token_hash` | VARCHAR(64) | NO | — | SHA-256 hex của raw token |
| `expires_at` | DATETIME | NO | — | Tạo lúc + 7 ngày |
| `last_active_at` | DATETIME | NO | CURRENT_TIMESTAMP | Cập nhật mỗi request |
| `ip_address` | VARCHAR(45) | YES | NULL | IPv4 hoặc IPv6 |
| `user_agent` | VARCHAR(500) | YES | NULL | |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_sessions_token (session_token_hash)`
- `KEY idx_sessions_user_id (user_id)`
- `KEY idx_sessions_expires_at (expires_at)` — dùng cho cleanup job

**Foreign Keys:**
- `CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`

---

### 2.3 `email_verifications`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `user_id` | BIGINT UNSIGNED | NO | — | FK → users.id |
| `token_hash` | VARCHAR(64) | NO | — | SHA-256 của raw token gửi qua email |
| `expires_at` | DATETIME | NO | — | Tạo lúc + 24 giờ (BR-04) |
| `used_at` | DATETIME | YES | NULL | NULL = chưa dùng |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_ev_token (token_hash)`
- `KEY idx_ev_user_id (user_id)`

**Foreign Keys:**
- `CONSTRAINT fk_ev_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`

---

### 2.4 `password_reset_tokens`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `user_id` | BIGINT UNSIGNED | NO | — | FK → users.id |
| `token_hash` | VARCHAR(64) | NO | — | SHA-256 của raw token gửi qua email |
| `expires_at` | DATETIME | NO | — | Tạo lúc + 1 giờ (BR-08) |
| `used_at` | DATETIME | YES | NULL | NULL = chưa dùng; BR-09: chỉ dùng 1 lần |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_prt_token (token_hash)`
- `KEY idx_prt_user_id (user_id)`

**Foreign Keys:**
- `CONSTRAINT fk_prt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`

---

### 2.5 `categories`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `user_id` | BIGINT UNSIGNED | YES | NULL | NULL = default category (BR-26, BR-27) |
| `name` | VARCHAR(50) | NO | — | BR-23: max 50 ký tự |
| `color` | VARCHAR(7) | NO | '#6B7280' | Hex color, ví dụ `#3B82F6` |
| `icon` | VARCHAR(50) | NO | 'circle' | Icon identifier |
| `is_default` | BOOLEAN | NO | FALSE | TRUE = 1 trong 6 default categories |
| `sort_order` | TINYINT UNSIGNED | NO | 0 | Thứ tự hiển thị |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |
| `updated_at` | DATETIME | NO | CURRENT_TIMESTAMP ON UPDATE | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `KEY idx_categories_user_id (user_id)`

**Foreign Keys:**
- `CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`

**Notes:**
- BR-25 (tên không trùng, không phân biệt hoa thường) được enforce ở Service layer
- Query lấy categories của user: `WHERE user_id IS NULL OR user_id = ?`

**Default categories (seed data):**

| name | color | icon | sort_order |
|---|---|---|---|
| Làm việc | #3B82F6 | briefcase | 1 |
| Di chuyển | #F59E0B | car | 2 |
| Ăn uống | #10B981 | utensils | 3 |
| Giải trí | #8B5CF6 | gamepad | 4 |
| Mạng xã hội | #EF4444 | smartphone | 5 |
| Ngủ nghỉ | #6B7280 | moon | 6 |

---

### 2.6 `activities`

| Column | Type | Nullable | Default | Ghi chú |
|---|---|---|---|---|
| `id` | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PK |
| `user_id` | BIGINT UNSIGNED | NO | — | FK → users.id |
| `category_id` | BIGINT UNSIGNED | NO | — | FK → categories.id |
| `name` | VARCHAR(100) | NO | — | BR-13: max 100 ký tự |
| `note` | TEXT | YES | NULL | BR-14: max 500 ký tự (enforce app layer) |
| `started_at` | DATETIME | NO | — | Lưu UTC |
| `ended_at` | DATETIME | NO | — | Lưu UTC; BR-15: ended_at > started_at |
| `duration_minutes` | SMALLINT UNSIGNED | NO | — | Computed: (ended_at - started_at) / 60; lưu để tránh tính lại |
| `deleted_at` | DATETIME | YES | NULL | NULL = active; NOT NULL = trong Trash (BR-22) |
| `created_at` | DATETIME | NO | CURRENT_TIMESTAMP | UTC |
| `updated_at` | DATETIME | NO | CURRENT_TIMESTAMP ON UPDATE | UTC |

**Indexes:**
- `PRIMARY KEY (id)`
- `KEY idx_activities_user_date (user_id, started_at)` — query chính theo ngày
- `KEY idx_activities_user_deleted (user_id, deleted_at)` — phân biệt active vs trash
- `KEY idx_activities_category (category_id)`

**Foreign Keys:**
- `CONSTRAINT fk_activities_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- `CONSTRAINT fk_activities_category FOREIGN KEY (category_id) REFERENCES categories(id)`

**Notes:**
- Soft delete: `deleted_at IS NULL` = active, `deleted_at IS NOT NULL` = trash
- Mọi query business đều thêm `WHERE deleted_at IS NULL`
- Trash query: `WHERE deleted_at IS NOT NULL AND user_id = ?`
- Auto cleanup: `WHERE deleted_at < NOW() - INTERVAL 30 DAY`
- Overlap detection query:
  ```sql
  SELECT id FROM activities
  WHERE user_id = ?
    AND deleted_at IS NULL
    AND id != ?  -- exclude self khi edit
    AND started_at < ? -- new ended_at
    AND ended_at > ?   -- new started_at
  ```

---

## 3. Migration Strategy

**Tool:** `golang-migrate/migrate`

**Convention đặt tên file:**
```
migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_create_sessions.up.sql
├── 000002_create_sessions.down.sql
├── 000003_create_tokens.up.sql
├── 000003_create_tokens.down.sql
├── 000004_create_categories.up.sql
├── 000004_create_categories.down.sql
├── 000005_create_activities.up.sql
├── 000005_create_activities.down.sql
└── 000006_seed_default_categories.up.sql
```

**Makefile commands:**
```bash
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback 1 step
make migrate-create name=add_column_x  # Tạo file migration mới
```

---

## 4. Query Patterns Quan Trọng

### Lấy activities theo ngày (Dashboard)
```sql
SELECT a.*, c.name as category_name, c.color, c.icon
FROM activities a
JOIN categories c ON a.category_id = c.id
WHERE a.user_id = ?
  AND a.deleted_at IS NULL
  AND a.started_at >= ?   -- 00:00:00 UTC của ngày (đã convert từ UTC+7)
  AND a.started_at < ?    -- 00:00:00 UTC của ngày hôm sau
ORDER BY a.started_at ASC
```

### Daily Summary (Report)
```sql
SELECT c.name, c.color, c.icon,
       SUM(a.duration_minutes) as total_minutes,
       COUNT(a.id) as activity_count
FROM activities a
JOIN categories c ON a.category_id = c.id
WHERE a.user_id = ?
  AND a.deleted_at IS NULL
  AND a.started_at >= ? AND a.started_at < ?
GROUP BY a.category_id, c.name, c.color, c.icon
ORDER BY total_minutes DESC
```

### Weekly Summary (Report)
```sql
SELECT
  DATE(CONVERT_TZ(a.started_at, '+00:00', '+07:00')) as activity_date,
  c.id as category_id,
  c.name, c.color,
  SUM(a.duration_minutes) as total_minutes
FROM activities a
JOIN categories c ON a.category_id = c.id
WHERE a.user_id = ?
  AND a.deleted_at IS NULL
  AND a.started_at >= ? AND a.started_at < ?
GROUP BY activity_date, a.category_id, c.name, c.color
ORDER BY activity_date ASC
```

### Trash List
```sql
SELECT a.*, c.name as category_name, c.color
FROM activities a
JOIN categories c ON a.category_id = c.id
WHERE a.user_id = ?
  AND a.deleted_at IS NOT NULL
ORDER BY a.deleted_at DESC
```

---

## 5. Index Strategy

| Index | Mục đích | Lý do |
|---|---|---|
| `users.email` | Login lookup | Unique, mỗi request auth |
| `sessions.session_token_hash` | Auth per request | Lookup mỗi request, phải cực nhanh |
| `sessions.expires_at` | Cleanup job | Scan toàn bộ bảng để tìm expired |
| `activities(user_id, started_at)` | Dashboard + Report queries | Query chính, full table scan nếu thiếu |
| `activities(user_id, deleted_at)` | Trash queries | Phân biệt active vs deleted |
| `categories.user_id` | Load categories của user | JOIN thường xuyên |
