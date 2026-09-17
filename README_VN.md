# ainovel-cli (Công cụ Tự Động Sáng Tác Tiểu Thuyết Bằng AI)

**ainovel-cli** là hệ thống AI sáng tác tiểu thuyết dài tập hoàn toàn tự động. Kiến trúc vận hành kết hợp giữa **Động cơ điều phối định tính (Deterministic Engine)** bằng mã nguồn chuẩn xác và **Mô hình ngôn ngữ lớn (LLM)** được kích hoạt đúng lúc đúng chỗ: Engine chịu trách nhiệm điều phối luồng công việc dựa trên trạng thái thực tế đến 3 Tác tử sáng tác độc lập (**Architect** / **Writer** / **Editor**), đồng thời tham vấn **Arbiter** (Trọng tài) khi cần đưa ra các phán quyết ngữ nghĩa then chốt.

Hệ thống cho phép chuyển hóa từ một câu ý tưởng ban đầu thành một tác phẩm tiểu thuyết hoàn chỉnh hàng trăm chương một cách mạch lạc, tự động và không cần sự can thiệp thủ công liên tục của con người.

---

## 🌟 Tính Năng Nổi Bật

- **Động Cơ Định Tính + Hệ Đa Tác Tử (Multi-Agent) Phối Hợp**: Engine điều phối 3 tác tử sáng tác chuyên biệt (**Architect** - Kiến trúc sư / **Writer** - Người viết / **Editor** - Biên tập viên) dựa trên bảng trạng thái thực tế. Vòng lặp điều phối chính hoàn toàn không tiêu tốn token LLM, đảm bảo tính ổn định và khả năng kiểm thử toàn diện.
- **Trọng Tài Ngữ Nghĩa Có Thể Kiểm Toán (Auditable Arbiter)**: Các quyết định mang tính ngữ nghĩa như lựa chọn kịch bản lập dàn ý, phân loại can thiệp từ người dùng, giải quyết xung đột hay gỡ nút thắt bế tắc... đều do Arbiter xử lý trong một lệnh gọi đơn lẻ. Mọi phán quyết đều được lưu vết để dễ dàng phát lại (replay) và kiểm thử hồi quy.
- **Khôi Phục Tiến Độ Cấp Bước (Step-level Checkpoint)**: Sau mỗi thao tác công cụ thành công, hệ thống lập tức lưu lại điểm kiểm soát (Checkpoint). Khi gặp sự cố mất mạng hoặc dừng tiến trình, hệ thống có thể khôi phục chính xác từng phân đoạn: `plan` (lập kế hoạch) ➔ `draft` (viết nháp) ➔ `check` (kiểm tra nhất quán) ➔ `commit` (nộp bản thảo).
- **Lập Dàn Ý Cuộn Phân Tầng (Rolling Arc/Volume Planning)**: Giải quyết triệt để vấn đề "dàn ý rỗng" ở các bộ truyện dài. Ban đầu hệ thống chỉ dựng khung tổng thể cho 2 Quyển đầu và chi tiết cho Hồi (Arc) 1. Khi tiến độ viết đạt mốc, Architect sẽ tiếp tục mở rộng các Hồi tiếp theo dựa trên tóm tắt diễn biến thực tế và hồ sơ nhân vật cập nhật.
- **Truy Hồi Thông Minh Các Chương Liên Quan**: Trong quá trình viết từng chương, hệ thống tự động đối chiếu và gợi ý các dữ kiện lịch sử liên quan theo 4 trục: *Phục bút (Foreshadowing), Sự xuất hiện của nhân vật, Biến chuyển trạng thái, Thay đổi mối quan hệ*, kết hợp với gợi mở của chương kế tiếp nhằm duy trì sự liền mạch cho tác phẩm từ 200 đến 500+ chương.
- **Chiến Lược Quản Lý Ngữ Cảnh Tự Thích Ứng (Adaptive Context)**: Tự động điều chỉnh linh hoạt giữa Ngữ cảnh đầy đủ, Cửa sổ trượt (Sliding Window) và Tóm tắt phân tầng dựa trên dung lượng tác phẩm, giúp kiểm soát ngân sách token tối ưu.
- **Thẩm Định Chất Lượng 7 Chiều (7-Dimensional Review)**: Editor đánh giá bản thảo dựa trên 7 tiêu chí cốt lõi: *Nhất quán thiết lập thế giới, Tính logic trong hành vi nhân vật, Nhịp điệu tình tiết (Pacing), Dòng chảy tự sự, Mạng lưới phục bút, Độ hấp dẫn của móc câu (Hook), và Chất lượng thẩm mỹ văn chương*. Tiêu chí thẩm mỹ được phân tích sâu qua: *Độ chân thực của miêu tả, Thủ pháp nghệ thuật, Khí chất thoại nhân vật, Độ tinh tế của từ ngữ, Sức truyền cảm* (kèm trích dẫn bằng chứng từ bản thảo).
- **Can Thiệp Tức Thời (Real-time Steering)**: Người dùng có thể nhập trực tiếp ý kiến điều chỉnh cốt truyện ngay khi hệ thống đang chạy. Trọng tài Arbiter sẽ phân tích phạm vi tác động và đưa các chương cần sửa vào hàng đợi làm lại mà không làm gián đoạn toàn bộ tiến trình.
- **Cơ Chế Duyệt Từng Chương Tùy Chọn (Step-by-Step Review Mode)**: Mặc định hệ thống tự động sáng tác 100%. Khi cần kiểm soát chi tiết, người dùng có thể kích hoạt `/review on`. Lệnh `/next` sẽ giải phóng chính xác 1 chương mới; các lượt viết lại do lỗi hay gián đoạn mạng sẽ không làm hao hụt hạn ngạch cấp phép.
- **Giao Diện Kép (TUI & Headless)**: Hỗ trợ giao diện Terminal trực quan sinh động (TUI) để theo dõi và điều khiển tại chỗ, đồng thời hỗ trợ chế độ chạy ngầm (Headless) cho máy chủ, VPS, NAS hoặc các đường ống CI/CD.
- **Hỗ Trợ Đa Dạng Nhà Cung Cấp LLM**: Tương thích mượt mà với OpenRouter, Anthropic (Claude), Google Gemini, OpenAI, DeepSeek, Ollama...

---

## 🏗️ Kiến Trúc Tổng Thể

Triết lý nền tảng: **Mặt Phẳng Trạng Thái Định Tính (Deterministic State) kết hợp Mặt Phẳng Ngữ Nghĩa Tự Chủ (Autonomous Semantics)**.

```text
┌──────────────────────────────────────────────────────────────┐
│                    Host / Engine (Định tính)                 │
│  Đọc Store ➔ Route ➔ Thực thi trực tiếp Worker ➔ Lặp lại     │
│  Xử lý Can thiệp / Bế tắc / Lỗi ➔ Tham vấn Trọng tài Arbiter │
└─────┬────────────┬────────────┬──────────────┬───────────────┘
      │            │            │              │
  ┌───▼────┐   ┌───▼────┐   ┌───▼────┐    ┌────▼────┐
  │Architect│  │ Writer │   │ Editor │    │ Arbiter │
  │(LLM Loop│  │(LLM Loop│  │(LLM Loop│   │(LLM Func│
  └───┬─────┘  └───┬────┘   └───┬────┘    └─────────┘
      └────────────┼────────────┘
                   │ Gọi Công cụ (Tools: IO + Checkpoint)
┌──────────────────▼───────────────────────────────────────────┐
│                            Store                             │
│  Tiến độ / Checkpoints / Dàn ý / Bản thảo / Thế giới / ...   │
└──────────────────────────────────────────────────────────────┘
```

### Phân Định Nhiệm Vụ Các Tác Tử

| Tác tử | Trách nhiệm chính | Công cụ (Tools) sử dụng |
|---|---|---|
| **Arbiter** (Trọng tài) | Đưa ra các phán quyết ngữ nghĩa: Khởi tạo kịch bản, phân loại can thiệp người dùng, chỉ định phương án khi Worker thất bại hoặc bế tắc. | Không gọi công cụ IO (Thực thi hàm LLM đơn lẻ, trả về dữ liệu cấu trúc). |
| **Architect** (Kiến trúc sư) | Xây dựng tiền đề (Premise), đại cương (Outline), thiết lập nhân vật, bối cảnh và quy tắc thế giới. | `novel_context`, `save_book`, `save_foundation` |
| **Writer** (Người viết) | Tự chủ tiếp nhận kế hoạch, đọc tư liệu, viết nháp, kiểm tra tính nhất quán và hoàn thiện bản thảo chương. | `novel_context`, `read_chapter`, `plan_chapter`, `draft_chapter`, `check_consistency`, `commit_chapter` |
| **Editor** (Biên tập viên) | Đọc bản thảo gốc, thẩm định chuyên sâu ở 2 tầng: Cấu trúc tình tiết & Thẩm mỹ văn học; tóm tắt Hồi/Quyển. | `novel_context`, `read_chapter`, `save_review`, `save_arc_summary`, `save_volume_summary` |

---

## 🔄 Quy Trình Sáng Tác Tiêu Chuẩn

```text
Ý tưởng người dùng ➔ Arbiter chọn kịch bản ➔ Architect dựng khung & Hồi 1 ➔ Writer viết từng chương ➔ Editor thẩm định Hồi
                                (Lưu Store)                                           ↑                     │
                                                                                      ├── Yêu cầu viết lại ◄┘
                                                                                      │
                                                                          Architect mở rộng Hồi/Quyển tiếp
                                                                          (Dựa trên Tóm tắt + Dữ liệu nhân vật)
```

### 6 Bước Chuẩn Hóa Của Writer Cho Mỗi Chương:
1. `novel_context`: Nạp gói ngữ cảnh (tóm tắt chương trước, mạng lưới phục bút, trạng thái nhân vật, văn phong chỉ định, dữ kiện liên quan).
2. `read_chapter`: Đọc lại phần đuôi của chương trước để bắt nhịp văn phong và cảm xúc liền mạch.
3. `plan_chapter`: Lập hợp đồng chương: mục tiêu, nhịp kịch tính, xung đột trung tâm, diễn biến tâm lý.
4. `draft_chapter`: Chấp bút sáng tác toàn bộ văn xuôi cho chương.
5. `check_consistency`: Đối chiếu bản thảo với các thiết lập thế giới và ràng buộc nhân vật.
6. `commit_chapter`: Nộp bản thảo chính thức, cập nhật dòng thời gian và lưu vết thực tế xuống đĩa.

---

## ⚡ Hướng Dẫn Khởi Chạy Nhanh

```bash
# Chạy trực tiếp mã nguồn bằng Go
go run ./cmd/ainovel-cli

# Hoặc biên dịch thành file nhị phân
go build -o ainovel ./cmd/ainovel-cli
./ainovel
```
