# User Requirements Specification (URS)
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-17
**Tác giả:** BA-SA Agent
**Dựa trên:** Project Vision v1.0, Business Use Cases v1.1

---

## 1. Introduction

### 1.1 Mục đích
Tài liệu này mô tả những gì người dùng cần từ hệ thống Time Tracker — viết bằng ngôn ngữ của người dùng, không phải ngôn ngữ kỹ thuật. URS là cầu nối giữa nhu cầu thực tế và Software Requirements Specification (SRS).

### 1.2 Phạm vi
Bao gồm toàn bộ chức năng của ứng dụng web Time Tracker giai đoạn MVP:
- Xác thực người dùng (đăng ký, đăng nhập, quản lý mật khẩu)
- Ghi lại và quản lý hoạt động hàng ngày
- Phân loại hoạt động theo danh mục
- Xem báo cáo thời gian (ngày / tuần / tháng)
- Khôi phục hoạt động đã xóa (Trash)
- Export dữ liệu

### 1.3 Đối tượng đọc
- Product Owner, BA, Developer, Tester
- Không yêu cầu kiến thức kỹ thuật để đọc tài liệu này

### 1.4 Định nghĩa
| Thuật ngữ | Định nghĩa |
|---|---|
| **Hoạt động** | Một sự kiện người dùng đã thực hiện, có thời gian bắt đầu và kết thúc |
| **Danh mục** | Nhóm phân loại hoạt động (vd: Làm việc, Giải trí) |
| **Danh mục mặc định** | 6 danh mục được tạo sẵn cho mọi tài khoản mới |
| **Danh mục tùy chỉnh** | Danh mục do người dùng tự tạo thêm |
| **Trash** | Kho chứa tạm các hoạt động đã xóa, trước khi bị xóa vĩnh viễn |
| **Dashboard** | Trang chính hiển thị hoạt động của ngày hôm nay |
| **Session** | Phiên đăng nhập của người dùng |
| **UTC+7** | Múi giờ Việt Nam — toàn bộ thời gian hiển thị theo múi giờ này |

### 1.5 Tài liệu tham chiếu
- Project Vision Document v1.0
- Business Use Cases Document v1.1
- System Architecture Document v1.0
- Database Schema Document v1.0

---

## 2. Bối cảnh nghiệp vụ

### 2.1 Vấn đề cần giải quyết
Người dùng kết thúc ngày với cảm giác "hôm nay mình đã làm gì?" mà không có câu trả lời rõ ràng. Thời gian bị tiêu tán vô hình — vào mạng xã hội, chờ đợi, di chuyển — mà không có dữ liệu thực tế để nhìn lại và cải thiện.

**Ba vấn đề cốt lõi:**
1. Không biết trong ngày mình đã làm những gì và trong bao lâu
2. Không phân biệt được thời gian "thật sự làm việc" vs. thời gian lãng phí
3. Không có dữ liệu để nhìn lại và hình thành thói quen tốt hơn

### 2.2 Tầm nhìn giải pháp
Một ứng dụng web đơn giản để người dùng **ghi lại, phân loại và nhìn lại** toàn bộ hoạt động trong ngày — hoạt động như một **gương phản chiếu** thời gian thực tế, không phải công cụ lập kế hoạch.

### 2.3 Người dùng mục tiêu

| Nhóm | Đặc điểm | Nhu cầu chính |
|---|---|---|
| **Nhân viên văn phòng** | 8h/ngày tại văn phòng, nhiều cuộc họp, hay bị gián đoạn | Biết thời gian làm việc thực tế bao nhiêu |
| **Sinh viên** | Tự quản lý lịch học, dễ xao nhãng | Nhận ra mình mất bao nhiêu giờ vào mạng xã hội |
| **Freelancer** | Làm nhiều dự án, cần track giờ để báo cáo/tính phí | Công cụ đơn giản ghi nhận thời gian theo task |

---

## 3. User Requirements

Mỗi User Requirement (UR) mô tả một nhu cầu của người dùng ở mức độ WHAT — không đề cập HOW.

---

### UR-01: Đăng ký tài khoản
**Nhu cầu:** Người dùng muốn tạo tài khoản riêng để dữ liệu thời gian của họ được lưu trữ riêng tư, chỉ họ mới xem được.

**Điều kiện chấp nhận:**
- Người dùng có thể tạo tài khoản chỉ với email và mật khẩu
- Tài khoản cần xác thực qua email trước khi sử dụng
- Sau khi tạo tài khoản, 6 danh mục mặc định được cung cấp sẵn
- Hệ thống không để lộ thông tin về email đã tồn tại hay chưa

**Nguồn:** UC-01, US-01, G4 (thói quen bền vững)

---

### UR-02: Đăng nhập và đăng xuất
**Nhu cầu:** Người dùng muốn truy cập an toàn vào tài khoản của mình và kết thúc phiên làm việc khi xong.

**Điều kiện chấp nhận:**
- Đăng nhập bằng email và mật khẩu
- Sau khi đăng nhập thành công, vào thẳng Dashboard (hoạt động hôm nay)
- Session tồn tại 7 ngày, không cần đăng nhập lại nếu vẫn dùng thường xuyên
- Sau 5 lần đăng nhập sai, tài khoản bị khóa tạm 15 phút
- Đăng xuất hủy session hoàn toàn và an toàn

**Nguồn:** UC-02, UC-03, US-02, BR-05, BR-06

---

### UR-03: Lấy lại mật khẩu khi quên
**Nhu cầu:** Khi quên mật khẩu, người dùng muốn có cách khôi phục quyền truy cập mà không cần liên hệ hỗ trợ.

**Điều kiện chấp nhận:**
- Nhập email → nhận link đặt lại mật khẩu qua email
- Link chỉ có hiệu lực 1 giờ và dùng được 1 lần duy nhất
- Sau khi đổi mật khẩu thành công, tất cả phiên đăng nhập cũ bị hủy
- Hệ thống không xác nhận email có tồn tại hay không (bảo mật)

**Nguồn:** UC-04, US-03, BR-08, BR-09, BR-10, BR-11

---

### UR-04: Đổi mật khẩu khi đang đăng nhập
**Nhu cầu:** Người dùng muốn thay đổi mật khẩu định kỳ vì lý do bảo mật.

**Điều kiện chấp nhận:**
- Nhập mật khẩu hiện tại để xác nhận trước khi đổi
- Mật khẩu mới phải khác mật khẩu hiện tại
- Sau khi đổi, các phiên đăng nhập khác bị hủy (trừ phiên hiện tại)

**Nguồn:** UC-05, BR-12

---

### UR-05: Log hoạt động nhanh (CORE)
**Nhu cầu:** Người dùng muốn ghi lại một hoạt động vừa thực hiện trong vòng dưới 10 giây — nếu quá lâu họ sẽ bỏ qua và không log.

**Điều kiện chấp nhận:**
- Nút "Log hoạt động" luôn nổi bật và dễ tìm ở mọi trang
- Form log chỉ yêu cầu: tên hoạt động, danh mục, thời gian bắt đầu, thời gian kết thúc
- Thời gian bắt đầu được điền sẵn = thời điểm hiện tại (người dùng chỉnh sửa nếu cần)
- Ghi chú là trường tùy chọn
- Sau khi submit, hoạt động xuất hiện ngay trong danh sách — không reload trang
- Có thể log hoạt động trong quá khứ (tối đa 2 năm trước)
- Không thể log hoạt động ở tương lai
- Nếu hoạt động trùng thời gian với hoạt động khác, hệ thống cảnh báo nhưng vẫn cho phép lưu

**Nguồn:** UC-06, US-04, US-05, US-06, G1, BR-13–BR-18

---

### UR-06: Xem danh sách hoạt động trong ngày
**Nhu cầu:** Người dùng muốn xem tổng quan những gì họ đã làm trong một ngày, có thể điều hướng qua các ngày.

**Điều kiện chấp nhận:**
- Dashboard mặc định hiển thị hoạt động của ngày hôm nay
- Mỗi hoạt động hiển thị: tên, danh mục (màu sắc), thời gian bắt đầu – kết thúc, tổng thời gian
- Hiển thị tóm tắt thời gian theo từng danh mục ở đầu trang
- Có thể điều hướng sang ngày trước / ngày sau
- Nếu ngày đang xem chưa có hoạt động, hiển thị trạng thái trống với gợi ý hành động

**Nguồn:** UC-07, BR-19, BR-20

---

### UR-07: Chỉnh sửa hoạt động đã log
**Nhu cầu:** Người dùng muốn sửa thông tin hoạt động khi nhập sai.

**Điều kiện chấp nhận:**
- Click vào hoạt động → chọn "Chỉnh sửa" → form hiện ra với dữ liệu hiện tại điền sẵn
- Chỉnh sửa xong → submit → danh sách cập nhật ngay lập tức
- Người dùng chỉ được sửa hoạt động của chính mình

**Nguồn:** UC-08, US-07, BR-21

---

### UR-08: Xóa hoạt động (với khả năng phục hồi)
**Nhu cầu:** Người dùng muốn xóa hoạt động nhập sai, nhưng vẫn muốn có cơ hội phục hồi nếu xóa nhầm.

**Điều kiện chấp nhận:**
- Chọn "Xóa" → hệ thống hỏi xác nhận trước khi xóa
- Hoạt động được chuyển vào Trash (không xóa vĩnh viễn ngay)
- Hoạt động trong Trash tự xóa vĩnh viễn sau 30 ngày
- Hoạt động trong Trash không xuất hiện trong báo cáo

**Nguồn:** UC-09, US-08, BR-22, BR-35

---

### UR-09: Xem và quản lý danh mục
**Nhu cầu:** Người dùng muốn xem 6 danh mục mặc định và có thể thêm danh mục phù hợp với nhu cầu riêng.

**Điều kiện chấp nhận:**
- Xem được toàn bộ danh mục (mặc định + tùy chỉnh) với tên, màu sắc, icon
- Thêm danh mục tùy chỉnh với tên, màu sắc và icon từ bộ icon có sẵn
- Tên danh mục không được trùng (không phân biệt hoa thường)
- Danh mục mặc định không thể sửa hay xóa
- Sửa được danh mục tùy chỉnh (tên, màu, icon)
- Xóa danh mục tùy chỉnh: nếu đang được dùng, phải chọn danh mục thay thế trước khi xóa

**Nguồn:** UC-10–UC-13, US-09, US-10, BR-23–BR-28

---

### UR-10: Xem báo cáo ngày
**Nhu cầu:** Người dùng muốn nhìn lại một ngày cụ thể: đã làm gì, mỗi nhóm chiếm bao nhiêu thời gian.

**Điều kiện chấp nhận:**
- Chọn ngày → xem tổng thời gian theo từng danh mục (giờ:phút và %)
- Biểu đồ tròn (pie chart) phân bổ thời gian theo danh mục
- Danh sách đầy đủ các hoạt động trong ngày theo thứ tự thời gian
- Nếu tổng thời gian > 24 giờ (do hoạt động chồng nhau), hiển thị cảnh báo

**Nguồn:** UC-14, US-11, BR-29, BR-30

---

### UR-11: Xem báo cáo tuần
**Nhu cầu:** Người dùng muốn nhận ra xu hướng sử dụng thời gian trong tuần — ngày nào làm nhiều, ngày nào ít.

**Điều kiện chấp nhận:**
- Chọn tuần muốn xem (tuần này / tuần trước / chọn ngày bất kỳ)
- Biểu đồ cột thể hiện thời gian theo ngày, phân theo danh mục (stacked bar)
- Tổng thời gian mỗi danh mục trong cả tuần
- Ngày có nhiều thời gian làm việc nhất / ít nhất trong tuần

**Nguồn:** UC-15, US-12, BR-31

---

### UR-12: Xem báo cáo tháng
**Nhu cầu:** Người dùng muốn nhìn tổng quan cả tháng để thấy xu hướng dài hạn.

**Điều kiện chấp nhận:**
- Chọn tháng / năm muốn xem
- Tổng thời gian theo từng danh mục trong tháng (giờ + %)
- Biểu đồ tròn phân bổ danh mục
- Biểu đồ xu hướng theo tuần trong tháng
- Ngày hoạt động nhiều nhất / ít nhất

**Nguồn:** UC-19, BR-36

---

### UR-13: Khôi phục hoạt động đã xóa (Trash)
**Nhu cầu:** Nếu người dùng xóa nhầm, họ muốn có cơ hội lấy lại hoạt động trong 30 ngày.

**Điều kiện chấp nhận:**
- Xem danh sách hoạt động trong Trash, có số ngày còn lại trước khi xóa vĩnh viễn
- Khôi phục một hoạt động → nó trở lại danh sách chính đúng vị trí theo thời gian
- Có thể xóa vĩnh viễn từng mục hoặc xóa toàn bộ Trash
- Trash rỗng thì hiển thị trạng thái trống

**Nguồn:** UC-17, BR-37, BR-38

---

### UR-14: Export dữ liệu
**Nhu cầu:** Người dùng muốn lưu trữ hoặc phân tích dữ liệu bên ngoài ứng dụng.

**Điều kiện chấp nhận:**
- Chọn khoảng thời gian (ngày bắt đầu – ngày kết thúc) và định dạng (CSV hoặc PDF)
- File được tải xuống ngay, không gửi qua email
- CSV: đầy đủ các trường (ngày, tên, danh mục, giờ bắt đầu, giờ kết thúc, tổng phút, ghi chú)
- PDF: summary theo danh mục + danh sách hoạt động định dạng đẹp
- Tối đa 1 năm dữ liệu mỗi lần export
- Hoạt động trong Trash không được export
- Thời gian trong file theo UTC+7

**Nguồn:** UC-18, BR-39–BR-43

---

### UR-15: Quản lý hồ sơ cá nhân
**Nhu cầu:** Người dùng muốn cập nhật thông tin cá nhân và xem tài khoản của mình.

**Điều kiện chấp nhận:**
- Xem và chỉnh sửa tên hiển thị
- Email được hiển thị nhưng không thể thay đổi
- Timezone cố định UTC+7 (Việt Nam)
- Tên hiển thị không được để trống, tối đa 50 ký tự

**Nguồn:** UC-16, BR-32–BR-34

---

## 4. Ràng buộc và Giả định

### 4.1 Ràng buộc người dùng
- Người dùng phải tự log thủ công — hệ thống không tự động theo dõi hoạt động
- Cần kết nối internet để sử dụng (ứng dụng web)
- Cần trình duyệt web hiện đại (Chrome, Firefox, Safari)
- Timezone cố định UTC+7; không hỗ trợ đổi timezone

### 4.2 Giả định
- Người dùng sẵn sàng dành <10 giây để log mỗi hoạt động
- Người dùng có tài khoản email để đăng ký và nhận thông báo
- Ứng dụng phục vụ cá nhân — không cần tính năng chia sẻ hay cộng tác

### 4.3 Ngoài phạm vi (Out of Scope)
Những gì ứng dụng **KHÔNG** làm trong giai đoạn MVP:
- Lập kế hoạch / to-do list / scheduling
- Tính năng team / collaboration
- Tích hợp công cụ bên thứ ba (Slack, Jira, Google Calendar)
- AI gợi ý hoặc coaching năng suất
- Mobile app native (iOS/Android)
- Xóa tài khoản

---

## 5. Ưu tiên User Requirements

| ID | User Requirement | Ưu tiên | Sprint |
|---|---|---|---|
| UR-01 | Đăng ký tài khoản | Must Have | 1 |
| UR-02 | Đăng nhập và đăng xuất | Must Have | 1 |
| UR-03 | Lấy lại mật khẩu | Must Have | 1 |
| UR-04 | Đổi mật khẩu | Should Have | 1 |
| UR-05 | Log hoạt động nhanh | Must Have | 2 |
| UR-06 | Xem danh sách hoạt động ngày | Must Have | 2 |
| UR-07 | Chỉnh sửa hoạt động | Must Have | 2 |
| UR-08 | Xóa hoạt động | Must Have | 2 |
| UR-09 | Quản lý danh mục | Should Have | 3 |
| UR-10 | Báo cáo ngày | Must Have | 4 |
| UR-11 | Báo cáo tuần | Should Have | 4 |
| UR-12 | Báo cáo tháng | Could Have | 4 |
| UR-13 | Trash & khôi phục | Should Have | 5 |
| UR-14 | Export dữ liệu | Could Have | 6 |
| UR-15 | Hồ sơ cá nhân | Should Have | 6 |
