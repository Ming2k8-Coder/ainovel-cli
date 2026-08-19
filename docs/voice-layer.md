# Tầng Văn Phong & Tiêu Chuẩn Thẩm Mỹ Văn Học (Voice Layer)

Tầng văn phong định nghĩa các tiêu chuẩn về giọng văn, nhịp điệu và chất lượng biểu đạt nghệ thuật cho các tác tử Writer (sáng tác) và Editor (thẩm định).

---

## 1. Ba Cấp Độ Áp Dụng Văn Phong (3-Tier Voice Overrides)

Hệ thống cho phép tùy biến văn phong linh hoạt theo thứ tự ưu tiên từ chung đến riêng:

```text
[Cấp 1: Hệ Thống (assets/voice.md)]
      │
      ▼ (Ghi đè theo Thể loại)
[Cấp 2: Thể Loại Văn Học (assets/styles/<style>.md)]
      │ (Kỳ ảo, Khoa học viễn tưởng, Trinh thám, Kiếm hiệp...)
      ▼ (Ghi đè theo Tác phẩm)
[Cấp 3: Tác Phẩm Cụ Thể (<ThưMụcTruyện>/style/voice.md)]
        (Văn phong độc bản dành riêng cho bộ truyện này)
```

---

## 2. Tiêu Chuẩn Loại Bỏ "Văn Phong AI" (De-AI Tone Principles)

Nhằm đảm bảo tác phẩm có chất văn tự nhiên, giàu cảm xúc và không bị rập khuôn theo các thói quen hành văn phổ biến của mô hình AI, hệ thống áp dụng các quy chuẩn thẩm mỹ nghiêm ngặt:

1. **Tránh Cấu Trúc Đệm Rập Khuôn**: Loại bỏ triệt để các mẫu câu sáo rỗng thường gặp của AI (như *"không phải... mà là..."*, *"trong một khoảnh khắc ngưng đọng..."*, *"không khí bỗng chốc trở nên..."*).
2. **Biến Hóa Cấu Trúc & Nhịp Điệu Câu**: Kết hợp linh hoạt giữa câu văn ngắn dồn dập trong cảnh hành động và câu văn dài gợi cảm xúc trong cảnh tự sự; tránh việc toàn bộ các câu đều có độ dài đều đều tẻ nhạt.
3. **Mở Đầu Trực Diện, Không Tóm Tắt Lại**: Tuyệt đối không mở đầu chương mới bằng việc tóm tắt lại các sự kiện vừa diễn ra ở chương trước; đi thẳng vào hành động và diễn biến tiếp nối.

---

## 3. Chuyển Đổi Phong Cách Prompt Nhanh & Kết Hợp A/B Feedback (`ainovel-cli style`)

Hệ thống cung cấp cơ chế **Chuyển Đổi Prompt Linh Hoạt** giúp tác giả đổi phong cách văn phong tức thời hoặc tự động đúc kết phong cách mới từ sở thích phản hồi A/B:

### a) Chuyển Đổi Phong Cách Nhanh (Quick Style Switch):
```bash
# Xem danh sách tất cả các bộ phong cách sẵn có (cyberpunk, xianxia, mystery, fantasy...)
ainovel-cli style --list --dir ./novel

# Chuyển nhanh sang bộ phong cách mới (Văn phong được nạp tự động vào System Prompt)
ainovel-cli style --switch "cyberpunk" --dir ./novel
```

### b) Đúc Kết Phong Cách Mới Từ Phản Hồi A/B (Create Style from A/B):
Tự động phân tích các quyết định A/B của tác giả và lưu thành một bộ phong cách hoàn toàn mới:
```bash
ainovel-cli style --create-from-ab "my-ab-cyberpunk" --dir ./novel
```

### c) Tinh Chỉnh & Sửa Lỗi Phong Cách Hiện Có (Tune/Fix Existing Style):
Bổ sung các chỉ thị ưu tiên (`ProseDirectives`) và ràng buộc cấm kỵ (`Taboos`) từ A/B vào một bộ phong cách hiện có:
```bash
ainovel-cli style --tune "cyberpunk" --dir ./novel
```

