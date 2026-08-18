# Tiến Hóa Mặt Phẳng Điều Khiển: Engine + Arbiter

> **Trạng thái**: Đã hoàn thiện và đưa vào sản xuất. Toàn bộ cơ chế Coordinator cũ dựa trên vòng lặp LLM kéo dài đã được xóa bỏ để chuyển sang kiến trúc Engine định tính kết hợp cùng Trọng tài ngữ nghĩa Arbiter.
>
> *(Chi tiết thiết kế hiện tại xem tại README_VN.md và `docs/engine-rfc.md`).*

---

## 1. Động Cơ Thay Đổi: Loại Bỏ Những Giả Định Lỗi Thời

Giả định ban đầu của dự án -- *"Duy trì một phiên LLM thường trực chạy vòng lặp dài để tự động điều phối toàn bộ cuốn sách"* -- đã bộc lộ nhiều điểm nghẽn nghiêm trọng khi mở rộng quy mô lên hàng trăm chương:

1. **Lãng phí chi phí token**: Trong hơn 90% số lần gọi điều phối, LLM chỉ làm nhiệm vụ chuyển tiếp nguyên văn chỉ thị từ bảng trạng thái mà không tạo thêm giá trị sáng tạo nào.
2. **Phát sinh mã vá lỗi phức tạp**: Để duy trì một phiên LLM không bị ngắt bất thường, hệ thống đã phải gánh thêm nhiều lớp bọc cồng kềnh (cơ chế chống dừng đột ngột, ép lệnh lặp lại, huấn luyện hành vi gượng gạo, đặt giới hạn lượt hội thoại lên tới 100.000).
3. **Thiếu tính xác định và khó kiểm thử**: Khó có thể kiểm thử hồi quy hoặc tái hiện lỗi một cách nhất quán khi mọi quyết định bị trộn lẫn trong một phiên chat dài dặc.

Toàn bộ các tiểu hệ thống mới của dự án (`import`, `simulation`, `userrules`) đều đã chứng minh hiệu quả vượt trội khi áp dụng mô hình: **"Mã nguồn Host trực tiếp điều phối + LLM đóng vai trò các hàm chức năng chuyên biệt"**. Việc cải tổ luồng chính nhằm đưa toàn bộ hệ thống về một chuẩn thống nhất và tin cậy.

---

## 2. Mô Hình Mục Tiêu Chuẩn Xác

```text
Điểm khởi chạy (Entry)
  ↓
Host (Vỏ bọc Engine tuần tự)
  ├─ Đọc dữ liệu từ Store ➔ flow.Route ➔ Thực thi trực tiếp Worker tương ứng
  ├─ Gặp tình huống ngữ nghĩa mở ➔ Triệu hồi hàm Trọng tài Arbiter
  └─ Cập nhật sự kiện / Giám sát ngân sách token / Kiểm soát điểm dừng
  ↓
Các tác tử Workers (Architect / Writer / Editor tự chủ hoàn toàn trong nhiệm vụ được giao)
  ↓
Các Công cụ (Tools) ➔ Ghi dữ liệu thực tế vào Store (Nguồn chân lý duy nhất)
```

### Phân Định Trách Nhiệm Rõ Ràng:
- **`flow.Route`**: Chịu trách nhiệm cho mọi bước đi tiếp theo có thể xác định qua tra bảng trạng thái.
- **`Arbiter`**: Chịu trách nhiệm cho các phán quyết ngữ nghĩa có phạm vi và cấu trúc rõ ràng.
- **`Workers`**: Chịu trách nhiệm cho quá trình sáng tạo văn học không giới hạn trong từng phân đoạn.
- **`Engine`**: Chịu trách nhiệm thực thi các quyết định, tuyệt đối không can thiệp vào thẩm định nghệ thuật.
- **`Observer / Diag`**: Chịu trách nhiệm quan sát và ghi nhận dữ liệu, không can thiệp luồng điều khiển.

---

## 3. Các Kịch Bản Tham Vấn Trọng Tài (Arbiter Scenarios)

| Kịch bản | Thời điểm kích hoạt | Nhiệm vụ chính |
|---|---|---|
| `plan_start` | Khởi tạo tác phẩm mới | Lựa chọn kịch bản lập dàn ý ngắn/dài tập phù hợp và mở rộng đề tài ban đầu. |
| `intervention` | Người dùng can thiệp | Phân loại yêu cầu của tác giả (thay đổi luật, điều chỉnh mạch truyện, yêu cầu viết lại). |
| `worker_failure` | Tác tử báo lỗi không thể tự phục hồi | Phân tích nguyên nhân và chỉ định hành động (Thử lại / Đổi hướng / Tạm dừng). |
| `deadlock` | Hệ thống lặp đi lặp lại cùng một nhiệm vụ | Đánh giá bế tắc (tối đa 3 lần tham vấn, sau 5 lần sẽ dừng an toàn chờ con người can thiệp). |

---

## 4. Giá Trị Đạt Được

- **Tiết kiệm tối đa chi phí**: Cắt giảm 100% chi phí LLM trung gian cho các bước chuyển giao chương.
- **Khả năng kiểm thử toàn diện**: Toàn bộ các phán quyết ngữ nghĩa được lưu vết và có thể phát lại (replay) để đánh giá hồi quy offline.
- **Phản hồi can thiệp tức thì**: Người dùng có thể can thiệp ngay lập tức mà không phải chờ đợi qua các chu kỳ trễ của vòng lặp cũ.
- **Đơn giản hóa mã nguồn**: Cắt giảm hơn 1.500 dòng mã bọc lót dư thừa, loại bỏ hoàn toàn các lỗi xung đột phiên làm việc.
