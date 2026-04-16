# System Architecture Document
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-16
**Tác giả:** SA Agent
**Dựa trên:** Business Use Cases v1.1 / Project Plan v1.0

---

## 1. Architecture Overview

### Pattern: Monolithic Layered Architecture (SSR)

**Lý do chọn:**
- Solo developer — monolith đơn giản hơn, không cần quản lý inter-service communication
- SSR (Server-Side Rendering) — phù hợp với Go templates + HTMX, không cần build pipeline FE phức tạp
- Layered (Handler → Service → Repository) — tách biệt rõ concerns, dễ test từng layer độc lập
- Không có lý do kỹ thuật để dùng microservices với scale cá nhân

**Design Principles:**
- **Separation of Concerns:** Handler xử lý HTTP, Service xử lý business logic, Repository xử lý DB
- **Dependency Inversion:** Service và Repository dùng interface — dễ mock khi test
- **Fail Fast:** Validate input sớm nhất có thể (ở handler layer)
- **Explicit over Implicit:** Error handling rõ ràng theo phong cách Go (`if err != nil`)
- **Security by Default:** Auth middleware bảo vệ toàn bộ route, whitelist thay vì blacklist

---

## 2. System Context Diagram

```
                    ┌─────────────────────────────┐
                    │         INTERNET             │
                    └──────────┬──────────────────┘
                               │ HTTPS
                    ┌──────────▼──────────────────┐
                    │    Web Browser (User)        │
                    │  Chrome / Firefox / Safari   │
                    └──────────┬──────────────────┘
                               │ HTTPS / HTML
                    ┌──────────▼──────────────────┐
                    │   TIME TRACKER APPLICATION  │
                    │   (Go + Gin + Templates)    │
                    └──┬──────────────────────┬───┘
                       │                      │
           ┌───────────▼───┐      ┌───────────▼───────────┐
           │  MySQL 8.x    │      │  Email Service         │
           │  (primary DB) │      │  (Resend SMTP)         │
           └───────────────┘      └───────────────────────┘

External Actors:
  - User (Browser): Tương tác qua HTML forms + HTMX requests
  - Email Service: Gửi verification email, reset password
  - Railway.app: Hosting platform (app + DB)
```

---

## 3. High-Level Architecture

### Layered Architecture

```
┌────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
│  Go HTML Templates + Tailwind CSS + HTMX + Alpine.js       │
│  (web/templates/**/*.html + web/static/)                   │
├────────────────────────────────────────────────────────────┤
│                    HANDLER LAYER                            │
│  Gin Router + HTTP Handlers                                 │
│  Trách nhiệm: Parse request, validate input, gọi service,  │
│               render template / trả JSON partial           │
│  (internal/handler/)                                        │
├────────────────────────────────────────────────────────────┤
│                    MIDDLEWARE LAYER                         │
│  Auth, CSRF, Rate Limit, Request Logger, Error Handler     │
│  (internal/middleware/)                                     │
├────────────────────────────────────────────────────────────┤
│                    SERVICE LAYER                            │
│  Business Logic, Use Case implementation, BR enforcement   │
│  (internal/service/)                                        │
├────────────────────────────────────────────────────────────┤
│                    REPOSITORY LAYER                         │
│  Database queries, GORM models, data mapping               │
│  (internal/repository/)                                     │
├────────────────────────────────────────────────────────────┤
│                    INFRASTRUCTURE LAYER                     │
│  MySQL (GORM), Email (gomail), Scheduler (robfig/cron),    │
│  Config (godotenv)                                          │
└────────────────────────────────────────────────────────────┘
```

### Request Flow (Happy Path)

```
Browser
  │ POST /activities (HTMX request)
  ▼
Gin Router
  │
  ▼
Auth Middleware ──► [401 nếu chưa login]
  │
  ▼
CSRF Middleware ──► [403 nếu token không hợp lệ]
  │
  ▼
Handler: ActivityHandler.Create()
  ├── Bind & validate form input
  ├── [400 nếu validation fail → trả partial HTML lỗi]
  │
  ▼
Service: ActivityService.Create()
  ├── Áp dụng Business Rules (BR-13 → BR-18)
  ├── Kiểm tra overlap → chuẩn bị warning nếu có
  │
  ▼
Repository: ActivityRepository.Insert()
  ├── GORM insert vào MySQL
  │
  ▼
Service trả về created activity + warning flag
  │
  ▼
Handler render HTML partial (HTMX swap vào DOM)
  │
  ▼
Browser cập nhật UI không reload trang
```

---

## 4. Component Breakdown

### 4.1 Handler Layer (`internal/handler/`)

| File | Trách nhiệm | UC liên quan |
|---|---|---|
| `auth.go` | Register, Login, Logout, Forgot/Reset password | UC-01 → UC-05 |
| `activity.go` | Log, List, Edit, Delete (soft) activity | UC-06 → UC-09 |
| `category.go` | CRUD category | UC-10 → UC-13 |
| `report.go` | Daily, Weekly, Monthly summary | UC-14, UC-15, UC-19 |
| `trash.go` | List trash, Restore, Permanent delete | UC-17 |
| `export.go` | Export CSV, PDF | UC-18 |
| `profile.go` | View, Edit profile | UC-16 |

**Input:** `*gin.Context`
**Output:** Rendered HTML template hoặc HTML partial (HTMX)

### 4.2 Service Layer (`internal/service/`)

| File | Trách nhiệm |
|---|---|
| `auth_service.go` | Password hashing, token generation, session management, BR-01 → BR-12 |
| `activity_service.go` | Activity CRUD, overlap detection, soft delete, BR-13 → BR-22 |
| `category_service.go` | Category CRUD, default seeding, reassign on delete, BR-23 → BR-28 |
| `report_service.go` | Aggregation logic, time calculation, BR-29 → BR-31 |
| `export_service.go` | CSV/PDF generation, BR-39 → BR-43 |
| `scheduler_service.go` | Trash auto-cleanup job (BR-38) |

**Interface pattern:**
```go
type ActivityService interface {
    Create(userID uint, req CreateActivityRequest) (*Activity, *Warning, error)
    List(userID uint, date time.Time) ([]Activity, error)
    Update(userID uint, id uint, req UpdateActivityRequest) (*Activity, error)
    SoftDelete(userID uint, id uint) error
}
```

### 4.3 Repository Layer (`internal/repository/`)

| File | Trách nhiệm |
|---|---|
| `user_repo.go` | User CRUD, find by email |
| `session_repo.go` | Session create, find, invalidate, cleanup |
| `activity_repo.go` | Activity CRUD, date-range queries, overlap check |
| `category_repo.go` | Category CRUD, list by user |
| `token_repo.go` | Email verification + password reset tokens |

**Interface pattern:**
```go
type ActivityRepository interface {
    Insert(ctx context.Context, activity *model.Activity) error
    FindByUserAndDate(ctx context.Context, userID uint, date time.Time) ([]model.Activity, error)
    FindOverlapping(ctx context.Context, userID uint, start, end time.Time, excludeID *uint) ([]model.Activity, error)
    SoftDelete(ctx context.Context, userID uint, id uint) error
}
```

### 4.4 Middleware (`internal/middleware/`)

| Middleware | Chức năng |
|---|---|
| `auth.go` | Kiểm tra session cookie hợp lệ, inject user vào context |
| `csrf.go` | Validate CSRF token cho POST/PUT/DELETE |
| `rate_limit.go` | Giới hạn request theo IP (đặc biệt cho auth routes) |
| `logger.go` | Log request/response |
| `error_handler.go` | Bắt panic, trả error page phù hợp |

### 4.5 Scheduler (`internal/scheduler/`)

- **Trash Cleanup Job:** Chạy mỗi ngày lúc 2:00 AM UTC+7 — xóa vĩnh viễn activities có `deleted_at` > 30 ngày

---

## 5. Data Architecture

### Storage Strategy

| Loại dữ liệu | Storage | Lý do |
|---|---|---|
| User data, Activities, Categories | MySQL 8.x | Relational data, ACID transactions |
| Session data | MySQL (`user_sessions` table) | Simple, không cần Redis cho scale hiện tại |
| Static assets (CSS, JS) | File system → serve bởi Gin | Đơn giản, không cần CDN |
| Uploaded files | Không có trong scope | N/A |

### Data Flow

```
User Input (HTML Form)
    │
    ▼ (HTTP POST — form-urlencoded hoặc JSON)
Handler: Bind → Validate
    │
    ▼
Service: Apply Business Rules
    │
    ▼
Repository: GORM → MySQL
    │
    ▼ (Query result)
Service: Map to Domain Model
    │
    ▼
Handler: Pass to Template Context
    │
    ▼ (HTML response)
Browser: HTMX swap vào DOM
```

### Timezone Strategy
- Tất cả thời gian **lưu dạng UTC** trong MySQL (`DATETIME` không có timezone)
- Application layer convert UTC ↔ UTC+7 khi đọc/ghi
- Template hiển thị UTC+7 cho người dùng

---

## 6. Integration Architecture

### Authentication Flow (Session-based)

```
Login Request
    │
    ▼
Validate credentials (bcrypt compare)
    │
    ▼
Generate secure random session token (32 bytes)
Hash token (SHA-256) → lưu vào DB
    │
    ▼
Set httpOnly cookie: session_token=<raw_token>; HttpOnly; Secure; SameSite=Strict
    │
    ▼
Per request: cookie → hash → lookup DB → inject user to context
```

**Lý do dùng session thay vì JWT:**
- App SSR — không cần stateless API
- Có thể revoke session ngay lập tức (logout, đổi mật khẩu)
- Đơn giản hơn, không có token expiry edge cases

### Email Integration

```
Service cần gửi email (verification, reset password)
    │
    ▼
EmailService.Send(to, subject, htmlBody)
    │
    ▼
gomail → SMTP (Resend API)
```

### Export Integration

```
Export Request (date range + format)
    │
    ▼
ExportService.GenerateCSV() → encoding/csv → []byte
ExportService.GeneratePDF() → HTML template → WeasyPrint/chromedp → []byte
    │
    ▼
Handler: gin.Context.Data(200, "application/csv", data)
         gin.Context.Data(200, "application/pdf", data)
```

---

## 7. Infrastructure & Deployment

### Development Environment

```
Local Machine
├── Go 1.22+
├── MySQL 8.x (local hoặc Docker)
├── Air (hot reload cho Go)
└── Makefile commands
```

### Production Environment (Railway.app)

```
Railway.app
├── Go Service (auto-deploy từ GitHub push)
│   ├── Build: go build ./cmd/server
│   └── Run: ./server
└── MySQL Plugin (managed MySQL 8.x)
```

### CI/CD Pipeline

```
git push → GitHub
    │
    ▼
GitHub Actions:
    ├── go vet
    ├── go test ./...
    └── (nếu pass) → Railway auto-deploy
```

### Docker (local dev option)

```yaml
# docker-compose.yml — MySQL local
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: time_tracker
    ports:
      - "3306:3306"
```

---

## 8. Security Architecture

### Authentication & Authorization

| Layer | Biện pháp |
|---|---|
| Password | bcrypt (cost=12) |
| Session token | Crypto random 32 bytes → SHA-256 hash lưu DB |
| Cookie | `HttpOnly`, `Secure`, `SameSite=Strict` |
| Session expiry | 7 ngày (BR-06), cleanup job hàng ngày |
| Route protection | Auth middleware whitelist (chỉ `/login`, `/register`, `/forgot-password` là public) |

### CSRF Protection

- Double Submit Cookie pattern hoặc Synchronizer Token Pattern
- Áp dụng cho tất cả form POST/PUT/DELETE

### Input Validation & Injection Prevention

| Threat | Mitigation |
|---|---|
| SQL Injection | GORM parameterized queries — không dùng raw string concat |
| XSS | Go `html/template` auto-escape HTML bằng mặc định |
| Path Traversal | Không có file upload trong scope |
| Mass Assignment | Bind vào DTO struct có tag `form:""` tường minh |

### Rate Limiting

| Endpoint | Giới hạn |
|---|---|
| POST /login | 10 req/phút per IP |
| POST /register | 5 req/phút per IP |
| POST /forgot-password | 3 req/phút per IP |
| Các route khác | 100 req/phút per IP |

### Brute Force Protection

- 5 lần đăng nhập sai → khóa tài khoản 15 phút (BR-05)
- Lưu `failed_login_attempts` + `locked_until` trong bảng `users`

### Sensitive Data

| Dữ liệu | Xử lý |
|---|---|
| Password | Chỉ lưu bcrypt hash, không bao giờ log |
| Session token | Chỉ lưu SHA-256 hash trong DB |
| Reset/verify token | Chỉ lưu SHA-256 hash, raw token chỉ gửi qua email |
| Email | Lưu plaintext (cần thiết cho gửi email) |

---

## 9. Architecture Decision Records (ADRs)

| ADR | Quyết định | Lý do | Trade-offs |
|---|---|---|---|
| ADR-01 | Monolith thay vì Microservices | Solo dev, personal scale | Khó horizontal scale sau này (chấp nhận được) |
| ADR-02 | SSR (Go Templates + HTMX) thay vì SPA | BE dev background, không cần build FE | Ít interactive hơn SPA (chấp nhận được với HTMX) |
| ADR-03 | Session-based auth thay vì JWT | SSR app, cần revoke ngay | Stateful — phụ thuộc DB cho mỗi request |
| ADR-04 | MySQL thay vì PostgreSQL | User đã quen, ổn định | PostgreSQL có nhiều advanced feature hơn |
| ADR-05 | Gin thay vì Echo/Fiber | Phổ biến nhất, nhiều tài liệu + AI examples | Không có lý do kỹ thuật đặc biệt |
| ADR-06 | GORM thay vì sqlx/sqlc | Solo dev, ORM nhanh hơn cho CRUD đơn giản | Ít control hơn raw SQL, performance thấp hơn một chút |
| ADR-07 | Railway.app thay vì VPS | Zero-config deploy, free tier đủ dùng | Less control, vendor lock-in nhẹ |
| ADR-08 | golang-migrate thay vì GORM AutoMigrate | Migration files versioned trong git, production-safe | Thêm bước manual khi deploy lần đầu |

---

## 10. Non-Functional Requirements Mapping

| NFR | Yêu cầu | Giải pháp |
|---|---|---|
| **Performance** | Log activity < 10 giây UX | HTMX partial update, không reload; index DB đúng chỗ |
| **Security** | Không leak data user khác | Auth middleware + userID check trong mọi query |
| **Reliability** | App không crash khi lỗi | Recovery middleware, graceful error pages |
| **Maintainability** | Có thể đọc code sau 1 tháng | Layered arch, interface-based, không over-engineer |
| **Testability** | Unit test từng layer | Interface + DI, dùng mock cho Repository layer |
| **Observability** | Biết khi có lỗi | Structured logging, error middleware |
| **Timezone** | Hiển thị đúng UTC+7 | Convert ở application layer, lưu UTC trong DB |
