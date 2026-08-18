# ainovel-cli (Động cơ Sáng tác Tiểu thuyết AI Tự động)

Động cơ sáng tác tiểu thuyết dài tập tự động hoàn toàn bằng AI. Mã nguồn vận hành dựa trên hệ thống cơ sở mã định tính (Deterministic Engine) phối hợp với các mô hình ngôn ngữ lớn (LLM) được triệu hồi chính xác tại các điểm ra quyết định: Engine thực hiện định tuyến dựa trên dữ liệu thực tế để điều phối 3 Tác nhân Sáng tác Độc lập (**Architect** / **Writer** / **Editor**), đồng thời kích hoạt **Arbiter** khi cần đưa ra các phán quyết ngữ nghĩa.

Từ một câu ý tưởng ban đầu đến một bộ tiểu thuyết hoàn chỉnh, toàn bộ quy trình vận hành tự động không cần con người can thiệp.

## 🌟 Tính Năng Nổi Bật

- **Động cơ Định tính + Đa Tác nhân (Multi-Agent) Hợp tác**: Engine điều phối 3 đại lý sáng tác độc lập (**Architect** / **Writer** / **Editor**) dựa trên bảng quyết định thực tế. Vòng lặp chính hoàn toàn KHÔNG tốn chi phí LLM, hành vi có thể kiểm thử toàn diện.
- **Phán quyết Ngữ nghĩa Có thể Kiểm toán (Auditable Arbiter)**: Các quyết định lựa chọn KTS kịch bản, phân loại can thiệp người dùng, giải quyết bế tắc... đều do Arbiter thực hiện trong một lần gọi đơn. Mỗi phán quyết được lưu vết xuống đĩa để phát lại (replay). Đơn giản, ổn định, nói KHÔNG với điều phối phức tạp.
- **Khôi phục 断点 (Checkpoint) Chính xác cấp Bước (Step-level)**: Sau mỗi công cụ thực thi thành công, hệ thống ghi ngay checkpoint. Khi bị ngắt/crash, khả năng khôi phục chính xác từng bước: `plan` -> `draft` -> `check` -> `commit`.
- **Lập kế hoạch Cuộn 2 Lớp (Rolling Arc/Volume Planning)**: Tiểu thuyết dài tập không còn bị lập kế hoạch "rỗng" một lần cho tất cả. Ban đầu chỉ lập khung cho 2 Quyển đầu + chi tiết Arc (Tập) 1. Các Arc/Quyển tiếp theo sẽ được Architect mở rộng khi tiến độ viết chạm mốc, luôn tham chiếu tóm tắt trước đó và trạng thái nhân vật.
- **Gợi ý Thông minh Chương Liên quan**: Mỗi khi viết một chương, hệ thống tự động đối chiếu và gợi ý các chương lịch sử liên quan theo 4 chiều: *Phục bút (Foreshadowing), Nhân vật xuất hiện, Thay đổi trạng thái, Quan hệ*, kết hợp với hé lộ chương tiếp theo để đảm bảo tính liên tục cho bộ truyện 500+ chương.
- **Chiến lược Ngữ cảnh Tự điều chỉnh (Adaptive Context)**: Tự động chuyển đổi giữa Ngữ cảnh toàn bộ / Cửa sổ trượt (Sliding Window) / Tóm tắt phân tầng dựa trên tổng số chương, hỗ trợ truyện siêu dài 500+ chương.
- **Đánh giá Chất lượng 7 Chiều (7-Dimensional Review)**: Editor kiểm duyệt dựa trên 7 chiều: *Tính nhất quán cài đặt, Hành vi nhân vật, Nhịp điệu (Pacing), Mạch tự sự, Phục bút, Hook câu khách, và Chất lượng thẩm mỹ*. Chiều thẩm mỹ chia nhỏ thành 5 tiêu chí: *Chất cảm miêu tả, Thủ pháp tự sự, Phân độ thoại, Chất lượng từ ngữ, Sức lay động cảm xúc* (mỗi mục bắt buộc dẫn chứng đoạn văn gốc).
- **Can thiệp Người dùng Rơle (Real-time Steering)**: Trong quá trình AI đang viết, bạn có thể gõ trực tiếp yêu cầu chỉnh sửa vào khung nhập mà không cần bấm tạm dừng. Hệ thống tự đánh giá phạm vi ảnh hưởng và viết lại các chương chịu tác động.
- **Nghiệm thu Từng chương Tùy chọn (Optional Step-by-Step Review)**: Mặc định chạy tự động 100%. Khi cần kiểm soát tinh tế, bật `/review on`. Mỗi lệnh `/next` chỉ cho phép ra 1 chương mới, làm lại (rework) hay ngắt mạng không làm tiêu tốn nhầm lượt cấp phép.
- **Cổng vào Đôi: TUI (Giao diện dòng lệnh trực quan) + Headless (Chạy ngầm)**: Vừa có thể quan sát/can thiệp trực quan trên Terminal, vừa có thể treo chạy ngầm liên tục trên Server, NAS hoặc CI.
- **Hỗ trợ Đa LLM**: Tự do chuyển đổi giữa OpenRouter, Anthropic (Claude), Google Gemini, OpenAI, DeepSeek, Ollama...

---

## 🏗️ Kiến Trúc Hệ Thống

Triết lý thiết kế cốt lõi: **Tầng Thực tế Định tính (Deterministic), Tầng Ngữ nghĩa Tự chủ (Autonomous)**.

```text
┌─────────────────────────────────────────────────────────┐
│              Host / Engine (Định tính)                  │
│  Đọc Store → Route → Chạy trực tiếp Worker → Vòng lặp   │
│  Kích hoạt Phán quyết / Phân loại / Bế tắc → Gọi Arbiter │
└────┬──────────┬──────────┬─────────────┬────────────────┘
     │          │          │             │
 ┌───▼────┐ ┌───▼───┐ ┌────▼────┐   ┌────▼────┐
 │Architect│ │Writer │ │ Editor  │   │ Arbiter │
 │(LLM Loop│ │(LLM Loop│ │(LLM Loop│   │(LLM Func│
 └───┬────┘ └───┬───┘ └────┬────┘   └─────────┘
     └──────────┼──────────┘
                │ Gọi Công cụ Tool (IO + Checkpoint)
┌───────────────▼─────────────────────────────────────────┐
│                       Store                             │
│  Progress / Checkpoint / Outline / Drafts / ...         │
└─────────────────────────────────────────────────────────┘
```

### Phân công Vai trò Tác nhân (Agents)

| Vai trò | Trách nhiệm | Công cụ (Tools) sử dụng |
|---|---|---|
| **Arbiter** | Phán quyết ngữ nghĩa: Chọn KTS kịch bản ban đầu, phân loại can thiệp người dùng, tìm lối thoát khi thất bại/bế tắc. | Không (Gọi LLM đơn lẻ, trả về quyết định cấu trúc). |
| **Architect** | Tạo tên tác phẩm, giới thiệu, tiền đề (premise), đại cương (outline), hồ sơ nhân vật, quy tắc thế giới. | `novel_context`, `save_book`, `save_foundation` |
| **Writer** | Tự chủ hoàn thành ý tưởng, sáng tác, tự duyệt và nộp một chương. | `novel_context`, `read_chapter`, `plan_chapter`, `draft_chapter`, `check_consistency`, `commit_chapter` |
| **Editor** | Đọc bản thảo gốc, thẩm định ở 2 tầng: Cấu trúc câu chuyện & Thẩm mỹ văn học. | `novel_context`, `read_chapter`, `save_review`, `save_arc_summary`, `save_volume_summary` |

---

## 🔄 Quy Trình Viết Bài (Writing Workflow)

```text
Yêu cầu người dùng → Arbiter Chọn KTS → Architect Lập khung + Arc 1 → Writer Viết từng chương → Editor Kiểm duyệt Arc
                         (Lưu đĩa)                                          ↑                      │
                                                                           ├── Viết lại/Mài dũa ◄──┘
                                                                           │
                                                                    Architect Mở rộng Arc/Quyển tiếp
                                                                   (Tham chiếu Tóm tắt + Ảnh chụp nhân vật)
```

### Quy trình chuẩn của Writer cho từng chương:
1. `novel_context`: Tải ngữ cảnh (tóm tắt chương trước, phục bút, trạng thái nhân vật, quy tắc văn phong, gợi ý chương liên quan).
2. `read_chapter`: Đọc lại đoạn văn trước để bắt đúng văn phong và nhịp điệu.
3. `plan_chapter`: Phác thảo mục tiêu, mâu thuẫn xung đột, tuyến cảm xúc chương này.
4. `draft_chapter`: Viết toàn bộ nội dung chương.
5. `check_consistency`: Đổi chiếu dữ liệu kiểm tra tính nhất quán (bắt buộc thực hiện sau draft).
6. `commit_chapter`: Nộp bản thảo cuối cùng, ghi dữ liệu thực tế xuống đĩa.

---

## 📁 Hướng Dẫn Việt Hóa Cấu Trúc Mã Nguồn (Localizing the Codebase)

Để Việt hóa toàn bộ kho mã nguồn này cho dự án của bạn (đặc biệt là Vũ trụ **SANBAKA**), hãy chú ý các thư mục quan trọng sau:

1. **Thư mục Prompt hệ thống (`assets/prompts/`)**:
   - `architect-long.md` & `architect-short.md`: Việt hóa câu lệnh cho Kiến trúc sư kịch bản.
   - `writer.md`: Việt hóa câu lệnh hướng dẫn Tác giả viết văn xuôi (thêm văn phong Cyberpunk/Hard Sci-Fi tại đây).
   - `editor.md`: Việt hóa tiêu chuẩn kiểm duyệt 7 chiều.
   - `arbiter-*.md`: Việt hóa câu lệnh phán quyết trọng tài.

2. **Thư mục Văn phong & Quy tắc (`assets/styles/` & `assets/voice.md`)**:
   - Chứa định nghĩa giọng văn, quy tắc từ ngữ cấm/khuyên dùng.

3. **Giao diện TUI (`internal/entry/tui/`)**:
   - Chứa các chuỗi hiển thị bảng điều khiển Terminal.

---

## ⚡ Bắt Đầu Nhanh

```bash
# Chạy dự án bằng Go
go run ./cmd/ainovel

# Hoặc build thành file thực thi
go build -o ainovel ./cmd/ainovel
./ainovel
```
