# Thiết Kế Bộ Nhớ Đệm Lời Nhắc (Prompt Cache Design)

Cơ chế đệm lời nhắc (Prompt Caching) được thiết kế nhằm tối ưu hóa chi phí token và giảm độ trễ phản hồi khi sáng tác tiểu thuyết dài tập với các mô hình hỗ trợ bộ nhớ đệm (như Anthropic Claude, OpenAI, DeepSeek).

---

## 1. Nguyên Tắc Cốt Lõi: Bảo Toàn Tính Bất Biến Của Tiền Tố (Prefix Invariance)

Để nhà cung cấp mô hình nhận diện và tái sử dụng bộ nhớ đệm thành công, **chuỗi byte ở phần đầu của yêu cầu (Prompt Prefix) phải hoàn toàn giống hệt nhau giữa các lượt gọi liên tiếp**:

1. **Sắp Xếp Công Cụ Định Tính**: Toàn bộ định nghĩa (Schema) và mô tả của các công cụ Tool luôn được sắp xếp theo một thứ tự cố định, tuyệt đối không thay đổi ngẫu nhiên.
2. **Lịch Sử Hội Thoại Chỉ Nối Tiếp (Append-only History)**: Không chỉnh sửa hoặc chèn xen kẽ nội dung vào các tin nhắn đã gửi trong quá khứ.
3. **Nội Dung Động Nằm Ở Đuôi (Dynamic Content at Tail)**: Các chỉ thị mới, dữ liệu ngữ cảnh vừa cập nhật hoặc câu hỏi hiện tại luôn được đặt ở phần cuối cùng của thông điệp.

---

## 2. Cấu Trúc Khóa Định Danh Bộ Nhớ Đệm (Cache Key)

Hệ thống tuân thủ nguyên tắc cách ly rõ ràng: **Một Tác Phẩm - Một Gốc, Một Vai Trò - Một Định Danh, Mỗi Lần Khởi Chạy - Tăng Tuần Tự**.

- Cú pháp chuẩn: `nvl-<MãBămTácPhẩm>-<TênVaiTrò>#<SốThứTựLượtChạy>`
- Ví dụ: `nvl-9f8a2c-writer#12`
