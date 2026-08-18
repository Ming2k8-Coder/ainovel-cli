# Tái Cấu Trúc Định Tuyến Theo Luồng Trạng Thái (Flow-Driven Routing)

Tài liệu này ghi lại quá trình cải tiến hệ thống điều phối của `ainovel-cli`: chuyển đổi từ việc phụ thuộc vào vòng lặp ngôn ngữ tự nhiên kéo dài sang cơ chế định tuyến định tính dựa trên hàm thuần `flow.Route`.

---

## 1. Nguyên Tắc Định Tuyến Tra Bảng Bằng Hàm Thuần

- **Chữ ký hàm `flow.Route(state State) -> *Instruction`**:
  - Nhận vào ảnh chụp trạng thái thực tế từ Store.
  - Không thực hiện bất kỳ thao tác đọc ghi đĩa (I/O) nào bên trong hàm.
  - Trả về con trỏ chỉ thị (`Instruction`) rõ ràng cho biết tác tử nào tiếp quản và thực hiện nhiệm vụ gì.
- **Kiểm Thử Vét Cạn**: Hàm được kiểm thử hồi quy trên hơn 120.000 tổ hợp trạng thái khác nhau, đảm bảo mọi tình huống ranh giới đều có chỉ thị định tính chuẩn xác và loại bỏ hoàn toàn chi phí LLM cho việc điều hướng.
