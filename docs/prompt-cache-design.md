# Thiết Kế Bộ Nhớ Đệm Prompt (Prompt Cache Design)

Tối ưu hóa chi phí và tốc độ phản hồi cho tiểu thuyết dài tập bằng cơ chế Prompt Caching của các nhà cung cấp mô hình (Anthropic Claude, OpenAI, DeepSeek).

## 1. Nguyên Tắc Cốt Lõi

Để tận dụng hiệu quả Prompt Cache, **tiền tố byte của yêu cầu (Prompt Prefix) phải giữ nguyên ổn định giữa các lượt gọi**:

1. **Thứ tự Công cụ Định tính**: Danh sách các Description và Schema của công cụ Tool luôn được sắp xếp theo thứ tự cố định.
2. **Lịch sử Chỉ Thêm Không Sửa (Append-only History)**: Không chỉnh sửa các tin nhắn đã gửi trong quá khứ.
3. **Nội Dung Động Ở Đuôi (Dynamic Content at Tail)**: Các chỉ thị mới hoặc thông tin thời gian chạy được nối vào cuối tin nhắn, tuyệt đối không chèn ngược lại đầu hội thoại.

## 2. Cấu Trúc Khóa Cache

Khóa Cache tuân thủ nguyên tắc: **Một Sách - Một Gốc, Một Vai Trò - Một Tên, Một Phiên - Một Khóa**.
- Định dạng: `nvl-<HashSách>-<VaiTrò>#<STT_Spawn>`
