# Quản Lý Quy Tắc Người Dùng Thời Gian Chạy (User Rules Runtime)

Tài liệu này mô tả cách hệ thống nạp, chuẩn hóa và thực thi các quy tắc tùy chỉnh do tác giả thiết lập riêng cho tác phẩm (đặt tại `~/.ainovel/rules/*.md` hoặc trong thư mục dự án `./.ainovel/rules/*.md`).

---

## 1. Hai Nhóm Quy Tắc Chính

Hệ thống phân tách quy tắc người dùng thành 2 cơ chế thực thi độc lập:

1. **Ràng Buộc Cơ Khí Cứng (`structured`)**:
   - `forbidden_chars`: Các ký tự đặc biệt bị cấm xuất hiện trong văn bản.
   - `forbidden_phrases`: Các từ ngữ, thành ngữ cấm kỵ.
   - `fatigue_words`: Các từ ngữ mệt mỏi/lặp lại quá nhiều bị giới hạn tần suất xuất hiện.
   - *Được kiểm tra tự động bằng mã nguồn khi nộp bản thảo (`commit_chapter`). Nếu vi phạm, hệ thống sẽ ghi nhận vào báo cáo thẩm định.*

2. **Sở Thích Sáng Tác Ngôn Ngữ Tự Nhiên (`preferences`)**:
   - Các chỉ thị về tính cách nhân vật, giọng điệu đặc trưng, nguyên tắc thế giới do tác giả mô tả bằng văn xuôi.
   - *Được tự động đưa vào `working_memory.user_rules` để các tác tử Architect, Writer và Editor đọc và tuân thủ trong suốt quá trình sáng tác.*

---

## 2. Thứ Tự Ưu Tiên Tuyệt Đối

> **Nguyên tắc bất di bất dịch: Ý chí và quy tắc của người dùng luôn có quyền ưu tiên CAO NHẤT.**
> Khi sở thích của người dùng có sự khác biệt so với các thiết lập mặc định của hệ thống, hệ thống sẽ hoàn toàn tuân theo chỉ dẫn của người dùng.
