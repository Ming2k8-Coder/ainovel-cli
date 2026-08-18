// Package authorship cung cấp hệ thống lưu vết quyền tác giả (Human Authorship Ledger)
// và chứng minh quyền sở hữu trí tuệ (Copyright Proof) cho các tác phẩm đồng sáng tác cùng AI.
//
// Thiết kế nhằm tối thiểu hóa thao tác của người dùng (Minimal Human Effort) nhưng vẫn thỏa mãn
// các tiêu chuẩn khắt khe về "Sự đóng góp sáng tạo có ý thức của con người" (Human Creative Control & Selection)
// theo quy định của Cục Bản quyền Tác giả (USCO, WIPO, Công ước Bern, Luật SHTT Việt Nam).
package authorship

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ContributionType định nghĩa các hình thức đóng góp sáng tạo của con người.
type ContributionType string

const (
	// TypePremiseSeed: Ý tưởng cốt lõi, tiền đề nghệ thuật do con người đặt ra.
	TypePremiseSeed ContributionType = "premise_seed"
	// TypeCreativeSteer: Can thiệp định hướng tình tiết, chuyển hướng cốt truyện hoặc thiết lập nhân vật.
	TypeCreativeSteer ContributionType = "creative_steer"
	// TypeOutlineCuration: Thẩm định, phê duyệt hoặc tinh chỉnh đại cương phân tầng.
	TypeOutlineCuration ContributionType = "outline_curation"
	// TypeCreativeChoice: Lựa chọn giữa các nhánh cốt truyện/phương án kịch tính do AI đề xuất.
	TypeCreativeChoice ContributionType = "creative_choice"
	// TypeHumanEdit: Trực tiếp chỉnh sửa, gọt giũa câu từ trong bản thảo.
	TypeHumanEdit ContributionType = "human_edit"
	// TypeChapterApproval: Phê duyệt nghiệm thu chương tại điểm chốt (Chapter Advance Gate).
	TypeChapterApproval ContributionType = "chapter_approval"
)

// Record lưu giữ một bằng chứng đóng góp sáng tạo cụ thể của con người.
type Record struct {
	ID          string           `json:"id"`
	Type        ContributionType `json:"type"`
	Author      string           `json:"author"`
	Chapter     int              `json:"chapter,omitempty"`
	Description string           `json:"description"`
	InputText   string           `json:"input_text,omitempty"`
	Hash        string           `json:"hash"`
	Timestamp   string           `json:"timestamp"`
}

// Ledger quản lý việc ghi nhận và truy xuất sổ cái quyền tác giả.
type Ledger struct {
	dir string
	mu  sync.RWMutex
}

// NewLedger khởi tạo Ledger trong thư mục dự án.
func NewLedger(dir string) *Ledger {
	return &Ledger{dir: dir}
}

func (l *Ledger) filePath() string {
	return filepath.Join(l.dir, "meta", "authorship_ledger.jsonl")
}

// RecordContribution ghi nhận một đóng góp sáng tạo mới của con người vào sổ cái bất biến.
func (l *Ledger) RecordContribution(contribType ContributionType, author, desc, inputText string, chapter int) (*Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if strings.TrimSpace(author) == "" {
		author = "Human Author"
	}
	now := time.Now().Format(time.RFC3339)

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s:%s:%d:%s", contribType, author, desc, inputText, chapter, now)))
	hashStr := hex.EncodeToString(hasher.Sum(nil))

	rec := Record{
		ID:          fmt.Sprintf("auth-%s", hashStr[:12]),
		Type:        contribType,
		Author:      author,
		Chapter:     chapter,
		Description: desc,
		InputText:   inputText,
		Hash:        hashStr,
		Timestamp:   now,
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("marshal authorship record: %w", err)
	}

	p := l.filePath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("create meta dir: %w", err)
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open authorship ledger: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("append authorship record: %w", err)
	}

	return &rec, nil
}

// LoadRecords đọc toàn bộ lịch sử đóng góp sáng tạo từ sổ cái.
func (l *Ledger) LoadRecords() ([]Record, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	p := l.filePath()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var records []Record
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var rec Record
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		records = append(records, rec)
	}

	return records, nil
}

// Summary thống kê các chỉ số đóng góp của con người.
type Summary struct {
	TotalInterventions int            `json:"total_interventions"`
	PremiseSeeds       int            `json:"premise_seeds"`
	CreativeSteers     int            `json:"creative_steers"`
	OutlineCurations   int            `json:"outline_curations"`
	CreativeChoices    int            `json:"creative_choices"`
	HumanEdits         int            `json:"human_edits"`
	ChapterApprovals   int            `json:"chapter_approvals"`
	FirstContribution  string         `json:"first_contribution"`
	LastContribution   string         `json:"last_contribution"`
	Authors            map[string]int `json:"authors"`
}

// GetSummary tính toán bản tổng kết thống kê đóng góp của con người.
func (l *Ledger) GetSummary() (Summary, error) {
	records, err := l.LoadRecords()
	if err != nil {
		return Summary{}, err
	}

	sum := Summary{
		TotalInterventions: len(records),
		Authors:            make(map[string]int),
	}

	if len(records) > 0 {
		sum.FirstContribution = records[0].Timestamp
		sum.LastContribution = records[len(records)-1].Timestamp
	}

	for _, r := range records {
		sum.Authors[r.Author]++
		switch r.Type {
		case TypePremiseSeed:
			sum.PremiseSeeds++
		case TypeCreativeSteer:
			sum.CreativeSteers++
		case TypeOutlineCuration:
			sum.OutlineCurations++
		case TypeCreativeChoice:
			sum.CreativeChoices++
		case TypeHumanEdit:
			sum.HumanEdits++
		case TypeChapterApproval:
			sum.ChapterApprovals++
		}
	}

	return sum, nil
}

// GenerateCopyrightReport tạo báo cáo pháp lý chứng minh quyền tác giả (Copyright & Authorship Dossier).
func (l *Ledger) GenerateCopyrightReport(bookTitle, primaryAuthor string) (string, error) {
	records, err := l.LoadRecords()
	if err != nil {
		return "", err
	}
	summary, err := l.GetSummary()
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(primaryAuthor) == "" {
		primaryAuthor = "Tác giả / Nhà sáng tạo chính"
	}
	if strings.TrimSpace(bookTitle) == "" {
		bookTitle = "Tác phẩm Không Tên"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# HỒ SƠ CHỨNG MINH QUYỀN TÁC GIẢ & SỞ HỮU TRÍ TUỆ\n\n")
	fmt.Fprintf(&b, "> **Tác phẩm**: 《%s》\n", bookTitle)
	fmt.Fprintf(&b, "> **Tác giả chủ quản**: %s\n", primaryAuthor)
	fmt.Fprintf(&b, "> **Ngày xuất báo cáo**: %s\n", time.Now().Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(&b, "> **Tiêu chuẩn áp dụng**: US Copyright Office (Compendium III) / Công ước Bern / Luật SHTT Việt Nam\n\n", )
	fmt.Fprintf(&b, "---\n\n")

	fmt.Fprintf(&b, "## 1. Tuyên Bố Quyền Tác Giả Của Con Người (Human Authorship Statement)\n\n")
	fmt.Fprintf(&b, "Tác phẩm này được sáng tạo thông qua quá trình **Con người định hướng và tuyển chọn nghệ thuật (Human Creative Direction, Selection & Arrangement)** với sự hỗ trợ công nghệ từ hệ thống `ainovel-cli`:\n\n")
	fmt.Fprintf(&b, "- **Ý tưởng & Tiền đề nghệ thuật**: Do tác giả con người trực tiếp hình thành và thiết lập.\n")
	fmt.Fprintf(&b, "- **Cấu trúc & Tuyến tính cốt truyện**: Được tác giả phê duyệt và định hình qua từng Hồi/Quyển.\n")
	fmt.Fprintf(&b, "- **Kiểm soát chất lượng & Thẩm mỹ**: Tác giả trực tiếp kiểm soát cổng phát hành và điều chỉnh tình tiết.\n\n")

	fmt.Fprintf(&b, "## 2. Thống Kê Điểm Chạm Sáng Tạo (Human Touchpoints Summary)\n\n")
	fmt.Fprintf(&b, "| Chỉ số đo lường | Số lượng ghi nhận | Diễn giải pháp lý |\n")
	fmt.Fprintf(&b, "|---|---|---|\n")
	fmt.Fprintf(&b, "| **Tổng số can thiệp sáng tạo** | **%d** | Chuỗi hành vi có ý thức của con người |\n", summary.TotalInterventions)
	fmt.Fprintf(&b, "| Tiền đề & Ý tưởng gốc | %d | Nguồn cảm hứng sáng tạo ban đầu |\n", summary.PremiseSeeds)
	fmt.Fprintf(&b, "| Can thiệp định hướng tình tiết | %d | Quyền kiểm soát mạch truyện (Creative Steering) |\n", summary.CreativeSteers)
	fmt.Fprintf(&b, "| Thẩm định & Phê duyệt đại cương | %d | Tuyển chọn và sắp đặt cấu trúc tác phẩm |\n", summary.OutlineCurations)
	fmt.Fprintf(&b, "| Lựa chọn ngã rẽ cốt truyện | %d | Quyết định nghệ thuật (Artistic Choice) |\n", summary.CreativeChoices)
	fmt.Fprintf(&b, "| Chỉnh sửa bản thảo trực tiếp | %d | Đóng góp trực tiếp vào văn bản cuối cùng |\n", summary.HumanEdits)
	fmt.Fprintf(&b, "| Nghiệm thu cổng từng chương | %d | Kiểm duyệt và cho phép xuất bản (Curation Gate) |\n\n", summary.ChapterApprovals)

	fmt.Fprintf(&b, "## 3. Nhật Ký Bằng Chứng Đóng Góp Chi Tiết (Audit Trail)\n\n")
	if len(records) == 0 {
		fmt.Fprintf(&b, "*Chưa có bản ghi nào được ghi nhận trong sổ cái.*\n")
	} else {
		for i, r := range records {
			chStr := ""
			if r.Chapter > 0 {
				chStr = fmt.Sprintf(" (Chương %d)", r.Chapter)
			}
			fmt.Fprintf(&b, "### %d. [%s] %s%s\n", i+1, r.Type, r.Description, chStr)
			fmt.Fprintf(&b, "- **Thời điểm**: `%s`\n", r.Timestamp)
			fmt.Fprintf(&b, "- **Mã chứng thực (SHA-256)**: `%s`\n", r.Hash)
			if r.InputText != "" {
				fmt.Fprintf(&b, "- **Nội dung can thiệp của tác giả**:\n```text\n%s\n```\n", strings.TrimSpace(r.InputText))
			}
			fmt.Fprintf(&b, "\n")
		}
	}

	fmt.Fprintf(&b, "---\n\n")
	fmt.Fprintf(&b, "*(Hồ sơ này có thể dùng làm tài liệu đính kèm khi đăng ký bản quyền tác phẩm văn học tại Cục Bản quyền tác giả hoặc cơ quan quản lý sở hữu trí tuệ tương đương).*\n")

	return b.String(), nil
}
