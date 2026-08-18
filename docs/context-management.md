# Hệ Thống Quản Lý Ngữ Cảnh & Ký Ức (Context Management)

Tài liệu này trình bày chi tiết về kiến trúc quản lý ngữ cảnh và phân bổ ký ức trong `ainovel-cli`, bao gồm:

- Lý do cần một hệ thống quản lý ngữ cảnh chuyên biệt cho tiểu thuyết dài tập.
- Nguồn gốc và cấu trúc của dữ liệu ngữ cảnh.
- Các tầng nén, phục hồi và chuyển giao ký ức giữa các phiên làm việc.
- Cơ chế ước lượng và tối ưu hóa ngân sách Token cho văn bản tiếng Việt và chữ tượng hình (CJK).

---

## 1. Triết Lý Thiết Kế: Ký Ức Tự Sự Có Cấu Trúc

Khác với các ứng dụng chatbot hội thoại thông thường (chỉ cần lưu chuỗi tin nhắn hỏi-đáp nối tiếp nhau), việc sáng tác tiểu thuyết dài tập từ vài chục đến hàng trăm chương đòi hỏi một phương pháp quản lý hoàn toàn khác biệt:

1. Cửa sổ ngữ cảnh (Context Window) của các mô hình LLM luôn có giới hạn và chi phí tăng nhanh theo dung lượng.
2. Tiểu thuyết không thể duy trì chất lượng bằng "lịch sử chat thô", mà bắt buộc phải dựa vào **ký ức tự sự có cấu trúc**.
3. Khi nén ngữ cảnh hoặc chuyển chương, tác tử Writer tuyệt đối không được đánh mất các thông tin sống còn: trạng thái nhân vật, mạng lưới phục bút, mục tiêu của Hồi, các quy tắc văn phong, và các chỉ thị sửa đổi từ khâu biên tập.
4. Khi tiếp tục viết sau khi khởi động lại, hệ thống không thể kỳ vọng LLM "tự nhớ những gì đã nói", mà phải xây dựng lại ngữ cảnh tươi mới từ nguồn dữ liệu thực tế được lưu trên đĩa.

Do đó, `ainovel-cli` triển khai mô hình **"Ký Ức Phân Tầng" (Layered Memory Architecture)**:

```text
┌─────────────────────────────────────────────────────────────┐
│ 1. Ký ức Ngắn hạn (Short-term Memory)                       │
│    Các tin nhắn trao đổi gần nhất trong phiên làm việc.     │
├─────────────────────────────────────────────────────────────┤
│ 2. Ký ức Trung hạn (Medium-term Memory)                     │
│    Tóm tắt diễn biến chương gần nhất, tóm tắt Hồi/Quyển.    │
├─────────────────────────────────────────────────────────────┤
│ 3. Ký ức Dài hạn (Long-term Memory)                         │
│    Hồ sơ nhân vật, quan hệ, phục bút, thế giới trong Store. │
├─────────────────────────────────────────────────────────────┤
│ 4. Gói Khôi Phục Động (Dynamic Restore Pack)                │
│    Tái tạo ngữ cảnh chuẩn bị viết từ novel_context.         │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Đường Ống Nén Ngữ Cảnh Chi Phí Thấp (Compression Pipeline)

Khi dung lượng tin nhắn trong phiên của Writer vượt quá ngưỡng an toàn của cửa sổ ngữ cảnh, hệ thống sẽ kích hoạt đường ống nén tự động theo thứ tự ưu tiên từ chi phí thấp đến cao:

```text
[1. ToolResultMicrocompact] ➔ [2. LightTrim] ➔ [3. StoreSummaryCompact] ➔ [4. FullSummary]
   Rút gọn kết quả công cụ       Cắt ngắn văn bản cũ     Nén trực tiếp từ đĩa       Tóm tắt bằng LLM
   (Chi phí: 0 token)            (Chi phí: 0 token)      (Không tốn LLM / Zero-LLM) (Tầng bảo hiểm cuối)
```

1. **ToolResultMicrocompact**: Thay thế toàn bộ kết quả trả về của các công cụ gọi trước đó bằng các văn bản giữ chỗ ngắn gọn, loại bỏ rác dữ liệu không còn giá trị điều hướng.
2. **LightTrim**: Thu gọn các đoạn văn bản quá dài từ các vòng lặp trước, chỉ giữ lại phần đầu và phần đuôi để tiết kiệm không gian.
3. **StoreSummaryCompact (Chiến lược độc quyền cho Writer)**: Thay vì phải gọi một LLM đắt đỏ để tóm tắt lại toàn bộ lịch sử trò chuyện (dễ làm trôi các chi tiết logic quan trọng), hệ thống đọc trực tiếp tóm tắt chương, tóm tắt Hồi, hồ sơ nhân vật và sổ cái phục bút từ Store trên đĩa để thay thế lịch sử cũ. **Phương pháp này tiêu tốn 0 token LLM và bảo toàn 100% độ chính xác của cốt truyện.**
4. **FullSummary**: Tầng bảo hiểm cuối cùng khi các bước trên không đủ hạ thấp dung lượng. Hệ thống gọi LLM tóm tắt lại toàn bộ diễn biến dựa trên Prompt chuyên biệt cho tiểu thuyết.

---

## 3. Cơ Chế Định Lượng Token Chuẩn Xác Cho Tiếng Việt & CJK

Các ngôn ngữ có dấu thanh như tiếng Việt và chữ tượng hình (CJK) thường chiếm số lượng token cao hơn đáng kể so với ký tự ASCII tiếng Anh thông thường.

Nhằm tránh tình trạng ước tính thấp hơn thực tế dẫn đến việc tràn ngữ cảnh ngoài ý muốn, hệ thống tích hợp bộ đo lường thông minh:
- Văn bản chứa ký tự tiếng Việt / CJK / Phi-ASCII: Tính theo hệ số `số_ký_tự × 1.5`
- Văn bản thuần ASCII: Tính theo hệ số `số_byte / 4`

Cơ chế này đảm bảo ngưỡng kích hoạt nén luôn diễn ra kịp thời và chính xác trước khi mô hình chạm trần cửa sổ ngữ cảnh.
