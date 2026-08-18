# Hướng Dẫn Quản Lý Ngữ Cảnh (Context Management) `ainovel-cli`

Tài liệu này giải thích hệ thống quản lý ngữ cảnh hiện tại của `ainovel-cli`, bao gồm:

- Tại sao cần quản lý ngữ cảnh
- Ngữ cảnh đến từ đâu
- Cách nén, khôi phục và chuyển giao ngữ cảnh khi vận hành
- Giá trị, điều kiện kích hoạt và kịch bản áp dụng của từng chiến lược
- Nơi cần kiểm tra đầu tiên khi gặp sự cố

---

## 1. Mục Tiêu Thiết Kế

Hệ thống quản lý ngữ cảnh của dự án này không dành cho trò chuyện (chat) thông thường, mà phục vụ riêng cho sáng tác tiểu thuyết. Nó giải thích và xử lý đồng thời nhiều bài toán:

1. Trò chuyện dài sẽ vượt quá cửa sổ ngữ cảnh (Context Window) của mô hình.
2. Sáng tác tiểu thuyết cần giữ lại không phải là "lịch sử trò chuyện thuần túy", mà là **ký ức tự sự có cấu trúc**.
3. Writer sau khi nén ngữ cảnh không được làm mất trạng thái nhân vật, phục bút, kế hoạch chương, quy tắc văn phong, và các điểm cần sửa sau thẩm định.
4. Khi khôi phục viết tiếp, không thể giả định mô hình vẫn "nhớ những gì đã nói trước đó", mà phải ưu tiên dựa vào các công cụ bền vững trên đĩa.

Do đó, chúng tôi áp dụng phương án **"Ký Ức Phân Tầng"**:

- **Ký ức ngắn hạn**: Phần đuôi tin nhắn gần nhất được giữ lại.
- **Ký ức trung hạn**: `ContextSummary` được tạo ra từ nén ngữ cảnh.
- **Ký ức dài hạn**: Các công cụ dữ liệu có cấu trúc trong Store của dự án.
- **Ký ức khôi phục**: Handoff / Restore Pack / `novel_context`.

---

## 2. Kiến Trúc Phân Tầng

### 2.1 Các Tầng Chính

1. `agentcore/context`: Quản lý ngân sách ngữ cảnh chung, đường ống chiến lược, khung nén/khôi phục.
2. `internal/tools/novel_context`: Đảm nhận việc lắp ráp dữ liệu có cấu trúc từ tiểu thuyết thành ngữ cảnh khả dụng cho lượt hiện tại.
3. Nén nhanh dựa trên Store dành riêng cho Writer.
4. `writer_restore`: Bổ sung gói khôi phục sau nén để đảm bảo Writer tiếp tục viết mượt mà.

### 2.2 Luồng Dữ Liệu Nén (Compression Pipeline)

Khi ngữ cảnh vượt quá ngưỡng token, hệ thống nén theo thứ tự chi phí từ thấp đến cao:

```text
ToolResultMicrocompact → LightTrim → StoreSummaryCompact → FullSummary
    Dọn kết quả công cụ cũ      Cắt ngắn văn bản dài      Nén từ Store (Zero LLM)    Dùng LLM tóm tắt bọc lót
```

---

## 3. Các Chiến Lược Nén Cốt Lõi

1. **ToolResultMicrocompact**: Thay thế các kết quả công cụ (Tool Results) cũ bằng văn bản chiếm chỗ ngắn gọn. Giảm tiếng ồn quy trình.
2. **LightTrim**: Cắt ngắn các khối văn bản quá dài (như chương bản thảo cũ), giữ lại đầu và đuôi.
3. **StoreSummaryCompact**: Dành riêng cho Writer. Thay vì gọi LLM tóm tắt lịch sử chat, lấy trực tiếp tóm tắt chương, tóm tắt Arc, hồ sơ nhân vật, phục bút từ Store đĩa để thay thế lịch sử cũ. **Chi phí LLM bằng 0**, không làm trôi dữ liệu.
4. **FullSummary**: Bước cuối cùng khi các tầng trên không đủ. Gọi LLM tóm tắt lại hội thoại với Prompt chuyên biệt cho tiểu thuyết (giữ lại trạng thái nhân vật, phục bút, phản hồi biên tập).

---

## 4. Tốc Độ Và Token CJK (Tiếng Trung / Tiếng Việt)

Hệ thống tự động phát hiện văn bản CJK / Phi-ASCII để ước tính token chính xác:
- Văn bản tiếng Trung / Tiếng Việt / CJK: `runes × 1.5`
- Văn bản ASCII: `bytes / 4`

Tránh việc ước tính thấp hơn thực tế gây trễ kích hoạt nén ngữ cảnh.

---

## 5. Tóm Lược

Khi gặp sự cố mất ký ức hoặc trôi logic của Writer:
1. Kiểm tra `novel_context` xem có bơm đủ `chapter_plan` và `previous_tail` không.
2. Kiểm tra `StoreSummaryCompact` có kích hoạt không hay bị rơi xuống `FullSummary`.
3. Kiểm tra các công cụ trên đĩa (`characters.json`, `foreshadow_ledger`, `state_changes.jsonl`).
