# Hướng Dẫn Đăng Ký & Bảo Hộ Bản Quyền Tác Phẩm AI Sáng Tác

> **Mục tiêu**: Hướng dẫn chi tiết chiến lược pháp lý giúp tác giả xác lập, chứng minh và bảo hộ 100% quyền tác giả (Copyright Ownership) đối với các bộ tiểu thuyết được đồng sáng tạo cùng `ainovel-cli` với **thời gian và công sức bỏ ra ít nhất**.

---

## 1. Cơ Sở Pháp Lý Quốc Tế & Việt Nam

### 1.1 Nguyên Tắc Cốt Lõi Về Quyền Tác Giả & AI
Theo hướng dẫn của **Cục Bản quyền Tác giả Hoa Kỳ (USCO - Policy on AI-Generated Content 2023)**, **Tổ chức Sở hữu Trí tuệ Thế giới (WIPO)** và **Luật Sở hữu Trí tuệ Việt Nam (Luật số 50/2005/QH11 sửa đổi 2022)**:

1. **AI không phải là chủ thể pháp lý**: AI không thể đứng tên tác giả hay sở hữu bản quyền.
2. **Nguyên tắc đóng góp sáng tạo có ý thức (Modicum of Creativity)**: Bản quyền bảo hộ cho **yếu tố con người** trong tác phẩm. Nếu con người giữ vai trò là "Tổng đạo diễn" — chỉ đạo ý tưởng, tuyển chọn kết cấu, điều hướng tình tiết và nghiệm thu sản phẩm cuối cùng — tác phẩm sẽ được cấp bản quyền đầy đủ dưới danh mục **Tác phẩm phái sinh hoặc Tác phẩm đồng sáng tạo có sự hỗ trợ của công nghệ**.
3. **Tiêu chuẩn "Tuyển chọn và Sắp đặt" (Selection & Arrangement)**: Tương tự như việc biên tập viên tuyển chọn các bài viết thành một tuyển tập, việc con người duyệt dàn ý, lựa chọn tình tiết và điều phối các tác tử AI viết từng chương là một hành vi sáng tạo có bản quyền hợp pháp (Precedent: *Feist Publications v. Rural Telephone Service*).

---

## 2. Chiến Lược "Tối Thiểu Thao Tác - Tối Đa Pháp Lý" Với `ainovel-cli`

Để biến cuốn tiểu thuyết thành tác phẩm có bản quyền hợp pháp mà chỉ tốn từ **5 đến 10 phút cho toàn bộ cuốn truyện 100 chương**, bạn chỉ cần thực hiện đúng quy trình 3 bước sau:

### Bước 1: Khởi Tạo Đề Tài Bằng Ý Tưởng Riêng (1-2 Phút)
- Đừng chỉ nhập đề tài chung chung như *"Hãy viết một truyện tiên hiệp"*.
- Hãy đưa vào **các yếu tố sáng tạo độc bản** của bạn: tên nhân vật chính, bối cảnh đặc thù, quy tắc ma thuật độc đáo hoặc xung đột cốt lõi (Ví dụ: *"Vũ trụ Sanbaka 1956 tại ĐHBK Hà Nội, hệ thống tính điểm rèn luyện DRL kết hợp công nghệ Cyberpunk"*).
- `ainovel-cli` sẽ tự động ký mã băm SHA-256 đoạn văn bản này vào `authorship_ledger.jsonl` dưới danh mục `premise_seed`.

### Bước 2: Bật Chế Độ Duyệt Từng Chương (`/review on`) & Can Thiệp Tối Thiểu
- Kích hoạt chế độ `/review on`.
- Sau mỗi chương, bạn chỉ cần gõ `/next` để cấp phép cho AI viết chương tiếp theo (tốn 1 giây/chương).
- Thỉnh thoảng (khoảng 3-5 chương một lần), hãy gõ 1 câu can thiệp ngắn (`steer`) để định hướng nhân vật (Ví dụ: *"Cho nhân vật A tiết lộ bí mật ở cuối chương này"*).
- Mọi thao tác này đều được ghi nhận vào sổ cái như là bằng chứng của quyền kiểm soát nghệ thuật (Creative Direction & Curation).

### Bước 3: Xuất Báo Cáo Hồ Sơ Quyền Tác Giả Khi Hoàn Thành
- Khi sách viết xong, chạy lệnh xuất hồ sơ:
  ```bash
  ainovel-cli copyright --author "Nguyễn Văn A" --output ./meta/COPYRIGHT_REPORT.md
  ```
- Tệp báo cáo này chứa đầy đủ:
  - Bảng thống kê toàn bộ các can thiệp sáng tạo của bạn.
  - Danh sách mã băm SHA-256 liên kết với từng chương bản thảo.
  - Tuyên bố quyền kiểm soát nghệ thuật của con người theo chuẩn USCO Compendium III.

---

## 3. Mẫu Khai Báo Hồ Sơ Khi Nộp Cục Bản Quyền Tác Giả

Khi điền tờ khai đăng ký quyền tác giả tại Cục Bản quyền tác giả (hoặc nộp online trên cổng Dịch vụ công Quốc gia):

1. **Mục Tên Tác Giả**: Điền tên thật/bút danh của bạn.
2. **Mục Hình Thức Tác Phẩm**: Chọn *Tác phẩm văn học (Tiểu thuyết / Truyện dài)*.
3. **Mục Mô Tả Quá Trình Sáng Tạo (Bảo hộ phần đóng góp của con người)**:
   > *"Tác phẩm được tác giả trực tiếp lên ý tưởng, xây dựng cốt truyện, thiết lập nhân vật và điều phối quy trình chấp bút tự động với sự hỗ trợ công nghệ từ phần mềm mã nguồn mở ainovel-cli. Toàn bộ cấu trúc đại cương, các bước ngoặt tình tiết và bản thảo cuối cùng đều do tác giả trực tiếp thẩm định, phê duyệt và tuyển chọn nghệ thuật theo hồ sơ lưu vết đính kèm."*
4. **Tài liệu đính kèm**: Bản in tác phẩm + Tệp `COPYRIGHT_AUTHORSHIP_REPORT.md` (xuất từ `ainovel-cli`).

---

## 4. Tóm Lược 3 Nguyên Tắc Vàng

1. **Luôn có Hạt giống Ý tưởng của con người** (Không bao giờ để AI tự bịa 100% đề tài mà không có prompt của bạn).
2. **Luôn có ít nhất một vài lần Can thiệp / Bấm duyệt `/next`** (Thể hiện quyền kiểm soát tác phẩm).
3. **Luôn lưu trữ tệp `authorship_ledger.jsonl`** cùng bản thảo gốc để làm bằng chứng pháp lý bất khả xâm phạm.
