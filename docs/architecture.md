# Kiến Trúc Thời Gian Chạy (Runtime Architecture) `ainovel-cli`

> Tầng Thực tế Định tính, Tầng Ngữ nghĩa Tự chủ: Một Engine định tính chuỗi, ba Worker tự chủ, một vài hàm Arbiter theo nhu cầu, và một Tầng Thực tế Hệ thống Tập tin.
>
> 2026-07-12 Hoàn tất thay thế mặt phẳng điều khiển: Vòng lặp dài Coordinator LLM nghỉ hưu, nhường chỗ cho Engine (vòng lặp định tính) + Arbiter (hàm phán quyết ngữ nghĩa) tiếp quản. Xem hồ sơ thiết kế và đánh giá tại `docs/engine-arbiter.md`, RFC tại `docs/engine-rfc.md`.

---

## 1. Mục tiêu (Theo thứ tự ưu tiên)

1. **Độ ổn định**: Nhập một câu ý tưởng, viết hoàn chỉnh và ổn định toàn bộ cuốn tiểu thuyết (200~500 chương). Không tự động gián đoạn giữa chừng vì sự cố kiến trúc.
2. **Khả năng lặp lại chất lượng**: Prompt / Tài liệu tham khảo / Tiêu chí đánh giá / Chiến lược ngữ cảnh có thể điều chỉnh độc lập, không làm ảnh hưởng đến kiến trúc.
3. **Khả năng khôi phục**: Sau khi bị crash, mất mạng, tạm dừng, có thể tiếp tục ngay từ điểm checkpoint gần nhất.
4. **Khả năng quan sát (Observability)**: Có thể truy vấn tiến độ, sản phẩm, thời gian của từng chương và từng bước (step).

"Ổn định" là tiền đề, "Chất lượng" là tầng trên. Mỗi quyết định kiến trúc ưu tiên phục vụ độ ổn định.

---

## 2. Các Nguyên Tắc Cốt Lõi

### 2.1 Tam phân pháp: Quyết định được phân loại theo bản chất

- **Chuyển trạng thái có thể liệt kê → Mã nguồn (Code)**. "Viết xong một chương thì giao cho ai tiếp theo" là đọc thực tế tra bảng: Hàm thuần `flow.Route` + kiểm thử quy cách vét cạn hàng vạn tổ hợp, tỷ lệ lỗi tiệm cận 0, chi phí LLM bằng 0.
- **Phán đoán ngữ nghĩa có ranh giới rõ ràng → Hàm LLM (Arbiter)**. Chọn KTS kịch bản, phân loại can thiệp người dùng, tìm lối thoát khi thất bại/bế tắc: Dữ liệu thực tế vào, quyết định cấu trúc ra, kiểm tra cơ khí bọc lót, mỗi phán quyết ghi đĩa có thể phát lại.
- **Sáng tác mở → Vòng lặp LLM (Worker)**. Trong phạm vi một chương, một lần đánh giá, một lần lập kế hoạch, Architect / Writer / Editor hoàn toàn tự chủ.

Sự đối xứng giữa hai mặt phẳng là kỷ luật xuyên suốt -- bất kỳ điểm quyết định mới nào trong tương lai cũng phải tuân theo hình thái này, không phát minh mô hình mới:

```text
Mặt phẳng Định tính:  flow.LoadState   → flow.Route     → Instruction   (Kiểm thử quy cách vét cạn)
Mặt phẳng Ngữ nghĩa:  arbiter.Collect* → arbiter.Decide* → XxxDecision   (decisions.jsonl + đánh giá hồi quy)
                      └── Thu thập thực tế(IO) ──┘└── Cốt lõi(Tự tái tạo offline) ──┘└── Engine thực thi ──┘
```

### 2.2 Công cụ là giao diện duy nhất của tầng thực tế

Mọi tương tác với hệ thống tập tin, Progress, Checkpoint đều do công cụ (Tools) hoàn thành. Tập tin đơn lẻ sử dụng cơ chế thay thế nguyên tử `temp + fsync + rename`; ghi nối tiếp đa tập tin không giả mạo giao dịch cơ sở dữ liệu: Nộp chương sử dụng Saga `PendingCommit` bền vững, ghi cấu trúc sử dụng tái tạo phát lại đẳng idempotent và công khai lỗi rõ ràng. Mỗi bước đều phải kiểm tra lỗi; chỉ các quy trình đã lưu vết ý định khôi phục mới hứa hẹn khôi phục sau khi khởi động lại.

### 2.3 Tầng quan sát chỉ quan sát

UI, chẩn đoán, nhật ký sự kiện đều là các thành phần tiêu thụ thụ động được chiếu từ luồng sự kiện / công cụ chỉ đọc. Đọc thực tế, không tạo ra thực tế, không ảnh hưởng đến luồng điều khiển.

Dữ liệu quan sát được chia làm ba tầng nghiêm ngặt: `agentcore.ProgressPayload` là tầng truyền tải, văn bản lỗi phải đầy đủ và không chứa chiến lược cắt gọt UI; `host.Event.Summary` là ngữ nghĩa hiển thị ngắn, `Detail` là chẩn đoán đầy đủ; nhật ký tập tin ưu tiên ghi `Detail` đầy đủ, TUI chỉ đọc `Summary` và cắt gọt theo độ rộng terminal khi hiển thị.

### 2.4 Tầng thực tế phẳng

Chỉ có ba loại thực tế:
- **Progress** -- Chỉ số tiến độ (viết đến chương mấy, danh sách chờ viết lại)
- **Checkpoint** -- Ghi chép tiến trình cấp bước (plan / draft / commit / review / arc_summary)
- **Artifact** -- Nội dung chương, đại cương, nhân vật, tóm tắt...

Không đưa vào các trừu tượng như WorkflowInstance / TaskInstance / Command.

### 2.5 Bốn quy tắc thép

**Quy tắc thép 1: Công cụ chỉ trả về thực tế, không trả về chỉ thị chuyển giao điều phối**.
**Quy tắc thép 2: Định tuyến quy trình do Flow Router đảm nhận, thực thi do Engine đảm nhận**.
**Quy tắc thép 3: Phán quyết ngữ nghĩa đi qua Arbiter, mỗi phán quyết phải lưu vết xuống đĩa**.
**Quy tắc thép 4: Mã hóa cứng các ranh giới rành rọt, không mã hóa cứng các phán đoán ngữ nghĩa không thể liệt kê**.

---

## 3. Toàn Cảnh Kiến Trúc

```text
[Entry: TUI / headless]
        │ prompt / steer
[Host Vỏ bọc]
   ├── observer            Báo cáo tiến độ Worker + Engine phát sự kiện → UI/Nhật ký
   ├── engine              Vòng lặp định tính: LoadState → Route → Kiểm tra trước → Chạy Worker → Ranh giới lính gác
   ├── Luồng can thiệp     Steer/Continue → Arbiter phán quyết → Thực thi hành động
   └── usage / Ngân sách / Điểm dừng / Quản lý Model
        │ Gọi chương trình subagent.Runner.Run
[architect_short/long · writer · editor] (Mỗi Agent độc lập context + model)
        │ Gọi công cụ Tool
[Tools]  novel_context · read_chapter · plan_chapter · draft_chapter · edit_chapter
         check_consistency · commit_chapter · save_review · save_arc_summary
         save_volume_summary · save_foundation
        │ Đơn tập tin nguyên tử + Tái tạo đẳng idempotent
[Store: Hệ thống tập tin (tmp + rename)]
   Progress · Checkpoints · Outline · Drafts · Summaries · Characters · World
   · Signals · Decisions (Nhật ký phán quyết) · Bảng phản hồi · Vi phạm
```

---

## 4. Hướng Dẫn Vận Hành & Khôi Phục

1. **Tạo mới**: Nhập yêu cầu ban đầu ➔ `plan_start` Arbiter phán quyết ➔ Khởi tạo `RunMeta` ➔ Engine bắt đầu chạy.
2. **Khôi phục (Resume)**: Khi bị ngắt/crash, khởi động lại ➔ Đọc `Store` ➔ `Route` tính toán bước tiếp theo từ trạng thái thực tế. Không cần khôi phục session LLM.
3. **Can thiệp người dùng (Steer)**: Người dùng nhập yêu cầu ➔ `intervention` Arbiter phán quyết ➔ Thêm vào hàng đợi ➔ Thực thi tại ranh giới chương tiếp theo.
