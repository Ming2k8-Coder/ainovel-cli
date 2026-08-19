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

## 9. Quản Lý Phiên Bản Git & Quy Trình Phát Triển Tác Phẩm (`git`)

Tự động hóa và quản lý quy trình phát triển tác phẩm theo mô hình Git chuyên nghiệp:

### a) Phân Nhánh & Chuyển Nhánh (Branching & Checkout):
```bash
# Xem danh sách các nhánh
ainovel-cli git --list-branches --dir ./novel

# Tạo nhánh Git mới
ainovel-cli git --create-branch "feat/arc2-rebellion" --dir ./novel

# Chuyển sang nhánh chỉ định
ainovel-cli git --checkout "main" --dir ./novel
```

### b) Quản Lý Nhiệm Vụ Cốt Truyện (Doing Issues):
```bash
# Tạo một Issue nhiệm vụ mới
ainovel-cli git --create-issue "Đại tu diễn biến trận đánh Hồi 2" --desc "Chỉnh sửa nhịp kịch tính và phục bút" --dir ./novel

# Xem danh sách Issues
ainovel-cli git --list-issues --dir ./novel

# Bắt đầu làm Issue (tạo nhánh feat/ mới)
ainovel-cli git --start-issue "ISSUE-1" --make-feat --dir ./novel

# Bắt đầu làm Issue trực tiếp trên main (nếu không tạo nhánh feat)
ainovel-cli git --start-issue "ISSUE-1" --make-feat=false --dir ./novel
```

### c) Thẩm Định & Gộp Nhánh (Pull Request Review & Merge):
```bash
# Mở một Pull Request mới để kiểm duyệt
ainovel-cli git --create-pr "Đại tu Hồi 2" --source "feat/issue-1-dai-tu" --target "main" --dir ./novel

# Xem danh sách PRs
ainovel-cli git --list-prs --dir ./novel

# Biên tập viên / Tác giả Thẩm định & Phê duyệt PR
ainovel-cli git --review-pr "PR-1" --approve --comment "Bản thảo đạt chuẩn thẩm mỹ, nhịp kịch tính tốt" --dir ./novel

# Gộp nhánh PR đã duyệt vào main
ainovel-cli git --merge-pr "PR-1" --dir ./novel
```

---

## 10. Phân Tích Đặc Tả & Viết Tiếp Truyện Có Sẵn (`continue`)

Cho phép phân tích một cuốn tiểu thuyết hoặc tệp đặc tả dự án sẵn có (như các thư mục server trong workspace `nexus-universe`), tự động trích xuất Hồ sơ Thế giới, Nhân vật, Dàn ý và Phong cách văn phong, sau đó thiết lập dự án viết tiếp từ Chương N+1:

```bash
# Phân tích thư mục truyện có sẵn trong workspace và thiết lập dự án viết tiếp
ainovel-cli continue --source "../servers/HUST-sanbaka" --dir ./novel
```

---

## 11. Bộ Điều Khiển Ngôn Ngữ Tự Nhiên (`cmd` / `do`)

Cho phép tác giả ra lệnh bằng bất kỳ câu thoại tiếng Việt hoặc tiếng Anh tự nhiên nào. Hệ thống tự động phân tích ý định (Intent Recognition) và điều phối công việc tương ứng trong `ainovel-cli`:

```bash
# Ra lệnh bằng ngôn ngữ tự nhiên
ainovel-cli cmd "Phân tích truyện HUST-sanbaka trong workspace rồi tiếp tục viết"

ainovel-cli cmd "Mở WebUI trên cổng 3000 và xuất báo cáo bản quyền"

ainovel-cli cmd "Chuyển phong cách prompt sang cyberpunk"
```

---

## 12. Quản Lý & Truy Vết Vũ Trụ Truyện Quy Mô MCU (`universe` / `mcu`)

Giải quyết bài toán quản lý các bộ truyện chung vũ trụ (Shared Universe) siêu lớn và phức tạp (tương tự Marvel Cinematic Universe - MCU, Star Wars, SANBAKA Universe):

### a) Truy Vết Lịch Sử Thực Thể Xuyên Các Bộ Truyện (MCU Tracing):
Theo vết một nhân vật, bảo vật hoặc công nghệ qua nhiều bộ truyện song song để biết nhân vật đó đã xuất hiện ở đâu, năm nào, thuộc bộ truyện nào:
```bash
ainovel-cli universe --trace "Dương Hùng" --dir ./novel
```

### b) Ghi Nhận Vết Sự Kiện Mới Trong Vũ Trụ (Universe Event Record):
```bash
ainovel-cli universe --record "Nhân vật A" --series "HUST-sanbaka" --chapter 12 --year 1956 --desc "Nhận chip Quantum" --loc "ĐHBK Hà Nội" --dir ./novel
```

### c) Thẩm Định Xung Đột Định Luật Vũ Trụ (Canon Continuity Audit):
Tự động phát hiện các nghịch lý thời gian (Time Paradoxes) hoặc sự vi phạm định luật tối cao của vũ trụ giữa các bộ truyện:
```bash
ainovel-cli universe --check-canon --dir ./novel
```





