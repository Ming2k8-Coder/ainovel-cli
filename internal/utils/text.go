package utils

// ThinkingSep là ký hiệu phân tách giữa phần suy nghĩ (thinking) và phần nội dung chính.
// observer chèn ký hiệu này trước đoạn suy nghĩ, TUI dựa vào đây để chuyển đổi phong cách hiển thị (render style).
const ThinkingSep = "\x02"

// TruncateRunes cắt ngắn chuỗi theo số rune và thêm dấu ba chấm '…'; nếu độ dài <= n thì giữ nguyên chuỗi.
func TruncateRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}
