# Khả Năng Quan Sát (Observability) & Tiểu Hệ Thống Chẩn Đoán

Hệ thống sử dụng gói nội bộ `internal/diag` làm hạ tầng quan sát, giám sát và chẩn đoán duy nhất cho toàn bộ quá trình sáng tác của `ainovel-cli`.

---

## 1. Ba Tầng Dữ Liệu Quan Sát Độc Lập

Để đảm bảo vừa cung cấp thông tin ngắn gọn, trực quan cho người dùng theo dõi trên màn hình Terminal, vừa lưu giữ đầy đủ dấu vết kỹ thuật phục vụ việc gỡ lỗi chuyên sâu, hệ thống phân chia dữ liệu thành 3 tầng:

1. **Tầng Truyền Tải Dữ Liệu Thô (`agentcore.ProgressPayload`)**: Thu thập toàn bộ các sự kiện chi tiết theo thời gian thực từ vòng lặp Worker (lệnh gọi công cụ, chuỗi suy luận, các lượt thử lại).
2. **Tầng Trình Diễn & Nhật Ký (`host.Event.Summary` & `Detail`)**:
   - `Summary`: Cung cấp dòng trạng thái ngắn gọn, súc tích để hiển thị trên giao diện TUI.
   - `Detail`: Lưu giữ đầy đủ thông tin ngữ cảnh, dữ liệu lỗi và stack trace vào tệp nhật ký trên đĩa.
3. **Màn Hình Chẩn Đoán Trực Quan (`/diag Screen`)**: Bảng điều khiển tích hợp sẵn trong TUI, cho phép phân tích theo thời gian thực các rủi ro về độ mệt mỏi từ ngữ, cảnh báo bế tắc cốt truyện, hoặc các vi phạm quy tắc sáng tác.

---

## 2. Kỷ Luật Quan Sát Bất Biến

> **Quy tắc cốt lõi: Tiểu hệ thống quan sát chỉ ghi nhận và cảnh báo — Tuyệt đối KHÔNG BAO GIỜ tự ý sửa đổi dữ liệu trong Store hoặc can thiệp vào logic điều phối của Engine.**

Mọi thay đổi đối với trạng thái tác phẩm đều phải đi qua các Công cụ (Tools) có thẩm quyền và tuân theo sự phán quyết của Engine/Arbiter.
