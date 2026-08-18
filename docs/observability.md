# Khả Năng Quan Sát (Observability) & Chẩn Đoán (`internal/diag`)

Dự án sử dụng gói `internal/diag` làm tiểu hệ thống quan sát duy nhất cho động cơ sáng tác.

## 1. Ba Tầng Dữ Liệu Quan Sát

1. `agentcore.ProgressPayload`: Tầng truyền tải sự kiện trực tiếp (raw events).
2. `host.Event.Summary` & `Detail`: Hiển thị tóm tắt ngắn cho TUI và lưu vết chẩn đoán chi tiết vào tập tin nhật ký.
3. `/diag` Screen: Giao diện chẩn đoán trực quan trong TUI, phân tích các rủi ro văn phong, bế tắc cốt truyện hoặc vi phạm cài đặt.

## 2. Kỷ Luật Quan Sát

> **Tiểu hệ thống chẩn đoán có thể đưa ra cảnh báo và gợi ý, nhưng KHÔNG BAO GIỜ tự tay sửa đổi dữ liệu hay thay đổi luồng điều phối của Engine.**

Mọi sửa đổi dữ liệu đều phải thông qua các công cụ Tool chuẩn và sự phán quyết của Engine/Arbiter.
