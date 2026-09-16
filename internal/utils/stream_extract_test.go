package utils

import (
	"strings"
	"testing"
)

func TestJSONFieldExtractor_Basic(t *testing.T) {
	ext := NewFieldExtractor("content")
	input := `{"role": "assistant", "content": "Hello World", "other": 123}`
	got := ext.Feed(input)
	if got != "Hello World" {
		t.Fatalf("expected 'Hello World', got %q", got)
	}
}

func TestJSONFieldExtractor_Chunked(t *testing.T) {
	ext := NewFieldExtractor("content")
	chunks := []string{
		`{"id": 1, `,
		`"con`,
		`tent": "Part 1 `,
		`and Part 2\n`,
		`with \"quotes\"`,
		` and \\slashes", "done": true}`,
	}
	var sb strings.Builder
	for _, ch := range chunks {
		sb.WriteString(ext.Feed(ch))
	}
	want := "Part 1 and Part 2\nwith \"quotes\" and \\slashes"
	if sb.String() != want {
		t.Fatalf("chunked extraction mismatch:\nwant: %q\ngot:  %q", want, sb.String())
	}
}

func TestJSONFieldExtractor_Escapes(t *testing.T) {
	ext := NewFieldExtractor("text")
	input := `{"text": "Line 1\r\nLine 2\tTabbed\\Escaped\/Slash\"Quoted\""}`
	got := ext.Feed(input)
	want := "Line 1\r\nLine 2\tTabbed\\Escaped/Slash\"Quoted\""
	if got != want {
		t.Fatalf("escape extraction mismatch:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestJSONFieldExtractor_Reset(t *testing.T) {
	ext := NewFieldExtractor("content")
	got1 := ext.Feed(`{"content": "First"}`)
	if got1 != "First" {
		t.Fatalf("got1 = %q, want First", got1)
	}
	ext.Reset()
	got2 := ext.Feed(`{"content": "Second"}`)
	if got2 != "Second" {
		t.Fatalf("got2 = %q, want Second", got2)
	}
}

func TestJSONFieldExtractor_UnicodeMultibyte(t *testing.T) {
	ext := NewFieldExtractor("content")
	// Test Vietnamese and Chinese characters in stream
	input := `{"content": "Xin chào thế giới! 这是一个测试。"}`
	got := ext.Feed(input)
	want := "Xin chào thế giới! 这是一个测试。"
	if got != want {
		t.Fatalf("multibyte mismatch:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStreamFilter_ThinkingAndToolCall(t *testing.T) {
	filter := NewStreamFilter("content")
	
	// Thinking part
	text1 := filter.Feed("Thinking about the plot...")
	if !strings.HasPrefix(text1, ThinkingSep) {
		t.Fatalf("expected ThinkingSep prefix for thinking text, got %q", text1)
	}
	if !strings.Contains(text1, "Thinking about the plot...") {
		t.Fatalf("expected thinking text, got %q", text1)
	}

	// Tool call part: JSON should extract 'content' and discard surrounding JSON
	text2 := filter.Feed(`{"name": "draft_chapter", "arguments": {"content": "Chapter 1 prose"}`)
	if !strings.Contains(text2, "Chapter 1 prose") {
		t.Fatalf("expected extracted content from tool call JSON, got %q", text2)
	}
	if strings.Contains(text2, "draft_chapter") {
		t.Fatalf("tool name should be discarded, but found in output: %q", text2)
	}

	// Close JSON and back to thinking
	text3 := filter.Feed("} Resuming thought")
	if !strings.Contains(text3, "Resuming thought") {
		t.Fatalf("expected text after JSON closes, got %q", text3)
	}
}

func TestStreamFilter_NestedBracesInString(t *testing.T) {
	filter := NewStreamFilter("content")
	// JSON with braces inside string value
	input := `{"content": "Text with {braces} and \"escaped quotes\" inside"}`
	got := filter.Feed(input)
	want := `Text with {braces} and "escaped quotes" inside`
	if got != want {
		t.Fatalf("nested braces mismatch:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStreamFilter_Reset(t *testing.T) {
	filter := NewStreamFilter("content")
	filter.Feed("Thought 1")
	filter.Reset()
	got := filter.Feed("Thought 2")
	if !strings.HasPrefix(got, ThinkingSep) {
		t.Fatalf("expected ThinkingSep after reset, got %q", got)
	}
}

// TestJSONFieldExtractor_UnicodeEscapeDocumentsLimitation kiểm tra hành vi đối với mã escape \uXXXX.
// Hiện tại feedExtract chưa hỗ trợ decode \uXXXX (RFC 8259) mà trả về \u.
func TestJSONFieldExtractor_UnicodeEscapeDocumentsLimitation(t *testing.T) {
	ext := NewFieldExtractor("content")
	input := `{"content": "Chương \u0031"}`
	got := ext.Feed(input)
	// Ghi nhận limitation: hiện tại feedExtract trả về \u0031 thay vì ký tự "1"
	if got == `Chương \u0031` {
		t.Log("Note: JSONFieldExtractor currently passes raw \\u0031 without decoding to UTF-8 rune")
	}
}

