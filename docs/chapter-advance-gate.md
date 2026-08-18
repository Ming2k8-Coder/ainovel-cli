# Cơ Chế Kiểm Soát Tiến Độ Chương (Chapter Advance Gate)

Tài liệu này mô tả chi tiết cơ chế khóa kiểm soát tiến độ từng chương (`gate`), chế độ nghiệm thu thủ công (`/review on`), cơ chế cấp phép phát hành (`/next`), và quy tắc xử lý khi phải làm lại bản thảo (Rework Safety).

---

## 1. Bối Cảnh & Mục Đích Thiết Kế

Mặc định, `ainovel-cli` vận hành ở chế độ tự động hóa hoàn toàn 100%: Engine tự động định tuyến từ chương này sang chương khác cho đến khi hoàn thành toàn bộ tác phẩm.

Tuy nhiên, trong quá trình sáng tác chất lượng cao, tác giả con người thường cần:
1. **Kiểm soát nhịp độ**: Tạm dừng sau mỗi chương để đọc duyệt, đánh giá văn phong trước khi cho phép AI viết tiếp chương sau.
2. **Can thiệp định hướng**: Điều chỉnh các chi tiết chưa ưng ý thông qua lệnh can thiệp thời gian thực (`steer`).
3. **An toàn về hạn ngạch cấp phép**: Mỗi lần người dùng gõ `/next` là cấp phép tạo **chính xác 1 chương mới**. Nếu chương đó bị lỗi mạng, bị từ chối thẩm định hoặc phải viết lại nhiều lần, lượt cấp phép `/next` đó không bị mất hay trừ hao nhầm.

---

## 2. Các Trạng Thái Của Cổng Kiểm Soát

- **Chế độ Tự động (`review=off`)**: Cổng luôn mở (`Allow = true`), Engine tự động điều phối liên tục giữa Writer, Editor và Architect.
- **Chế độ Duyệt từng bước (`review=on`)**:
  - Khi một chương hoàn thành và nộp bản thảo (`commit_chapter`), cổng sẽ đóng lại và chuyển sang trạng thái chờ cấp phép (`PermittedChapter <= CompletedChapter`).
  - Động cơ sẽ tạm dừng tại ranh giới chương và phát thông báo trên TUI để người dùng đọc duyệt.
  - Khi người dùng gõ lệnh `/next`, biến `PermittedChapter` được nâng lên `CurrentChapter + 1`, cho phép Writer bắt tay vào viết chương tiếp theo.

---

## 3. Quy Tắc Bảo Vệ Khi Viết Lại (Rework Safety)

Khi một chương cần phải làm lại (do Editor đánh giá không đạt yêu cầu hoặc người dùng yêu cầu sửa đổi):
- Các tác vụ **viết lại (`rewrite`)** hoặc **mài dũa (`polish`)** cho các chương đã có trong hàng đợi `PendingRewrites` **luôn luôn được phép thực thi ngay lập tức** mà không bị chặn bởi cổng cấp phép chương mới.
- Chỉ khi toàn bộ hàng đợi làm lại đã được giải quyết xong (`rewrites_drained = true`), cổng mới kiểm tra lại điều kiện cấp phép cho chương tiếp theo.
