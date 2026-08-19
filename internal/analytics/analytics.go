// Package analytics cung cấp bảng điều khiển thống kê chi phí Token, tốc độ sáng tác
// và tỷ lệ đóng góp của tác giả con người so với AI.
package analytics

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/cryptoaudit"
	"github.com/voocel/ainovel-cli/internal/store"
)

// Report chứa toàn bộ chỉ số thống kê của dự án.
type Report struct {
	BookTitle        string  `json:"book_title"`
	TotalChapters    int     `json:"total_chapters"`
	TotalWordCount   int     `json:"total_word_count"`
	AvgWordsPerChapter int   `json:"avg_words_per_chapter"`
	TotalTouchpoints int     `json:"total_touchpoints"`
	HumanEffortRatio float64 `json:"human_effort_ratio"` // Tỷ lệ can thiệp nghệ thuật của con người
	EstimatedTokens  int     `json:"estimated_tokens"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// GenerateReport tính toán toàn bộ chỉ số phân tích.
func GenerateReport(dir string) (*Report, error) {
	st := store.NewStore(dir)
	bookTitle := filepath.Base(dir)
	if book, err := st.Book.Load(); err == nil && book != nil && book.Title != "" {
		bookTitle = book.Title
	}

	progress, _ := st.Progress.Load()
	totalCh := 0
	totalWords := 0
	if progress != nil {
		totalCh = len(progress.CompletedChapters)
		totalWords = progress.TotalWordCount
	}

	avgWords := 0
	if totalCh > 0 {
		avgWords = totalWords / totalCh
	}

	ledger := authorship.NewLedger(dir)
	records, _ := ledger.LoadRecords()
	touchpoints := len(records)

	// Ước tính Token (Trung bình 1 từ tiếng Việt ~ 1.8 tokens cho LLM + Prompt context overhead)
	estimatedTokens := (totalWords * 2) + (totalCh * 8000)
	// Ước tính chi phí trung bình (Giả định giá hỗn hợp $1.5 / 1 triệu tokens cho Claude 3.5 / DeepSeek V3)
	estimatedCost := (float64(estimatedTokens) / 1000000.0) * 1.5

	humanRatio := 0.0
	if totalCh > 0 {
		humanRatio = (float64(touchpoints) / float64(totalCh)) * 100.0
	}

	return &Report{
		BookTitle:          bookTitle,
		TotalChapters:      totalCh,
		TotalWordCount:     totalWords,
		AvgWordsPerChapter: avgWords,
		TotalTouchpoints:   touchpoints,
		HumanEffortRatio:   humanRatio,
		EstimatedTokens:    estimatedTokens,
		EstimatedCostUSD:   estimatedCost,
	}, nil
}

// PrintReport in ra màn hình bảng phân tích đẹp mắt.
func PrintReport(dir string) error {
	rep, err := GenerateReport(dir)
	if err != nil {
		return err
	}

	dossier, _ := cryptoaudit.BuildDossier(dir, rep.BookTitle, "Human Author")

	fmt.Println("📊 BẢNG PHÂN TÍCH THỐNG KÊ CHI PHÍ & HIỆU SUẤT SÁNG TÁC (ANALYTICS DASHBOARD)")
	fmt.Println("=======================================================================")
	fmt.Printf("• Tên tác phẩm: 《%s》\n", rep.BookTitle)
	fmt.Printf("• Tổng số chương đã hoàn thành: %d chương\n", rep.TotalChapters)
	fmt.Printf("• Tổng số từ (Word Count): %d từ (Trung bình %d từ/chương)\n", rep.TotalWordCount, rep.AvgWordsPerChapter)
	fmt.Printf("• Số lần can thiệp sáng tạo của con người (Touchpoints): %d lần\n", rep.TotalTouchpoints)
	fmt.Printf("• Chỉ số kiểm soát nghệ thuật của con người (Human Control Index): %.1f%% \n", rep.HumanEffortRatio)
	fmt.Printf("• Phân lượng Token ước tính đã tiêu thụ: ~%d Tokens\n", rep.EstimatedTokens)
	fmt.Printf("• Chi phí API ước tính (LLM API Cost): ~$%.3f USD\n", rep.EstimatedCostUSD)

	if dossier != nil && dossier.MerkleRoot != "" {
		fmt.Printf("• Mã băm chứng thực số Merkle Root: 0x%s\n", dossier.MerkleRoot)
	}

	fmt.Println("=======================================================================")
	return nil
}
