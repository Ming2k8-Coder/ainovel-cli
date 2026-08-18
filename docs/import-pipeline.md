# Quy Trình Biên Dịch Ngữ Nghĩa Nhập Tác Phẩm (`/import`)

Tính năng nhập tiểu thuyết từ nguồn bên ngoài (`/import`) cho phép phân tích một tệp văn bản thô (định dạng TXT hoặc EPUB) và tự động biên dịch thành cấu trúc dữ liệu chuẩn của `ainovel-cli`, giúp AI có thể nắm bắt toàn bộ bối cảnh và tiếp tục viết các chương tiếp theo một cách mạch lạc.

---

## 1. Đường Ống Biên Dịch 5 Giai Đoạn (Compilation Pipeline)

```text
Văn bản thô (Tệp TXT / EPUB)
   ↓
[1. Tiếp Nhận (Ingest)]
   Đọc tệp, làm sạch mã hóa ký tự và đánh chỉ mục từng khối văn bản.
   ↓
[2. Phân Đoạn (Segment)]
   Tự động nhận diện ranh giới chương, chia tách hồi/quyển, phần mở đầu và kết thúc.
   ↓
[3. Phân Tích (Analyze)]
   Quét sâu từng chương để trích xuất sự kiện chính, hồ sơ nhân vật, mạng lưới phục bút và loại móc câu.
   ↓
[4. Tổng Hợp (Synthesize)]
   Tổng hợp toàn bộ tác phẩm: Đúc kết tiền đề (Premise), dựng đại cương phân tầng và la bàn định hướng (Compass).
   ↓
[5. Xuất Bản (Publish)]
   Ghi toàn bộ cấu trúc vào Store để các tác tử Architect, Writer và Editor sẵn sàng tiếp quản sáng tác.
```

---

## 2. Các Đặc Tính Kỹ Thuật Nổi Bật

- **Tính Đẳng Biến (Idempotency) & Tự Khôi Phục**: Mỗi giai đoạn biên dịch đều được gắn chặt với mã băm (fingerprint) của dữ liệu đầu vào. Nếu quá trình phân tích bị gián đoạn giữa chừng, hệ thống có thể tiếp tục xử lý chính xác tại bước dở dang mà không phải phân tích lại từ đầu.
- **Bảo Toàn Chân Lý Văn Bản**: Hệ thống chỉ trích xuất các dữ kiện thực sự xuất hiện trong văn bản gốc, tuyệt đối không suy diễn hay tự bịa đặt các chi tiết ngoài lề vào hồ sơ thế giới.
