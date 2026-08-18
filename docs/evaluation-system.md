# Hệ Thống Thẩm Định & Đánh Giá Offline (Evaluation System)

Hệ thống đánh giá ngoại tuyến (`eval`) được thiết kế nhằm mục đích đo lường, phân tích và so sánh chất lượng sáng tác văn học của các mô hình AI một cách khoa học mà không phụ thuộc vào giao diện hiển thị thời gian thực.

---

## 1. Các Mục Tiêu Cốt Lõi

1. **Thử Nghiệm A/B Prompt & Văn Phong (A/B Testing)**: Cho phép so sánh trực quan hiệu quả biểu đạt giữa các phiên bản System Prompt, bộ hướng dẫn văn phong (Voice), hoặc các tham số nhiệt độ (Temperature) khác nhau trên cùng một bộ đề văn học mẫu.
2. **Kiểm Thử Hồi Quy (Regression Testing)**: Đảm bảo rằng mọi đợt nâng cấp mã nguồn, thay đổi cấu trúc công cụ hay tinh chỉnh prompt không làm suy giảm chất lượng hành văn, không làm phá vỡ tính nhất quán của cốt truyện và không gây ra hiện tượng lặp từ ngữ.
3. **Phát Lại & Đánh Giá Phán Quyết Của Trọng Tài (Arbiter Replay)**: Đọc lại toàn bộ lịch sử các phán quyết ngữ nghĩa từ tập tin nhật ký `decisions.jsonl` để đánh giá độ chính xác, tính hợp lý và sự ổn định của Trọng tài Arbiter qua từng phiên bản mô hình.

---

## 2. Hướng Dẫn Sử Dụng Lệnh `eval`

Hệ thống cung cấp công cụ dòng lệnh chuyên biệt để tự động hóa việc đánh giá hàng loạt:

```bash
# Chạy bộ thẩm định chất lượng trên tập dữ liệu chuẩn
ainovel-cli eval --dataset ./evals/dataset.json --output ./evals/results/

# Đánh giá so sánh văn phong giữa hai mô hình
ainovel-cli eval --compare --model-a claude-3-5-sonnet --model-b gpt-4o --prompt-set ./evals/prompts/
```

Kết quả thẩm định sẽ được xuất bản dưới dạng báo cáo tổng hợp kèm các chỉ số thống kê chi tiết về độ trôi phục bút, tần suất từ mệt mỏi và điểm thẩm mỹ văn học 7 chiều.
