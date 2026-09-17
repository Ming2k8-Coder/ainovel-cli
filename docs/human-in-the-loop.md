# Cơ Chế Con Người Trong Vòng Lặp (Human-In-The-Loop: Tối Thiểu Thao Tác - Tối Đa Bản Quyền)

Tài liệu này trình bày giải pháp kiến trúc **Human-In-The-Loop (HITL)** của `ainovel-cli`, được thiết kế với phương châm:

> **"Người dùng can thiệp ít nhất (Minimal Effort), nhưng vẫn xác lập đầy đủ căn cứ pháp lý để giữ 100% bản quyền tác phẩm (Maximum Copyright Protection)."**

---

## 1. Bài Toán & Triết Lý Thiết Kế

Theo quy định pháp lý quốc tế (Cục Bản quyền Tác giả Hoa Kỳ - USCO, Tổ chức Sở hữu Trí tuệ Thế giới - WIPO, Công ước Bern) và Luật Sở hữu Trí tuệ Việt Nam:
- **Tác phẩm do AI tạo ra 100% tự động** mà không có sự kiểm soát hoặc đóng góp sáng tạo có ý thức của con người thì **KHÔNG ĐƯỢC CẤP BẢN QUYỀN** (Public Domain / Không có tác giả hợp pháp).
- **Tác phẩm có sự tham gia của con người (AI-Assisted Work)** sẽ **ĐƯỢC BẢO HỘ BẢN QUYỀN TOÀN DIỆN** nếu tác giả con người thực hiện các quyền sau:
  1. **Định hướng sáng tạo ban đầu (Creative Direction)**: Thiết lập đề tài, nhân vật, xung đột gốc.
  2. **Tuyển chọn và sắp đặt nghệ thuật (Selection & Arrangement)**: Phê duyệt, lựa chọn cấu trúc đại cương hoặc các ngã rẽ cốt truyện.
  3. **Kiểm soát chất lượng và cấp phép xuất bản (Curation & Gate Approval)**: Trực tiếp nghiệm thu và cho phép phát hành từng chương.

Thay vì bắt người dùng phải gõ từng dòng văn bản tốn hàng trăm giờ, `ainovel-cli` xây dựng cơ chế **"Điểm Chạm Tối Giản" (Micro-Touchpoints)**: Mỗi hành động chỉ tốn vài giây nhưng được tự động ghi vết và ký số để cấu thành bằng chứng quyền tác giả không thể chối cãi.

---

## 2. Mô Hình 4 Điểm Chạm Sáng Tạo Tối Giản

```text
[1. Hạt Giống Ý Tưởng (Premise Seed)]
      │ Người dùng nhập 1 câu / 1 đoạn ý tưởng ban đầu (Mất 30 giây)
      ▼
   Architect Lập Khung & Hồi 1
      │
[2. Phê Duyệt Cấu Trúc (Outline Curation)]
      │ Người dùng bấm Enter duyệt đại cương / chỉnh 1 chi tiết nhỏ (Mất 5 giây)
      ▼
   Writer Sáng Tác Từng Chương (Hoàn toàn tự động)
      │
[3. Can Thiệp Tức Thời (Creative Steering - Tùy chọn)]
      │ Gõ 1 câu chuyển hướng khi đang viết: "Cho nhân vật B xuất hiện ở viện bảo tàng" (Mất 10 giây)
      ▼
   Editor Thẩm Định Bản Thảo
      │
[4. Cổng Cấp Phép Chương (Advance Gate)]
      │ Người dùng gõ /next duyệt chương (Mất 1 giây)
      ▼
[Sổ Cái Quyền Tác Giả (authorship_ledger.jsonl)] ➔ Xuất Hồ Sơ Bản Quyền (COPYRIGHT_REPORT.md)
```

---

## 3. Hệ Thống Sổ Cái Quyền Tác Giả (`authorship_ledger.jsonl`)

Mỗi khi người dùng thực hiện bất kỳ hành động nào trong 4 điểm chạm trên, hệ thống tự động:
1. Ghi lại nội dung văn bản can thiệp, định danh tác giả và số chương liên quan.
2. Đóng dấu thời gian chuẩn RFC3339.
3. Tạo mã băm mật mã học **SHA-256** liên kết chặt chẽ với bản thảo chương đó.

Tệp này được lưu trữ vĩnh viễn tại `meta/authorship_ledger.jsonl` và không thể làm giả hay đảo ngược.

---

## 4. Xuất Báo Cáo Bản Quyền Tác Giả

Bất cứ lúc nào trong quá trình sáng tác hoặc sau khi hoàn tất tác phẩm, người dùng có thể xuất hồ sơ bản quyền:

```bash
# Xuất hồ sơ bản quyền kèm chứng thực mã băm SHA-256
ainovel-cli copyright --author "Tên Của Bạn"
```

Hệ thống sẽ tổng hợp toàn bộ các can thiệp sáng tạo thành một tài liệu chuẩn mực `COPYRIGHT_AUTHORSHIP_REPORT.md` để đính kèm vào hồ sơ nộp Cục Bản quyền tác giả.
