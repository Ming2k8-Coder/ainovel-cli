# Cổng Tiến Độ Chương (Chapter Advance Gate)

> Trạng thái: Đã triển khai  
> Ngày: 2026-07-14  
> Giải quyết: Nghiệm thu từng chương, tạm dừng an toàn sau khi can thiệp, cấp phép chương chính xác khi khôi phục sau sự cố (crash recovery).

## 1. Tại Sao Cần Thiết?

Rủi ro cốt lõi của sáng tác tự động dài tập không phải là tiêu tốn thêm một lần gọi API, mà là trong lúc người dùng đang đọc nghiệm thu, hệ thống vẫn tiếp tục viết chương mới và đưa các tóm tắt, trạng thái nhân vật cũng như phản hồi đại cương (dựa trên cốt truyện cũ) vào nguồn thực tế tiếp theo. Việc xóa chương viết thừa không thể tự động hoàn tác các trạng thái phái sinh này, khiến người dùng mất niềm tin vào quy trình sáng tác.

Dự án vẫn định vị mặc định là "tự động hoàn thành liên tục sau khi nhận mục tiêu". Do đó, không biến việc xác nhận từng chương thành mặc định toàn cục. Hệ thống cung cấp 2 chính sách rõ ràng:

- `auto`: Chế độ mặc định, tự động đẩy tiến độ liên tục.
- `review`: Chế độ nghiệm thu từng chương do người dùng chủ động chọn, mỗi chương mới đều cần một giấy phép (permit) chính xác.

## 2. Ranh Giới Nhiệm Vụ

| Vấn đề | Thuộc về | Lý do |
|---|---|---|
| Chế độ hiện tại có phải nghiệm thu từng chương | RunMeta / Host | Ý định vận hành lâu dài của người dùng |
| Chương nào đã được cấp phép | RunMeta / Gate | Thực tế cơ khí có thể xác minh và khôi phục |
| Tiếp theo chạy Worker nào | `flow.Route` | Suy luận từ hàm thuần thực tế sáng tác |
| Lệnh có bắt đầu một chương mới tiến về phía trước không | `flow.StartsForwardChapter` | Phán đoán cơ khí phân loại |
| "Sửa xong cho tôi xem" có cần tạm dừng không | Arbiter | Phán đoán ngữ nghĩa ngôn ngữ tự nhiên |
| Khi nào kích hoạt tạm dừng | `ChapterAdvanceGate` | Thực thi định tính đối với ý định một lần |

## 3. Chế Độ Vận Hành

- `/review on`: Chuyển sang chế độ nghiệm thu từng chương.
- `/review off`: Chuyển về chế độ tự động.
- `/next`: Cấp phép cho chương tiếp theo (`NextChapter()`) và khởi chạy Engine.

## 4. Quy Tắc Thép
- Cổng tiến độ chỉ áp dụng cho việc **bắt đầu một chương mới**. Việc sửa đổi (rewrite), thẩm định (review), lập kế hoạch đại cương (outline) không bị chặn bởi cổng này.
- Khi ngắt máy/crash, giấy phép được đối chiếu chặt chẽ với trạng thái thực tế để đảm bảo không bị tiêu tốn nhầm cho chương tiếp theo.
