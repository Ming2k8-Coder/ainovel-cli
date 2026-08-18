# Tổng Hợp Sơ Đồ Hệ Thống & Đường Ống Kiến Trúc (Architecture & Pipeline Diagrams)

Tài liệu này tổng hợp toàn bộ các sơ đồ **Mermaid** trực quan hóa kiến trúc, quy trình sáng tác, máy trạng thái, đường ống nén ngữ cảnh và cơ chế sở hữu trí tuệ của `ainovel-cli`.

---

## 1. Kiến Trúc Tổng Thể Hai Mặt Phẳng (Dual-Plane Runtime Architecture)

Sự phân tách giữa **Mặt Phẳng Trạng Thái Định Tính (Deterministic Engine Plane)** và **Mặt Phẳng Ngữ Nghĩa Tự Chủ (Autonomous Semantics Plane)**.

```mermaid
flowchart TB
    subgraph UI["Giao Diện & Điểm Khởi Chạy"]
        TUI["TUI Dashboard (Terminal)"]
        Headless["Headless Mode (Server/CI)"]
    end

    subgraph HostEngine["Mặt Phẳng Điều Phối Định Tính (Deterministic Host / Engine)"]
        Engine["Engine Loop (Vòng lặp chuỗi đơn)"]
        Router["flow.Route(State)<br/>(Hàm thuần tra bảng 0 Token)"]
        Gate["Chapter Advance Gate<br/>(Khóa kiểm soát tiến độ)"]
        Observer["Observer & Telemetry<br/>(Giám sát sự kiện)"]
        Budget["Budget Sentinel<br/>(Giám sát ngân sách Token)"]
    end

    subgraph Semantics["Mặt Phẳng Ngữ Nghĩa (Autonomous Semantics)"]
        Arbiter["Trọng tài Arbiter<br/>(LLM-as-a-Function: plan_start, steer, failure, deadlock)"]
        Architect["Architect Agent<br/>(Dàn ý, nhân vật, thế giới)"]
        Writer["Writer Agent<br/>(Sáng tác văn xuôi từng chương)"]
        Editor["Editor Agent<br/>(Thẩm định 7 chiều & tóm tắt)"]
    end

    subgraph StoreLayer["Tầng Dữ Liệu Thực Tế Bất Biến (FileSystem Store)"]
        Progress["meta/progress.json"]
        Checkpoints["checkpoints/*.json"]
        Drafts["drafts/ & chapters/"]
        World["world/ & characters.json"]
        Authorship["meta/authorship_ledger.jsonl"]
        Decisions["meta/decisions.jsonl"]
    end

    UI --> HostEngine
    Engine --> Router
    Router -->|Tra bảng trạng thái| Engine
    Engine -->|Chỉ thị công việc| Architect
    Engine -->|Chỉ thị công việc| Writer
    Engine -->|Chỉ thị công việc| Editor
    Engine -->|Tình huống ngữ nghĩa mở| Arbiter
    Arbiter -->|Phán quyết cấu trúc| Engine
    Gate -->|Cho phép / Tạm dừng| Engine

    Architect -->|Gọi Tools: IO + Checkpoints| StoreLayer
    Writer -->|Gọi Tools: IO + Checkpoints| StoreLayer
    Editor -->|Gọi Tools: IO + Checkpoints| StoreLayer
    StoreLayer -.->|LoadState facts| Router
    StoreLayer -.->|CollectInterventionFacts| Arbiter
    Engine -.->|Cập nhật sự kiện| Observer
    Observer -.-> UI
```

---

## 2. Quy Trình Sáng Tác Toàn Vòng Đời (End-to-End Novel Lifecycle)

Quy trình từ ý tưởng ban đầu đến tác phẩm hoàn chỉnh hàng trăm chương.

```mermaid
sequenceDiagram
    autonumber
    actor User as Tác Giả (Human)
    participant Host as Host / Engine
    participant Arbiter as Trọng tài Arbiter
    participant Architect as Kiến trúc sư Architect
    participant Writer as Người viết Writer
    participant Editor as Biên tập viên Editor
    participant Store as Hệ thống Store (Đĩa)

    User->>Host: Nhập đề tài & ý tưởng ban đầu
    Host->>Store: Ghi nhận Premise Seed vào Authorship Ledger
    Host->>Arbiter: Tham vấn plan_start (Chọn kịch bản & mở rộng yêu cầu)
    Arbiter-->>Host: Quyết định kịch bản (architect_long / architect_short)
    Host->>Architect: Phân công xây dựng đại cương & Hồi 1
    Architect->>Store: Lưu premise.md, characters.json, outline.json
    
    loop Từng Chương Bản Thảo
        Host->>Writer: Chỉ thị viết Chương N
        Writer->>Store: Nạp novel_context, lập kế hoạch, viết nháp, commit_chapter
        Store-->>Host: Trả về trạng thái chương hoàn tất
        
        opt Nếu bật chế độ /review on
            Host->>Host: Chapter Advance Gate kích hoạt (Chờ cấp phép)
            User->>Host: Lệnh /next duyệt chương
            Host->>Store: Ký duyệt Chapter Approval vào Authorship Ledger
        end
    end

    opt Cuối Hồi (Arc End) hoặc Cuối Quyển (Volume End)
        Host->>Editor: Thẩm định Hồi & Tạo tóm tắt
        Editor->>Store: Lưu save_review (scope=arc), save_arc_summary
        alt Phát hiện lỗi nghiêm trọng
            Editor->>Store: Đưa các chương lỗi vào PendingRewrites
            Host->>Writer: Chỉ thị viết lại / mài dũa các chương lỗi
        else Đạt chuẩn chất lượng
            Host->>Architect: Mở rộng Hồi tiếp theo (expand_arc / append_volume)
            Architect->>Store: Cập nhật layered_outline.json
        end
    end

    Host->>Store: Kiểm tra điều kiện hoàn thành toàn thư (MarkComplete)
    Host->>User: Xuất bản báo cáo hoàn thành & Chứng thư bản quyền
```

---

## 3. Đường Ống 6 Bước Chấp Bút Mỗi Chương Của Writer (Writer 6-Step Pipeline)

Mỗi chương được tác tử Writer tự chủ thực hiện qua 6 bước nghiêm ngặt kèm điểm kiểm soát Checkpoint.

```mermaid
flowchart TD
    Start(["Nhận lệnh viết Chương N"]) --> Step1["1. novel_context<br/>Nạp tóm tắt, phục bút, nhân vật, văn phong, gợi ý liên quan"]
    Step1 --> Step2["2. read_chapter<br/>Đọc lại 800 từ cuối chương N-1 để bắt đúng nhịp văn"]
    Step2 --> Step3["3. plan_chapter<br/>Lập hợp đồng chương: mục tiêu, nhịp kịch tính, cảm xúc, xung đột"]
    Step3 --> Step4["4. draft_chapter<br/>Chấp bút sáng tác toàn bộ văn xuôi cho chương"]
    Step4 --> Step5["5. check_consistency<br/>Đối chiếu bản thảo với thiết lập thế giới và ràng buộc nhân vật"]
    Step5 --> CheckPass{"Vượt qua<br/>kiểm tra?"}
    CheckPass -- "Không đạt" --> Step4
    CheckPass -- "Đạt" --> Step6["6. commit_chapter<br/>Nộp bản thảo chính thức, lưu Saga PendingCommit, cập nhật timeline"]
    Step6 --> Done(["Hoàn thành chương N"])

    classDef step fill:#e1f5fe,stroke:#0288d1,stroke-width:2px;
    classDef finish fill:#e8f5e9,stroke:#388e3c,stroke-width:2px;
    class Step1,Step2,Step3,Step4,Step5,Step6 step;
    class Done finish;
```

---

## 4. Đường Ống Thẩm Định 7 Chiều Của Editor (Editor 7-Dimensional Quality Gate)

Editor thẩm định bản thảo trên 2 tầng: **Cấu trúc tình tiết** và **Thẩm mỹ văn học**.

```mermaid
flowchart LR
    subgraph Input["Đầu Vào Bản Thảo"]
        Draft["Chương hoàn thành<br/>(chapters/NN.md)"]
        Context["Hồ sơ thế giới & Quy tắc<br/>(world_rules, style_rules)"]
    end

    subgraph Dimensions["7 Chiều Đánh Giá Chất Lượng"]
        D1["1. Nhất quán bối cảnh<br/>(World Consistency)"]
        D2["2. Tính cách nhân vật<br/>(Character Logic)"]
        D3["3. Nhịp điệu tình tiết<br/>(Pacing & Tension)"]
        D4["4. Mạch tự sự<br/>(Narrative Cohesion)"]
        D5["5. Mạng lưới phục bút<br/>(Foreshadowing)"]
        D6["6. Độ giật gân móc câu<br/>(Hook Quality)"]
        D7["7. Thẩm mỹ văn học<br/>(Aesthetic Excellence)"]
    end

    subgraph Aesthetic["Chi Tiết Thẩm Mỹ 5 Tiêu Chí"]
        A1["• Chất cảm miêu tả"]
        A2["• Thủ pháp nghệ thuật"]
        A3["• Khí chất thoại nhân vật"]
        A4["• Độ tinh tế từ ngữ"]
        A5["• Sức lay động cảm xúc"]
    end

    subgraph Outcome["Kết Luận Thẩm Định"]
        Pass["✅ Đạt yêu cầu<br/>(Tiếp tục chương sau)"]
        Rework["⚠️ Cần làm lại<br/>(Vào hàng đợi PendingRewrites)"]
    end

    Input --> Dimensions
    D7 -.-> Aesthetic
    Dimensions -->|Tổng hợp điểm số| Outcome
```

---

## 5. Cơ Chế Quản Lý Ngữ Cảnh & Nén Zero-LLM (Adaptive Context Pipeline)

Quy trình nén ngữ cảnh tự động từ chi phí thấp đến cao khi dung lượng token vượt ngưỡng an toàn.

```mermaid
flowchart TD
    Trigger(["Token hội thoại > Ngưỡng Context Budget"]) --> Stage1["1. ToolResultMicrocompact<br/>Rút gọn kết quả công cụ cũ thành text giữ chỗ<br/><b>Chi phí: 0 Token</b>"]
    Stage1 --> Check1{"Dung lượng<br/>đã an toàn?"}
    Check1 -- "Có" --> Safe(["Tiếp tục sáng tác"])
    Check1 -- "Chưa" --> Stage2["2. LightTrim<br/>Cắt ngắn các đoạn văn bản quá dài từ vòng trước<br/><b>Chi phí: 0 Token</b>"]
    
    Stage2 --> Check2{"Dung lượng<br/>đã an toàn?"}
    Check2 -- "Có" --> Safe
    Check2 -- "Chưa" --> Stage3["3. StoreSummaryCompact (Độc quyền Writer)<br/>Nạp trực tiếp tóm tắt chương & hồ sơ từ Store trên đĩa<br/><b>Chi phí: 0 Token LLM</b> - 100% chính xác"]
    
    Stage3 --> Check3{"Dung lượng<br/>đã an toàn?"}
    Check3 -- "Có" --> Safe
    Check3 -- "Chưa" --> Stage4["4. FullSummary (Bọc lót cuối cùng)<br/>Gọi LLM tóm tắt lại diễn biến với Prompt chuyên biệt"]
    Stage4 --> Safe
```

---

## 6. Cơ Chế Human-In-The-Loop & Sổ Cái Quyền Tác Giả Merkle Tree

Mô hình "Tối thiểu thao tác - Tối đa bằng chứng bản quyền" và cây Merkle Root.

```mermaid
flowchart TD
    subgraph HumanActions["4 Điểm Chạm Sáng Tạo Của Con Người (Micro-Touchpoints)"]
        Seed["1. Ý tưởng ban đầu (Premise Seed)"]
        OutlineCurate["2. Duyệt khung dàn ý (Outline Curation)"]
        ABSelect["3. Lựa chọn phương án A/B (Creative Choice)"]
        GateApprove["4. Cấp phép chương /next (Chapter Approval)"]
    end

    subgraph Ledger["Sổ Cái Bất Biến (meta/authorship_ledger.jsonl)"]
        Rec1["Bản ghi 1: SHA-256 (Seed)"]
        Rec2["Bản ghi 2: SHA-256 (Outline)"]
        Rec3["Bản ghi 3: SHA-256 (A/B Choice)"]
        Rec4["Bản ghi 4: SHA-256 (Approval)"]
    end

    subgraph MerkleTree["Cây Merkle Tree Toàn Tác Phẩm"]
        Node1["Hash(Ch01 + Rec1)"]
        Node2["Hash(Ch02 + Rec2)"]
        Node3["Hash(Ch03 + Rec3)"]
        Node4["Hash(Ch04 + Rec4)"]
        Parent1["Hash(Node1 + Node2)"]
        Parent2["Hash(Node3 + Node4)"]
        MerkleRoot["👑 MERKLE ROOT: 0x7a8f...9b2c<br/>(Đại diện cho toàn bộ tác phẩm)"]
    end

    subgraph OutputLegal["Hồ Sơ Pháp Lý Đăng Ký Bản Quyền"]
        Doc1["COPYRIGHT_AUTHORSHIP_REPORT.md<br/>(Chuẩn USCO Compendium III / WIPO)"]
        Doc2["LEGAL_AFFIDAVIT.md<br/>(Chứng thư mã băm số học)"]
    end

    HumanActions --> Ledger
    Ledger --> MerkleTree
    Node1 --> Parent1
    Node2 --> Parent1
    Node3 --> Parent2
    Node4 --> Parent2
    Parent1 --> MerkleRoot
    Parent2 --> MerkleRoot
    MerkleTree --> OutputLegal
```

---

## 7. Vòng Lặp Phản Hồi Lời Nhắc A/B (A/B Selection & Prompt Feedback Loop)

Quy trình tự động học gu thẩm mỹ của tác giả và tiêm trực tiếp vào Lời nhắc Hệ thống.

```mermaid
flowchart LR
    A["Phương án A<br/>Thoại sắc bén / Nhịp nhanh"] --> UserSelect{"Tác giả chọn<br/>(Chỉ mất 2s)"}
    B["Phương án B<br/>Tâm lý / Thuyết minh"] --> UserSelect
    UserSelect -->|Chọn A| Record["Lưu Choice vào ab_feedback.jsonl<br/>+ Ghi nhận Authorship Ledger"]
    Record --> Synthesizer["Bộ Phân Tích & Tổng Hợp Sở Thích<br/>(SynthesizePreferences)"]
    Synthesizer --> Directives["Đúc kết Chỉ Thị Ưu Tiên & Ràng Buộc Cấm Kỵ<br/>(ProseDirectives & Taboos)"]
    Directives --> VoiceFile["Cập nhật style/voice.md"]
    VoiceFile --> InjectContext["Tiêm tự động vào novel_context của Writer"]
    InjectContext --> NextChapters["Các chương tiếp theo tự động thích ứng chuẩn gu"]
```

---

## 8. Máy Trạng Thái Toàn Hệ Thống (State Machine Transitions)

Các giai đoạn sáng tác (`Phase`) và luồng vận hành (`Flow State`).

```mermaid
stateDiagram-v2
    [*] --> PhaseInit: Khởi tạo tác phẩm
    PhaseInit --> PhasePremise: Thiết lập đề tài
    PhasePremise --> PhaseOutline: Lập đại cương
    PhaseOutline --> PhaseWriting: Bắt đầu viết văn xuôi
    
    state PhaseWriting {
        [*] --> FlowWriting: Sáng tác chương mới
        FlowWriting --> FlowReviewing: Chờ thẩm định
        FlowReviewing --> FlowRewriting: Phát hiện lỗi -> Viết lại
        FlowReviewing --> FlowPolishing: Lỗi nhẹ -> Mài dũa câu từ
        FlowRewriting --> FlowWriting: Đã giải quyết hết hàng đợi
        FlowPolishing --> FlowWriting: Đã làm mịn xong
        FlowWriting --> FlowSteering: Người dùng can thiệp thời gian thực
        FlowSteering --> FlowWriting: Arbiter xử lý xong can thiệp
    }

    PhaseWriting --> PhaseComplete: Toàn bộ tác phẩm hoàn thành
    PhaseComplete --> [*]: Xuất bản & Đăng ký bản quyền
    PhaseComplete --> PhaseWriting: Reopen (Mở lại để đại tu/viết tiếp quyển mới)
```
