# Tài Liệu Đặc Tả Kỹ Thuật (RFC): Động Cơ Thực Thi Trực Tiếp Worker

> **Mục đích**: Giải quyết 4 vấn đề kỹ thuật trọng tâm khi chuyển dịch từ vòng lặp Coordinator sang Động cơ tuần tự định tính (Engine).
> *(Tham chiếu thiết kế kiến trúc tại `docs/engine-arbiter.md`).*

---

## 1. Cơ Chế Khởi Chạy Tác Tử (Worker Execution Plane)

Hàm `Runner.Run(agent, task)` sẽ khởi tạo một vòng lặp `agentcore.AgentLoop` hoàn chỉnh cho mỗi nhiệm vụ. Động cơ Engine gọi trực tiếp hàm này, đồng thời kế thừa nguyên vẹn toàn bộ các cấu hình đã thiết lập:
- Mô hình chuyên biệt cho từng vai trò kèm cơ chế tự động chuyển vùng dự phòng (Failover).
- Khóa định danh bộ nhớ đệm Prompt (`prompt cache key`) tăng tuần tự.
- Mức độ suy luận tư duy (`ThinkingLevel`).
- Bộ ghi nhận chi phí (`UsageRecorder`) và nhật ký phiên (`SessionLogger`).
- Bộ quản lý ngữ cảnh phân tầng (`ContextManagerFactory`) và gói khôi phục (`RestorePack`).

Toàn bộ kết quả định kiểu và chuỗi lỗi sẽ được trả về trực tiếp cho Engine mà không cần phải trải qua các bước giải mã JSON trung gian hay phán đoán kết quả công cụ phức tạp.

---

## 2. Vòng Đời Của Động Cơ (Engine Lifecycle)

- Động cơ vận hành dưới dạng một chuỗi tuần tự đơn (Single goroutine loop).
- Khi ngữ cảnh bị hủy (`ctx.Done()` hoặc lệnh Cancel), hệ thống sẽ lập tức tạm dừng an toàn và lan truyền tín hiệu vào vòng lặp Worker; các điểm kiểm soát (Checkpoints) đã lưu trên đĩa đảm bảo không làm thất thoát dữ liệu.
- Khi người dùng gửi lệnh tiếp tục (`Resume` hoặc `Continue`), hệ thống sẽ khởi tạo một chu kỳ mới và tính toán lại bước đi từ trạng thái thực tế trong Store.

---

## 3. Phân Loại Và Xử Lý Lỗi (Deterministic Error Taxonomy)

Hệ thống phân cấp xử lý lỗi rõ ràng, ưu tiên các biện pháp định tính trước khi gọi đến LLM:

1. **Lỗi tạm thời (Mạng gián đoạn, Giới hạn tốc độ, Luồng dữ liệu bị treo)**: Được tác tử tự động thử lại nội bộ (tối đa 7 lần) mà không làm thoát vòng lặp.
2. **Lỗi thực thi nghiêm trọng từ Worker**: Engine sẽ tự động thử lại cùng một chỉ thị 1 lần duy nhất. Nếu vẫn thất bại, hệ thống sẽ chuyển dữ liệu lỗi cho Trọng tài Arbiter (`worker_failure`) để ra phán quyết: Thử lại lần nữa (`retry`), Đổi hướng sang tác tử khác (`reroute`), hoặc Tạm dừng an toàn (`abort`).
3. **Lỗi cú pháp / Sai tham số / Tác tử không tồn tại**: Tạm dừng ngay lập tức và phát thông báo lỗi cho người dùng (đây là lỗi mã nguồn, việc thử lại tự động là vô nghĩa).

---

## 4. Giao Thức Chống Bế Tắc (Deadlock Prevention Protocol)

Sau mỗi vòng lặp, hệ thống ghi nhận khóa nhận diện nhiệm vụ `instructionKey(Agent + Chapter/Task)`. Nếu sau khi thực thi, bộ định tuyến `flow.Route` vẫn tiếp tục sinh ra cùng một khóa nhiệm vụ cũ (nghĩa là tác vụ trước đó chưa thỏa mãn điều kiện chuyển bước), biến đếm lặp `repeats` sẽ tăng thêm 1:

- **Khi `repeats == 3`**: Tự động tham vấn Trọng tài Arbiter để tìm giải pháp gỡ nút thắt bế tắc.
- **Khi `repeats == 5`**: Kích hoạt cơ chế ngắt mạch an toàn (Hard abort) — Tạm dừng tiến trình và gửi thông báo khẩn cấp để người dùng trực tiếp xử lý.
