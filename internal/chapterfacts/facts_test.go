package chapterfacts

import (
	"testing"

	"github.com/voocel/ainovel-cli/internal/domain"
)

func validFacts() domain.ChapterFacts {
	return domain.ChapterFacts{
		Title:      "Chương 1: Khởi đầu",
		Summary:    "Nhân vật chính bước vào thế giới mới.",
		Characters: []string{"Nam chính", "Sư phụ"},
		KeyEvents:  []string{"Thức tỉnh sức mạnh"},
		TimelineEvents: []domain.TimelineEvent{
			{Time: "Sáng sớm", Event: "Rời khỏi thôn", Characters: []string{"Nam chính"}},
		},
		ForeshadowUpdates: []domain.ForeshadowUpdate{
			{ID: "fs-1", Action: "plant", Description: "Ngọc bội phát sáng bí ẩn"},
		},
		RelationshipChanges: []domain.RelationshipEntry{
			{CharacterA: "Nam chính", CharacterB: "Sư phụ", Relation: "Thầy trò"},
		},
		StateChanges: []domain.StateChange{
			{Entity: "Nam chính", Field: "Cảnh giới", OldValue: "Phàm nhân", NewValue: "Luyện Khí"},
		},
		CastIntros: []domain.CastIntro{
			{Name: "Sư phụ", BriefRole: "Người dẫn đường tu tiên"},
		},
	}
}

func TestProperties(t *testing.T) {
	propsNoFeedback := Properties(false)
	if len(propsNoFeedback) == 0 {
		t.Fatal("expected non-empty properties")
	}

	propsWithFeedback := Properties(true)
	if len(propsWithFeedback) != len(propsNoFeedback)+1 {
		t.Fatalf("expected feedback to add 1 property: without=%d with=%d", len(propsNoFeedback), len(propsWithFeedback))
	}
}

func TestValidate_Success(t *testing.T) {
	f := validFacts()
	if err := Validate(f); err != nil {
		t.Fatalf("expected valid facts to pass, got: %v", err)
	}
}

func TestValidate_RequiredCoreFields(t *testing.T) {
	f := validFacts()
	f.Title = ""
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty Title")
	}

	f = validFacts()
	f.Summary = "   "
	if err := Validate(f); err == nil {
		t.Fatal("expected error on whitespace Summary")
	}

	f = validFacts()
	f.KeyEvents = nil
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty KeyEvents")
	}

	f = validFacts()
	f.Characters = []string{"Alice", "  "}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on whitespace character item")
	}
}

func TestValidate_TimelineEvents(t *testing.T) {
	f := validFacts()
	f.TimelineEvents = []domain.TimelineEvent{
		{Time: "", Event: "Event 1"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty timeline time")
	}

	f.TimelineEvents = []domain.TimelineEvent{
		{Time: "Noon", Event: "Event 1", Characters: []string{""}},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty timeline character")
	}
}

func TestValidate_ForeshadowUpdates(t *testing.T) {
	f := validFacts()
	f.ForeshadowUpdates = []domain.ForeshadowUpdate{
		{ID: "", Action: "plant", Description: "desc"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty foreshadow ID")
	}

	f.ForeshadowUpdates = []domain.ForeshadowUpdate{
		{ID: "fs-1", Action: "plant", Description: ""},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on plant without description")
	}

	f.ForeshadowUpdates = []domain.ForeshadowUpdate{
		{ID: "fs-1", Action: "unknown_action"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on invalid foreshadow action")
	}

	f.ForeshadowUpdates = []domain.ForeshadowUpdate{
		{ID: "fs-1", Action: "advance"},
		{ID: "fs-2", Action: "resolve"},
	}
	if err := Validate(f); err != nil {
		t.Fatalf("expected advance/resolve without description to be valid, got %v", err)
	}
}

func TestValidate_RelationshipChanges(t *testing.T) {
	f := validFacts()
	f.RelationshipChanges = []domain.RelationshipEntry{
		{CharacterA: "Alice", CharacterB: "Alice", Relation: "Friend"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on relationship relating character to itself")
	}

	f.RelationshipChanges = []domain.RelationshipEntry{
		{CharacterA: "", CharacterB: "Bob", Relation: "Friend"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on missing CharacterA")
	}

	// Edge case: CharacterA and CharacterB are identical except for surrounding whitespace
	f.RelationshipChanges = []domain.RelationshipEntry{
		{CharacterA: "Alice", CharacterB: "Alice ", Relation: "Self"},
	}
	if err := Validate(f); err == nil {
		t.Errorf("ISSUE: Validate allows relating character to itself when whitespace differs: %q vs %q", "Alice", "Alice ")
	}
}

func TestValidate_StateChangesAndCast(t *testing.T) {
	f := validFacts()
	f.StateChanges = []domain.StateChange{
		{Entity: "Hero", Field: "", NewValue: "Val"},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty state change field")
	}

	f = validFacts()
	f.CastIntros = []domain.CastIntro{
		{Name: "Sidekick", BriefRole: "  "},
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty cast brief_role")
	}
}

func TestValidate_HooksAndFeedback(t *testing.T) {
	f := validFacts()
	f.HookType = "invalid_hook_type_xyz"
	if err := Validate(f); err == nil {
		t.Fatal("expected error on invalid hook_type")
	}

	f = validFacts()
	f.DominantStrand = "invalid_strand_xyz"
	if err := Validate(f); err == nil {
		t.Fatal("expected error on invalid dominant_strand")
	}

	f = validFacts()
	f.Feedback = &domain.OutlineFeedback{
		Deviation:  "",
		Suggestion: "Do something",
	}
	if err := Validate(f); err == nil {
		t.Fatal("expected error on empty feedback deviation")
	}
}
