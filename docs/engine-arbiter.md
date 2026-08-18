# Tiến Hóa Mặt Phẳng Điều Khiển: Engine + Arbiter (Gỡ bỏ Vòng Lặp Dài Coordinator)

> Trạng thái (2026-07-14 v6): **Đã hoàn thành mã nguồn** -- Engine và Arbiter đã đi vào hoạt động, Coordinator và toàn bộ các thành phần đi kèm đã được xóa bỏ hoàn toàn.
> **Kiến trúc hiện tại**: Xem tại README phần Kiến trúc và `docs/engine-rfc.md`.

---

## 1. Động Cơ: Các Giả Định Đã Cũ Kỹ

Giả định ban đầu của dự án -- "Một câu Prompt, một LLM常驻 (thường trực) chạy vòng lặp dài điều khiển toàn bộ cuốn sách" -- đã lỗi thời. 

Sau đợt tái cấu trúc Hybrid, quyền quyết định thực tế nằm ở `flow.Route`. Coordinator trong vòng lặp chính 90% lần gọi chỉ làm nhiệm vụ **chuyển tiếp nguyên văn**. Vì "duy trì một phiên LLM không nên dừng", hệ thống đã phải gánh thêm một đống mã bọc lót phức tạp:
1. StopGuard + kỹ thuật văn bản động blockMessage
2. Giao thức lệnh trùng lặp Dispatcher ("Hạ lệnh lần thứ N")
3. Huấn luyện hành vi `coordinator.md`
4. Cổng chặn giai đoạn completePhaseGate / writerExpandedChapterGate
5. MaxTurns = 100,000

Các tiểu hệ thống mới của dự án (`import` / `simulation` / `cocreate` / `userrules`) đều đã chuyển sang mô hình **"Host trực tiếp điều phối + LLM đóng vai trò hàm"**. Giải pháp này đồng bộ toàn bộ luồng chính về cùng một mô hình chuẩn đã được kiểm chứng.

---

## 2. Mô Hình Mục Tiêu

```text
Entry
  ↓
Host (Vỏ bọc EngineLoop)
  ├─ Đọc Store → flow.Route → Chạy trực tiếp Worker
  ├─ Bối cảnh ngữ nghĩa rõ ràng → Gọi hàm Arbiter
  └─ Sự kiện / Ngân sách / Điểm dừng / Thông báo
  ↓
Workers (architect / writer / editor - Tự chủ nội bộ, giữ lại Checkpoint-delta Guard)
  ↓
Tools → Store (Nguồn thực tế duy nhất)
```

**Phân công trách nhiệm:**
- `Route`: Quản lý mọi bước tiếp theo có thể tra bảng.
- `Arbiter`: Quản lý phán đoán ngữ nghĩa có ranh giới rõ ràng.
- `Worker`: Quản lý sáng tác mở.
- `Engine`: Thực thi quyết định, không tham gia phán đoán văn học.
- `Observer/Diag`: Chỉ quan sát.

---

## 3. Các Kịch Bản Arbiter

| Kịch bản | Kích hoạt | Ghi chú |
|---|---|---|
| `plan_start` | Khởi động sách mới | Chọn KTS short/long + mở rộng yêu cầu quá ngắn |
| `intervention` | Can thiệp người dùng | Truy vấn / Quy tắc / Điều chỉnh cốt truyện / Viết lại chương cũ |
| `worker_failure` | Worker báo lỗi và mã định tính không có lối thoát | Thử lại (retry) / Chuyển hướng (reroute) / Tạm dừng (abort) |
| `deadlock` | Lệnh tương tự lặp lại liên tục không có tiến triển | 3 lần hỏi Arbiter, 5 lần tạm dừng cứng |

---

## 4. Tóm Lược Tác Động

- **Chi phí LLM mỗi chương**: Tiết kiệm 1 lần gọi chuyển tiếp ở mỗi ranh giới chương.
- **Khả năng kiểm thử phán quyết**: Phát lại offline từng kịch bản + đánh giá hồi quy.
- **Phản hồi can thiệp**: Tức thì ở ranh giới lệnh.
- **Độ phức tạp**: Cắt giảm 1,500+ dòng mã dư thừa.
