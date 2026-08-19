# Hướng Dẫn Các Tính Năng Nâng Cao (Advanced Features)

Nhánh `human-in-the-loop` bổ sung một bộ tính năng nâng cao toàn diện nhằm tối ưu hóa quyền sở hữu trí tuệ, mở rộng khả năng sáng tạo đa vũ trụ và tự động hóa quy trình xuất bản.

---

## 1. Chứng Thư Bản Quyền Mật Mã Học Merkle Tree (`audit-proof`)

Nhằm biến toàn bộ cuốn tiểu thuyết thành một bằng chứng pháp lý có giá trị chứng cứ tuyệt đối trước tòa án sở hữu trí tuệ và cơ quan cấp bản quyền quốc tế:
- Hệ thống tự động tính toán mã băm SHA-256 cho từng chương bản thảo và từng can thiệp sáng tạo của tác giả.
- Xây dựng **Cây Merkle Tree** và tính toán **Mã gốc Merkle Root** (`0x...`) duy nhất đại diện cho toàn bộ tác phẩm.
- Bất kỳ sự thay đổi dù chỉ 1 ký tự trong tương lai đều sẽ làm sai lệch Merkle Root, chứng minh rằng tác giả đã hoàn thiện và sở hữu tác phẩm tại thời điểm đóng dấu.

```bash
# Xuất chứng thư bản quyền số học Merkle Root
ainovel-cli audit-proof --author "Nguyễn Văn A" --title "Sanbaka Universe" --dir ./novel
```

---

## 2. Rẽ Nhánh Đa Vũ Trụ / Cốt Truyện Song Song (Multiverse Branching)

Cho phép tác giả "Fork" tác phẩm tại bất kỳ chương hoặc hồi nào để khám phá các kịch bản "Nếu như..." (What-If) mà không làm xáo trộn mạch truyện chính:
- Nhánh A: Nhân vật chính hy sinh ở Hồi 2 để cứu bạn bè.
- Nhánh B: Nhân vật chính mở khóa công nghệ Cyberpunk bí mật và sống sót.

Toàn bộ các nhánh được quản lý độc lập trong Store và có thể chuyển đổi dễ dàng.

---

## 3. Xuất Bản Đa Định Dạng Tự Động (`publish`)

Không cần phải định dạng hay chuyển đổi thủ công, hệ thống hỗ trợ xuất bản tác phẩm ngay lập tức:

### Xuất bản MDX (Cho Astro / Next.js / sanbaka-web):
```bash
# Tự động tạo các tệp chapter-001.mdx có Frontmatter chuẩn SEO
ainovel-cli publish --format mdx --dir ./novel --out ./dist/mdx
```

### Xuất bản Toàn Thư Markdown (Cho Sách điện tử / Kindle / Ebook):
```bash
# Gom toàn bộ tác phẩm thành một tệp Markdown hoàn chỉnh duy nhất
ainovel-cli publish --format full-md --dir ./novel --out ./dist
```

---

## 4. Máy Chủ Xem Truyện Web Cục Bộ (`serve`)

Khởi chạy máy chủ Web Reader với giao diện Dark Mode & Glassmorphic hiện đại tại `http://localhost:8080`, cho phép tác giả xem và đọc các chương truyện ngay trên trình duyệt:

```bash
ainovel-cli serve --dir ./novel --port 8080
```

---

## 5. Mạng Lưới Quản Lý Phục Bút & Bí Ẩn (`foreshadow`)

Theo dõi các manh mối đã gài, các nút thắt cốt truyện và các câu hỏi bí ẩn chưa giải quyết:

```bash
# Gài một chi tiết phục bút mới vào Chương 3
ainovel-cli foreshadow --add "Nhân vật A giấu chìa khóa vạn năng trong viện bảo tàng" --chapter 3 --category mystery --dir ./novel

# Đánh dấu tháo gỡ phục bút tại Chương 12
ainovel-cli foreshadow --resolve "fs-12345" --chapter 12 --note "Nhân vật B tìm thấy chìa khóa" --dir ./novel
```

---

## 6. Ma Trận Chỉ Số Tâm Lý & Quan Hệ Nhân Vật (`char-matrix`)

Theo dõi mức độ tin tưởng, căng thẳng và trạng thái liên minh giữa các nhân vật chính:

```bash
# Cập nhật quan hệ giữa hai nhân vật
ainovel-cli char-matrix -a "Nguyễn Văn A" -b "Trần Thị B" --trust 75 --tension 20 --alliance ally --chapter 5 --dir ./novel
```

---

## 7. Bảng Phân Tích Chi Phí Token & Hiệu Suất (`stats`)

Xuất bảng phân tích toàn diện về lượng Token đã tiêu thụ, chi phí API ước tính (USD), và chỉ số kiểm soát nghệ thuật của tác giả:

```bash
ainovel-cli stats --dir ./novel
```

---

## 8. Bảng Điều Khiển Trực Quan Trình Duyệt (`webui`)

Khởi chạy bảng điều khiển tương tác trực quan cao cấp trên trình duyệt (WebUI Dashboard) tại `http://localhost:3000`, hoạt động song song cùng với TUI (Terminal) và Headless:

```bash
ainovel-cli webui --dir ./novel --port 3000
```

---

## 9. Quản Lý Phiên Bản Git Tự Động (`git`)

Tự động hóa toàn bộ việc sao lưu, gắn nhãn mốc lịch sử và commit từng chương bản thảo vào Git Version Control:

```bash
# Kiểm tra trạng thái Git repository
ainovel-cli git --dir ./novel

# Khởi tạo Git repo cho tác phẩm
ainovel-cli git --init --dir ./novel

# Tạo nhãn mốc hoàn thành Quyển 1
ainovel-cli git --tag "vol1-complete" --dir ./novel
```


