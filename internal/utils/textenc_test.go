package utils

import (
	"bytes"
	"testing"
)

func TestDecodeText_ValidUTF8(t *testing.T) {
	input := []byte("Hello, thế giới! Tiếng Việt có dấu.")
	got := DecodeText(input)
	if got != string(input) {
		t.Fatalf("expected %q, got %q", string(input), got)
	}
}

func TestDecodeText_StripBOM(t *testing.T) {
	input := append([]byte("\uFEFF"), []byte("Content without BOM")...)
	got := DecodeText(input)
	if got != "Content without BOM" {
		t.Fatalf("expected BOM to be stripped, got %q", got)
	}
}

func TestDecodeText_GBK(t *testing.T) {
	// GB18030 / GBK bytes for "中文小说"
	// 中: 0xD6 0xD0, 文: 0xCE 0xC4, 小: 0xD0 0xA1, 说: 0xCB 0xB5
	gbkBytes := []byte{0xD6, 0xD0, 0xCE, 0xC4, 0xD0, 0xA1, 0xCB, 0xB5}
	got := DecodeText(gbkBytes)
	want := "中文小说"
	if got != want {
		t.Fatalf("GBK decode mismatch: want %q, got %q", want, got)
	}
}

func TestDecodeText_Empty(t *testing.T) {
	if got := DecodeText([]byte{}); got != "" {
		t.Fatalf("expected empty string for empty input, got %q", got)
	}
}

func TestDecodeText_PreservesNewlines(t *testing.T) {
	input := []byte("Line 1\r\nLine 2\nLine 3\r")
	got := DecodeText(input)
	if !bytes.Equal([]byte(got), input) {
		t.Fatalf("expected preserved newlines, got %q", got)
	}
}
