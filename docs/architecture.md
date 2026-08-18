# Kiến Trúc Thời Gian Chạy (Runtime Architecture) `ainovel-cli`

> **Mặt Phẳng Trạng Thái Định Tính, Mặt Phẳng Ngữ Nghĩa Tự Chủ**: Kết hợp giữa một Động cơ tuần tự định tính (Engine), ba Tác tử sáng tác tự chủ (Workers), các hàm Trọng tài ngữ nghĩa theo nhu cầu (Arbiter), và Tầng lưu trữ dữ liệu thực tế duy nhất trên hệ thống tập tin (Store).
>
> *(Kiến trúc cập nhật tháng 07/2026: Đã loại bỏ hoàn toàn vòng lặp Coordinator kéo dài để chuyển sang mô hình Engine + Arbiter hiện đại, ổn định và không tốn chi phí LLM trung gian).*

---

## 1. Mục Tiêu Thiết Kế (Theo Thứ Tự Ưu Tiên)

1. **Tính Ổn Định Tuyệt Đối**: Tiếp nhận một ý tưởng ban đầu, sáng tác liên tục và bền bỉ toàn bộ tác phẩm dài từ 200 đến 500+ chương mà không bị sập nguồn hay đứt gãy giữa chừng do lỗi kiến trúc.
2. **Khả Năng Tái Lặp & Tinh Chỉnh Chất Lượng**: Cho phép điều chỉnh độc lập hệ thống Prompt, tài liệu tham khảo văn phong, tiêu chí thẩm định chất lượng hoặc chiến lược phân bổ token mà không làm xáo trộn luồng điều phối chung.
3. **Khả Năng Tự Phục Hồi Sau Sự Cố**: Sau khi ứng dụng bị ngắt đột ngột, đứt kết nối mạng hoặc dừng tiến trình, hệ thống có thể khôi phục và tiếp tục công việc ngay tại điểm kiểm soát (Checkpoint) gần nhất mà không làm sai lệch dữ liệu.
4. **Khả Năng Quan Sát & Đo Lường (Observability)**: Dễ dàng theo dõi tiến độ, xem lại bản thảo, kiểm toán chi phí token và thời gian thực thi của từng phân đoạn công việc.

> *"Độ ổn định là nền móng, chất lượng tác phẩm là tầng cao."* Mọi quyết định kỹ thuật đều phải ưu tiên phục vụ sự vận hành bền bỉ của hệ thống.

---

## 2. Các Nguyên Tắc Kiến Trúc Cốt Lõi

### 2.1 Tam Phân Pháp: Phân Loại Quyết Định Theo Đúng Bản Chất

Hệ thống phân chia ranh giới xử lý thành ba nhóm rõ ràng:

- **Chuyển trạng thái có thể liệt kê ➔ Mã nguồn định tính (Code)**: Các quyết định như *"Sau khi hoàn thành chương này thì giao cho ai làm tiếp"* được xử lý bằng hàm thuần `flow.Route` thông qua việc đọc trạng thái từ Store. Toàn bộ các kịch bản chuyển trạng thái đều được kiểm thử vét cạn, tỷ lệ lỗi tiệm cận 0 và chi phí LLM bằng 0.
- **Phán đoán ngữ nghĩa có ranh giới rõ ràng ➔ Hàm LLM đơn lẻ (Arbiter)**: Các nhiệm vụ như chọn kịch bản dàn ý, phân loại can thiệp từ người dùng, tìm lối thoát khi Worker thất bại... tiếp nhận dữ liệu thực tế và trả về cấu trúc quyết định định dạng chuẩn. Mọi quyết định đều được ghi vết xuống đĩa (`decisions.jsonl`) để phục vụ phát lại offline.
- **Sáng tác mở ➔ Vòng lặp tác tử tự chủ (Worker Loop)**: Trong phạm vi một chương cụ thể, một lần thẩm định hoặc một đợt lập dàn ý, các tác tử Architect, Writer và Editor hoàn toàn tự chủ quyết định việc đọc tài liệu, lên ý tưởng và hoàn thiện văn bản.

```text
Mặt phẳng Định tính:  flow.LoadState   ➔ flow.Route     ➔ Instruction   (Kiểm thử quy cách vét cạn)
Mặt phẳng Ngữ nghĩa:  arbiter.Collect* ➔ arbiter.Decide* ➔ XxxDecision   (Lưu vết decisions.jsonl)
                      └── Thu thập dữ liệu (IO) ──┘└── Quyết định LLM ──┘└── Engine thực thi ──┘
```

---

### 2.2 Công Cụ Là Giao Diện Duy Nhất Của Tầng Dữ Liệu

Mọi thao tác đọc ghi với hệ thống tập tin, tiến độ (`Progress`) và điểm kiểm soát (`Checkpoint`) đều phải thông qua các Công cụ (Tools) được chuẩn hóa:
- Tập tin đơn lẻ áp dụng cơ chế ghi nguyên tử `temp + fsync + rename` để tránh ghi đè dở dang khi mất điện/crash.
- Ghi nhiều tập tin liên đới áp dụng mô hình giao dịch Saga (`PendingCommit`) có thể phục hồi.
- Mọi bước thực thi đều kiểm tra điều kiện tiên quyết và điều kiện hậu nghiệm chặt chẽ.

---

### 2.3 Tầng Quan Sát Tách Biệt Tuyệt Đối

Giao diện người dùng (TUI), nhật ký chẩn đoán (`diag`) và bộ theo dõi sự kiện (`observer`) chỉ đóng vai trò tiêu thụ dữ liệu thụ động từ luồng sự kiện:
- Chỉ đọc dữ liệu, không tự ý thay đổi dữ liệu trong Store.
- Không can thiệp hoặc làm thay đổi logic điều phối của Engine.
- Phân tách rõ: dữ liệu chi tiết cho nhật ký (`Detail`) và thông tin ngắn gọn, súc tích cho màn hình hiển thị (`Summary`).

---

### 2.4 Bốn Quy Tắc Bất Biến

1. **Quy tắc 1**: Công cụ chỉ trả về dữ liệu thực tế, tuyệt đối không trả về chỉ thị điều phối luồng.
2. **Quy tắc 2**: Định tuyến công việc do `Flow Router` phụ trách, việc thực thi do `Engine` đảm nhận.
3. **Quy tắc 3**: Các phán quyết ngữ nghĩa đều thông qua `Arbiter` và bắt buộc phải lưu vết xuống đĩa.
4. **Quy tắc 4**: Chỉ mã hóa cứng các quy tắc logic xác định, tuyệt đối không mã hóa cứng các phán đoán văn học mở.

---

## 3. Sơ Đồ Khối Kiến Trúc

```text
[Giao diện TUI / Chế độ Headless]
        │ Lệnh người dùng / Can thiệp
[Host / Engine Điều Phối]
   ├── observer            Thu thập tiến độ Worker + Phát sự kiện ra màn hình/nhật ký
   ├── engine              Vòng lặp định tính: LoadState ➔ Route ➔ Kiểm tra ➔ Chạy Worker ➔ Điểm chốt
   ├── luồng can thiệp     Tiếp nhận can thiệp ➔ Arbiter phán quyết ➔ Đưa vào hàng đợi
   └── usage / ngân sách   Giám sát hạn ngạch token, điểm dừng và quản lý mô hình
        │ Kích hoạt subagent.Runner.Run
[Architect · Writer · Editor] (Mỗi tác tử có ngữ cảnh và mô hình độc lập)
        │ Gọi các Công cụ (Tools)
[Tools]  novel_context · read_chapter · plan_chapter · draft_chapter · edit_chapter
         check_consistency · commit_chapter · save_review · save_arc_summary · save_foundation
        │ Ghi nguyên tử vào hệ thống tập tin
[Store: Hệ Thống Tập Tin]
   Progress · Checkpoints · Outline · Drafts · Summaries · Characters · World
   · Signals · Decisions (Nhật ký phán quyết) · Phản hồi dàn ý · Vi phạm quy tắc
```

---

## 4. Vòng Đời Vận Hành

1. **Khởi tạo (Start)**: Người dùng nhập đề tài ➔ `plan_start` Arbiter chọn kịch bản ➔ Khởi tạo dữ liệu nền tảng ➔ Engine bắt đầu vòng lặp sáng tác.
2. **Khôi phục (Resume)**: Khởi động lại sau sự cố ➔ Engine đọc `Store` ➔ `flow.Route` tính toán chỉ thị tiếp theo từ dữ liệu thực tế mà không cần phục hồi phiên chat cũ của LLM.
3. **Can thiệp (Steer)**: Người dùng nhập phản hồi ➔ `intervention` Arbiter phân loại phạm vi ➔ Ghi nhận hành động vào hàng đợi ➔ Thực thi an toàn tại ranh giới chương tiếp theo.
