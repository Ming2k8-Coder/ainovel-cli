package abselect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voocel/ainovel-cli/internal/authorship"
)

func TestABSelectFeedbackLoop(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "abselect_test_*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ledger := authorship.NewLedger(tempDir)
	mgr := NewManager(tempDir, ledger)

	// 1. Ghi nhận lựa chọn A/B thứ nhất (Chọn A - đối thoại sắc bén, bỏ B - thuyết minh dài dòng)
	choice1 := Choice{
		ID:      "ab-test-1",
		Chapter: 1,
		Scope:   "dialogue_tone",
		OptionA: Variant{
			Key:         "A",
			Title:       "Đoạn đối thoại gay cấn trong phòng thí nghiệm",
			Description: "Nhân vật tranh luận dồn dập về bản thiết kế",
			Traits:      []string{"dialogue_sharp", "fast_paced"},
		},
		OptionB: Variant{
			Key:         "B",
			Title:       "Đoạn độc thoại giải thích bối cảnh lịch sử",
			Description: "Thuyết minh chi tiết về lịch sử trường Bách Khoa",
			Traits:      []string{"verbose_exposition", "slow_pacing"},
		},
		Selected:   "A",
		UserReason: "Thích nhịp đối thoại nhanh và sắc bén hơn.",
	}
	if err := mgr.RecordChoice(choice1); err != nil {
		t.Fatalf("RecordChoice 1 failed: %v", err)
	}

	// 2. Ghi nhận lựa chọn A/B thứ hai (Chọn A - chi tiết kỹ thuật khoa học, bỏ B - sáo rỗng)
	choice2 := Choice{
		ID:      "ab-test-2",
		Chapter: 2,
		Scope:   "chapter_draft",
		OptionA: Variant{
			Key:         "A",
			Title:       "Bản thảo nhấn mạnh chi tiết cơ khí Cyberpunk",
			Description: "Mô tả vi mạch bán dẫn và động cơ phản lực",
			Traits:      []string{"hard_scifi_detail", "dialogue_sharp"},
		},
		OptionB: Variant{
			Key:         "B",
			Title:       "Bản thảo tình cảm học đường thông thường",
			Description: "Tập trung vào cảm xúc lãng mạn sáo rỗng",
			Traits:      []string{"melodrama_cliche"},
		},
		Selected:   "A",
		UserReason: "Đúng chất Hard Sci-Fi kỹ thuật Sanbaka hơn.",
	}
	if err := mgr.RecordChoice(choice2); err != nil {
		t.Fatalf("RecordChoice 2 failed: %v", err)
	}

	// 3. Kiểm tra đồng bộ sang Authorship Ledger
	records, err := ledger.LoadRecords()
	if err != nil {
		t.Fatalf("LoadRecords failed: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 authorship records from AB choices, got %d", len(records))
	}
	if records[0].Type != authorship.TypeCreativeChoice {
		t.Fatalf("expected TypeCreativeChoice, got %s", records[0].Type)
	}

	// 4. Tổng hợp sở thích Prompt & Voice
	synth, err := mgr.SynthesizePreferences()
	if err != nil {
		t.Fatalf("SynthesizePreferences failed: %v", err)
	}

	if len(synth.PreferredTraits) == 0 || len(synth.AvoidedTraits) == 0 {
		t.Fatalf("expected preferred and avoided traits, got %+v", synth)
	}

	// 5. Tự động áp dụng vào style/voice.md
	if err := mgr.ApplyFeedbackToProject(); err != nil {
		t.Fatalf("ApplyFeedbackToProject failed: %v", err)
	}

	voiceContent, err := os.ReadFile(filepath.Join(tempDir, "style", "voice.md"))
	if err != nil {
		t.Fatalf("read voice.md failed: %v", err)
	}
	voiceStr := string(voiceContent)
	if !strings.Contains(voiceStr, "Định Hướng Hành Văn Ưu Tiên") || !strings.Contains(voiceStr, "Các Yếu Tố Cần Tránh") {
		t.Fatalf("voice.md missing expected sections: %s", voiceStr)
	}
}
