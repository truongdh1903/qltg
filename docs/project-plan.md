# Project Plan
**Dự án:** Quản Lý Thời Gian (Time Tracker)
**Phiên bản:** 1.0
**Ngày lập:** 2026-04-16
**PM:** AI Agent (Claude)
**Dựa trên:** Business Use Cases v1.1

---

## 1. Project Summary

| Thông tin | Chi tiết |
|---|---|
| **Tên dự án** | Quản Lý Thời Gian (Time Tracker) |
| **Start date** | 2026-04-17 (thứ Sáu) |
| **Target go-live** | ~2026-06-13 (8 tuần) |
| **Methodology** | Lightweight Scrum — Sprint 1 tuần |
| **Team** | 1 người (Full-stack solo dev + AI pair programming) |
| **Capacity** | ~8 giờ/ngày × 5 ngày/tuần = 40 giờ/sprint |
| **Platform** | Web application |
| **Scope** | 19 Use Cases / 5 modules (Auth, Activity, Category, Report, Export/Trash) |

> **Lưu ý PM:** Tech stack chưa được chốt — cần quyết định trong Sprint 0. Tech stack sẽ ảnh hưởng trực tiếp đến estimate của Sprint 1 trở đi.

---

## 2. Phases & Milestones

| Phase | Tên | Thời gian | Milestone | Deliverables |
|---|---|---|---|---|
| **P0** | Foundation | Apr 17–23 (1 tuần) | ✦ M0: Project Ready | System Architecture, DB Schema, API Design, UI Wireframes, Project scaffold |
| **P1** | Core MVP | Apr 24 – May 14 (3 tuần) | ✦ M1: MVP Internal | Auth hoàn chỉnh, Log/Edit/Delete activity, Category management |
| **P2** | Full Features | May 15 – Jun 4 (3 tuần) | ✦ M2: Feature Complete | Reports (Day/Week/Month), Trash, Export, Profile |
| **P3** | QA & Hardening | Jun 5–11 (1 tuần) | ✦ M3: Release Candidate | Integration test, Security review, Bug fixes, UI polish |
| **P4** | Deploy | Jun 12–13 (2 ngày) | ✦ M4: Go Live | Production deploy, monitoring |

---

## 3. Work Breakdown Structure (WBS)

### P0 — Foundation

**Epic 1: Documentation & Design**
| Task | Effort | Owner |
|---|---|---|
| System Architecture Document | 1 ngày | Dev + AI |
| Database Schema / ERD | 0.5 ngày | Dev + AI |
| API Design Document (endpoints, request/response) | 1 ngày | Dev + AI |
| UI Wireframes — low fidelity (5 màn hình chính) | 1 ngày | Dev + AI |

**Epic 2: Project Setup**
| Task | Effort | Owner |
|---|---|---|
| Chốt tech stack FE + BE + DB | 0.5 ngày | Dev |
| Tạo repository, cấu trúc thư mục chuẩn | 0.5 ngày | Dev |
| Project scaffold FE (routing, layout, component base) | 0.5 ngày | Dev + AI |
| Project scaffold BE (folder structure, middleware, error handler) | 0.5 ngày | Dev + AI |
| Database setup + migration tool | 0.5 ngày | Dev |
| CI/CD cơ bản (lint + test tự động khi push) | 0.5 ngày | Dev + AI |

> **Buffer P0:** 0.5 ngày cho phát sinh ngoài dự kiến

---

### P1 — Core MVP

**Epic 3: Authentication (Sprint 1 — Apr 24–30)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Register API + validation | UC-01 | 0.5 ngày | Dev + AI |
| Email verification flow + gửi email | UC-01 | 0.5 ngày | Dev + AI |
| Login API + session/JWT | UC-02 | 0.5 ngày | Dev + AI |
| Logout + hủy token phía server | UC-03 | 0.25 ngày | Dev + AI |
| Forgot password + gửi email reset | UC-04 | 0.5 ngày | Dev + AI |
| Reset password flow + invalidate links cũ | UC-04 | 0.5 ngày | Dev + AI |
| Change password (khi đang login) | UC-05 | 0.25 ngày | Dev + AI |
| Auth middleware / guard (protect routes) | — | 0.5 ngày | Dev + AI |
| Rate limiting + brute force protection (BR-05) | — | 0.5 ngày | Dev + AI |
| UI: Register / Login / Forgot password pages | — | 1 ngày | Dev + AI |
| Unit tests: Auth module | — | 1 ngày | Dev + AI |
> **Buffer Sprint 1:** 0.25 ngày

---

**Epic 4: Activity Logging — CORE (Sprint 2 — May 1–7)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Log activity API + validation (BR-13 → BR-18) | UC-06 | 0.75 ngày | Dev + AI |
| Overlap detection + toast warning | UC-06 BR-18 | 0.5 ngày | Dev + AI |
| List activities by date API | UC-07 | 0.25 ngày | Dev + AI |
| Edit activity API + validation | UC-08 | 0.5 ngày | Dev + AI |
| Soft delete (move to Trash) API | UC-09 | 0.5 ngày | Dev + AI |
| Dashboard UI: daily activity list + navigation ngày | UC-07 | 1 ngày | Dev + AI |
| UI: Quick log form (target < 10s per log) | UC-06 | 0.75 ngày | Dev + AI |
| UI: Edit inline / modal | UC-08 | 0.5 ngày | Dev + AI |
| Timezone handling UTC ↔ UTC+7 (BR-34) | — | 0.5 ngày | Dev + AI |
| Unit tests: Activity module | — | 0.75 ngày | Dev + AI |
> **Buffer Sprint 2:** 0.25 ngày

---

**Epic 5: Category Management (Sprint 3 — May 8–14)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Seed 6 default categories khi tạo account mới | UC-10 | 0.25 ngày | Dev + AI |
| List categories API | UC-10 | 0.25 ngày | Dev + AI |
| Add custom category API + validation | UC-11 | 0.5 ngày | Dev + AI |
| Edit custom category API | UC-12 | 0.25 ngày | Dev + AI |
| Delete category + reassign activities flow (BR-28) | UC-13 | 0.75 ngày | Dev + AI |
| UI: Category management page | UC-10–13 | 1 ngày | Dev + AI |
| Unit tests: Category module | — | 0.5 ngày | Dev + AI |
> **Còn lại Sprint 3** (~2 ngày): Bắt đầu sớm Epic 6 (Daily Report)

---

### P2 — Full Features

**Epic 6: Reports (Sprint 4 — May 15–21)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Daily Summary API (tổng theo category + list) | UC-14 | 0.5 ngày | Dev + AI |
| Weekly Summary API (stacked data by day) | UC-15 | 0.75 ngày | Dev + AI |
| Monthly Summary API (aggregated data) | UC-19 | 0.75 ngày | Dev + AI |
| UI: Daily report + Pie chart | UC-14 | 1 ngày | Dev + AI |
| UI: Weekly report + Stacked bar chart | UC-15 | 1 ngày | Dev + AI |
| UI: Monthly report + trend chart | UC-19 | 1 ngày | Dev + AI |
| Unit tests: Report module | — | 0.75 ngày | Dev + AI |
> **Buffer Sprint 4:** 0.25 ngày

---

**Epic 7: Trash Management (Sprint 5 — May 22–28)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Trash list API (filtered, sorted by delete date) | UC-17 | 0.5 ngày | Dev + AI |
| Restore activity API | UC-17 | 0.25 ngày | Dev + AI |
| Permanent delete API | UC-17 | 0.25 ngày | Dev + AI |
| Scheduled cleanup job (auto delete sau 30 ngày — BR-38) | UC-17 | 0.5 ngày | Dev + AI |
| UI: Trash page | UC-17 | 0.75 ngày | Dev + AI |

**Epic 8: Export (Sprint 5 — May 22–28)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| Export CSV API (BR-39 → BR-43) | UC-18 | 0.75 ngày | Dev + AI |
| Export PDF API | UC-18 | 1 ngày | Dev + AI |
| UI: Export page (chọn date range + format) | UC-18 | 0.5 ngày | Dev + AI |

**Epic 9: Profile (Sprint 5 — May 22–28)**
| Task | UC | Effort | Owner |
|---|---|---|---|
| View / edit profile API (tên, timezone) | UC-16 | 0.5 ngày | Dev + AI |
| UI: Profile / Settings page | UC-16 | 0.5 ngày | Dev + AI |
> **Buffer Sprint 5:** 0.5 ngày

---

**Epic 10: Integration & UX Polish (Sprint 6 — May 29 – Jun 4)**
| Task | Effort | Owner |
|---|---|---|
| Integration tests toàn bộ flow (end-to-end test chính) | 1.5 ngày | Dev + AI |
| Security review checklist (OWASP Top 10) | 1 ngày | Dev + AI |
| Fix các lỗi phát hiện từ testing | 1 ngày | Dev |
| UI/UX polish: responsive, loading states, empty states | 1 ngày | Dev + AI |
| Performance: query optimization, frontend bundle | 0.5 ngày | Dev + AI |
> **Buffer Sprint 6:** 1 ngày

---

### P3 — Deploy

**Epic 11: Deployment (Jun 5–11)**
| Task | Effort | Owner |
|---|---|---|
| Production environment setup (hosting, DB) | 0.5 ngày | Dev |
| Environment variables + secrets management | 0.25 ngày | Dev |
| Domain + SSL/HTTPS setup | 0.25 ngày | Dev |
| CI/CD pipeline hoàn chỉnh (auto deploy khi merge main) | 0.5 ngày | Dev + AI |
| Smoke test trên production | 0.25 ngày | Dev |
| Basic monitoring (uptime, error alerts) | 0.25 ngày | Dev + AI |
| **Go Live** | — | — |

---

## 4. Sprint Plan

| Sprint | Tuần | Mục tiêu | Key Deliverables | Capacity |
|---|---|---|---|---|
| **Sprint 0** | Apr 17–23 | Foundation sẵn sàng | SA Doc, DB Schema, API Design, Wireframes, Project scaffold | 40h |
| **Sprint 1** | Apr 24–30 | Auth hoàn chỉnh | Register, Login, Logout, Forgot/Reset password, auth tests | 40h |
| **Sprint 2** | May 1–7 | Log activity hoạt động | Log/Edit/Delete activity, Dashboard UI, overlap warning | 40h |
| **Sprint 3** | May 8–14 | Category + bắt đầu Report | Category CRUD, daily report (API + UI một phần) | 40h |
| **Sprint 4** | May 15–21 | Reports hoàn chỉnh | Daily / Weekly / Monthly report với charts | 40h |
| **Sprint 5** | May 22–28 | Feature complete | Trash, Export CSV/PDF, Profile | 40h |
| **Sprint 6** | May 29 – Jun 4 | Release candidate | Integration test, Security review, Bug fix, UI polish | 40h |
| **Sprint 7** | Jun 5–11 | Go live | Production deploy, CI/CD, monitoring, smoke test | 40h |

**Target go-live: 2026-06-13** *(8 tuần kể từ start)*

---

## 5. Resource Plan

| Role | Người | Involvement | Phase |
|---|---|---|---|
| Developer (BE + FE) | Bạn (solo) | 100% | Toàn dự án |
| AI Pair Programmer | Claude | On-demand | Toàn dự án |
| PM / BA / SA / QA | Claude (agent) | On-demand | Theo phase |

> **Ghi chú:** Solo dev + AI là mô hình "vibe coding" — AI đảm nhận boilerplate, scaffold, test template; dev quyết định architecture và review output.

---

## 6. Risk Register

| ID | Rủi ro | Xác suất | Tác động | Mitigation | Owner |
|---|---|---|---|---|---|
| R-01 | Scope creep — muốn thêm tính năng giữa sprint | Cao | Trung bình | Ghi vào backlog, không đưa vào sprint đang chạy; review mỗi sprint | Dev |
| R-02 | FE/UI mất nhiều thời gian hơn dự kiến (BE background) | Cao | Cao | Dùng UI component library có sẵn; AI generate UI code | Dev |
| R-03 | Tech stack chưa chốt gây delay Sprint 0 | Trung bình | Cao | Quyết định tech stack **ngay ngày đầu** Sprint 0 — không nghiên cứu quá 2 giờ | Dev |
| R-04 | PDF export phức tạp hơn dự kiến | Trung bình | Thấp | Dùng library có sẵn (puppeteer / pdfkit); timeboxed 1 ngày, nếu quá thì defer | Dev |
| R-05 | Mất động lực giữa chừng | Thấp | Cao | Dùng app của mình hàng ngày từ Sprint 2; demo sau mỗi sprint | Dev |
| R-06 | Security issue phát hiện muộn | Thấp | Cao | Áp dụng security checklist từ Sprint 1 (auth); không để dồn đến Sprint 6 | Dev |
| R-07 | Email service (verification, reset) có vấn đề | Thấp | Trung bình | Chọn service đơn giản (Resend / SendGrid free tier); test sớm Sprint 1 | Dev |

---

## 7. Communication Plan

*(Solo project — tự quản lý)*

| Hoạt động | Tần suất | Mục đích | Công cụ |
|---|---|---|---|
| **Daily standup với AI** | Mỗi ngày | Xác nhận task ngày hôm nay, unblock nếu bị stuck | Chat với Claude |
| **Sprint Review** | Cuối mỗi sprint (thứ Năm hoặc thứ Sáu) | Kiểm tra deliverables đã xong, demo với chính mình | Tự test |
| **Sprint Retrospective** | Cuối mỗi sprint | Gì tốt? Gì cần cải thiện? Estimate có đúng không? | Ghi note ngắn |
| **Sprint Planning** | Đầu sprint mới (thứ Hai) | Xác nhận backlog, commit task tuần này | Checklist |
| **Backlog Grooming** | Khi cần (mid-sprint) | Thêm/bỏ/ưu tiên lại task | Backlog doc |

---

## 8. Dependencies & Critical Path

### Dependencies chính

```
[SA Document] → [DB Schema] → [Project Setup] → [Sprint 1: Auth]
                                                        ↓
[API Design]  ──────────────────────────────→  [Sprint 2: Activity]
                                                        ↓
                                               [Sprint 3: Category]
                                                        ↓
                                               [Sprint 4: Reports]
                                                        ↓
                                               [Sprint 5: Trash/Export/Profile]
                                                        ↓
                                               [Sprint 6: QA/Security]
                                                        ↓
                                               [Sprint 7: Deploy → Go Live]
```

### Critical Path
**SA Doc → DB Schema → Project Setup → Auth → Activity → Category → Reports → QA → Deploy**

Bất kỳ delay nào trên critical path đều ảnh hưởng trực tiếp đến go-live date.

### Nút thắt cần chú ý
| Điểm | Rủi ro | Hành động |
|---|---|---|
| Chốt tech stack (Sprint 0, ngày 1) | Delay cả dự án nếu kéo dài | Timeboxed 2 giờ, quyết định và đi |
| Auth middleware (Sprint 1) | Mọi module sau đều depend vào | Không sang Sprint 2 nếu auth chưa xong |
| Timezone handling (Sprint 2) | Bug ẩn, khó fix muộn | Xử lý đúng ngay từ đầu, test kỹ |
| Integration test (Sprint 6) | Phát hiện bug cross-module | Không skip, dù có deadline |

---

## 9. Definition of Done (DoD)

Một task/story được coi là **DONE** khi:

- [ ] Code đã được viết và tự review
- [ ] Unit test đã viết và pass
- [ ] Không có lỗi console / warning nghiêm trọng
- [ ] Chạy được trên local environment
- [ ] API đúng với API Design Document (nếu có)
- [ ] Không vi phạm Business Rules liên quan
- [ ] Acceptance Criteria đã pass (tự test theo Gherkin)

---

## 10. Backlog ưu tiên thấp (defer nếu cần)

Những tính năng này có trong scope nhưng có thể defer sang v1.1 nếu timeline bị trễ:

| Feature | UC | Lý do có thể defer |
|---|---|---|
| Export PDF | UC-18 | CSV đã đủ dùng; PDF phức tạp hơn |
| Monthly Report | UC-19 | Daily + Weekly đã đủ MVP |
| Scheduled cleanup job | UC-17 BR-38 | Có thể cleanup thủ công ban đầu |
| Rate limiting nâng cao | BR-05 | Basic rate limit là đủ cho MVP |
