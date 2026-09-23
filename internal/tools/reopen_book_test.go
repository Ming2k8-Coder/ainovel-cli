package tools

import (
	"testing"

	"github.com/voocel/ainovel-cli/internal/domain"
	"github.com/voocel/ainovel-cli/internal/store"
)

// completedBook tạo một cuốn tiểu thuyết N chương đã hoàn thành (phase=complete, CompletedChapters=1..n).
func completedBook(t *testing.T, n int) *store.Store {
	t.Helper()
	s := store.NewStore(t.TempDir())
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := s.Progress.Init(n); err != nil {
		t.Fatalf("InitProgress: %v", err)
	}
	for ch := 1; ch <= n; ch++ {
		if err := s.Progress.MarkChapterComplete(ch, 100, "", ""); err != nil {
			t.Fatalf("MarkChapterComplete(%d): %v", ch, err)
		}
	}
	if err := s.Progress.MarkComplete(); err != nil {
		t.Fatalf("MarkComplete: %v", err)
	}
	return s
}

func TestReopenBookReopensCompletedBook(t *testing.T) {
	s := completedBook(t, 3)

	if err := ReopenBook(s, []int{3, 1}, "Dọn dẹp các ký tự đặc biệt"); err != nil {
		t.Fatalf("ReopenBook: %v", err)
	}

	p, _ := s.Progress.Load()
	if p.Phase != domain.PhaseWriting {
		t.Errorf("phase = %s, want writing", p.Phase)
	}
	if p.Flow != domain.FlowRewriting {
		t.Errorf("flow = %s, want rewriting", p.Flow)
	}
	if len(p.PendingRewrites) != 2 || p.PendingRewrites[0] != 3 || p.PendingRewrites[1] != 1 {
		t.Errorf("PendingRewrites = %v, want [3 1] (giữ nguyên thứ tự đưa vào hàng đợi)", p.PendingRewrites)
	}

	if cp := s.Checkpoints.LatestByStep(domain.GlobalScope(), "reopen"); cp == nil {
		t.Error("expected a 'reopen' checkpoint")
	}
}

func TestReopenBookRejectsNonCompleteBook(t *testing.T) {
	// Sách đang viết (chưa hoàn thành) không thể reopen
	s := store.NewStore(t.TempDir())
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := s.Progress.Init(5); err != nil {
		t.Fatalf("InitProgress: %v", err)
	}
	if err := s.Progress.MarkChapterComplete(1, 100, "", ""); err != nil { // phase→writing
		t.Fatalf("MarkChapterComplete: %v", err)
	}
	if err := ReopenBook(s, []int{1}, ""); err == nil {
		t.Fatal("expected reopen to be rejected when phase != complete")
	}
}

func TestReopenBookRejectsUnwrittenChapters(t *testing.T) {
	s := completedBook(t, 3)

	// Chương 5 không tồn tại → từ chối (thuộc dạng viết tiếp/vượt giới hạn, cần điều chỉnh độ dài)
	if err := ReopenBook(s, []int{2, 5}, ""); err == nil {
		t.Fatal("expected reopen to be rejected for unwritten chapter")
	}
	// chapters rỗng → từ chối
	if err := ReopenBook(s, nil, ""); err == nil {
		t.Fatal("expected reopen to be rejected for empty chapters")
	}
}
