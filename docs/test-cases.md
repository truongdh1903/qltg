# Test Case Specification
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-17
**Dựa trên:** SRS v1.0, URS v1.0, Business Use Cases v1.1

---

## Quy ước

### Loại test
| Ký hiệu | Loại | Mô tả |
|---|---|---|
| **Unit** | Unit Test | Test một function/method độc lập, mock dependency |
| **Integration** | Integration Test | Test nhiều layer phối hợp (Handler → Service → Repository → DB) |
| **E2E** | End-to-End | Test từ browser đến DB, không mock |

### Mức độ ưu tiên
| Ký hiệu | Ý nghĩa |
|---|---|
| **P1** | Blocker — Hệ thống không dùng được nếu fail |
| **P2** | Critical — Tính năng core bị ảnh hưởng |
| **P3** | Major — Tính năng phụ bị ảnh hưởng |
| **P4** | Minor — UX nhỏ, không ảnh hưởng chức năng |

### Kết quả mong đợi
- **PASS**: Kết quả thực tế khớp với kết quả mong đợi
- **FAIL**: Kết quả thực tế không khớp

---

## Module 1: Authentication

### TC-001: Đăng ký — email hợp lệ
- **SR**: SR-F-001, SR-F-004
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Email chưa tồn tại trong DB
- **Steps**:
  1. POST `/register` với `email=newuser@example.com`, `password=Password1`, `confirm_password=Password1`
- **Expected Result**: HTTP 302 redirect; tài khoản tạo với `status=unverified`; 1 email xác thực được gửi; token lưu trong bảng `email_verification_tokens`

---

### TC-002: Đăng ký — email sai định dạng
- **SR**: SR-F-001
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/register` với `email=notanemail`, `password=Password1`
- **Expected Result**: HTTP 422; response chứa lỗi validation trường `email`; không tạo record nào trong DB

---

### TC-003: Đăng ký — mật khẩu không đủ mạnh
- **SR**: SR-F-002
- **Type**: Unit | **Priority**: P2
- **Test cases con** (parameterized):

| Input mật khẩu | Lý do fail | Expected |
|---|---|---|
| `short1A` (7 ký tự) | < 8 ký tự | Lỗi validation |
| `alllowercase1` | Không có chữ hoa | Lỗi validation |
| `NoNumberHere` | Không có số | Lỗi validation |
| `Password1` | Hợp lệ | Pass |

- **Expected Result**: Các input fail trả về lỗi rõ ràng; `Password1` không có lỗi

---

### TC-004: Đăng ký — email đã tồn tại
- **SR**: SR-F-003
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Email `existing@example.com` đã có trong DB
- **Steps**:
  1. POST `/register` với `email=existing@example.com`
- **Expected Result**: Hiển thị thông báo lỗi chung (không nói "Email đã tồn tại"); không tạo tài khoản mới; HTTP status không tiết lộ sự tồn tại của email

---

### TC-005: Đăng ký — xác thực email hợp lệ
- **SR**: SR-F-005
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Tài khoản `unverified`, token còn hiệu lực (< 24h)
- **Steps**:
  1. GET `/verify-email?token=<valid_token>`
- **Expected Result**: Tài khoản chuyển sang `status=active`; 6 danh mục mặc định được tạo (SR-F-006); redirect đến trang đăng nhập

---

### TC-006: Đăng ký — link xác thực hết hạn
- **SR**: SR-F-005
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Token tạo > 24h trước
- **Steps**:
  1. GET `/verify-email?token=<expired_token>`
- **Expected Result**: Hiển thị trang thông báo hết hạn + nút "Gửi lại email xác thực"; tài khoản vẫn `unverified`

---

### TC-007: Tạo 6 danh mục mặc định sau xác thực
- **SR**: SR-F-006
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Xác thực email thành công
- **Steps**:
  1. Truy vấn DB: `SELECT * FROM categories WHERE user_id = :new_user_id AND is_default = true`
- **Expected Result**: Đúng 6 danh mục: `Làm việc`, `Di chuyển`, `Ăn uống`, `Giải trí`, `Mạng xã hội`, `Ngủ nghỉ`

---

### TC-008: Đăng nhập — thành công
- **SR**: SR-F-007, SR-F-010
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Tài khoản `active`, email đã xác thực
- **Steps**:
  1. POST `/login` với `email=user@example.com`, `password=Password1`
- **Expected Result**: HTTP 302 redirect đến Dashboard; cookie `session_id` được set (HttpOnly=true, Secure=true, SameSite=Lax); session lưu trong DB với `expires_at = NOW() + 7 days`

---

### TC-009: Đăng nhập — sai mật khẩu
- **SR**: SR-F-007
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/login` với mật khẩu sai
- **Expected Result**: Thông báo lỗi "Email hoặc mật khẩu không đúng" — không chỉ rõ trường nào sai; không tạo session

---

### TC-010: Đăng nhập — khóa sau 5 lần sai
- **SR**: SR-F-008
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/login` với mật khẩu sai 5 lần liên tiếp
  2. Lần thứ 6: POST `/login` với mật khẩu đúng
- **Expected Result**: Sau lần 5: tài khoản bị khóa, thông báo thời gian mở khóa; lần 6: vẫn bị từ chối dù mật khẩu đúng, hiển thị thời gian còn lại

---

### TC-011: Đăng nhập — tài khoản chưa xác thực email
- **SR**: SR-F-009
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Tài khoản `status=unverified`
- **Steps**:
  1. POST `/login` với thông tin đúng
- **Expected Result**: Từ chối đăng nhập; hiển thị tùy chọn "Gửi lại email xác thực"

---

### TC-012: Đăng xuất
- **SR**: SR-F-011
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Người dùng đang đăng nhập
- **Steps**:
  1. POST `/logout`
  2. Thử GET `/app/dashboard` với cookie cũ
- **Expected Result**: Cookie bị xóa; session bị hủy trong DB; bước 2 redirect về `/login`

---

### TC-013: Quên mật khẩu — gửi link
- **SR**: SR-F-012, SR-F-013
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. POST `/forgot-password` với `email=user@example.com` (email tồn tại)
  2. POST `/forgot-password` với `email=ghost@example.com` (email không tồn tại)
- **Expected Result**: Cả 2 case đều trả về thông báo giống nhau "Nếu email tồn tại, bạn sẽ nhận được email"; chỉ case 1 gửi email thật; token lưu trong DB với `expires_at = NOW() + 1 hour`

---

### TC-014: Quên mật khẩu — đặt lại thành công
- **SR**: SR-F-014
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Token hợp lệ, chưa dùng, chưa hết hạn
- **Steps**:
  1. GET `/reset-password?token=<valid_token>` → form hiện ra
  2. POST `/reset-password` với `token`, `password=NewPass1`, `confirm=NewPass1`
- **Expected Result**: Mật khẩu cập nhật; toàn bộ session cũ bị xóa trong DB; redirect đến trang đăng nhập

---

### TC-015: Quên mật khẩu — token hết hạn
- **SR**: SR-F-015
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Token tạo > 1 giờ trước
- **Steps**:
  1. GET `/reset-password?token=<expired_token>`
- **Expected Result**: Thông báo "Link đã hết hạn" + form yêu cầu link mới; không thể đặt lại mật khẩu

---

### TC-016: Quên mật khẩu — token đã dùng
- **SR**: SR-F-015
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Token đã được dùng 1 lần
- **Steps**:
  1. POST `/reset-password` lần 2 với cùng token
- **Expected Result**: Thông báo "Link không còn hợp lệ"; không thay đổi mật khẩu

---

### TC-017: Đổi mật khẩu — thành công
- **SR**: SR-F-016, SR-F-018
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Người dùng đang đăng nhập, có 2 session hoạt động
- **Steps**:
  1. POST `/app/settings/change-password` với `current=Password1`, `new=NewPass2`, `confirm=NewPass2`
- **Expected Result**: Mật khẩu cập nhật; session khác bị hủy; session hiện tại vẫn hợp lệ

---

### TC-018: Đổi mật khẩu — sai mật khẩu hiện tại
- **SR**: SR-F-016
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/app/settings/change-password` với `current=WrongPass`
- **Expected Result**: Lỗi "Mật khẩu hiện tại không đúng"; không thay đổi mật khẩu

---

### TC-019: Đổi mật khẩu — mật khẩu mới trùng cũ
- **SR**: SR-F-017
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. POST `/app/settings/change-password` với `new=Password1` (trùng mật khẩu hiện tại)
- **Expected Result**: Lỗi "Mật khẩu mới không được trùng mật khẩu hiện tại"

---

### TC-020: Auth Middleware — route protected
- **SR**: SR-F-071
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. GET `/app/dashboard` không có cookie session
  2. GET `/app/activities` không có cookie session
- **Expected Result**: Cả 2 redirect về `/login`

---

### TC-021: Session sliding window
- **SR**: SR-F-072
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Session tạo 6 ngày trước, vẫn còn hiệu lực
- **Steps**:
  1. Thực hiện bất kỳ request nào đến `/app/*`
- **Expected Result**: `expires_at` trong DB được cập nhật = `NOW() + 7 days`

---

## Module 2: Activity Management

### TC-022: Log hoạt động — thành công
- **SR**: SR-F-019, SR-F-023
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Người dùng đang đăng nhập
- **Steps**:
  1. POST `/app/activities` với `name=Đọc sách`, `category_id=1`, `start_time=2026-04-17T08:00`, `end_time=2026-04-17T09:00`
- **Expected Result**: HTTP 200 (HTMX partial); hoạt động lưu vào DB với đúng `user_id`; danh sách cập nhật mà không reload trang

---

### TC-023: Log hoạt động — pre-fill thời gian hiện tại
- **SR**: SR-F-020
- **Type**: E2E | **Priority**: P3
- **Steps**:
  1. Click nút "Log hoạt động"
- **Expected Result**: Trường `start_time` trong form tự động điền = thời gian hiện tại (UTC+7), chính xác đến phút

---

### TC-024: Log hoạt động — tên rỗng
- **SR**: SR-F-021
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/app/activities` với `name=` (rỗng)
- **Expected Result**: HTTP 422; lỗi "Tên hoạt động không được để trống"; không lưu vào DB

---

### TC-025: Log hoạt động — tên vượt 100 ký tự
- **SR**: SR-F-021
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. Validate service với `name` = 101 ký tự
- **Expected Result**: Lỗi validation "Tên hoạt động tối đa 100 ký tự"

---

### TC-026: Log hoạt động — ghi chú vượt 500 ký tự
- **SR**: SR-F-021
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. Validate service với `note` = 501 ký tự
- **Expected Result**: Lỗi validation "Ghi chú tối đa 500 ký tự"

---

### TC-027: Log hoạt động — thời gian kết thúc ≤ bắt đầu
- **SR**: SR-F-021
- **Type**: Unit | **Priority**: P2
- **Test cases con** (parameterized):

| start_time | end_time | Expected |
|---|---|---|
| 08:00 | 07:59 | Lỗi "Thời gian kết thúc phải sau thời gian bắt đầu" |
| 08:00 | 08:00 | Lỗi |
| 08:00 | 08:01 | Pass |

---

### TC-028: Log hoạt động — thời gian tương lai
- **SR**: SR-F-021
- **Type**: Unit | **Priority**: P2
- **Steps**:
  1. Validate với `start_time` = `NOW() + 1 hour`
- **Expected Result**: Lỗi "Không thể log hoạt động ở thời điểm tương lai"

---

### TC-029: Log hoạt động — quá khứ hợp lệ (trong 2 năm)
- **SR**: SR-F-021
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. Validate với `start_time` = `NOW() - 1 year`
  2. Validate với `start_time` = `NOW() - 2 years - 1 day`
- **Expected Result**: Case 1: Pass. Case 2: Lỗi "Chỉ được log hoạt động trong vòng 2 năm trở lại"

---

### TC-030: Log hoạt động — chồng thời gian (overlap)
- **SR**: SR-F-022
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Đã có hoạt động 08:00–09:00 trong ngày
- **Steps**:
  1. POST activity với `start_time=08:30`, `end_time=10:00`
- **Expected Result**: Hoạt động được lưu thành công; hiển thị toast notification cảnh báo overlap; không bị chặn

---

### TC-031: Xem danh sách hoạt động — isolation per user
- **SR**: SR-F-025
- **Type**: Integration | **Priority**: P1
- **Preconditions**: User A và User B đều có hoạt động ngày hôm nay
- **Steps**:
  1. Đăng nhập User A, GET `/app/dashboard`
- **Expected Result**: Chỉ thấy hoạt động của User A, không thấy hoạt động User B

---

### TC-032: Xem danh sách hoạt động — sắp xếp theo thời gian
- **SR**: SR-F-026
- **Type**: Integration | **Priority**: P3
- **Preconditions**: Có 3 hoạt động: 10:00, 08:00, 14:00
- **Steps**:
  1. GET `/app/dashboard`
- **Expected Result**: Thứ tự hiển thị: 08:00 → 10:00 → 14:00

---

### TC-033: Xem danh sách — ngày không có hoạt động
- **SR**: SR-F-030
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. GET `/app/dashboard?date=2020-01-01` (ngày không có data)
- **Expected Result**: Hiển thị trạng thái rỗng với thông điệp gợi ý hành động; không có lỗi

---

### TC-034: Chỉnh sửa hoạt động — thành công
- **SR**: SR-F-031, SR-F-035
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Hoạt động ID=1 thuộc user hiện tại
- **Steps**:
  1. PUT `/app/activities/1` với `name=Tên mới`, các trường khác giữ nguyên
- **Expected Result**: DB cập nhật; danh sách refresh (HTMX partial); tên mới hiển thị

---

### TC-035: Chỉnh sửa hoạt động — không phải của mình
- **SR**: SR-F-033
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Hoạt động ID=99 thuộc User B; User A đang đăng nhập
- **Steps**:
  1. PUT `/app/activities/99` (User A)
- **Expected Result**: HTTP 403 Forbidden; không thay đổi gì trong DB

---

### TC-036: Xóa hoạt động — soft delete
- **SR**: SR-F-037
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Hoạt động ID=5 tồn tại, thuộc user hiện tại
- **Steps**:
  1. DELETE `/app/activities/5` (sau khi xác nhận)
- **Expected Result**: `activities.deleted_at` được set = NOW(); hoạt động biến khỏi danh sách chính; record vẫn còn trong DB

---

### TC-037: Xóa hoạt động — không xuất hiện trong báo cáo
- **SR**: SR-F-039
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Hoạt động ID=5 đã soft-deleted (có `deleted_at`)
- **Steps**:
  1. GET `/app/reports/daily?date=<ngày của activity 5>`
- **Expected Result**: Hoạt động ID=5 không xuất hiện trong báo cáo; tổng thời gian không tính activity đó

---

### TC-038: Xóa hoạt động — không phải của mình
- **SR**: SR-F-038
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. DELETE `/app/activities/99` (thuộc user khác)
- **Expected Result**: HTTP 403; `deleted_at` của activity 99 vẫn NULL

---

---

## Module 3: Category Management

### TC-039: Xem danh mục
- **SR**: SR-F-040
- **Type**: Integration | **Priority**: P2
- **Preconditions**: User có 6 danh mục mặc định + 2 tùy chỉnh
- **Steps**:
  1. GET `/app/categories`
- **Expected Result**: Hiển thị 8 danh mục; mặc định không có nút sửa/xóa; tùy chỉnh có cả 2 nút

---

### TC-040: Thêm danh mục — thành công
- **SR**: SR-F-041
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. POST `/app/categories` với `name=Thể thao`, `color=#FF5733`, `icon=dumbbell`
- **Expected Result**: Danh mục mới lưu vào DB với đúng `user_id`; hiển thị trong danh sách

---

### TC-041: Thêm danh mục — tên trùng (case-insensitive)
- **SR**: SR-F-042
- **Type**: Unit | **Priority**: P2
- **Preconditions**: Đã có danh mục `Thể thao`
- **Test cases con**:

| Input tên | Expected |
|---|---|
| `Thể thao` | Lỗi "Tên danh mục đã tồn tại" |
| `thể thao` | Lỗi "Tên danh mục đã tồn tại" |
| `THỂ THAO` | Lỗi "Tên danh mục đã tồn tại" |
| `Thể Thao2` | Pass |

---

### TC-042: Thêm danh mục — tên vượt 50 ký tự
- **SR**: SR-F-041
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. Validate với `name` = 51 ký tự
- **Expected Result**: Lỗi "Tên danh mục tối đa 50 ký tự"

---

### TC-043: Sửa danh mục mặc định — bị từ chối
- **SR**: SR-F-043
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Category ID=1 là danh mục mặc định (`is_default=true`)
- **Steps**:
  1. PUT `/app/categories/1` với tên mới
- **Expected Result**: HTTP 403; không thay đổi gì

---

### TC-044: Xóa danh mục mặc định — bị từ chối
- **SR**: SR-F-043
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. DELETE `/app/categories/1` (danh mục mặc định)
- **Expected Result**: HTTP 403; danh mục vẫn tồn tại trong DB

---

### TC-045: Xóa danh mục tùy chỉnh — chưa được dùng
- **SR**: SR-F-046
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Danh mục ID=10 tùy chỉnh, không có hoạt động nào dùng
- **Steps**:
  1. DELETE `/app/categories/10` (xác nhận dialog)
- **Expected Result**: Danh mục bị xóa khỏi DB

---

### TC-046: Xóa danh mục tùy chỉnh — đang được dùng
- **SR**: SR-F-045
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Danh mục ID=10 đang được dùng bởi 5 hoạt động
- **Steps**:
  1. DELETE `/app/categories/10`
  2. Hệ thống yêu cầu chọn danh mục thay thế
  3. Chọn danh mục ID=2 làm thay thế
  4. Xác nhận
- **Expected Result**: 5 hoạt động cập nhật sang `category_id=2`; danh mục ID=10 bị xóa

---

### TC-047: Xóa danh mục — không phải của mình
- **SR**: SR-F-043
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. DELETE `/app/categories/99` (thuộc user khác)
- **Expected Result**: HTTP 403

---

---

## Module 4: Reports

### TC-048: Báo cáo ngày — tính toán đúng
- **SR**: SR-F-047, SR-F-048
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Ngày 2026-04-17 có 3 hoạt động:
  - "Làm việc" 08:00–10:00 (120 phút)
  - "Ăn uống" 12:00–12:30 (30 phút)
  - "Làm việc" 14:00–16:00 (120 phút)
- **Steps**:
  1. GET `/app/reports/daily?date=2026-04-17`
- **Expected Result**:
  - Tổng: 270 phút (4h30m)
  - "Làm việc": 240 phút (88.9%)
  - "Ăn uống": 30 phút (11.1%)
  - Biểu đồ tròn hiển thị đúng tỷ lệ
  - Danh sách 3 hoạt động theo thứ tự thời gian

---

### TC-049: Báo cáo ngày — cảnh báo > 24 giờ
- **SR**: SR-F-049
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Ngày có các hoạt động chồng nhau, tổng > 1440 phút
- **Steps**:
  1. GET `/app/reports/daily?date=<ngày đó>`
- **Expected Result**: Báo cáo vẫn hiển thị; có cảnh báo "Tổng thời gian vượt 24 giờ do các hoạt động trùng nhau"

---

### TC-050: Báo cáo ngày — không có dữ liệu
- **SR**: SR-F-050
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. GET `/app/reports/daily?date=2020-01-01`
- **Expected Result**: Hiển thị trạng thái rỗng; không có lỗi; không có chart rỗng gây lỗi JS

---

### TC-051: Báo cáo tuần — tính đúng khoảng thứ Hai đến Chủ Nhật
- **SR**: SR-F-051
- **Type**: Unit | **Priority**: P2
- **Steps**:
  1. Gọi service `GetWeekRange(date=2026-04-17)` (Thứ Sáu)
- **Expected Result**: Trả về `start=2026-04-13 (Thứ Hai)`, `end=2026-04-19 (Chủ Nhật)`

---

### TC-052: Báo cáo tuần — stacked bar chart data
- **SR**: SR-F-052
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Có dữ liệu đủ 7 ngày trong tuần
- **Steps**:
  1. GET `/app/reports/weekly?date=2026-04-17`
- **Expected Result**: Response chứa dữ liệu cho đủ 7 ngày (Thứ Hai đến Chủ Nhật); ngày không có data trả về 0 cho tất cả danh mục; xác định đúng ngày nhiều/ít nhất

---

### TC-053: Báo cáo tháng — boundary ngày đầu/cuối tháng
- **SR**: SR-F-054
- **Type**: Unit | **Priority**: P2
- **Steps**:
  1. Gọi service `GetMonthRange(year=2026, month=2)` (tháng 2)
- **Expected Result**: `start=2026-02-01`, `end=2026-02-28`

---

### TC-054: Báo cáo tháng — dữ liệu theo tuần
- **SR**: SR-F-055
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. GET `/app/reports/monthly?year=2026&month=4`
- **Expected Result**: Biểu đồ xu hướng có 4–5 điểm dữ liệu (tuần 1 đến tuần 4/5 của tháng 4); tổng theo danh mục chính xác

---

---

## Module 5: Trash

### TC-055: Xem Trash — hiển thị đúng
- **SR**: SR-F-057
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Có 3 hoạt động trong Trash với `deleted_at` khác nhau
- **Steps**:
  1. GET `/app/trash`
- **Expected Result**: Hiển thị 3 mục; sắp xếp theo `deleted_at` giảm dần; mỗi mục hiển thị số ngày còn lại (ví dụ "còn 28 ngày")

---

### TC-056: Khôi phục hoạt động từ Trash
- **SR**: SR-F-058
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Hoạt động ID=5 trong Trash, `start_time=2026-04-10T09:00`
- **Steps**:
  1. POST `/app/trash/5/restore`
- **Expected Result**: `deleted_at` set thành NULL; hoạt động xuất hiện lại ở ngày 2026-04-10 đúng vị trí; không còn trong Trash

---

### TC-057: Xóa vĩnh viễn 1 mục
- **SR**: SR-F-059
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. DELETE `/app/trash/5` (xác nhận dialog)
- **Expected Result**: Record bị DELETE khỏi DB (hard delete); không thể khôi phục

---

### TC-058: Empty Trash
- **SR**: SR-F-060
- **Type**: Integration | **Priority**: P3
- **Preconditions**: Có 5 mục trong Trash
- **Steps**:
  1. DELETE `/app/trash` với action=`empty` (xác nhận dialog)
- **Expected Result**: Toàn bộ 5 mục bị DELETE khỏi DB; Trash hiển thị trạng thái rỗng

---

### TC-059: Scheduler — tự động xóa Trash sau 30 ngày
- **SR**: SR-F-061
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Có hoạt động với `deleted_at = NOW() - 31 days` và `deleted_at = NOW() - 29 days`
- **Steps**:
  1. Chạy manually hàm `CleanupTrash()`
- **Expected Result**: Mục 31 ngày bị DELETE; mục 29 ngày vẫn còn trong DB

---

### TC-060: Trash không tính vào báo cáo
- **SR**: SR-F-039
- **Type**: Integration | **Priority**: P1
- (Xem TC-037 — đã cover)

---

---

## Module 6: Export

### TC-061: Export CSV — thành công
- **SR**: SR-F-062, SR-F-064, SR-F-067
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Có 10 hoạt động trong khoảng 2026-04-01 đến 2026-04-17
- **Steps**:
  1. GET `/app/export?start=2026-04-01&end=2026-04-17&format=csv`
- **Expected Result**: HTTP 200; `Content-Type: text/csv`; `Content-Disposition: attachment; filename=...`; file CSV có header đúng: `date,activity_name,category,start_time,end_time,total_minutes,note`; đúng 10 rows; thời gian theo UTC+7

---

### TC-062: Export PDF — thành công
- **SR**: SR-F-065, SR-F-067
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. GET `/app/export?start=2026-04-01&end=2026-04-17&format=pdf`
- **Expected Result**: HTTP 200; `Content-Type: application/pdf`; file PDF có summary danh mục + danh sách hoạt động; thời gian theo UTC+7

---

### TC-063: Export — khoảng thời gian > 365 ngày
- **SR**: SR-F-063
- **Type**: Integration | **Priority**: P2
- **Steps**:
  1. GET `/app/export?start=2024-01-01&end=2026-01-02` (> 365 ngày)
- **Expected Result**: Lỗi "Chỉ export tối đa 1 năm dữ liệu mỗi lần"; không tạo file

---

### TC-064: Export — khoảng thời gian không có dữ liệu
- **SR**: SR-F-067
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. GET `/app/export?start=2020-01-01&end=2020-01-31&format=csv`
- **Expected Result**: Thông báo "Không có dữ liệu trong khoảng thời gian này"; không trả file rỗng

---

### TC-065: Export — loại trừ hoạt động trong Trash
- **SR**: SR-F-066
- **Type**: Integration | **Priority**: P1
- **Preconditions**: Ngày 2026-04-10 có 3 hoạt động, 1 trong Trash
- **Steps**:
  1. GET `/app/export?start=2026-04-10&end=2026-04-10&format=csv`
- **Expected Result**: File CSV chỉ có 2 rows (không tính mục trong Trash)

---

---

## Module 7: Profile

### TC-066: Xem hồ sơ
- **SR**: SR-F-068
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. GET `/app/settings/profile`
- **Expected Result**: Hiển thị tên hiển thị (input editable), email (read-only), timezone "UTC+7" (read-only)

---

### TC-067: Cập nhật tên hiển thị — thành công
- **SR**: SR-F-069
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. PUT `/app/settings/profile` với `display_name=Nguyễn Văn A`
- **Expected Result**: DB cập nhật `display_name`; thông báo thành công

---

### TC-068: Cập nhật tên hiển thị — để trống
- **SR**: SR-F-069
- **Type**: Unit | **Priority**: P3
- **Steps**:
  1. Validate với `display_name=` (rỗng)
- **Expected Result**: Lỗi "Tên hiển thị không được để trống"

---

### TC-069: Cập nhật tên hiển thị — vượt 50 ký tự
- **SR**: SR-F-069
- **Type**: Unit | **Priority**: P4
- **Steps**:
  1. Validate với `display_name` = 51 ký tự
- **Expected Result**: Lỗi "Tên hiển thị tối đa 50 ký tự"

---

### TC-070: Thay đổi email — bị từ chối
- **SR**: SR-F-070
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. PUT `/app/settings/profile` với `email=newemail@example.com`
- **Expected Result**: Trường email bị bỏ qua hoàn toàn; email trong DB không thay đổi; hoặc HTTP 400 nếu field không được phép

---

---

## Module 8: Security

### TC-071: SQL Injection — login form
- **SR**: SR-NF-009
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. POST `/login` với `email=admin'--`, `password=anything`
  2. POST `/login` với `email='; DROP TABLE users; --`
- **Expected Result**: Trả về lỗi validation email thông thường; không thực thi SQL nào ngoài query parameterized; DB không bị ảnh hưởng

---

### TC-072: XSS — tên hoạt động
- **SR**: SR-NF-010
- **Type**: E2E | **Priority**: P1
- **Steps**:
  1. Tạo hoạt động với `name=<script>alert('xss')</script>`
  2. Xem Dashboard
- **Expected Result**: Text hiển thị là `<script>alert('xss')</script>` (escaped); không có popup; không có script chạy

---

### TC-073: Truy cập dữ liệu user khác — direct URL
- **SR**: SR-NF-012
- **Type**: Integration | **Priority**: P1
- **Preconditions**: User A (ID=1), User B (ID=2) cùng đăng nhập
- **Steps**:
  1. User B thử GET `/app/activities?date=2026-04-17` khi đang đăng nhập
  2. User B thử PUT `/app/activities/1` (activity thuộc User A)
- **Expected Result**: Bước 1: Chỉ thấy dữ liệu của User B; Bước 2: HTTP 403

---

### TC-074: Session cookie attributes
- **SR**: SR-NF-007
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. Đăng nhập thành công
  2. Inspect response headers
- **Expected Result**: `Set-Cookie` header có: `HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`

---

### TC-075: CSRF protection
- **SR**: SR-NF-011
- **Type**: Integration | **Priority**: P1
- **Steps**:
  1. POST `/app/activities` từ origin khác (simulated cross-site request) không có CSRF token
- **Expected Result**: HTTP 403 hoặc request bị từ chối; không tạo hoạt động

---

---

## Module 9: Edge Cases

### TC-076: Hoạt động qua nửa đêm (EC-01)
- **SR**: SRS Section 7 EC-01
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. POST activity với `start_time=2026-04-17T23:30`, `end_time=2026-04-18T00:30`
- **Expected Result**: Hoạt động được lưu; gán vào ngày `2026-04-17` (ngày của `start_time`); hiển thị cảnh báo nhẹ "Hoạt động kéo qua ngày hôm sau"

---

### TC-077: Tổng thời gian > 24h do overlap (EC-02)
- **SR**: SR-F-049, SRS Section 7 EC-02
- **Type**: Integration | **Priority**: P2
- **Preconditions**: Ngày 2026-04-17 có 3 hoạt động chồng nhau, tổng = 1500 phút
- **Steps**:
  1. GET `/app/reports/daily?date=2026-04-17`
- **Expected Result**: Báo cáo hiển thị tổng = 1500 phút; có cảnh báo rõ ràng; không crash; pie chart vẫn render được

---

### TC-078: Export khoảng thời gian không có dữ liệu (EC-04)
- (Xem TC-064 — đã cover)

---

### TC-079: Xóa danh mục, chọn danh mục mặc định làm thay thế (EC-05)
- **SR**: SR-F-045
- **Type**: Integration | **Priority**: P3
- **Steps**:
  1. Xóa danh mục tùy chỉnh "Thể thao" đang có 3 hoạt động
  2. Chọn "Làm việc" (danh mục mặc định) làm thay thế
- **Expected Result**: 3 hoạt động cập nhật sang "Làm việc"; "Thể thao" bị xóa — không có lỗi

---

### TC-080: Timezone — lưu UTC, hiển thị UTC+7
- **SR**: SR-NF-019
- **Type**: Unit | **Priority**: P1
- **Steps**:
  1. Tạo hoạt động với `start_time=2026-04-17T08:00` (UTC+7 = 01:00 UTC)
  2. Truy vấn trực tiếp DB
  3. Xem trên Dashboard
- **Expected Result**: DB lưu `start_time=2026-04-17T01:00:00Z`; Dashboard hiển thị `08:00`

---

---

## Tổng hợp Test Cases

| Module | Số TC | P1 | P2 | P3 | P4 |
|---|---|---|---|---|---|
| Authentication | 21 (TC-001–021) | 9 | 8 | 4 | 0 |
| Activity Management | 17 (TC-022–038) | 5 | 8 | 4 | 0 |
| Category Management | 9 (TC-039–047) | 3 | 4 | 2 | 0 |
| Reports | 7 (TC-048–054) | 1 | 4 | 2 | 0 |
| Trash | 6 (TC-055–060) | 2 | 3 | 1 | 0 |
| Export | 5 (TC-061–065) | 1 | 3 | 1 | 0 |
| Profile | 5 (TC-066–070) | 1 | 0 | 3 | 1 |
| Security | 5 (TC-071–075) | 5 | 0 | 0 | 0 |
| Edge Cases | 5 (TC-076–080) | 1 | 1 | 3 | 0 |
| **Tổng** | **80** | **28** | **31** | **20** | **1** |

---

## Thứ tự thực thi đề xuất

**Sprint 1 (Auth):** TC-001 → TC-021, TC-074, TC-075
**Sprint 2 (Activity):** TC-022 → TC-038, TC-071, TC-072, TC-073, TC-080
**Sprint 3 (Category):** TC-039 → TC-047
**Sprint 4 (Reports):** TC-048 → TC-054
**Sprint 5 (Trash):** TC-055 → TC-060, TC-059
**Sprint 6 (Export + Profile):** TC-061 → TC-070
**Regression:** TC-076 → TC-080 (edge cases chạy sau mỗi sprint liên quan)
