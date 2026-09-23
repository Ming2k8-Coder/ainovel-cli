package utils

import "strings"

// JSONFieldExtractor trích xuất giá trị chuỗi của một trường được chỉ định từ các mảnh JSON dạng stream.
//
// Khi LLM sinh tool call ở chế độ streaming, các tham số có thể đến từng đoạn nhỏ (OpenAI/Anthropic)
// hoặc đến cùng một lúc (Gemini). Trình trích xuất này dùng máy trạng thái (state machine) quét từng ký tự,
// khi phát hiện trường (key) mục tiêu sẽ trích xuất chuỗi giá trị và xử lý thoát chuỗi (escape) JSON.
type JSONFieldExtractor struct {
	key      string // Mục tiêu khớp, ví dụ `"content"` hoặc `"task"`
	state    extractState
	matchPos int
	escape   bool
	buf      strings.Builder
}

type extractState int

const (
	stateScan    extractState = iota // Đang quét, tìm kiếm key mục tiêu
	stateColon                       // Đã khớp key, chờ dấu hai chấm và dấu ngoặc kép mở
	stateExtract                     // Đang trích xuất nội dung chuỗi giá trị
)

func NewFieldExtractor(fieldName string) *JSONFieldExtractor {
	return &JSONFieldExtractor{key: `"` + fieldName + `"`}
}

// Feed xử lý một đoạn delta, trả về chuỗi văn bản đã trích xuất được (có thể rỗng).
func (e *JSONFieldExtractor) Feed(delta string) string {
	e.buf.Reset()
	for _, r := range delta {
		switch e.state {
		case stateScan:
			e.feedScan(r)
		case stateColon:
			e.feedColon(r)
		case stateExtract:
			e.feedExtract(r)
		}
	}
	return e.buf.String()
}

func (e *JSONFieldExtractor) feedScan(r rune) {
	if e.matchPos < len(e.key) && byte(r) == e.key[e.matchPos] {
		e.matchPos++
		if e.matchPos == len(e.key) {
			e.state = stateColon
			e.matchPos = 0
		}
		return
	}
	e.matchPos = 0
	if byte(r) == e.key[0] {
		e.matchPos = 1
	}
}

func (e *JSONFieldExtractor) feedColon(r rune) {
	switch r {
	case ':', ' ', '\t':
		// Bỏ qua khoảng trắng và dấu hai chấm
	case '"':
		e.state = stateExtract
		e.escape = false
	default:
		e.state = stateScan
		e.matchPos = 0
		if byte(r) == e.key[0] {
			e.matchPos = 1
		}
	}
}

func (e *JSONFieldExtractor) feedExtract(r rune) {
	if e.escape {
		e.escape = false
		switch r {
		case 'n':
			e.buf.WriteByte('\n')
		case 't':
			e.buf.WriteByte('\t')
		case 'r':
			e.buf.WriteByte('\r')
		case '"', '\\', '/':
			e.buf.WriteRune(r)
		default:
			e.buf.WriteByte('\\')
			e.buf.WriteRune(r)
		}
		return
	}
	switch r {
	case '\\':
		e.escape = true
	case '"':
		e.state = stateScan
		e.matchPos = 0
	default:
		e.buf.WriteRune(r)
	}
}

// Reset đặt lại trạng thái (được gọi khi bắt đầu lượt tin nhắn LLM mới).
func (e *JSONFieldExtractor) Reset() {
	e.state = stateScan
	e.matchPos = 0
	e.escape = false
}

// StreamFilter phân biệt giữa câu trả lời văn bản thông thường của SubAgent và lệnh gọi công cụ (tool call) dạng JSON.
// Câu trả lời dạng văn bản được đánh dấu là nội dung suy nghĩ (thêm tiền tố ThinkingSep); còn JSON tool call chỉ trích xuất trường được chỉ định.
//
// Cơ chế nhận diện: khi gặp '{' sẽ chuyển sang chế độ JSON (theo dõi độ sâu ngoặc nhọn),
// khi độ sâu về 0 thì quay lại chế độ văn bản.
type StreamFilter struct {
	fieldExt   *JSONFieldExtractor
	mode       filterMode
	braceDepth int
	inString   bool // Đang ở trong chuỗi JSON (ngoặc nhọn không được tính)
	escJSON    bool // Thoát chuỗi bên trong chuỗi JSON
	thinking   bool // Đang ở trong đoạn văn bản suy nghĩ
	buf        strings.Builder
}

type filterMode int

const (
	filterText filterMode = iota // Trả về văn bản trực tiếp
	filterJSON                   // Gọi công cụ JSON, trích xuất trường mục tiêu
)

func NewStreamFilter(fieldName string) *StreamFilter {
	return &StreamFilter{fieldExt: NewFieldExtractor(fieldName)}
}

// Feed xử lý một đoạn delta, trả về văn bản có thể hiển thị.
// Văn bản thông thường được xuất trực tiếp; giá trị trường mục tiêu trong JSON được trích xuất; các cấu trúc JSON khác bị loại bỏ.
func (f *StreamFilter) Feed(delta string) string {
	f.buf.Reset()
	for _, r := range delta {
		switch f.mode {
		case filterText:
			if r == '{' {
				f.thinking = false
				f.mode = filterJSON
				f.braceDepth = 1
				f.inString = false
				f.escJSON = false
				f.fieldExt.Reset()
				f.feedExtractor(r)
			} else {
				if !f.thinking {
					f.thinking = true
					f.buf.WriteString(ThinkingSep)
				}
				f.buf.WriteRune(r)
			}
		case filterJSON:
			f.feedExtractor(r)
			f.trackBraces(r)
		}
	}
	return f.buf.String()
}

// feedExtractor chuyển từng ký tự cho fieldExt, kết quả trích xuất được ghi vào buf.
func (f *StreamFilter) feedExtractor(r rune) {
	if text := f.fieldExt.Feed(string(r)); text != "" {
		f.buf.WriteString(text)
	}
}

// trackBraces theo dõi độ sâu ngoặc nhọn JSON, khi độ sâu về 0 sẽ chuyển lại về chế độ văn bản.
func (f *StreamFilter) trackBraces(r rune) {
	if f.escJSON {
		f.escJSON = false
		return
	}
	if f.inString {
		switch r {
		case '\\':
			f.escJSON = true
		case '"':
			f.inString = false
		}
		return
	}
	switch r {
	case '"':
		f.inString = true
	case '{':
		f.braceDepth++
	case '}':
		f.braceDepth--
		if f.braceDepth <= 0 {
			f.mode = filterText
		}
	}
}

// Reset đặt lại trạng thái.
func (f *StreamFilter) Reset() {
	f.mode = filterText
	f.braceDepth = 0
	f.inString = false
	f.escJSON = false
	f.thinking = false
	f.fieldExt.Reset()
}
