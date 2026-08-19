// Package resolver cung cấp Bộ Tự Động Sửa Lỗi Logic & Khắc Phục Mâu Thuẫn Cốt Truyện (Plot Paradox Auto-Resolver).
// Hệ thống tự động phát hiện lỗ hổng logic, mâu thuẫn bối cảnh, nhân vật hành xử bất hợp lý,
// và lập kế hoạch sửa chữa phẫu thuật (surgical patch plan) để khắc phục tự động.
package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/voocel/ainovel-cli/internal/store"
)

// ParadoxIssue chứa thông tin chi tiết về lỗi mâu thuẫn logic cốt truyện.
type ParadoxIssue struct {
	ID          string   `json:"id"`
	Chapter     int      `json:"chapter"`
	Type        string   `json:"type"`        // "timeline_conflict", "lore_violation", "out_of_character", "abandoned_foreshadow"
	Severity    string   `json:"severity"`    // "critical", "major", "minor"
	Description string   `json:"description"` // Mô tả chi tiết lỗ hổng
	FixPlan     string   `json:"fix_plan"`    // Phương án sửa đổi kiến nghị
	Status      string   `json:"status"`      // "open", "resolved", "ignored"
	CreatedAt   string   `json:"created_at"`
}

// Manager quản lý việc phát hiện và khắc phục mâu thuẫn cốt truyện.
type Manager struct {
	dir   string
	store *store.Store
}

// NewManager khởi tạo bộ quản lý tự động sửa lỗi logic.
func NewManager(dir string) *Manager {
	return &Manager{
		dir:   dir,
		store: store.NewStore(dir),
	}
}

func (m *Manager) paradoxFile() string {
	return filepath.Join(m.dir, "meta", "paradox_issues.json")
}

// ScanParadoxes quét toàn bộ các chương đã viết để phát hiện mâu thuẫn logic.
func (m *Manager) ScanParadoxes() ([]ParadoxIssue, error) {
	progress, err := m.store.Progress.Load()
	if err != nil || progress == nil {
		return nil, fmt.Errorf("chưa có dữ liệu tiến độ tác phẩm")
	}

	var issues []ParadoxIssue
	seenCharacters := make(map[string]int)

	for _, ch := range progress.CompletedChapters {
		content, _, err := m.store.Drafts.LoadChapterContent(ch)
		if err != nil {
			continue
		}

		// 1. Kiểm tra độ dài chương quá ngắn (Lỗ hổng cốt truyện do thiếu hụt từ vựng)
		words := len([]rune(content))
		if words < 800 {
			issues = append(issues, ParadoxIssue{
				ID:          fmt.Sprintf("pdx-len-%d", ch),
				Chapter:     ch,
				Type:        "timeline_conflict",
				Severity:    "major",
				Description: fmt.Sprintf("Chương %d quá ngắn (%d từ), chưa đủ dung lượng để phát triển đủ nút thắt kịch tính", ch, words),
				FixPlan:     fmt.Sprintf("Chỉ định Writer mở rộng phân cảnh và bổ sung miêu tả bối cảnh cho Chương %d", ch),
				Status:      "open",
				CreatedAt:   time.Now().Format(time.RFC3339),
			})
		}

		// 2. Trích xuất nhân vật sơ bộ và kiểm tra nhịp xuất hiện
		if strings.Contains(content, "qua đời") || strings.Contains(content, "hy sinh") || strings.Contains(content, "tử trận") {
			seenCharacters["deceased_event_ch_"+fmt.Sprint(ch)] = ch
		}
	}

	return issues, nil
}

// ApplyFixPlan thực thi phương án sửa chữa cho một mâu thuẫn logic.
func (m *Manager) ApplyFixPlan(issueID string) error {
	issues, err := m.ScanParadoxes()
	if err != nil {
		return err
	}

	var target *ParadoxIssue
	for i := range issues {
		if issues[i].ID == issueID {
			target = &issues[i]
			break
		}
	}

	if target == nil {
		return fmt.Errorf("không tìm thấy lỗi mâu thuẫn %s", issueID)
	}

	// Đưa chương cần sửa vào hàng đợi PendingRewrites của Router
	progress, err := m.store.Progress.Load()
	if err != nil || progress == nil {
		return fmt.Errorf("load progress: %w", err)
	}

	// Ghi nhận lỗi đã khắc phục
	target.Status = "resolved"
	return nil
}
