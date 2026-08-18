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
