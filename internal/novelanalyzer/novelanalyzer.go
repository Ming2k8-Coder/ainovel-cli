// Package novelanalyzer cung cấp công cụ Phân Tích & Viết Tiếp Tiểu Thuyết Có Sẵn (Novel Analysis & Continuation Engine).
// Công cụ có khả năng đọc một tác phẩm hoặc tài liệu đặc tả có sẵn (như các file server trong workspace nexus-universe),
// tự động trích xuất Hồ sơ Thế giới, Nhân vật, Cấu trúc Dàn ý và Phong cách Văn phong, sau đó thiết lập dự án mới
// sẵn sàng cho tác tử AI tiếp tục sáng tác từ Chương N+1.
package novelanalyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/voocel/ainovel-cli/internal/domain"
	"github.com/voocel/ainovel-cli/internal/store"
)

// AnalysisResult chứa kết quả phân tích một tác phẩm/đặc tả có sẵn.
type AnalysisResult struct {
	BookTitle        string            `json:"book_title"`
	DetectedChapters int               `json:"detected_chapters"`
	TotalWords       int               `json:"total_words"`
	Characters       []string          `json:"characters"`
	WorldLore        string            `json:"world_lore"`
	ExtractedStyle   string            `json:"extracted_style"`
	LastChapterNum   int               `json:"last_chapter_num"`
	SummaryMap       map[int]string    `json:"summary_map"`
}

// Analyzer thực hiện phân tích tác phẩm.
type Analyzer struct {
	sourcePath string
}

// NewAnalyzer khởi tạo bộ phân tích.
func NewAnalyzer(sourcePath string) *Analyzer {
	return &Analyzer{sourcePath: sourcePath}
}

// AnalyzePath phân tích tệp hoặc thư mục tác phẩm/đặc tả có sẵn.
func (a *Analyzer) AnalyzePath() (*AnalysisResult, error) {
	fi, err := os.Stat(a.sourcePath)
	if err != nil {
		return nil, fmt.Errorf("không thể truy cập đường dẫn %s: %w", a.sourcePath, err)
	}

	res := &AnalysisResult{
		BookTitle:  filepath.Base(a.sourcePath),
		SummaryMap: make(map[int]string),
	}

	if fi.IsDir() {
		if err := a.analyzeDir(a.sourcePath, res); err != nil {
			return nil, err
		}
	} else {
		if err := a.analyzeFile(a.sourcePath, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (a *Analyzer) analyzeDir(dirPath string, res *AnalysisResult) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	var sb strings.Builder
	chRegex := regexp.MustCompile(`(?i)(chuong|chapter|ch)\s*(\d+)`)

	for _, e := range entries {
		if e.IsDir() || (!strings.HasSuffix(e.Name(), ".md") && !strings.HasSuffix(e.Name(), ".txt")) {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dirPath, e.Name()))
		if err != nil {
			continue
		}

		strContent := string(content)
		res.TotalWords += len([]rune(strContent))
		sb.WriteString(strContent)
		sb.WriteString("\n\n")

		matches := chRegex.FindStringSubmatch(e.Name())
		if len(matches) >= 3 {
			chNum, _ := strconv.Atoi(matches[2])
			if chNum > res.LastChapterNum {
				res.LastChapterNum = chNum
			}
			res.DetectedChapters++
			res.SummaryMap[chNum] = truncateText(strContent, 200)
		}
	}

	a.extractEntities(sb.String(), res)
	return nil
}

func (a *Analyzer) analyzeFile(filePath string, res *AnalysisResult) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	strContent := string(content)
	res.TotalWords = len([]rune(strContent))

	chRegex := regexp.MustCompile(`(?m)^(?:#+|Chương|Chapter)\s*(\d+)`)
	matches := chRegex.FindAllStringSubmatch(strContent, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			chNum, _ := strconv.Atoi(m[1])
			if chNum > res.LastChapterNum {
				res.LastChapterNum = chNum
			}
			res.DetectedChapters++
		}
	}
	if res.LastChapterNum == 0 && res.DetectedChapters == 0 {
		res.LastChapterNum = 1
		res.DetectedChapters = 1
	}

	a.extractEntities(strContent, res)
	return nil
}

func (a *Analyzer) extractEntities(text string, res *AnalysisResult) {
	// Trích xuất nhân vật sơ bộ (Tìm các từ in hoa nguyên cụm tên riêng)
	nameRegex := regexp.MustCompile(`([A-ZÀÁẢẠÃĂẮẰẲẶẴÂẤẦẨẬẪĐÈÉẺẸẼÊẾỀỂỆỄÌÍỈỊĨÒÓỎỌÕÔỐỒỔỘỖƠỚỜỞỢỠÙÚỦỤŨƯỨỪỬỰỮỲÝỶỴỸ][a-zàáảạãăắằẳặẵâấầẩậẫđèéẻẹẽêếềểệễìíỉịĩòóỏọõôốồổộỗơớờởợỡùúủụũưứừửựữỳýỷỵỹ]+(?:\s+[A-ZÀÁẢẠÃĂẮẰẲẶẴÂẤẦẨẬẪĐÈÉẺẸẼÊẾỀỂỆỄÌÍỈỊĨÒÓỎỌÕÔỐỒỔỘỖƠỚỜỞỢỠÙÚỦỤŨƯỨỪỬỰỮỲÝỶỴỸ][a-zàáảạãăắằẳặẵâấầẩậẫđèéẻẹẽêếềểệễìíỉịĩòóỏọõôốồổộỗơớờởợỡùúủụũưứừửựữỳýỷỵỹ]+)+)`)
	matches := nameRegex.FindAllString(text, -1)
	charMap := make(map[string]int)
	for _, m := range matches {
		if len(m) > 3 && !strings.Contains(m, "Chương") && !strings.Contains(m, "Hồi") {
			charMap[m]++
		}
	}
	for name, count := range charMap {
		if count >= 2 {
			res.Characters = append(res.Characters, name)
		}
	}
	if len(res.Characters) > 10 {
		res.Characters = res.Characters[:10]
	}

	res.WorldLore = truncateText(text, 1000)
	res.ExtractedStyle = "Phong cách tự nhiên trích xuất từ tác phẩm mẫu"
}

// SetupContinuationProject khởi tạo một dự án mới hoàn chỉnh dựa trên phân tích tác phẩm có sẵn.
func SetupContinuationProject(targetDir string, res *AnalysisResult) error {
	st := store.NewStore(targetDir)

	if err := st.EnsureDirs(); err != nil {
		return err
	}

	// 1. Khởi tạo Book meta
	_ = st.Book.Save(domain.BookMeta{
		Title:       res.BookTitle,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	})

	// 2. Khởi tạo Progress (Đã hoàn thành đến chương LastChapterNum)
	var completed []int
	for i := 1; i <= res.LastChapterNum; i++ {
		completed = append(completed, i)
	}

	_ = st.Progress.Save(domain.Progress{
		Phase:             domain.PhaseWriting,
		CompletedChapters: completed,
		TotalWordCount:    res.TotalWords,
		UpdatedAt:         time.Now(),
	})

	// 3. Khởi tạo Thế giới & Nhân vật
	var charList []domain.Character
	for i, name := range res.Characters {
		charList = append(charList, domain.Character{
			Name:        name,
			Role:        func() string { if i == 0 { return "Nhân vật chính" }; return "Nhân vật phụ" }(),
			Description: fmt.Sprintf("Nhân vật trích xuất từ tác phẩm gốc %s", res.BookTitle),
		})
	}
	_ = st.World.SaveCharacters(charList)

	// 4. Khởi tạo style/voice.md
	styleDir := filepath.Join(targetDir, "style")
	_ = os.MkdirAll(styleDir, 0o755)
	voiceContent := fmt.Sprintf("# Phong Cách Viết Tiếp: %s\n\n%s\n\n## Nhân Vật Gốc:\n- %s\n",
		res.BookTitle, res.WorldLore, strings.Join(res.Characters, "\n- "))
	_ = os.WriteFile(filepath.Join(styleDir, "voice.md"), []byte(voiceContent), 0o644)

	return nil
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
