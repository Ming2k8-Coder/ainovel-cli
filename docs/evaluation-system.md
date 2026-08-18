# Hệ Thống Đánh Giá Ngoại Tuyến (Offline Evaluation System) `ainovel-cli`

Hệ thống đánh giá (`eval`) được thiết kế để đo lường, so sánh và kiểm thử chất lượng sáng tác tiểu thuyết bằng AI mà không phụ thuộc vào trải nghiệm TUI thời gian thực.

## 1. Mục Tiêu

1. **Đánh giá A/B Prompt & Văn Phong**: So sánh hiệu quả của các System Prompt, bộ hướng dẫn văn phong (Voice), hoặc các cài đặt model khác nhau đối với cùng một bộ đề văn học.
2. **Kiểm Thử Hồi Quy (Regression Testing)**: Đảm bảo các đợt cập nhật mã nguồn hoặc thay đổi prompt không làm giảm chất lượng văn phong hoặc phá vỡ các quy tắc nhất quán cốt truyện.
3. **Phát Lại Phán Quyết Offline**: Phát lại và đánh giá các phán quyết ngữ nghĩa của Arbiter dựa trên file nhật ký `decisions.jsonl`.

## 2. Cách Sử Dụng Lệnh `eval`

```bash
# Chạy đánh giá offline từ dòng lệnh
ainovel-cli eval --dataset ./evals/dataset.json --output ./evals/results/
```
