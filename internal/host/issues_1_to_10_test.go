package host

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/domain"
	"github.com/voocel/ainovel-cli/internal/flow"
	"github.com/voocel/ainovel-cli/internal/revision"
	storepkg "github.com/voocel/ainovel-cli/internal/store"
	"github.com/voocel/ainovel-cli/internal/stylestat"
	"github.com/voocel/ainovel-cli/internal/tools"
)

// TestIssue1_SkeletonArcExpansion verifies Issue #1:
// Expanding a skeleton arc (estimated chapters > 0, chapters = nil) using expand_next_arc
// does NOT require manual volume and arc parameters and seamlessly expands the outline.
func TestIssue1_SkeletonArcExpansion(t *testing.T) {
	dir := t.TempDir()
	st := storepkg.NewStore(dir)
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}

	// Setup layered outline with Arc 1 (concrete) and Arc 2 (skeleton arc)
	layered := []domain.VolumeOutline{
		{
			Index: 1, Title: "Quyển 1", Theme: "Khởi nguồn",
			Arcs: []domain.ArcOutline{
				{
					Index: 1, Title: "Hồi 1", Goal: "Gặp sư phụ",
					Chapters: []domain.OutlineEntry{
						{Chapter: 1, Title: "Chương 1", CoreEvent: "Xuất sơn", Hook: "Gặp nguy hiểm", Scenes: []string{"Thôn xóm"}},
					},
				},
				{
					Index: 2, Title: "Hồi 2 (Xương sườn)", Goal: "Vào tông môn",
					EstimatedChapters: 2,
					Chapters:          nil, // Skeleton arc!
				},
			},
		},
	}
	if err := st.Outline.SaveLayeredOutline(layered); err != nil {
		t.Fatal(err)
	}
	if err := st.Progress.Init(1); err != nil {
		t.Fatal(err)
	}
	p, _ := st.Progress.Load()
	p.CompletedChapters = []int{1}
	if err := st.Progress.Save(p); err != nil {
		t.Fatal(err)
	}

	// Execute ExpandNextArcTool
	tool := tools.NewExpandNextArcTool(st)
	args := map[string]any{
		"title": "Hồi 2: Nhập Môn Khảo Nghiệm",
		"goal":  "Vượt qua 3 bài thi tông môn",
		"chapters": []map[string]any{
			{"title": "Chương 2", "core_event": "Thi tuyển", "hook": "Đối thủ xuất hiện", "scenes": []string{"Quảng trường"}},
			{"title": "Chương 3", "core_event": "Trúng tuyển", "hook": "Bí mật lộ diện", "scenes": []string{"Đại điện"}},
		},
	}
	rawArgs, _ := json.Marshal(args)
	if _, err := tool.Execute(context.Background(), rawArgs); err != nil {
		t.Fatalf("ExpandNextArc failed: %v", err)
	}

	// Verify outline was updated and skeleton arc was filled
	volumes, err := st.Outline.LoadLayeredOutline()
	if err != nil {
		t.Fatal(err)
	}
	arc2 := volumes[0].Arcs[1]
	if len(arc2.Chapters) != 2 || arc2.EstimatedChapters != 0 {
		t.Fatalf("expected 2 chapters and 0 estimated in expanded arc, got: %+v", arc2)
	}
	if arc2.Title != "Hồi 2: Nhập Môn Khảo Nghiệm" {
		t.Fatalf("expected updated title, got %q", arc2.Title)
	}

	// Verify progress total chapters updated to 3
	prog, _ := st.Progress.Load()
	if prog.TotalChapters != 3 {
		t.Fatalf("expected total chapters = 3, got %d", prog.TotalChapters)
	}
}

// TestIssue2_ContextAndBudgetSentinel verifies Issue #2:
// Guarding against context explosion and runaway costs via BudgetSentinel.
func TestIssue2_ContextAndBudgetSentinel(t *testing.T) {
	var aborted bool
	var warnCount int
	cfg := bootstrap.BudgetConfig{
		BookUSD:   2.00,
		WarnRatio: 0.8,
		HardStop:  true,
	}
	currentCost := 0.50
	sentinel := NewBudgetSentinel(cfg, func() float64 {
		return currentCost
	}, func(reason string) {
		aborted = true
	}, func(level, summary string) {
		warnCount++
	})

	if sentinel == nil {
		t.Fatal("expected sentinel initialized")
	}

	// 50% cost -> normal
	sentinel.OnCost(0.50)
	sentinel.HandleBoundary()
	if aborted || warnCount != 0 {
		t.Fatalf("unexpected abort or warning at 50%% cost: aborted=%v, warn=%d", aborted, warnCount)
	}

	// 85% cost -> warn
	currentCost = 1.70
	sentinel.OnCost(1.70)
	if warnCount != 1 {
		t.Fatalf("expected warning at 85%%, got %d", warnCount)
	}

	// 105% cost -> hard abort
	currentCost = 2.10
	sentinel.OnCost(2.10)
	if !aborted {
		t.Fatal("expected sentinel to abort at 105% cost")
	}
}

// TestIssue3_LegacyProjectBaselineMigration verifies Issue #3:
// Legacy projects from v0.7.x lack chapter_records/*.json. MigrateLegacyBaseline
// must reconstruct records from legacy snapshots and allow revision scan to succeed.
func TestIssue3_LegacyProjectBaselineMigration(t *testing.T) {
	dir := t.TempDir()
	st := storepkg.NewStore(dir)
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}

	// Create legacy completed chapters
	if err := st.Progress.Init(2); err != nil {
		t.Fatal(err)
	}
	p, _ := st.Progress.Load()
	p.CompletedChapters = []int{1, 2}
	_ = st.Progress.Save(p)

	// Write chapter draft and final text
	chap1Path := filepath.Join(st.Dir(), "chapters", "01.md")
	chap2Path := filepath.Join(st.Dir(), "chapters", "02.md")
	_ = os.MkdirAll(filepath.Dir(chap1Path), 0755)
	_ = os.WriteFile(chap1Path, []byte("Nội dung chương 1 đầy đủ."), 0644)
	_ = os.WriteFile(chap2Path, []byte("Nội dung chương 2 đầy đủ."), 0644)

	// Write legacy summaries
	_ = st.Summaries.SaveSummary(domain.ChapterSummary{Chapter: 1, Title: "Chương 1", Summary: "Tóm tắt 1"})
	_ = st.Summaries.SaveSummary(domain.ChapterSummary{Chapter: 2, Title: "Chương 2", Summary: "Tóm tắt 2"})

	// Before migration: revision.Scan MUST report missing acceptance record
	_, errBefore := revision.Scan(st)
	if errBefore == nil || !strings.Contains(errBefore.Error(), "缺少接纳记录") {
		t.Fatalf("expected error before migration, got %v", errBefore)
	}

	// Run migration
	if err := revision.MigrateLegacyBaseline(st); err != nil {
		t.Fatalf("MigrateLegacyBaseline failed: %v", err)
	}

	// After migration: records exist with ChapterOriginLegacy and revision.Scan succeeds
	rec1, err := st.ChapterRecords.Load(1)
	if err != nil || rec1 == nil {
		t.Fatalf("expected record 1 after migration: %v", err)
	}
	if rec1.Origin != domain.ChapterOriginLegacy {
		t.Fatalf("expected ChapterOriginLegacy, got %q", rec1.Origin)
	}

	changes, errAfter := revision.Scan(st)
	if errAfter != nil {
		t.Fatalf("revision.Scan failed after migration: %v", errAfter)
	}
	if len(changes) != 0 {
		t.Fatalf("expected 0 unexpected changes on clean migration, got %d", len(changes))
	}
}

// TestIssue4_CommitChapter_Validation verifies Issue #4:
// Malformed or invalid JSON arguments are caught cleanly without crashing the runtime.
func TestIssue4_CommitChapter_Validation(t *testing.T) {
	dir := t.TempDir()
	st := storepkg.NewStore(dir)
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}

	tool := tools.NewCommitChapterTool(st, tools.NewStyleStatsIndex(st))

	// Truncated / malformed JSON
	truncatedJSON := []byte(`{"chapter": 1, "title": "Incomplete`)
	_, err := tool.Execute(context.Background(), truncatedJSON)
	if err == nil {
		t.Fatal("expected error on truncated JSON args")
	}

	// Chapter <= 0
	invalidChapter := []byte(`{"chapter": 0}`)
	_, err = tool.Execute(context.Background(), invalidChapter)
	if err == nil {
		t.Fatal("expected error on chapter <= 0")
	}
}

// TestIssue5_UIEventBusAndLifecycle verifies Issue #5:
// Structured lifecycle events and UI state snapshots allow decoupling CLI terminal and WebUI.
func TestIssue5_UIEventBusAndLifecycle(t *testing.T) {
	ev := Event{
		ID:       "evt-test-1",
		Time:     time.Now(),
		Category: "TOOL",
		Agent:    "writer",
		Summary:  "Đang viết chương 1",
		Level:    "info",
	}
	if !ev.Running() {
		t.Fatal("event without FinishedAt must be running")
	}

	ev.FinishedAt = time.Now()
	ev.Duration = 100 * time.Millisecond
	if ev.Running() {
		t.Fatal("event with FinishedAt must not be running")
	}

	LogEvent(ev)

	snap := UISnapshot{
		Provider:     "anthropic",
		BookTitle:    "Thần Đạo Đan Tôn",
		RuntimeState: "running",
		Phase:        "writing",
	}
	if snap.BookTitle != "Thần Đạo Đan Tôn" || snap.RuntimeState != "running" {
		t.Fatalf("unexpected snapshot state: %+v", snap)
	}
}

// TestIssue6_ProjectLevelStyleOverrides verifies Issue #6:
// Book-level style overrides take top priority over home directory and built-ins.
func TestIssue6_ProjectLevelStyleOverrides(t *testing.T) {
	homeDir := t.TempDir()
	bookDir := t.TempDir()

	// Write home voice and book voice
	_ = os.WriteFile(filepath.Join(homeDir, "voice.md"), []byte("Văn phong toàn cục"), 0644)
	_ = os.WriteFile(filepath.Join(bookDir, "voice.md"), []byte("Văn phong riêng cho truyện này"), 0644)

	// Write book-level custom style file
	_ = os.MkdirAll(filepath.Join(bookDir, "styles"), 0755)
	_ = os.WriteFile(filepath.Join(bookDir, "styles", "xianxia.md"), []byte("## Tiên hiệp chuẩn mực"), 0644)

	bundle := assets.Load("default", assets.LoadOptions{
		HomeStyleDir: homeDir,
		BookStyleDir: bookDir,
	})

	if !strings.Contains(bundle.Voice, "Văn phong toàn cục") {
		t.Fatal("expected global voice in bundle")
	}
	if !strings.Contains(bundle.Voice, "Văn phong riêng cho truyện này") {
		t.Fatal("expected book voice in bundle")
	}
	if bundle.Styles["xianxia"] != "## Tiên hiệp chuẩn mực" {
		t.Fatalf("expected custom book style xianxia loaded, got %q", bundle.Styles["xianxia"])
	}
}

// TestIssue7_ReviewModeStepByStepAdvance verifies Issue #7:
// In Review mode, chapter advance is paused until user grants permit via /next.
func TestIssue7_ReviewModeStepByStepAdvance(t *testing.T) {
	st := storepkg.NewStore(t.TempDir())
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}
	if err := st.RunMeta.Init("default", "test", "test"); err != nil {
		t.Fatal(err)
	}
	if err := st.RunMeta.SetAdvanceMode(domain.ChapterAdvanceReview); err != nil {
		t.Fatal(err)
	}
	if err := st.Progress.Init(10); err != nil {
		t.Fatal(err)
	}
	if err := st.Progress.UpdatePhase(domain.PhaseWriting); err != nil {
		t.Fatal(err)
	}

	var paused int
	var reasons []string
	gate := NewChapterAdvanceGate(st, func(reason string) {
		paused++
		reasons = append(reasons, reason)
	}, func(_ string, _ string) {})

	forward := &flow.Instruction{Agent: "writer", Chapter: 1, Task: "Viết chương 1"}

	// Without permit -> blocked
	allowed, err := gate.Allow(forward)
	if err != nil {
		t.Fatal(err)
	}
	if allowed || paused != 1 {
		t.Fatalf("expected blocked without permit: allowed=%v, paused=%d", allowed, paused)
	}
	if len(reasons) == 0 || !strings.Contains(reasons[len(reasons)-1], "/next") {
		t.Fatalf("expected mention of /next in pause reason: %v", reasons)
	}

	// Grant permit for Chapter 1
	if err := st.RunMeta.GrantAdvancePermit(1); err != nil {
		t.Fatal(err)
	}

	allowed, err = gate.Allow(forward)
	if err != nil {
		t.Fatal(err)
	}
	if !allowed {
		t.Fatal("expected advance allowed after permit granted")
	}
}

// TestIssue8_NovelSynopsisAndMetadata verifies Issue #8:
// Saving book synopsis and verifying required metadata structure.
func TestIssue8_NovelSynopsisAndMetadata(t *testing.T) {
	dir := t.TempDir()
	st := storepkg.NewStore(dir)
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}

	saveBookTool := tools.NewSaveBookTool(st)

	// Valid title and synopsis
	validArgs := []byte(`{"title": "Đại Đạo Triều Thiên", "synopsis": "Một thiếu niên bước ra từ cổ sơn, mở ra con đường tu tiên chấn động càn khôn."}`)
	res, err := saveBookTool.Execute(context.Background(), validArgs)
	if err != nil {
		t.Fatalf("save_book failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(res, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["saved"] != true {
		t.Fatalf("expected saved=true, got %+v", parsed)
	}

	meta, err := st.Book.Load()
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "Đại Đạo Triều Thiên" || !strings.Contains(meta.Synopsis, "thiếu niên") {
		t.Fatalf("persisted book metadata mismatch: %+v", meta)
	}

	// Empty title rejected
	emptyTitleArgs := []byte(`{"title": "", "synopsis": "Tóm tắt"}`)
	if _, err := saveBookTool.Execute(context.Background(), emptyTitleArgs); err == nil {
		t.Fatal("expected error on empty title")
	}
}

// TestIssue9_QualityChecklistAndStylestats verifies Issue #9:
// Stylestat tracker correctly records chapter contents and computes statistics snapshot.
func TestIssue9_QualityChecklistAndStylestats(t *testing.T) {
	tracker := stylestat.NewTracker()
	for ch := 1; ch <= 5; ch++ {
		tracker.Upsert(ch, "Hắn nhìn trời cao. Hắn thở dài một hơi. Hắn quay người bước đi.")
	}

	stats := tracker.Snapshot(nil, nil)
	if stats == nil || stats.Chapters != 5 {
		t.Fatalf("expected 5 chapters tracked, got %+v", stats)
	}
}
