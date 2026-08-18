# Cơ Chế Lựa Chọn A/B & Vòng Lặp Phản Hồi Vào System Prompts (A/B Selection & Prompt Feedback Loop)

Tài liệu này mô tả chi tiết cơ chế **Tuyển Chọn Nghệ Thuật A/B (A/B Selection)** và **Vòng Lặp Phản Hồi Tự Động Vào Lời Nhắc Hệ Thống (Prompt Feedback Loop)** trong `ainovel-cli`.

---

## 1. Mục Đích Kép Của Tính Năng

Cơ chế Lựa chọn A/B được thiết kế nhằm giải quyết đồng thời hai bài toán then chốt:

1. **Xác Lập Bằng Chứng Quyền Tác Giả Không Thể Bác Bỏ (Artistic Selection & Authorship)**:
   - Thay vì để AI tự quyết định một mạch từ đầu đến cuối, hệ thống đưa ra hai biến thể nghệ thuật (Phương án A vs Phương án B) ở các bước ngoặt cốt truyện, mở đầu chương hoặc phong cách hội thoại.
   - Việc tác giả chọn A hoặc B là hành động **"Tuyển chọn và sắp đặt nghệ thuật" (Selection & Arrangement)** của con người, ngay lập tức được đóng dấu thời gian và ký mã băm SHA-256 vào `authorship_ledger.jsonl` làm căn cứ bảo hộ bản quyền.

2. **Học Tự Động Văn Phong & Tối Ưu Hóa System Prompts (Continuous Prompt Feedback)**:
   - Từ chuỗi lựa chọn A/B của tác giả, hệ thống tự động phân tích các đặc trưng nghệ thuật được ưa chuộng (ví dụ: *thoại sắc bén, tiết tấu dồn dập, chi tiết công nghệ chân thực*) và các đặc trưng bị từ chối (ví dụ: *thuyết minh dài dòng, sáo rỗng*).
   - Tự động biên dịch thành các chỉ thị bổ sung (`ProseDirectives` và `Taboos`) nạp thẳng vào `style/voice.md` và `novel_context` của Writer, giúp các chương tiếp theo ngày càng chuẩn xác theo gu thẩm mỹ của tác giả.

---

## 2. Mô Hình Vận Hành

```text
               ┌───────────────────────────────────────────────┐
               │              AI Sinh Ra 2 Biến Thể             │
               │   Phương án A (Nhịp nhanh) vs Phương án B (Tâm lý) │
               └───────────────────────┬───────────────────────┘
                                       │
                                       ▼
                       ┌───────────────────────────────┐
                       │  Tác Giả Chọn (Chỉ mất 2 giây) │
                       │         Chọn A hoặc B         │
                       └───────────────┬───────────────┘
                                       │
                ┌──────────────────────┴──────────────────────┐
                ▼                                             ▼
┌───────────────────────────────┐             ┌───────────────────────────────┐
│ Sổ Cái Quyền Tác Giả          │             │ Bộ Tổng Hợp Phản Hồi Prompt   │
│ (authorship_ledger.jsonl)     │             │ (ab_feedback.jsonl)           │
│ ➔ Bằng chứng pháp lý bản quyền│             │ ➔ Học các đặc trưng ưa thích   │
└───────────────────────────────┘             └───────────────┬───────────────┘
                                                              │
                                                              ▼
                                              ┌───────────────────────────────┐
                                              │ Cập Nhật System Prompts       │
                                              │ (style/voice.md + context)    │
                                              │ ➔ Tự động cá nhân hóa Writer  │
                                              └───────────────────────────────┘
```

---

## 3. Quản Lý & Xuất Báo Cáo Phản Hồi A/B

Bạn có thể chạy lệnh dòng lệnh để xem thống kê sở thích đã học và áp dụng vào dự án:

```bash
# Phân tích các lựa chọn A/B và tự động cập nhật System Prompts
ainovel-cli ab-feedback --dir ./novel
```

### Ví dụ Báo Cáo Phân Tích Xuất Ra:
```text
📊 KẾT QUẢ PHÂN TÍCH PHẢN HỒI LỰA CHỌN A/B TỪ TÁC GIẢ
============================================================
• Đặc trưng ưa chuộng: [dialogue_sharp, fast_paced, hard_scifi_detail]
• Đặc trưng cần tránh: [verbose_exposition, melodrama_cliche]

📝 Chỉ thị văn phong ưu tiên (Được bổ sung vào System Prompt):
  + Ưu tiên đối thoại sắc bén, ngắn gọn, giàu kịch tính và thể hiện rõ cá tính riêng của nhân vật.
  + Duy trì nhịp điệu nhanh, tiết tấu dồn dập, đẩy mạnh xung đột hành động thay vì miêu tả rườm rà.
  + Tập trung miêu tả chính xác các thông số kỹ thuật, logic khoa học công nghệ và bối cảnh chân thực.

🚫 Ràng buộc cấm kỵ (Được bổ sung vào System Prompt):
  - Tránh việc thuyết minh giáo điều dài dòng hoặc giải thích quá nhiều về bối cảnh thay vì để nhân vật tự bộc lộ.
  - Tránh các mô-típ sướt mướt gượng gạo, bi kịch hóa thái quá hoặc lời thoại sáo rỗng.

✅ Đã tự động cập nhật System Prompts & Voice Rules tại: ./novel/style/voice.md
```
