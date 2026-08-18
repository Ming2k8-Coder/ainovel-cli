# Quy Trình Biên Dịch Ngữ Nghĩa Nhập Tiểu Thuyết (`/import`)

Quy trình nhập tiểu thuyết từ bên ngoài (`/import`) giúp phân tích một tệp văn bản thô (TXT/EPUB) và biên dịch nó thành cấu trúc dữ liệu chuẩn của `ainovel-cli`.

## 1. Các Giai Đoạn Biên Dịch (Compilation Pipeline)

```text
Văn bản thô (Raw File)
   ↓
[1. Ingest]       → Đọc và đánh chỉ số dòng / khối văn bản.
   ↓
[2. Segment]      → Phân đoạn ngữ nghĩa (xác định ranh giới chương, quyển, phần mở đầu/kết thúc).
   ↓
[3. Analyze]      → Phân tích từng chương (rút trích sự kiện, nhân vật, phục bút, loại hook).
   ↓
[4. Synthesize]   → Tổng hợp toàn thư (định hình premise, lập đại cương phân tầng, compass câu chuyện).
   ↓
[5. Publish]      → Xuất bản thành Store chính thức để AI có thể tiếp tục viết tiếp.
```

## 2. Đặc Điểm

- **Tái Tạo Đẳng Idempotent**: Mỗi giai đoạn đều gắn với vân tay (fingerprint) của dữ liệu đầu vào. Khi bị ngắt giữa chừng, hệ thống có thể khôi phục và chạy tiếp chính xác tại giai đoạn đó.
- **Không Trôi Dữ Liệu**: Chỉ trích xuất các thực tế **có thật trong văn bản thô**, không bịa đặt thêm các chi tiết chưa viết.
