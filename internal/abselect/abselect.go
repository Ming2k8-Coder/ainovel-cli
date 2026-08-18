// Package abselect cung cấp hệ thống Lựa chọn A/B (A/B Selection) và Phản hồi Tự động
// vào Lời nhắc Hệ thống (System Prompts & Voice Rules).
//
// Cơ chế này phục vụ 2 mục tiêu chiến lược:
//  1. Quyền tác giả (Copyright Authorship): Con người thực hiện quyền "Tuyển chọn nghệ thuật" (Artistic Selection)
//     bằng việc chọn nhánh A hoặc B, được tự động ghi nhận vào Authorship Ledger.
//  2. Vòng lặp phản hồi thông minh (Prompt Feedback Loop): Từ các lựa chọn A/B của người dùng,
//     hệ thống tự động phân tích sở thích văn phong (thích thoại sắc bén, nhịp điệu nhanh, chi tiết kỹ thuật...)
//     và tự động cập nhật trực tiếp vào System Prompts & Quy tắc văn phong của tác phẩm.
package abselect

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/domain"
)

// Variant đại diện cho một phương án A hoặc B được tạo ra để người dùng lựa chọn.
type Variant struct {
	Key         string   `json:"key"`                   // "A" hoặc "B"
	Title       string   `json:"title"`                 // Tên tóm tắt phương án
	Description string   `json:"description"`           // Mô tả trọng tâm kịch bản/văn phong
	Content     string   `json:"content"`               // Nội dung văn bản mẫu hoặc dàn ý
	Traits      []string `json:"traits"`                // Các đặc trưng văn học (ví dụ: "dialogue_sharp", "fast_paced", "hard_scifi_detail")
	Model       string   `json:"model,omitempty"`       // Tên mô hình tạo ra (nếu thử nghiệm đa mô hình)
}

// Choice lưu trữ quyết định lựa chọn của người dùng.
type Choice struct {
	ID         string    `json:"id"`
	Chapter    int       `json:"chapter,omitempty"`
	Scope      string    `json:"scope"`                 // "premise", "outline_arc", "chapter_draft", "dialogue_tone"
	OptionA    Variant   `json:"option_a"`
	OptionB    Variant   `json:"option_b"`
	Selected   string    `json:"selected"`              // "A", "B", hoặc "custom"
	UserReason string    `json:"user_reason,omitempty"` // Ghi chú ngắn của tác giả (tùy chọn)
	Timestamp  string    `json:"timestamp"`
}

// Manager quản lý việc lưu trữ, phân tích và phản hồi các lựa chọn A/B.
type Manager struct {
	dir        string
	authorship *authorship.Ledger
	mu         sync.RWMutex
}

// NewManager khởi tạo bộ quản lý A/B selection.
func NewManager(dir string, ledger *authorship.Ledger) *Manager {
	return &Manager{
		dir:        dir,
		authorship: ledger,
	}
}

func (m *Manager) feedbackFile() string {
	return filepath.Join(m.dir, "meta", "ab_feedback.jsonl")
}

// RecordChoice lưu lại quyết định A/B của người dùng và đồng thời đồng bộ vào Authorship Ledger.
func (m *Manager) RecordChoice(c Choice) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.Timestamp == "" {
		c.Timestamp = time.Now().Format(time.RFC3339)
	}
	if c.ID == "" {
		c.ID = fmt.Sprintf("ab-%d-%s", time.Now().UnixNano(), c.Selected)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal ab choice: %w", err)
	}

	p := m.feedbackFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("create meta dir: %w", err)
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open ab feedback file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write ab feedback: %w", err)
	}

	// Tự động ghi nhận vào Authorship Ledger (Bằng chứng quyền tác giả)
	if m.authorship != nil {
		selectedVariant := c.OptionA
		if c.Selected == "B" {
			selectedVariant = c.OptionB
		}
		desc := fmt.Sprintf("Lựa chọn phương án %s (%s): %s", c.Selected, selectedVariant.Title, selectedVariant.Description)
		_, _ = m.authorship.RecordContribution(
			authorship.TypeCreativeChoice,
			"Author",
			desc,
			c.UserReason,
			c.Chapter,
		)
	}

	return nil
}

// LoadChoices nạp toàn bộ lịch sử lựa chọn A/B.
func (m *Manager) LoadChoices() ([]Choice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.feedbackFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var choices []Choice
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var c Choice
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			continue
		}
		choices = append(choices, c)
	}
	return choices, nil
}

// SynthesizedPromptFeedback chứa các quy tắc văn phong và sở thích được học tự động từ các lựa chọn A/B.
type SynthesizedPromptFeedback struct {
	PreferredTraits []string                  `json:"preferred_traits"`
	AvoidedTraits   []string                  `json:"avoided_traits"`
	ProseDirectives []string                  `json:"prose_directives"`
	Taboos          []string                  `json:"taboos"`
	AuthorStyle     domain.AuthorRevisionStyle `json:"author_style"`
}

// SynthesizePreferences phân tích các lựa chọn trong quá khứ để đúc kết thành chỉ thị Prompt.
func (m *Manager) SynthesizePreferences() (SynthesizedPromptFeedback, error) {
	choices, err := m.LoadChoices()
	if err != nil {
		return SynthesizedPromptFeedback{}, err
	}

	traitVotes := make(map[string]int)

	for _, c := range choices {
		var chosenTraits []string
		var rejectedTraits []string

		if c.Selected == "A" {
			chosenTraits = c.OptionA.Traits
			rejectedTraits = c.OptionB.Traits
		} else if c.Selected == "B" {
			chosenTraits = c.OptionB.Traits
			rejectedTraits = c.OptionA.Traits
		}

		for _, t := range chosenTraits {
			traitVotes[t]++
		}
		for _, t := range rejectedTraits {
			traitVotes[t]--
		}
	}

	res := SynthesizedPromptFeedback{
		AuthorStyle: domain.AuthorRevisionStyle{
			UpdatedAt: time.Now(),
		},
	}

	for trait, score := range traitVotes {
		if score > 0 {
			res.PreferredTraits = append(res.PreferredTraits, trait)
			directive := mapTraitToDirective(trait)
			if directive != "" {
				res.ProseDirectives = append(res.ProseDirectives, directive)
				res.AuthorStyle.Prose = append(res.AuthorStyle.Prose, directive)
			}
		} else if score < 0 {
			res.AvoidedTraits = append(res.AvoidedTraits, trait)
			taboo := mapTraitToTaboo(trait)
			if taboo != "" {
				res.Taboos = append(res.Taboos, taboo)
				res.AuthorStyle.Taboos = append(res.AuthorStyle.Taboos, taboo)
			}
		}
	}

	return res, nil
}

// ApplyFeedbackToProject tự động áp dụng các sở thích đã học vào tệp `style/voice.md` của tác phẩm.
func (m *Manager) ApplyFeedbackToProject() error {
	feedback, err := m.SynthesizePreferences()
	if err != nil {
		return err
	}

	if len(feedback.ProseDirectives) == 0 && len(feedback.Taboos) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("# Quy Tắc Văn Phong & Sở Thích Tác Giả (Học Tự Động Từ Lựa Chọn A/B)\n\n")
	b.WriteString("> Tài liệu này được cập nhật tự động dựa trên các quyết định nghệ thuật của tác giả qua các đợt A/B Selection.\n\n")

	if len(feedback.ProseDirectives) > 0 {
		b.WriteString("## 1. Định Hướng Hành Văn Ưu Tiên (Preferred Directives)\n\n")
		for _, d := range feedback.ProseDirectives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}

	if len(feedback.Taboos) > 0 {
		b.WriteString("## 2. Các Yếu Tố Cần Tránh (Taboos & Constraints)\n\n")
		for _, t := range feedback.Taboos {
			fmt.Fprintf(&b, "- %s\n", t)
		}
		b.WriteString("\n")
	}

	styleDir := filepath.Join(m.dir, "style")
	if err := os.MkdirAll(styleDir, 0o755); err != nil {
		return err
	}

	voicePath := filepath.Join(styleDir, "voice.md")
	return os.WriteFile(voicePath, []byte(b.String()), 0o644)
}

func mapTraitToDirective(trait string) string {
	switch trait {
	case "dialogue_sharp":
		return "Ưu tiên đối thoại sắc bén, ngắn gọn, giàu kịch tính và thể hiện rõ cá tính riêng của nhân vật."
	case "fast_paced":
		return "Duy trì nhịp điệu nhanh, tiết tấu dồn dập, đẩy mạnh xung đột hành động thay vì miêu tả rườm rà."
	case "hard_scifi_detail":
		return "Tập trung miêu tả chính xác các thông số kỹ thuật, logic khoa học công nghệ và bối cảnh chân thực."
	case "deep_psychology":
		return "Khắc họa sâu sắc diễn biến tâm lý, độc thoại nội tâm và sự giằng xé nội tại của nhân vật."
	case "atmospheric_mood":
		return "Xây dựng bầu không khí đậm chất điện ảnh, giàu chất cảm không gian và hình ảnh gợi cảm."
	default:
		return fmt.Sprintf("Ưu tiên đặc trưng phong cách: %s", trait)
	}
}

func mapTraitToTaboo(trait string) string {
	switch trait {
	case "verbose_exposition":
		return "Tránh việc thuyết minh giáo điều dài dòng hoặc giải thích quá nhiều về bối cảnh thay vì để nhân vật tự bộc lộ."
	case "slow_pacing":
		return "Tránh nhịp điệu chậm chạp, lê thê, phân cảnh thiếu xung đột thúc đẩy câu chuyện."
	case "melodrama_cliche":
		return "Tránh các mô-típ sướt mướt gượng gạo, bi kịch hóa thái quá hoặc lời thoại sáo rỗng."
	case "superficial_dialogue":
		return "Tránh các câu thoại vô thưởng vô phạt không mang lại thông tin hay chuyển biến tình cảm."
	default:
		return fmt.Sprintf("Hạn chế đặc trưng: %s", trait)
	}
}
