# Software Requirements Specification (SRS)
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-17
**Tác giả:** BA-SA Agent
**Dựa trên:** URS v1.0, Business Use Cases v1.1, System Architecture v1.0
**Chuẩn tham chiếu:** IEEE 830 (adapted)

---

## 1. Introduction

### 1.1 Mục đích
Tài liệu này định nghĩa chính xác những gì **hệ thống phải làm** — dành cho developer, tester và architect. Mỗi requirement là atomic, verifiable, traceable và consistent.

### 1.2 Phạm vi hệ thống
**Time Tracker** là ứng dụng web cho phép người dùng ghi lại, phân loại và phân tích cách sử dụng thời gian cá nhân hàng ngày. Hệ thống bao gồm:
- Backend: Go 1.22 + Gin (REST/SSR)
- Frontend: Go HTML Templates + HTMX + Tailwind CSS + Alpine.js
- Database: MySQL 8.x
- Email: Resend API (SMTP)
- Deployment: Railway.app

### 1.3 Định nghĩa & Viết tắt
| Ký hiệu | Ý nghĩa |
|---|---|
| SR-F-XXX | System Requirement — Functional |
| SR-NF-XXX | System Requirement — Non-Functional |
| UR-XX | User Requirement (từ URS v1.0) |
| BR-XX | Business Rule (từ Business Use Cases v1.1) |
| UC-XX | Use Case (từ Business Use Cases v1.1) |
| SHALL | Bắt buộc (mandatory) |
| SHOULD | Khuyến nghị (recommended) |
| MAY | Tùy chọn (optional) |
| UTC+7 | Asia/Bangkok — múi giờ Việt Nam |

### 1.4 Tài liệu tham chiếu
- User Requirements Specification v1.0
- Business Use Cases v1.1
- System Architecture Document v1.0
- Database Schema Document v1.0
- Project Vision Document v1.0

---

## 2. Overall Description

### 2.1 Bối cảnh hệ thống
Hệ thống là ứng dụng web độc lập (monolithic SSR). Người dùng tương tác qua trình duyệt web. Hệ thống giao tiếp với:
- **MySQL Database**: Lưu trữ toàn bộ dữ liệu
- **Resend Email API**: Gửi email xác thực và đặt lại mật khẩu
- **Client Browser**: Nhận HTML từ server, thực hiện HTMX requests

### 2.2 Chức năng tổng quan

| Module | Chức năng chính |
|---|---|
| **Auth** | Đăng ký, đăng nhập, session, quên/đổi mật khẩu |
| **Activity** | CRUD hoạt động, soft delete, log quá khứ |
| **Category** | Xem, thêm, sửa, xóa danh mục tùy chỉnh |
| **Report** | Daily / Weekly / Monthly summary + charts |
| **Trash** | Xem, khôi phục, xóa vĩnh viễn |
| **Export** | CSV và PDF theo khoảng thời gian |
| **Profile** | Xem và cập nhật thông tin cá nhân |
| **Scheduler** | Tự động xóa Trash và Session hết hạn |

### 2.3 Phân lớp người dùng (User Classes)

| Lớp | Mô tả | Quyền hạn |
|---|---|---|
| **Guest** | Người chưa đăng nhập | Chỉ xem landing page, form đăng ký/đăng nhập/quên mật khẩu |
| **User** | Người đã đăng nhập và xác thực email | Toàn bộ tính năng của ứng dụng, giới hạn trong dữ liệu của chính họ |
| **System** | Scheduler, background jobs | Xóa Trash hết hạn, dọn session cũ |

### 2.4 Ràng buộc thiết kế
- Không dùng JWT — authentication dựa trên session cookie (HttpOnly, Secure)
- Toàn bộ thời gian lưu vào DB ở dạng UTC; hiển thị ra ngoài theo UTC+7
- Không có SPA — Server-Side Rendering; dùng HTMX cho các partial update
- Không có shared data giữa các user — mọi query MUST filter theo `user_id`

---

## 3. Functional Requirements

---

### UC-01: Đăng ký tài khoản

#### SR-F-001
Hệ thống SHALL validate email đầu vào theo định dạng RFC 5322 trước khi tạo tài khoản.

**Acceptance Criteria:**
- Given người dùng nhập email `user@example.com` / When submit / Then validation pass
- Given người dùng nhập email `user@` / When submit / Then hiển thị lỗi inline "Email không hợp lệ", không gửi request

#### SR-F-002
Hệ thống SHALL validate mật khẩu: tối thiểu 8 ký tự, có ít nhất 1 chữ hoa, 1 chữ số.

**Acceptance Criteria:**
- Given `Password1` / When submit / Then pass
- Given `password` / When submit / Then lỗi "Mật khẩu cần ít nhất 8 ký tự, 1 chữ hoa và 1 chữ số"

#### SR-F-003
Hệ thống SHALL kiểm tra email đã tồn tại. Nếu trùng, SHALL hiển thị thông báo lỗi chung **không tiết lộ** email đã đăng ký hay chưa (BR-07).

**Acceptance Criteria:**
- Given email đã tồn tại / When submit / Then thông báo lỗi chung (không nói "Email đã tồn tại")

#### SR-F-004
Hệ thống SHALL tạo tài khoản ở trạng thái `unverified` và gửi email xác thực chứa token duy nhất.

#### SR-F-005
Hệ thống SHALL xử lý click link xác thực email:
- Nếu token hợp lệ và còn hiệu lực (< 24h): kích hoạt tài khoản, redirect đến trang đăng nhập
- Nếu token hết hạn: hiển thị tùy chọn gửi lại email xác thực
- Nếu token không tồn tại: hiển thị lỗi "Link không hợp lệ"

#### SR-F-006
Hệ thống SHALL tự động tạo 6 danh mục mặc định cho mỗi tài khoản mới được xác thực thành công:
`Làm việc`, `Di chuyển`, `Ăn uống`, `Giải trí`, `Mạng xã hội`, `Ngủ nghỉ`

---

### UC-02: Đăng nhập

#### SR-F-007
Hệ thống SHALL xác thực email và mật khẩu. Nếu thông tin sai, SHALL trả về thông báo lỗi chung "Email hoặc mật khẩu không đúng" mà không chỉ rõ trường nào sai (BR-07).

**Acceptance Criteria:**
- Given email sai / When submit / Then "Email hoặc mật khẩu không đúng"
- Given mật khẩu sai / When submit / Then "Email hoặc mật khẩu không đúng"

#### SR-F-008
Hệ thống SHALL khóa tài khoản 15 phút sau 5 lần đăng nhập sai liên tiếp (BR-05). Bộ đếm reset khi đăng nhập thành công.

**Acceptance Criteria:**
- Given lần thứ 5 sai / When submit / Then tài khoản bị khóa, thông báo thời gian mở khóa

#### SR-F-009
Hệ thống SHALL từ chối đăng nhập tài khoản chưa xác thực email, hiển thị tùy chọn gửi lại email xác thực (BR-03).

#### SR-F-010
Hệ thống SHALL tạo session cookie (HttpOnly, Secure, SameSite=Lax) với thời hạn 7 ngày kể từ lần hoạt động cuối (BR-06). Redirect đến Dashboard sau khi đăng nhập thành công.

---

### UC-03: Đăng xuất

#### SR-F-011
Hệ thống SHALL hủy session phía server và xóa cookie phía client khi người dùng đăng xuất. Token cũ SHALL bị vô hiệu hóa ngay lập tức.

**Acceptance Criteria:**
- Given người dùng đang đăng nhập / When click "Đăng xuất" / Then session bị hủy, redirect đến trang đăng nhập, truy cập lại trang protected thì bị chuyển hướng

---

### UC-04: Quên mật khẩu

#### SR-F-012
Hệ thống SHALL gửi email chứa link đặt lại mật khẩu khi nhận request quên mật khẩu. Link chứa token an toàn, có hiệu lực 1 giờ và chỉ dùng được 1 lần (BR-08, BR-09).

#### SR-F-013
Hệ thống SHALL luôn hiển thị thông báo "Nếu email tồn tại, bạn sẽ nhận được email" — bất kể email có tồn tại hay không (BR-11).

#### SR-F-014
Hệ thống SHALL cập nhật mật khẩu mới và hủy toàn bộ session cũ sau khi người dùng đặt lại mật khẩu thành công (BR-10).

#### SR-F-015
Hệ thống SHALL xử lý các trường hợp lỗi của link đặt lại mật khẩu:
- Token hết hạn → thông báo + form yêu cầu link mới
- Token đã dùng → thông báo "Link không còn hợp lệ"

---

### UC-05: Đổi mật khẩu

#### SR-F-016
Hệ thống SHALL yêu cầu xác nhận mật khẩu hiện tại trước khi cho phép đổi sang mật khẩu mới.

#### SR-F-017
Hệ thống SHALL từ chối nếu mật khẩu mới trùng với mật khẩu hiện tại (BR-12).

#### SR-F-018
Hệ thống SHALL hủy toàn bộ session khác (giữ lại session hiện tại) sau khi đổi mật khẩu thành công.

---

### UC-06: Log hoạt động mới

#### SR-F-019
Hệ thống SHALL cung cấp form log hoạt động với các trường bắt buộc: tên hoạt động, danh mục, thời gian bắt đầu, thời gian kết thúc. Trường ghi chú là tùy chọn.

#### SR-F-020
Hệ thống SHALL pre-fill thời gian bắt đầu = thời điểm hiện tại (theo UTC+7) khi form được mở.

#### SR-F-021
Hệ thống SHALL validate dữ liệu đầu vào:
- Tên hoạt động: không rỗng, tối đa 100 ký tự (BR-13)
- Ghi chú: tối đa 500 ký tự (BR-14)
- Thời gian kết thúc > thời gian bắt đầu (BR-15)
- Không cho phép thời gian ở tương lai (BR-17)
- Cho phép log quá khứ tối đa 2 năm (BR-16)

**Acceptance Criteria:**
- Given tên rỗng / When submit / Then lỗi inline "Tên hoạt động không được để trống", không gửi request
- Given thời gian kết thúc ≤ bắt đầu / When submit / Then lỗi "Thời gian kết thúc phải sau thời gian bắt đầu"
- Given thời gian ở tương lai / When submit / Then lỗi "Không thể log hoạt động ở tương lai"

#### SR-F-022
Hệ thống SHALL phát hiện hoạt động chồng thời gian với hoạt động khác cùng ngày. Nếu chồng, SHALL hiển thị toast notification cảnh báo nhưng vẫn lưu (BR-18).

#### SR-F-023
Hệ thống SHALL lưu hoạt động thành công và cập nhật danh sách ngay lập tức (HTMX partial update) — không reload toàn trang.

**Acceptance Criteria:**
- Given form hợp lệ / When submit / Then hoạt động xuất hiện trong danh sách mà không reload trang

---

### UC-07: Xem danh sách hoạt động trong ngày

#### SR-F-024
Hệ thống SHALL hiển thị danh sách hoạt động của ngày hôm nay theo mặc định khi người dùng vào Dashboard.

#### SR-F-025
Hệ thống SHALL chỉ truy vấn và hiển thị dữ liệu của người dùng đang đăng nhập (BR-19). Mọi query hoạt động MUST có điều kiện `WHERE user_id = :current_user_id`.

#### SR-F-026
Hệ thống SHALL sắp xếp danh sách hoạt động theo `start_time` tăng dần.

#### SR-F-027
Hệ thống SHALL hiển thị cho mỗi hoạt động: tên, danh mục (tên + màu sắc), thời gian bắt đầu, thời gian kết thúc, tổng thời gian (định dạng `Xh Ym`).

#### SR-F-028
Hệ thống SHALL hiển thị tóm tắt thời gian theo từng danh mục ở đầu trang (tổng giờ phút + %).

#### SR-F-029
Hệ thống SHALL cung cấp điều hướng ngày: nút "Hôm qua", "Hôm nay", "Ngày mai" và date picker.

#### SR-F-030
Hệ thống SHALL hiển thị trạng thái rỗng với thông điệp gợi ý nếu ngày đang xem chưa có hoạt động nào.

---

### UC-08: Chỉnh sửa hoạt động

#### SR-F-031
Hệ thống SHALL cho phép người dùng chỉnh sửa bất kỳ trường nào của hoạt động: tên, danh mục, thời gian bắt đầu, thời gian kết thúc, ghi chú.

#### SR-F-032
Hệ thống SHALL pre-fill form chỉnh sửa với giá trị hiện tại của hoạt động.

#### SR-F-033
Hệ thống SHALL kiểm tra authorization: người dùng chỉ được chỉnh sửa hoạt động thuộc `user_id` của chính mình (BR-21). Trả về HTTP 403 nếu vi phạm.

#### SR-F-034
Hệ thống SHALL áp dụng các validation tương tự SR-F-021 khi chỉnh sửa.

#### SR-F-035
Hệ thống SHALL cập nhật danh sách ngay lập tức sau khi chỉnh sửa thành công (HTMX partial update).

---

### UC-09: Xóa hoạt động (Soft Delete)

#### SR-F-036
Hệ thống SHALL hiển thị dialog xác nhận trước khi xóa, thông báo rõ "Hoạt động sẽ được chuyển vào Trash và tự xóa sau 30 ngày".

#### SR-F-037
Hệ thống SHALL thực hiện soft delete: set `deleted_at = NOW()`, không xóa record khỏi DB. Hoạt động biến mất khỏi danh sách chính ngay lập tức.

#### SR-F-038
Hệ thống SHALL kiểm tra authorization tương tự SR-F-033 cho thao tác xóa (BR-21).

#### SR-F-039
Hệ thống SHALL loại trừ hoạt động có `deleted_at IS NOT NULL` khỏi tất cả query danh sách hoạt động, báo cáo, và export (BR-35, BR-42).

---

### UC-10–UC-13: Quản lý danh mục

#### SR-F-040
Hệ thống SHALL hiển thị danh sách đầy đủ danh mục của người dùng: 6 danh mục mặc định + danh mục tùy chỉnh. Mỗi danh mục hiển thị: tên, màu sắc, icon, loại (mặc định / tùy chỉnh).

#### SR-F-041
Hệ thống SHALL cho phép tạo danh mục tùy chỉnh với: tên (bắt buộc, max 50 ký tự), màu sắc, icon từ bộ icon có sẵn (BR-23).

#### SR-F-042
Hệ thống SHALL kiểm tra tên danh mục không trùng với danh mục hiện có của người dùng (không phân biệt hoa thường) (BR-25).

**Acceptance Criteria:**
- Given đã có danh mục "Học tập" / When tạo danh mục "học tập" / Then lỗi "Tên danh mục đã tồn tại"

#### SR-F-043
Hệ thống SHALL từ chối sửa hoặc xóa danh mục mặc định (BR-26, BR-27). Trả về HTTP 403 nếu thử.

#### SR-F-044
Hệ thống SHALL cho phép chỉnh sửa danh mục tùy chỉnh: tên, màu sắc, icon. Áp dụng validation BR-25 khi đổi tên.

#### SR-F-045
Khi xóa danh mục tùy chỉnh đang được dùng bởi ít nhất 1 hoạt động:
- Hệ thống SHALL hiển thị cảnh báo và yêu cầu chọn danh mục thay thế
- Hệ thống SHALL cập nhật tất cả hoạt động liên quan sang danh mục thay thế trước khi xóa (BR-28)

#### SR-F-046
Khi xóa danh mục tùy chỉnh chưa được dùng bởi hoạt động nào, hệ thống SHALL xóa trực tiếp sau dialog xác nhận.

---

### UC-14: Báo cáo ngày (Daily Summary)

#### SR-F-047
Hệ thống SHALL tính tổng thời gian theo từng danh mục cho ngày được chọn:
`total_minutes = SUM(TIMESTAMPDIFF(MINUTE, start_time, end_time))` theo danh mục (BR-29).

#### SR-F-048
Hệ thống SHALL hiển thị:
- Tổng thời gian đã log trong ngày (định dạng Xh Ym)
- Thời gian theo từng danh mục: số giờ phút + phần trăm
- Biểu đồ tròn phân bổ thời gian theo danh mục
- Danh sách đầy đủ hoạt động theo thứ tự thời gian

#### SR-F-049
Hệ thống SHALL hiển thị cảnh báo khi tổng thời gian trong ngày vượt quá 1440 phút (24 giờ) do hoạt động chồng nhau (BR-30).

#### SR-F-050
Hệ thống SHALL hiển thị trạng thái rỗng nếu ngày được chọn không có hoạt động.

---

### UC-15: Báo cáo tuần (Weekly Summary)

#### SR-F-051
Hệ thống SHALL tính toán dữ liệu báo cáo theo tuần (Thứ Hai đến Chủ Nhật) (BR-31).

#### SR-F-052
Hệ thống SHALL hiển thị:
- Biểu đồ cột stacked bar: thời gian theo ngày trong tuần, phân theo danh mục
- Tổng thời gian mỗi danh mục trong cả tuần
- Ngày có nhiều thời gian làm việc nhất / ít nhất

#### SR-F-053
Hệ thống SHALL cung cấp điều hướng tuần: "Tuần trước", "Tuần này", "Tuần sau" và date picker.

---

### UC-19: Báo cáo tháng (Monthly Summary)

#### SR-F-054
Hệ thống SHALL tính toán dữ liệu báo cáo theo tháng dương lịch (ngày 1 đến ngày cuối tháng) (BR-36).

#### SR-F-055
Hệ thống SHALL hiển thị:
- Tổng thời gian theo từng danh mục trong tháng (giờ + %)
- Biểu đồ tròn phân bổ danh mục
- Biểu đồ xu hướng theo tuần trong tháng (4–5 điểm)
- Ngày hoạt động nhiều nhất / ít nhất trong tháng

#### SR-F-056
Hệ thống SHALL cung cấp điều hướng tháng: "Tháng trước", "Tháng này", "Tháng sau" và month picker.

---

### UC-17: Trash

#### SR-F-057
Hệ thống SHALL hiển thị danh sách hoạt động trong Trash của người dùng, sắp xếp theo `deleted_at` giảm dần. Mỗi mục hiển thị số ngày còn lại trước khi xóa vĩnh viễn.

#### SR-F-058
Hệ thống SHALL khôi phục hoạt động từ Trash: set `deleted_at = NULL`. Hoạt động xuất hiện lại đúng vị trí theo `start_time` trong ngày tương ứng (BR-37).

#### SR-F-059
Hệ thống SHALL cho phép xóa vĩnh viễn từng mục trong Trash sau dialog xác nhận.

#### SR-F-060
Hệ thống SHALL cho phép xóa toàn bộ Trash ("Empty Trash") sau dialog xác nhận. Không thể hoàn tác thao tác này.

#### SR-F-061
Hệ thống (Scheduler) SHALL tự động xóa vĩnh viễn các hoạt động có `deleted_at < NOW() - INTERVAL 30 DAY` mỗi ngày lúc 2:00 AM UTC+7 (BR-38).

---

### UC-18: Export dữ liệu

#### SR-F-062
Hệ thống SHALL cho phép người dùng chọn khoảng thời gian export (ngày bắt đầu – ngày kết thúc) và định dạng (CSV hoặc PDF).

#### SR-F-063
Hệ thống SHALL từ chối khoảng thời gian vượt quá 365 ngày (BR-41).

#### SR-F-064
Hệ thống SHALL tạo file CSV với các cột (BR-39):
`date, activity_name, category, start_time, end_time, total_minutes, note`
Thời gian theo UTC+7 (BR-43).

#### SR-F-065
Hệ thống SHALL tạo file PDF gồm (BR-40):
- Summary theo danh mục (tổng giờ, %)
- Danh sách hoạt động có định dạng đẹp (header, footer, phân trang nếu cần)
Thời gian theo UTC+7 (BR-43).

#### SR-F-066
Hệ thống SHALL loại trừ hoạt động trong Trash khỏi file export (BR-42).

#### SR-F-067
Hệ thống SHALL trả file export trực tiếp dưới dạng download response (Content-Disposition: attachment), không gửi qua email.

---

### UC-16: Hồ sơ cá nhân

#### SR-F-068
Hệ thống SHALL hiển thị: tên hiển thị (editable), email (read-only), timezone (hiển thị UTC+7, read-only).

#### SR-F-069
Hệ thống SHALL validate tên hiển thị: không rỗng, tối đa 50 ký tự (BR-33).

#### SR-F-070
Hệ thống SHALL từ chối bất kỳ thay đổi nào đến trường email (BR-32).

---

### Auth Middleware

#### SR-F-071
Hệ thống SHALL bảo vệ tất cả route trong nhóm `/app/*` bằng middleware xác thực session. Truy cập không có session hợp lệ SHALL bị redirect về trang đăng nhập.

#### SR-F-072
Hệ thống SHALL gia hạn session (sliding window) mỗi khi người dùng thực hiện request, giữ thời hạn 7 ngày kể từ lần hoạt động cuối (BR-06).

---

### Scheduler (Background Jobs)

#### SR-F-073
Hệ thống SHALL chạy job dọn Trash hàng ngày lúc 2:00 AM UTC+7: xóa vĩnh viễn (`DELETE`) các hoạt động có `deleted_at < NOW() - INTERVAL 30 DAY`.

#### SR-F-074
Hệ thống SHALL chạy job dọn session hàng ngày lúc 3:00 AM UTC+7: xóa các session đã hết hạn.

---

## 4. Non-Functional Requirements

| Mã | Loại | Mô tả | Chỉ số | Cách kiểm tra |
|---|---|---|---|---|
| SR-NF-001 | Performance | Thời gian phản hồi trang Dashboard | P95 < 500ms với ≤ 100 hoạt động/ngày | Load test với k6 |
| SR-NF-002 | Performance | Thời gian submit form log hoạt động | P95 < 300ms | Load test |
| SR-NF-003 | Performance | Thời gian tạo báo cáo tuần/tháng | P95 < 1000ms | Load test |
| SR-NF-004 | Performance | Thời gian tạo file export CSV/PDF | P95 < 3000ms | Load test |
| SR-NF-005 | Usability | Luồng log hoạt động | Hoàn thành trong < 10 giây (từ click đến xong) | Manual test |
| SR-NF-006 | Security | HTTPS | Toàn bộ traffic qua HTTPS; HTTP redirect sang HTTPS | SSL scan |
| SR-NF-007 | Security | Session cookie | HttpOnly=true, Secure=true, SameSite=Lax | Dev tools kiểm tra |
| SR-NF-008 | Security | Password hashing | Bcrypt cost factor ≥ 10 | Code review |
| SR-NF-009 | Security | SQL Injection | Sử dụng parameterized queries 100% | Code review + SAST |
| SR-NF-010 | Security | XSS | HTML escape toàn bộ user input trong template | Manual test + ZAP scan |
| SR-NF-011 | Security | CSRF | CSRF token trên tất cả form POST/PUT/DELETE | Penetration test |
| SR-NF-012 | Security | Authorization | Mọi query có điều kiện `user_id` — không thể đọc/sửa/xóa dữ liệu người khác | Unit test mỗi handler |
| SR-NF-013 | Reliability | Uptime | ≥ 99% trong giờ sử dụng bình thường | Railway.app monitoring |
| SR-NF-014 | Reliability | Scheduler | Background jobs chạy đúng lịch, không bị bỏ lỡ khi server restart | Log monitoring |
| SR-NF-015 | Maintainability | Code readability | Code có thể đọc lại và hiểu sau 1 tháng không cần giải thích | Peer review |
| SR-NF-016 | Maintainability | Test coverage | Unit test cho service layer ≥ 70% | `go test -cover` |
| SR-NF-017 | Compatibility | Browser support | Chrome 100+, Firefox 100+, Safari 15+ | Manual test |
| SR-NF-018 | Compatibility | Responsive | Sử dụng được trên màn hình ≥ 768px (tablet trở lên) | Manual test |
| SR-NF-019 | Data | Timezone | Lưu UTC trong DB, hiển thị UTC+7 ở tất cả giao diện | Unit test timezone conversion |
| SR-NF-020 | Data | Retention | Dữ liệu Trash tự xóa sau 30 ngày; session hết hạn sau 7 ngày không hoạt động | Integration test |

---

## 5. System Constraints & Assumptions

### 5.1 Ràng buộc hệ thống
- **Go 1.22**: Phiên bản ngôn ngữ tối thiểu
- **MySQL 8.x**: Phiên bản DB tối thiểu (dùng tính năng `JSON`, `generated columns` nếu cần)
- **Charset**: `utf8mb4` để hỗ trợ emoji và tiếng Việt đầy đủ
- **Deployment**: Railway.app — phải dùng environment variables cho tất cả secret
- **No file storage**: Không có hệ thống lưu file (avatar, attachment) — chỉ lưu DB
- **Single instance**: Không có horizontal scaling — một instance duy nhất

### 5.2 Giả định kỹ thuật
- Người dùng dùng múi giờ UTC+7 (Việt Nam) — timezone cố định, không thể thay đổi
- Server và DB cùng timezone setting UTC để tránh nhầm lẫn
- Email service (Resend) có uptime ≥ 99% — nếu email không gửi được, hệ thống ghi log lỗi và cho phép retry
- Không có CDN — static files phục vụ trực tiếp từ server
- Số người dùng giai đoạn MVP: 1–10 người (cá nhân + thân quen)

### 5.3 Phụ thuộc bên ngoài
| Phụ thuộc | Mục đích | Fallback |
|---|---|---|
| Resend API | Gửi email xác thực, reset mật khẩu | Log lỗi, hiển thị thông báo "Vui lòng thử lại" |
| Railway.app | Hosting và managed MySQL | N/A (deployment target cố định) |
| Tailwind CSS CDN | CSS framework | N/A (có thể bundle sau) |
| HTMX CDN | Partial page updates | N/A (có thể bundle sau) |

---

## 6. Traceability Matrix

### UR → SR Mapping

| User Requirement | System Requirements |
|---|---|
| UR-01 (Đăng ký) | SR-F-001, SR-F-002, SR-F-003, SR-F-004, SR-F-005, SR-F-006 |
| UR-02 (Đăng nhập/xuất) | SR-F-007, SR-F-008, SR-F-009, SR-F-010, SR-F-011, SR-F-071, SR-F-072 |
| UR-03 (Quên mật khẩu) | SR-F-012, SR-F-013, SR-F-014, SR-F-015 |
| UR-04 (Đổi mật khẩu) | SR-F-016, SR-F-017, SR-F-018 |
| UR-05 (Log hoạt động) | SR-F-019, SR-F-020, SR-F-021, SR-F-022, SR-F-023 |
| UR-06 (Xem danh sách ngày) | SR-F-024, SR-F-025, SR-F-026, SR-F-027, SR-F-028, SR-F-029, SR-F-030 |
| UR-07 (Chỉnh sửa hoạt động) | SR-F-031, SR-F-032, SR-F-033, SR-F-034, SR-F-035 |
| UR-08 (Xóa hoạt động) | SR-F-036, SR-F-037, SR-F-038, SR-F-039 |
| UR-09 (Quản lý danh mục) | SR-F-040, SR-F-041, SR-F-042, SR-F-043, SR-F-044, SR-F-045, SR-F-046 |
| UR-10 (Báo cáo ngày) | SR-F-047, SR-F-048, SR-F-049, SR-F-050 |
| UR-11 (Báo cáo tuần) | SR-F-051, SR-F-052, SR-F-053 |
| UR-12 (Báo cáo tháng) | SR-F-054, SR-F-055, SR-F-056 |
| UR-13 (Trash & khôi phục) | SR-F-057, SR-F-058, SR-F-059, SR-F-060, SR-F-061, SR-F-073, SR-F-074 |
| UR-14 (Export) | SR-F-062, SR-F-063, SR-F-064, SR-F-065, SR-F-066, SR-F-067 |
| UR-15 (Hồ sơ cá nhân) | SR-F-068, SR-F-069, SR-F-070 |

### SR → BR Mapping (Functional)

| System Requirement | Business Rule(s) |
|---|---|
| SR-F-001 | BR-01 |
| SR-F-002 | BR-02 |
| SR-F-003 | BR-07 |
| SR-F-005 | BR-04 |
| SR-F-006 | (implicit từ UC-01 postcondition) |
| SR-F-008 | BR-05 |
| SR-F-009 | BR-03 |
| SR-F-010 | BR-06 |
| SR-F-012 | BR-08, BR-09 |
| SR-F-013 | BR-11 |
| SR-F-014 | BR-10 |
| SR-F-017 | BR-12 |
| SR-F-021 | BR-13, BR-14, BR-15, BR-16, BR-17 |
| SR-F-022 | BR-18 |
| SR-F-025 | BR-19 |
| SR-F-033 | BR-21 |
| SR-F-037 | BR-22 |
| SR-F-039 | BR-35, BR-42 |
| SR-F-041 | BR-23 |
| SR-F-042 | BR-25 |
| SR-F-043 | BR-26, BR-27 |
| SR-F-045 | BR-28 |
| SR-F-047 | BR-29 |
| SR-F-049 | BR-30 |
| SR-F-051 | BR-31 |
| SR-F-054 | BR-36 |
| SR-F-058 | BR-37 |
| SR-F-061 | BR-38 |
| SR-F-063 | BR-41 |
| SR-F-064 | BR-39, BR-43 |
| SR-F-065 | BR-40, BR-43 |
| SR-F-066 | BR-42 |
| SR-F-068–070 | BR-32, BR-33, BR-34 |

---

## 7. Edge Cases cần xử lý

| ID | Tình huống | Yêu cầu hệ thống |
|---|---|---|
| EC-01 | Hoạt động qua nửa đêm (bắt đầu 23:30, kết thúc 00:30) | Ghi vào ngày của `start_time`. Hiển thị cảnh báo nhẹ "Hoạt động kéo qua ngày hôm sau". |
| EC-02 | Tổng thời gian trong ngày > 24h do overlapping | Vẫn lưu và tính đúng. Hiển thị cảnh báo trên trang báo cáo ngày (SR-F-049). |
| EC-03 | Người dùng mở 2 tab cùng lúc | HTMX load data theo từng request độc lập; không có conflict ở server-side. |
| EC-04 | Export với khoảng thời gian không có dữ liệu | Hiển thị thông báo "Không có dữ liệu trong khoảng thời gian này", không tạo file rỗng. |
| EC-05 | Xóa danh mục và chọn danh mục thay thế là danh mục mặc định | Cho phép — danh mục mặc định là lựa chọn hợp lệ. |
