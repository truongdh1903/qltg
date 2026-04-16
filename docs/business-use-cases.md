# Business Use Case Document
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.1
**Ngày:** 2026-04-16
**Tác giả:** BA Agent
**Dựa trên:** Project Vision Document v1.0
**Changelog v1.1:** Cập nhật Business Rules sau khi confirm Open Questions. Bổ sung UC-17 (Trash), UC-18 (Export), UC-19 (Monthly Report).

---

## 1. Actor List

| Actor | Loại | Mô tả |
|---|---|---|
| **Khách (Guest)** | Primary | Người dùng chưa đăng ký, chỉ xem được landing page |
| **Người dùng (User)** | Primary | Người đã đăng ký và đăng nhập — actor chính của toàn bộ hệ thống |
| **Hệ thống (System)** | Secondary | Thực hiện các hành động tự động: tính toán thống kê, validate dữ liệu |

---

## 2. Use Case Overview

```
[Khách]
  └── UC-01: Đăng ký tài khoản
  └── UC-02: Đăng nhập

[Người dùng]
  ├── Xác thực
  │     └── UC-03: Đăng xuất
  │     └── UC-04: Quên mật khẩu
  │     └── UC-05: Đổi mật khẩu
  │
  ├── Quản lý hoạt động (CORE)
  │     └── UC-06: Log hoạt động mới
  │     └── UC-07: Xem danh sách hoạt động trong ngày
  │     └── UC-08: Chỉnh sửa hoạt động
  │     └── UC-09: Xóa hoạt động
  │
  ├── Quản lý danh mục
  │     └── UC-10: Xem danh mục hoạt động
  │     └── UC-11: Thêm danh mục tùy chỉnh
  │     └── UC-12: Chỉnh sửa danh mục tùy chỉnh
  │     └── UC-13: Xóa danh mục tùy chỉnh
  │
  ├── Báo cáo
  │     └── UC-14: Xem báo cáo ngày (Daily Summary)
  │     └── UC-15: Xem báo cáo tuần (Weekly Summary)
  │     └── UC-19: Xem báo cáo tháng (Monthly Summary)
  │
  ├── Quản lý Trash
  │     └── UC-17: Xem và khôi phục hoạt động đã xóa
  │
  ├── Export dữ liệu
  │     └── UC-18: Export dữ liệu ra CSV / PDF
  │
  └── Hồ sơ cá nhân
        └── UC-16: Xem và chỉnh sửa thông tin cá nhân
```

---

## 3. Use Case Details

---

### UC-01: Đăng ký tài khoản

**Actor:** Khách (Guest)
**Mục tiêu:** Tạo tài khoản mới để sử dụng hệ thống

**Điều kiện tiên quyết:**
- Người dùng chưa có tài khoản
- Người dùng đang ở trang đăng ký

**Luồng chính (Happy Path):**
1. Khách điền email, mật khẩu, xác nhận mật khẩu
2. Hệ thống validate định dạng email và độ mạnh mật khẩu
3. Hệ thống kiểm tra email chưa tồn tại trong hệ thống
4. Hệ thống tạo tài khoản và gửi email xác thực
5. Khách kiểm tra email và click link xác thực
6. Hệ thống kích hoạt tài khoản và chuyển hướng đến trang đăng nhập

**Luồng thay thế / Ngoại lệ:**
- [A1] Email đã tồn tại → hiển thị thông báo lỗi, không tiết lộ email có tồn tại hay không (bảo mật)
- [A2] Mật khẩu không khớp → hiển thị lỗi inline tại field xác nhận
- [A3] Email không đúng định dạng → hiển thị lỗi inline
- [A4] Link xác thực hết hạn (sau 24h) → hiển thị tùy chọn gửi lại email
- [E1] Lỗi hệ thống khi gửi email → thông báo lỗi, cho phép thử lại

**Hậu điều kiện:**
- Tài khoản được tạo ở trạng thái "chưa xác thực"
- Sau khi xác thực email → trạng thái "hoạt động"
- 6 danh mục mặc định được tạo tự động cho tài khoản mới

**Business Rules:**
- BR-01: Email phải là định dạng hợp lệ (RFC 5322)
- BR-02: Mật khẩu tối thiểu 8 ký tự, có ít nhất 1 chữ hoa, 1 số
- BR-03: Tài khoản chưa xác thực email không thể đăng nhập
- BR-04: Link xác thực có hiệu lực 24 giờ

---

### UC-02: Đăng nhập

**Actor:** Khách (Guest)
**Mục tiêu:** Truy cập vào hệ thống với tài khoản đã có

**Điều kiện tiên quyết:**
- Người dùng đã có tài khoản và đã xác thực email

**Luồng chính (Happy Path):**
1. Khách nhập email và mật khẩu
2. Hệ thống xác thực thông tin đăng nhập
3. Hệ thống tạo session / token
4. Hệ thống chuyển hướng đến trang Dashboard (danh sách hoạt động hôm nay)

**Luồng thay thế / Ngoại lệ:**
- [A1] Sai email hoặc mật khẩu → thông báo lỗi chung ("Email hoặc mật khẩu không đúng"), không chỉ rõ cái nào sai (bảo mật)
- [A2] Tài khoản chưa xác thực email → thông báo và cho phép gửi lại email xác thực
- [A3] Tài khoản bị khóa (sau 5 lần sai) → thông báo và hướng dẫn mở khóa
- [E1] Lỗi hệ thống → thông báo lỗi chung, log lỗi phía server

**Hậu điều kiện:**
- Session hợp lệ được tạo
- Người dùng được chuyển đến Dashboard

**Business Rules:**
- BR-05: Sau 5 lần đăng nhập sai liên tiếp → khóa tài khoản tạm thời 15 phút
- BR-06: Session có thời hạn (cần làm rõ: bao lâu? xem Open Questions)
- BR-07: Không hiển thị chi tiết lỗi xác thực để tránh user enumeration attack

---

### UC-03: Đăng xuất

**Actor:** Người dùng (User)
**Mục tiêu:** Kết thúc phiên làm việc an toàn

**Điều kiện tiên quyết:** Người dùng đang đăng nhập

**Luồng chính:**
1. Người dùng click "Đăng xuất"
2. Hệ thống hủy session / token phía server
3. Hệ thống xóa token phía client
4. Hệ thống chuyển hướng về trang đăng nhập

**Hậu điều kiện:** Session bị hủy hoàn toàn, không thể tái sử dụng token cũ

---

### UC-04: Quên mật khẩu

**Actor:** Khách (Guest)
**Mục tiêu:** Lấy lại quyền truy cập khi quên mật khẩu

**Điều kiện tiên quyết:** Người dùng có tài khoản nhưng không nhớ mật khẩu

**Luồng chính:**
1. Khách nhập email tại trang "Quên mật khẩu"
2. Hệ thống gửi email chứa link đặt lại mật khẩu (dù email có tồn tại hay không — bảo mật)
3. Khách click link trong email
4. Hệ thống xác thực link còn hiệu lực
5. Khách nhập mật khẩu mới và xác nhận
6. Hệ thống cập nhật mật khẩu, hủy toàn bộ session cũ
7. Hệ thống chuyển hướng về trang đăng nhập

**Luồng thay thế:**
- [A1] Link hết hạn → thông báo và cho phép yêu cầu link mới
- [A2] Link đã được sử dụng → thông báo không hợp lệ

**Business Rules:**
- BR-08: Link đặt lại mật khẩu có hiệu lực 1 giờ
- BR-09: Mỗi link chỉ dùng được 1 lần
- BR-10: Sau khi đổi mật khẩu thành công, toàn bộ session cũ bị hủy
- BR-11: Luôn hiển thị thông báo "Nếu email tồn tại, bạn sẽ nhận được email" — không xác nhận email có tồn tại hay không

---

### UC-05: Đổi mật khẩu

**Actor:** Người dùng (User)
**Mục tiêu:** Thay đổi mật khẩu khi đang đăng nhập

**Điều kiện tiên quyết:** Người dùng đang đăng nhập

**Luồng chính:**
1. Người dùng vào phần Cài đặt → Đổi mật khẩu
2. Nhập mật khẩu hiện tại, mật khẩu mới, xác nhận mật khẩu mới
3. Hệ thống xác thực mật khẩu hiện tại
4. Hệ thống cập nhật mật khẩu mới
5. Hệ thống hủy toàn bộ session khác (giữ lại session hiện tại)

**Luồng thay thế:**
- [A1] Mật khẩu hiện tại sai → thông báo lỗi
- [A2] Mật khẩu mới trùng mật khẩu cũ → thông báo yêu cầu mật khẩu khác

**Business Rules:**
- BR-02: Áp dụng (mật khẩu mới phải đủ mạnh)
- BR-12: Mật khẩu mới không được trùng mật khẩu hiện tại

---

### UC-06: Log hoạt động mới ⭐ (CORE)

**Actor:** Người dùng (User)
**Mục tiêu:** Ghi lại một hoạt động đã thực hiện trong ngày

**Điều kiện tiên quyết:** Người dùng đang đăng nhập

**Luồng chính (Happy Path):**
1. Người dùng click nút "Log hoạt động" (nút nổi bật, dễ thấy)
2. Hệ thống hiển thị form nhanh với các trường: Tên hoạt động, Danh mục, Thời gian bắt đầu, Thời gian kết thúc, Ghi chú (tùy chọn)
3. Người dùng điền thông tin (hệ thống gợi ý thời gian bắt đầu = thời điểm hiện tại)
4. Người dùng chọn danh mục từ dropdown
5. Người dùng submit form
6. Hệ thống validate dữ liệu
7. Hệ thống lưu hoạt động và cập nhật danh sách ngay lập tức (không reload trang)
8. Hiển thị thông báo thành công ngắn gọn

**Luồng thay thế / Ngoại lệ:**
- [A1] Tên hoạt động bỏ trống → thông báo lỗi inline
- [A2] Không chọn danh mục → thông báo lỗi inline
- [A3] Thời gian kết thúc ≤ thời gian bắt đầu → thông báo lỗi inline
- [A4] Thời gian thuộc ngày khác (hôm qua, hôm kia) → vẫn cho phép log, hiển thị cảnh báo nhẹ
- [E1] Lỗi lưu dữ liệu → thông báo lỗi, giữ nguyên form để người dùng thử lại

**Hậu điều kiện:**
- Hoạt động được lưu vào database
- Danh sách hoạt động trong ngày được cập nhật
- Báo cáo ngày/tuần phản ánh dữ liệu mới

**Business Rules:**
- BR-13: Tên hoạt động tối đa 100 ký tự
- BR-14: Ghi chú tối đa 500 ký tự
- BR-15: Thời gian kết thúc phải sau thời gian bắt đầu
- BR-16: Cho phép log hoạt động trong quá khứ (tối đa bao nhiêu ngày? → xem Open Questions)
- BR-17: Không cho phép log hoạt động ở thời điểm tương lai
- BR-18: Các hoạt động được phép chồng nhau về thời gian (overlapping) — hệ thống cảnh báo nhưng không ngăn chặn (xem Open Questions)

---

### UC-07: Xem danh sách hoạt động trong ngày

**Actor:** Người dùng (User)
**Mục tiêu:** Xem toàn bộ hoạt động đã log trong một ngày cụ thể

**Điều kiện tiên quyết:** Người dùng đang đăng nhập

**Luồng chính:**
1. Người dùng vào trang Dashboard (mặc định hiển thị ngày hôm nay)
2. Hệ thống truy vấn và hiển thị danh sách hoạt động của ngày, sắp xếp theo thời gian bắt đầu tăng dần
3. Mỗi hoạt động hiển thị: tên, danh mục (màu/icon), thời gian bắt đầu–kết thúc, tổng thời gian
4. Hiển thị tổng thời gian theo từng danh mục ở đầu trang (summary nhanh)
5. Người dùng có thể điều hướng sang ngày trước / sau

**Luồng thay thế:**
- [A1] Không có hoạt động nào trong ngày → hiển thị trạng thái trống với gợi ý "Log hoạt động đầu tiên"

**Business Rules:**
- BR-19: Chỉ hiển thị dữ liệu của người dùng đang đăng nhập
- BR-20: Mặc định hiển thị ngày hiện tại theo timezone của người dùng (xem Open Questions)

---

### UC-08: Chỉnh sửa hoạt động

**Actor:** Người dùng (User)
**Mục tiêu:** Cập nhật thông tin của một hoạt động đã log

**Điều kiện tiên quyết:**
- Người dùng đang đăng nhập
- Hoạt động cần chỉnh sửa thuộc về người dùng này

**Luồng chính:**
1. Người dùng click vào hoạt động trong danh sách → chọn "Chỉnh sửa"
2. Hệ thống hiển thị form với dữ liệu hiện tại đã điền sẵn
3. Người dùng thay đổi thông tin cần thiết
4. Người dùng submit
5. Hệ thống validate và lưu thay đổi
6. Danh sách cập nhật ngay lập tức

**Luồng thay thế:**
- [A1] Validate thất bại → hiển thị lỗi inline, không lưu
- [A2] Người dùng hủy chỉnh sửa → đóng form, dữ liệu không thay đổi

**Business Rules:**
- BR-15, BR-13, BR-14: Áp dụng tương tự UC-06
- BR-21: Người dùng chỉ được chỉnh sửa hoạt động của chính mình (authorization)

---

### UC-09: Xóa hoạt động (Soft Delete)

**Actor:** Người dùng (User)
**Mục tiêu:** Xóa một hoạt động đã log — hoạt động được chuyển vào Trash, chưa xóa vĩnh viễn

**Điều kiện tiên quyết:**
- Người dùng đang đăng nhập
- Hoạt động cần xóa thuộc về người dùng này

**Luồng chính:**
1. Người dùng chọn "Xóa" trên một hoạt động
2. Hệ thống hiển thị dialog xác nhận ("Hoạt động sẽ được chuyển vào Trash và tự xóa sau 30 ngày")
3. Người dùng xác nhận
4. Hệ thống đánh dấu hoạt động là đã xóa (soft delete), ẩn khỏi danh sách chính và cập nhật ngay
5. Báo cáo không tính hoạt động trong Trash

**Luồng thay thế:**
- [A1] Người dùng hủy tại dialog → không xóa, đóng dialog

**Business Rules:**
- BR-22: Xóa là soft delete — hoạt động chuyển vào Trash, giữ 30 ngày rồi tự xóa vĩnh viễn
- BR-21: Người dùng chỉ được xóa hoạt động của chính mình
- BR-35: Hoạt động trong Trash không được tính vào báo cáo

---

### UC-10: Xem danh mục hoạt động

**Actor:** Người dùng (User)
**Mục tiêu:** Xem danh sách tất cả danh mục (mặc định + tùy chỉnh)

**Luồng chính:**
1. Người dùng vào trang Cài đặt → Danh mục
2. Hệ thống hiển thị 6 danh mục mặc định và danh mục tùy chỉnh (nếu có)
3. Mỗi danh mục hiển thị: tên, màu sắc, icon, loại (mặc định / tùy chỉnh)

**Danh mục mặc định (6):**
| # | Tên | Ý nghĩa |
|---|---|---|
| 1 | Làm việc | Công việc chuyên môn, học tập có chủ đích |
| 2 | Di chuyển | Đi lại, commute |
| 3 | Ăn uống | Ăn, uống, nghỉ giải lao |
| 4 | Giải trí | Xem phim, đọc sách, chơi game |
| 5 | Mạng xã hội | Facebook, TikTok, YouTube, lướt web |
| 6 | Ngủ nghỉ | Ngủ, nghỉ ngơi, nằm không |

---

### UC-11: Thêm danh mục tùy chỉnh

**Actor:** Người dùng (User)
**Mục tiêu:** Tạo danh mục mới phù hợp với nhu cầu cá nhân

**Luồng chính:**
1. Người dùng click "Thêm danh mục"
2. Nhập tên danh mục, chọn màu sắc, chọn icon (từ bộ icon có sẵn)
3. Hệ thống validate tên không trùng với danh mục hiện có
4. Hệ thống lưu danh mục mới

**Luồng thay thế:**
- [A1] Tên trùng với danh mục đã có (không phân biệt hoa thường) → thông báo lỗi

**Business Rules:**
- BR-23: Tên danh mục tối đa 50 ký tự
- BR-24: Mỗi người dùng tối đa 20 danh mục tùy chỉnh (xem Open Questions)
- BR-25: Tên danh mục không phân biệt hoa thường khi kiểm tra trùng lặp

---

### UC-12: Chỉnh sửa danh mục tùy chỉnh

**Actor:** Người dùng (User)
**Mục tiêu:** Cập nhật tên, màu sắc hoặc icon của danh mục tùy chỉnh

**Business Rules:**
- BR-26: Chỉ được chỉnh sửa danh mục tùy chỉnh, không được sửa danh mục mặc định
- BR-25: Áp dụng khi đổi tên

---

### UC-13: Xóa danh mục tùy chỉnh

**Actor:** Người dùng (User)
**Mục tiêu:** Xóa một danh mục tùy chỉnh không còn cần thiết

**Luồng chính:**
1. Người dùng chọn "Xóa" trên danh mục tùy chỉnh
2. Nếu danh mục đang được dùng bởi các hoạt động đã log: hệ thống hiển thị cảnh báo và yêu cầu chọn danh mục thay thế
3. Người dùng chọn danh mục thay thế → hệ thống cập nhật tất cả hoạt động liên quan → xóa danh mục
4. Nếu danh mục chưa được dùng: xóa trực tiếp sau khi xác nhận

**Business Rules:**
- BR-27: Không được xóa danh mục mặc định
- BR-28: Khi xóa danh mục đang được sử dụng, bắt buộc phải chọn danh mục thay thế (không được để hoạt động không có danh mục)

---

### UC-14: Xem báo cáo ngày (Daily Summary) ⭐

**Actor:** Người dùng (User)
**Mục tiêu:** Xem tổng kết thời gian sử dụng trong một ngày cụ thể

**Điều kiện tiên quyết:** Có ít nhất 1 hoạt động được log trong ngày đó

**Luồng chính:**
1. Người dùng vào trang Báo cáo → chọn ngày
2. Hệ thống tính toán và hiển thị:
   - Tổng thời gian đã log trong ngày
   - Thời gian theo từng danh mục (số giờ phút + % tổng)
   - Biểu đồ tròn (pie chart) phân bổ thời gian theo danh mục
   - Danh sách hoạt động theo thứ tự thời gian

**Luồng thay thế:**
- [A1] Ngày được chọn không có dữ liệu → hiển thị trạng thái trống

**Business Rules:**
- BR-29: Tổng thời gian = tổng (thời gian kết thúc − thời gian bắt đầu) của tất cả hoạt động trong ngày
- BR-30: Nếu các hoạt động chồng nhau (overlapping), hệ thống vẫn cộng dồn (không deduplicate) — hiển thị cảnh báo khi tổng > 24 giờ

---

### UC-15: Xem báo cáo tuần (Weekly Summary)

**Actor:** Người dùng (User)
**Mục tiêu:** Nhìn lại xu hướng sử dụng thời gian trong tuần

**Luồng chính:**
1. Người dùng vào trang Báo cáo → chọn tab Tuần
2. Người dùng điều hướng chọn tuần muốn xem (tuần trước / tuần này / chọn ngày)
3. Hệ thống hiển thị:
   - Biểu đồ cột: thời gian theo ngày trong tuần, phân theo danh mục (stacked bar chart)
   - Tổng thời gian mỗi danh mục trong tuần
   - Ngày có nhiều thời gian làm việc nhất / ít nhất

**Business Rules:**
- BR-31: Tuần tính từ Thứ Hai đến Chủ Nhật
- BR-29: Áp dụng cho từng ngày

---

### UC-19: Xem báo cáo tháng (Monthly Summary)

**Actor:** Người dùng (User)
**Mục tiêu:** Nhìn lại tổng quan sử dụng thời gian trong một tháng

**Điều kiện tiên quyết:** Có ít nhất 1 hoạt động trong tháng đó

**Luồng chính:**
1. Người dùng vào trang Báo cáo → chọn tab Tháng
2. Người dùng chọn tháng / năm muốn xem
3. Hệ thống hiển thị:
   - Tổng thời gian theo từng danh mục trong tháng (số giờ + %)
   - Biểu đồ tròn phân bổ danh mục
   - Biểu đồ xu hướng theo tuần trong tháng (4 điểm dữ liệu)
   - Ngày hoạt động nhiều nhất / ít nhất trong tháng

**Business Rules:**
- BR-36: Tháng tính theo lịch dương từ ngày 1 đến ngày cuối tháng
- BR-29: Áp dụng cho từng ngày trong tháng

---

### UC-17: Xem và khôi phục hoạt động đã xóa (Trash)

**Actor:** Người dùng (User)
**Mục tiêu:** Xem lại và phục hồi hoạt động bị xóa nhầm trong vòng 30 ngày

**Điều kiện tiên quyết:** Người dùng đang đăng nhập

**Luồng chính:**
1. Người dùng vào trang Trash (từ menu hoặc cài đặt)
2. Hệ thống hiển thị danh sách hoạt động đã xóa, sắp xếp theo ngày xóa mới nhất, kèm số ngày còn lại trước khi xóa vĩnh viễn
3. Người dùng chọn một hoạt động → click "Khôi phục"
4. Hệ thống chuyển hoạt động về trạng thái bình thường, xuất hiện lại trong danh sách chính

**Luồng thay thế:**
- [A1] Người dùng chọn "Xóa vĩnh viễn" trên một mục → dialog xác nhận → xóa permanent
- [A2] Người dùng chọn "Xóa tất cả" → dialog xác nhận → xóa toàn bộ Trash vĩnh viễn
- [A3] Trash rỗng → hiển thị trạng thái trống

**Business Rules:**
- BR-22: Hoạt động trong Trash tự xóa vĩnh viễn sau 30 ngày kể từ ngày xóa
- BR-35: Hoạt động trong Trash không tính vào báo cáo
- BR-37: Sau khi khôi phục, hoạt động xuất hiện đúng vị trí theo thời gian trong ngày tương ứng
- BR-38: Hệ thống chạy job tự động xóa vĩnh viễn các mục quá 30 ngày (scheduled cleanup)

---

### UC-18: Export dữ liệu

**Actor:** Người dùng (User)
**Mục tiêu:** Tải xuống dữ liệu hoạt động ra file để lưu trữ hoặc phân tích bên ngoài

**Điều kiện tiên quyết:** Người dùng đang đăng nhập, có ít nhất 1 hoạt động

**Luồng chính:**
1. Người dùng vào trang Export (từ menu hoặc trang Báo cáo)
2. Người dùng chọn: khoảng thời gian (ngày bắt đầu – ngày kết thúc), định dạng (CSV hoặc PDF)
3. Người dùng click "Export"
4. Hệ thống tạo file và trả về để tải xuống ngay (không gửi qua email)

**Luồng thay thế:**
- [A1] Khoảng thời gian không có dữ liệu → thông báo "Không có dữ liệu trong khoảng thời gian này"
- [A2] Khoảng thời gian quá lớn (> 1 năm) → thông báo giới hạn, yêu cầu chọn lại

**Business Rules:**
- BR-39: CSV chứa: ngày, tên hoạt động, danh mục, giờ bắt đầu, giờ kết thúc, tổng thời gian (phút), ghi chú
- BR-40: PDF chứa: summary theo danh mục + danh sách hoạt động có định dạng đẹp
- BR-41: Export tối đa 1 năm dữ liệu mỗi lần
- BR-42: Hoạt động trong Trash không được export
- BR-43: Thời gian trong file export theo UTC+7

---

### UC-16: Xem và chỉnh sửa thông tin cá nhân

**Actor:** Người dùng (User)
**Mục tiêu:** Quản lý thông tin hồ sơ cá nhân

**Luồng chính:**
1. Người dùng vào trang Cài đặt → Hồ sơ
2. Hệ thống hiển thị: tên hiển thị, email (không cho sửa), timezone
3. Người dùng thay đổi tên hoặc timezone
4. Hệ thống lưu và xác nhận

**Business Rules:**
- BR-32: Email không được phép thay đổi sau khi đăng ký
- BR-33: Tên hiển thị tối đa 50 ký tự, không được để trống
- BR-34: Timezone mặc định = timezone của trình duyệt tại thời điểm đăng ký

---

## 4. User Stories

### Module: Xác thực

| ID | User Story |
|---|---|
| US-01 | Với tư cách là **khách**, tôi muốn **đăng ký tài khoản bằng email và mật khẩu** để **bắt đầu sử dụng hệ thống** |
| US-02 | Với tư cách là **khách**, tôi muốn **đăng nhập vào hệ thống** để **truy cập dữ liệu cá nhân của mình** |
| US-03 | Với tư cách là **người dùng**, tôi muốn **đặt lại mật khẩu qua email** để **lấy lại quyền truy cập khi quên mật khẩu** |

### Module: Log hoạt động

| ID | User Story |
|---|---|
| US-04 | Với tư cách là **người dùng**, tôi muốn **log một hoạt động trong dưới 10 giây** để **không bị gián đoạn công việc đang làm** |
| US-05 | Với tư cách là **người dùng**, tôi muốn **chọn danh mục khi log hoạt động** để **phân loại thời gian rõ ràng** |
| US-06 | Với tư cách là **người dùng**, tôi muốn **log hoạt động trong quá khứ** để **bổ sung những gì tôi quên log lúc đó** |
| US-07 | Với tư cách là **người dùng**, tôi muốn **chỉnh sửa hoạt động đã log** để **sửa lỗi nhập liệu** |
| US-08 | Với tư cách là **người dùng**, tôi muốn **xóa hoạt động đã log** để **loại bỏ dữ liệu không chính xác** |

### Module: Danh mục

| ID | User Story |
|---|---|
| US-09 | Với tư cách là **người dùng**, tôi muốn **thêm danh mục tùy chỉnh** để **phân loại theo nhu cầu riêng của mình** |
| US-10 | Với tư cách là **người dùng**, tôi muốn **xóa danh mục không cần thiết và chuyển hoạt động sang danh mục khác** để **giữ danh sách gọn gàng** |

### Module: Báo cáo

| ID | User Story |
|---|---|
| US-11 | Với tư cách là **người dùng**, tôi muốn **xem tổng thời gian theo từng danh mục trong ngày** để **biết mình đã dùng thời gian vào đâu** |
| US-12 | Với tư cách là **người dùng**, tôi muốn **xem báo cáo theo tuần** để **nhận ra xu hướng và thay đổi thói quen** |

---

## 5. Acceptance Criteria (Trọng tâm)

### US-04: Log hoạt động nhanh

```
Given người dùng đang ở trang Dashboard
When người dùng click nút "Log hoạt động", điền tên, chọn danh mục, nhập giờ bắt đầu và kết thúc, rồi submit
Then hoạt động được lưu thành công và xuất hiện ngay trong danh sách mà không cần reload trang
```

```
Given người dùng đang điền form log hoạt động
When người dùng để trống tên hoạt động và submit
Then hệ thống hiển thị thông báo lỗi inline tại field tên, không gửi request lên server
```

```
Given người dùng nhập thời gian kết thúc sớm hơn thời gian bắt đầu
When người dùng submit form
Then hệ thống hiển thị lỗi "Thời gian kết thúc phải sau thời gian bắt đầu", không lưu dữ liệu
```

### US-11: Daily Summary

```
Given người dùng đã log ít nhất 1 hoạt động trong ngày hôm nay
When người dùng vào trang báo cáo ngày
Then hệ thống hiển thị tổng thời gian theo từng danh mục (giờ:phút) và biểu đồ tỷ lệ phần trăm
```

```
Given người dùng vào trang báo cáo một ngày chưa có dữ liệu
When trang được tải
Then hệ thống hiển thị trạng thái trống với thông điệp hướng dẫn
```

### US-01: Đăng ký tài khoản

```
Given người dùng điền form đăng ký với email hợp lệ và mật khẩu đủ mạnh
When người dùng submit form
Then hệ thống tạo tài khoản, gửi email xác thực, và hiển thị thông báo "Vui lòng kiểm tra email để xác thực tài khoản"
```

```
Given người dùng nhập email đã tồn tại trong hệ thống
When người dùng submit form đăng ký
Then hệ thống hiển thị thông báo lỗi chung không tiết lộ email có tồn tại hay không
```

---

## 6. Business Rules — Tổng hợp

| ID | Rule | Module |
|---|---|---|
| BR-01 | Email phải đúng định dạng RFC 5322 | Auth |
| BR-02 | Mật khẩu >= 8 ký tự, có chữ hoa và số | Auth |
| BR-03 | Tài khoản chưa xác thực email không thể đăng nhập | Auth |
| BR-04 | Link xác thực email hiệu lực 24 giờ | Auth |
| BR-05 | Khóa tài khoản 15 phút sau 5 lần đăng nhập sai | Auth |
| BR-06 | Session timeout = 7 ngày kể từ lần hoạt động cuối | Auth |
| BR-07 | Không tiết lộ thông tin xác thực chi tiết (user enumeration) | Auth |
| BR-08 | Link reset mật khẩu hiệu lực 1 giờ | Auth |
| BR-09 | Link reset mật khẩu chỉ dùng được 1 lần | Auth |
| BR-10 | Đổi mật khẩu hủy toàn bộ session cũ | Auth |
| BR-11 | Không xác nhận email tồn tại hay không trong luồng quên mật khẩu | Auth |
| BR-12 | Mật khẩu mới không được trùng mật khẩu hiện tại | Auth |
| BR-13 | Tên hoạt động tối đa 100 ký tự | Activity |
| BR-14 | Ghi chú tối đa 500 ký tự | Activity |
| BR-15 | Thời gian kết thúc > thời gian bắt đầu | Activity |
| BR-16 | Cho phép log hoạt động trong quá khứ tối đa 2 năm tính từ ngày hiện tại | Activity |
| BR-17 | Không log hoạt động ở thời điểm tương lai | Activity |
| BR-18 | Hoạt động chồng nhau được phép; hệ thống hiển thị toast notification cảnh báo nhẹ, không chặn submit | Activity |
| BR-19 | Chỉ hiển thị dữ liệu của người dùng đang đăng nhập | Activity |
| BR-20 | Timezone người dùng áp dụng cho toàn bộ hiển thị ngày giờ | Activity |
| BR-21 | Chỉ chỉnh sửa/xóa hoạt động của chính mình | Activity |
| BR-22 | Xóa hoạt động là soft delete — chuyển vào Trash, tự xóa vĩnh viễn sau 30 ngày | Activity |
| BR-23 | Tên danh mục tối đa 50 ký tự | Category |
| BR-24 | Không giới hạn cứng danh mục tùy chỉnh; nếu ảnh hưởng performance thì áp hard limit = 50, có cơ chế extend | Category |
| BR-25 | Tên danh mục không phân biệt hoa thường khi kiểm tra trùng | Category |
| BR-26 | Không được sửa danh mục mặc định | Category |
| BR-27 | Không được xóa danh mục mặc định | Category |
| BR-28 | Xóa danh mục đang dùng phải chọn danh mục thay thế | Category |
| BR-29 | Tổng thời gian = tổng (end - start) không deduplicate | Report |
| BR-30 | Cảnh báo khi tổng thời gian > 24 giờ/ngày | Report |
| BR-31 | Tuần tính từ Thứ Hai đến Chủ Nhật | Report |
| BR-35 | Hoạt động trong Trash không tính vào báo cáo | Report/Trash |
| BR-36 | Tháng tính theo lịch dương từ ngày 1 đến ngày cuối tháng | Report |
| BR-37 | Khôi phục hoạt động từ Trash → xuất hiện đúng vị trí theo thời gian gốc | Trash |
| BR-38 | Job tự động xóa vĩnh viễn các mục Trash quá 30 ngày (scheduled cleanup) | Trash |
| BR-39 | CSV export gồm: ngày, tên HĐ, danh mục, giờ bắt đầu, giờ kết thúc, tổng phút, ghi chú | Export |
| BR-40 | PDF export gồm: summary theo danh mục + danh sách HĐ có định dạng đẹp | Export |
| BR-41 | Export tối đa 1 năm dữ liệu mỗi lần | Export |
| BR-42 | Hoạt động trong Trash không được export | Export |
| BR-43 | Thời gian trong file export theo UTC+7 | Export |
| BR-32 | Email không được thay đổi sau khi đăng ký | Profile |
| BR-33 | Tên hiển thị tối đa 50 ký tự, không được để trống | Profile |
| BR-34 | Timezone cố định UTC+7; toàn bộ dữ liệu lưu dạng UTC trong DB, hiển thị ra ngoài theo UTC+7 | Profile |

---

## 7. Edge Cases & Open Questions

### Open Questions — ĐÃ RESOLVED ✅

| ID | Câu hỏi | Quyết định |
|---|---|---|
| OQ-01 | Session timeout bao lâu? | **7 ngày** → BR-06 |
| OQ-02 | Log hoạt động quá khứ tối đa bao lâu? | **2 năm** → BR-16 |
| OQ-03 | Cảnh báo khi hoạt động chồng nhau? | **Toast notification, không chặn** → BR-18 |
| OQ-04 | Có undo sau khi xóa không? | **Có — Trash 30 ngày** → BR-22, UC-17 |
| OQ-05 | Giới hạn số danh mục tùy chỉnh? | **Không giới hạn cứng; fallback = 50 + extend** → BR-24 |
| OQ-06 | Có báo cáo tháng không? | **Có** → UC-19 |
| OQ-07 | Có export CSV/PDF không? | **Có** → UC-18 |
| OQ-08 | Timezone xử lý thế nào? | **Cố định UTC+7; lưu UTC trong DB, hiển thị UTC+7** → BR-34 |

### Edge Cases đáng chú ý

| ID | Tình huống | Hành vi kỳ vọng |
|---|---|---|
| EC-01 | Log hoạt động bắt đầu 23:30 kết thúc 00:30 (qua nửa đêm) | Ghi vào ngày bắt đầu? Hay split? → Cần xác nhận |
| EC-02 | Người dùng log 25 giờ trong một ngày (do overlapping) | Cảnh báo nhưng vẫn lưu (BR-30) |
| EC-03 | Xóa tài khoản — dữ liệu xử lý như thế nào? | Chưa có UC xóa tài khoản — Out of scope? |
| EC-04 | Hai tab mở cùng lúc, log trên cả hai | Cần xử lý conflict/sync |
| EC-05 | Người dùng thay đổi timezone sau khi đã có dữ liệu | Hiển thị theo timezone mới hay giữ nguyên? |
