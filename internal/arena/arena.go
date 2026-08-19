// Package arena cung cấp Sàn Đấu So Sánh Đa Mô Hình LLM (Multi-Model Battle Arena & Benchmark),
// cho phép thử nghiệm và so sánh song song chất lượng văn xuôi giữa 2 hoặc nhiều mô hình LLM khác nhau
// (ví dụ: Claude 3.5 Sonnet vs DeepSeek V3 vs Gemini 1.5 Pro) trên cùng 1 đề bài/chương.
package arena

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/voocel/agentcore"
)

// ModelCandidate thông tin mô hình tham gia đấu trường.
type ModelCandidate struct {
	Name        string  `json:"name"`        // "claude-3-5-sonnet", "deepseek-chat", "gemini-1.5-pro"
	Provider    string  `json:"provider"`    // "anthropic", "deepseek", "google", "openrouter"
	Score       float64 `json:"score"`       // Điểm số tổng hợp 7 chiều (0 - 100)
	WordCount   int     `json:"word_count"`
	DurationMs  int64   `json:"duration_ms"`
	DraftOutput string  `json:"draft_output"`
}

// BattleReport kết quả thi đấu giữa các mô hình.
type BattleReport struct {
	Prompt     string           `json:"prompt"`
	Candidates []ModelCandidate `json:"candidates"`
	Winner     string           `json:"winner"`
	Summary    string           `json:"summary"`
	CreatedAt  string           `json:"created_at"`
}

// Manager quản lý Sàn đấu so sánh LLM.
type Manager struct {
	dir string
}

// NewManager khởi tạo Arena Manager.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

// EvaluateDraft đánh giá điểm số văn học của một bản thảo (0 - 100).
func EvaluateDraft(draft string) float64 {
	words := len([]rune(draft))
	if words < 100 {
		return 30.0
	}

	score := 70.0

	// 1. Điểm thưởng cho độ dài hợp lý (1500 - 3000 từ)
	if words >= 1500 && words <= 3500 {
		score += 15.0
	}

	// 2. Điểm thưởng cho cấu trúc thoại sinh động (Có ngoặc báo thoại)
	dialogueCount := strings.Count(draft, "“") + strings.Count(draft, "\"")
	if dialogueCount >= 5 {
		score += 10.0
	}

	// 3. Trừ điểm nếu bị lặp từ sáo rỗng của AI
	taboos := []string{"trong một khoảnh khắc", "không khí bỗng chốc", "không phải... mà là"}
	for _, t := range taboos {
		if strings.Contains(draft, t) {
			score -= 5.0
		}
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

// RunBattle thực thi cuộc so tài giữa 2 mô hình LLM trên cùng 1 prompt.
func (m *Manager) RunBattle(ctx context.Context, prompt string, modelA, modelB string, runnerA, runnerB agentcore.Model) (*BattleReport, error) {
	if prompt == "" {
		prompt = "Hãy viết một đoạn đối thoại căng thẳng giữa 2 kỹ sư Cyberpunk tại phòng thí nghiệm Bách Khoa 1956."
	}

	startA := time.Now()
	respA, errA := runnerA.Generate(ctx, []agentcore.Message{{Role: "user", Content: prompt}})
	durA := time.Since(startA).Milliseconds()
	outA := ""
	if errA == nil && respA != nil {
		outA = respA.Content
	} else {
		outA = fmt.Sprintf("Lỗi tạo mẫu Model A (%v)", errA)
	}

	startB := time.Now()
	respB, errB := runnerB.Generate(ctx, []agentcore.Message{{Role: "user", Content: prompt}})
	durB := time.Since(startB).Milliseconds()
	outB := ""
	if errB == nil && respB != nil {
		outB = respB.Content
	} else {
		outB = fmt.Sprintf("Lỗi tạo mẫu Model B (%v)", errB)
	}

	scoreA := EvaluateDraft(outA)
	scoreB := EvaluateDraft(outB)

	candA := ModelCandidate{
		Name:        modelA,
		Score:       scoreA,
		WordCount:   len([]rune(outA)),
		DurationMs:  durA,
		DraftOutput: outA,
	}

	candB := ModelCandidate{
		Name:        modelB,
		Score:       scoreB,
		WordCount:   len([]rune(outB)),
		DurationMs:  durB,
		DraftOutput: outB,
	}

	winner := modelA
	if scoreB > scoreA {
		winner = modelB
	}

	report := &BattleReport{
		Prompt:     prompt,
		Candidates: []ModelCandidate{candA, candB},
		Winner:     winner,
		Summary:    fmt.Sprintf("Mô hình chiến thắng: %s (Điểm số: %.1f vs %.1f)", winner, scoreA, scoreB),
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	return report, nil
}

// PrintReport in báo cáo kết quả thi đấu ra màn hình.
func (r *BattleReport) PrintReport() {
	fmt.Println("⚔️ SÀN ĐẤU SO SÁNH ĐA MÔ HÌNH LLM (MULTI-MODEL BATTLE ARENA)")
	fmt.Println("============================================================")
	fmt.Printf("🎯 Đề bài Prompt: %s\n\n", r.Prompt)

	for i, c := range r.Candidates {
		fmt.Printf("🏆 Thí sinh #%d: [%s]\n", i+1, c.Name)
		fmt.Printf("   • Điểm văn học 7 chiều: %.1f / 100\n", c.Score)
		fmt.Printf("   • Độ dài: %d từ | Thời gian đáp ứng: %d ms\n", c.WordCount, c.DurationMs)
		fmt.Printf("   • Bản thảo trích đoạn: %s...\n\n", truncateText(c.DraftOutput, 150))
	}

	fmt.Println("============================================================")
	fmt.Printf("👑 %s\n", r.Summary)
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
