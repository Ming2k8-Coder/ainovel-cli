# Tái Cấu Trúc Điều Hướng Theo Luồng (Flow-Driven Refactoring)

Tài liệu này ghi lại quá trình chuyển đổi hệ thống điều phối từ vòng lặp chuỗi sang điều hướng định tính dựa trên hàm thuần `flow.Route`.

## 1. Nguyên Tắc Tra Bảng Định Tính

- **`flow.Route(state) → *Instruction`**: Là một hàm thuần (pure function) nhận dữ liệu thực tế từ Store và trả về chỉ thị công việc tiếp theo cho Engine.
- Được kiểm thử quy cách vét cạn với hơn 120,000 tổ hợp trạng thái, đảm bảo tỷ lệ lỗi tiệm cận 0 và chi phí LLM bằng 0 cho việc định tuyến.
