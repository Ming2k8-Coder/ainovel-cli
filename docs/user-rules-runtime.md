# Thời Gian Chạy Quy Tắc Người Dùng (User Rules Runtime)

Tài liệu này mô tả cách hệ thống quản lý và áp dụng các quy tắc tùy chỉnh do người dùng thiết lập cho dự án tiểu thuyết (`~/.ainovel/rules/*.md` hoặc `./.ainovel/rules/*.md`).

## 1. Hai Loại Quy Tắc

1. **Ràng buộc Cơ khí (`structured`)**:
   - `forbidden_chars`: Các ký tự cấm xuất hiện trong bản thảo.
   - `forbidden_phrases`: Các cụm từ/từ ngữ bị cấm.
   - `fatigue_words`: Các từ ngữ mệt mỏi/lặp lại quá nhiều bị hạn chế tần suất.
   - *Được kiểm tra cứng tự động khi nộp chương (`commit_chapter`).*

2. **Sở thích Ngôn ngữ Tự nhiên (`preferences`)**:
   - Các yêu cầu về thiết lập nhân vật, văn phong, quy tắc thế giới do người dùng nhập vào.
   - *Được bơm trực tiếp vào `working_memory.user_rules` để các Agent tham chiếu khi sáng tác.*

## 2. Thứ Tự Ưu Tiên

> **Quy tắc và sở thích của người dùng luôn có ưu tiên CAO NHẤT.**
> Khi sở thích của người dùng xung đột với các thiết lập mặc định của hệ thống, hệ thống sẽ tuân theo sở thích của người dùng.
