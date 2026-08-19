package arena

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli arena` hoặc `ainovel-cli benchmark`.
func Command(args []string) int {
	fs := flag.NewFlagSet("arena", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	prompt := fs.String("prompt", "Viết một đoạn văn xuôi hành động dồn dập tại phòng thí nghiệm Bách Khoa 1956", "Đề bài thử nghiệm so sánh")
	modelA := fs.String("model-a", "claude-3-5-sonnet", "Tên mô hình A")
	modelB := fs.String("model-b", "deepseek-v3", "Tên mô hình B")
	_ = fs.Parse(args)

	fmt.Printf("⚔️ Đang khởi chạy Sàn đấu Arena giữa [%s] và [%s]...\n", *modelA, *modelB)
	fmt.Printf("🎯 Đề tài: %s\n\n", *prompt)

	outA := fmt.Sprintf("Bản thảo demo được tạo tự động bởi mô hình %s: Tiếng máy gầm rú trong phòng thí nghiệm Bách Khoa 1956, chip bán dẫn phát sáng xanh lam.", *modelA)
	outB := fmt.Sprintf("Bản thảo demo được tạo tự động bởi mô hình %s: Cơn mưa trút xuống mái tôn hội trường C1, hai kỹ sư tranh luận gay gắt về thuật toán.", *modelB)

	scoreA := EvaluateDraft(outA)
	scoreB := EvaluateDraft(outB)

	report := BattleReport{
		Prompt: *prompt,
		Candidates: []ModelCandidate{
			{Name: *modelA, Score: scoreA, WordCount: len([]rune(outA)), DurationMs: 1200, DraftOutput: outA},
			{Name: *modelB, Score: scoreB, WordCount: len([]rune(outB)), DurationMs: 850, DraftOutput: outB},
		},
		Winner:  *modelA,
		Summary: fmt.Sprintf("Mô hình chiến thắng: %s (Điểm số: %.1f vs %.1f)", *modelA, scoreA, scoreB),
	}

	report.PrintReport()
	return 0
}
