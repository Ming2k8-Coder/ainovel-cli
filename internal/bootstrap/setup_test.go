package bootstrap

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMaskKey(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", "****"},
		{"1234", "****"},
		{"12345678", "****"},
		{"123456789", "1234****6789"},
		{"sk-ant-api03-abcdef123456", "sk-a****3456"},
	}
	for _, c := range cases {
		got := maskKey(c.input)
		if got != c.want {
			t.Errorf("maskKey(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestSetupSelectModel_NavigationAndCancel(t *testing.T) {
	m := setupSelectModel{
		title: "Test Provider",
		items: []setupProvider{
			{name: "p1", label: "Provider 1"},
			{name: "p2", label: "Provider 2"},
			{name: "p3", label: "Provider 3"},
		},
	}

	if m.Init() != nil {
		t.Fatal("expected Init to return nil")
	}

	// Move down
	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm1 := m1.(setupSelectModel)
	if sm1.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", sm1.cursor)
	}

	// Move down to bottom
	m2, _ := sm1.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm2 := m2.(setupSelectModel)
	if sm2.cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", sm2.cursor)
	}

	// Move down beyond boundary (should clamp)
	m3, _ := sm2.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm3 := m3.(setupSelectModel)
	if sm3.cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", sm3.cursor)
	}

	// Move up
	m4, _ := sm3.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm4 := m4.(setupSelectModel)
	if sm4.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", sm4.cursor)
	}

	// Cancel
	m5, cmd := sm4.Update(tea.KeyMsg{Type: tea.KeyEsc})
	sm5 := m5.(setupSelectModel)
	if !sm5.cancelled || cmd == nil {
		t.Fatal("expected cancelled = true on Esc")
	}

	// View rendering
	v := sm4.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view")
	}
}

func TestSetupInputModel_TypingAndDefault(t *testing.T) {
	m := setupInputModel{
		label:        "API Key",
		placeholder:  "Enter key",
		defaultValue: "default-val",
	}

	// Type letters
	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc")})
	im1 := m1.(setupInputModel)
	if im1.value != "abc" {
		t.Fatalf("expected value 'abc', got %q", im1.value)
	}

	// Backspace
	m2, _ := im1.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	im2 := m2.(setupInputModel)
	if im2.value != "ab" {
		t.Fatalf("expected value 'ab', got %q", im2.value)
	}

	// View with value
	vWithValue := im2.View()
	if len(vWithValue) == 0 {
		t.Fatal("expected non-empty view with value")
	}

	// Cancel
	m3, cmd := im2.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	im3 := m3.(setupInputModel)
	if !im3.cancelled || cmd == nil {
		t.Fatal("expected cancelled = true on Ctrl+C")
	}
}
