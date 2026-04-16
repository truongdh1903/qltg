# Project Vision Document
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày:** 2026-04-16
**Giai đoạn:** Ý tưởng

---

## 1. Problem Statement

Hầu hết mọi người — từ nhân viên văn phòng, sinh viên đến freelancer — đều kết thúc một ngày với cảm giác **"hôm nay mình đã làm gì?"** mà không có câu trả lời rõ ràng.

Thời gian bị tiêu tán vô hình vào mạng xã hội, di chuyển, chờ đợi, các cuộc trò chuyện không cần thiết — nhưng không ai đo được con số thực tế. Kết quả là người ta liên tục cảm thấy bận rộn mà vẫn thiếu năng suất, không biết mình đang đầu tư thời gian đúng chỗ hay không.

**Vấn đề cụ thể:**
- Không biết trong ngày mình đã làm những gì và trong bao lâu
- Không phân biệt được thời gian "thật sự làm việc" vs. thời gian bị lãng phí
- Không có dữ liệu để nhìn lại và cải thiện thói quen

---

## 2. Vision Statement

> **Mọi người đều biết chính xác thời gian của mình đang đi về đâu — và có thể chủ động quyết định dùng nó tốt hơn.**

---

## 3. Mission Statement

Xây dựng một ứng dụng web đơn giản, trực quan giúp người dùng **ghi lại, phân loại và nhìn lại** toàn bộ hoạt động trong ngày — để hiểu rõ mình đang dùng thời gian như thế nào và dần hình thành thói quen sống có chủ đích hơn.

---

## 4. Target Users

| Phân khúc | Đặc điểm | Nỗi đau chính |
|---|---|---|
| **Nhân viên văn phòng** | 8h/ngày tại văn phòng, nhiều cuộc họp, hay bị gián đoạn | Không biết thời gian làm việc thực tế bao nhiêu |
| **Sinh viên** | Tự quản lý lịch học, dễ xao nhãng | Mất nhiều giờ vào mạng xã hội mà không hay |
| **Freelancer** | Làm nhiều dự án, cần track giờ để báo cáo / tính phí | Không có công cụ đơn giản để ghi nhận thời gian theo task |

---

## 5. Core Value Proposition

Không phải công cụ lập kế hoạch — đây là **gương phản chiếu thời gian thực tế** của bạn.

- **Ghi lại thực tế, không phải kế hoạch:** Người dùng log những gì họ *đã làm*, không phải những gì họ *dự định làm*
- **Phân loại trực quan:** Phân biệt rõ: làm việc / di chuyển / ăn uống / giải trí / mạng xã hội
- **Insight đơn giản:** Tổng kết cuối ngày/tuần bằng con số dễ hiểu — không phức tạp, không áp lực

---

## 6. Strategic Goals

1. **G1 — Trải nghiệm ghi chép nhanh:** Người dùng có thể log một hoạt động trong dưới 10 giây
2. **G2 — Phân loại thời gian rõ ràng:** Hệ thống có ít nhất 6 nhóm hoạt động mặc định có thể tùy chỉnh
3. **G3 — Báo cáo có giá trị:** Người dùng thấy được thống kê ngày / tuần đủ để rút ra insight cá nhân
4. **G4 — Thói quen bền vững:** Thiết kế khuyến khích người dùng quay lại mỗi ngày mà không cảm thấy bị ép buộc
5. **G5 — Nền tảng học tập:** Dự án phục vụ mục tiêu học vibe coding — code sạch, có thể mở rộng, dễ iterate

---

## 7. Success Metrics (KPIs)

| Chỉ số | Mục tiêu giai đoạn MVP |
|---|---|
| Thời gian log 1 hoạt động | < 10 giây |
| Số loại hoạt động mặc định | >= 6 categories |
| Người dùng quay lại > 3 ngày/tuần (bản thân) | Có thể duy trì 2 tuần liên tiếp |
| Tính năng báo cáo hoạt động | Hiển thị được daily summary |
| Chất lượng code | Có thể đọc lại và hiểu sau 1 tháng |

---

## 8. Out of Scope

Những gì dự án này **KHÔNG** giải quyết trong giai đoạn hiện tại:

- Lập kế hoạch / scheduling / to-do list (không phải Todoist hay Google Calendar)
- Tính năng team / collaboration / báo cáo cho quản lý
- Tích hợp với công cụ bên thứ ba (Slack, Jira, Google...)
- AI gợi ý / coaching về năng suất
- Tính năng monetization / subscription
- Ứng dụng mobile native (iOS/Android)

---

## 9. Assumptions & Risks

### Giả định nền tảng
- Người dùng sẵn sàng tự log thủ công (không tự động tracking)
- Web app là nền tảng đủ dùng cho giai đoạn đầu
- Dự án phục vụ mục tiêu học tập — không có áp lực thương mại

### Rủi ro chiến lược
| Rủi ro | Mức độ | Hướng xử lý |
|---|---|---|
| Người dùng lười log sau vài ngày | Cao | Thiết kế UX cực kỳ nhanh và không ma sát |
| Scope creep — muốn thêm tính năng liên tục | Trung bình | Giữ Out of Scope rõ ràng, iterate theo version |
| Kỹ năng kỹ thuật chưa đủ để hoàn thiện | Trung bình | Học qua làm, ưu tiên MVP chạy được trước |
| Mất động lực giữa chừng | Trung bình | Dùng chính app của mình hàng ngày để duy trì động lực |
