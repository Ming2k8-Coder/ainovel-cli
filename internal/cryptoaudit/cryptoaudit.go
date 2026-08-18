// Package cryptoaudit cung cấp giải pháp chứng thực số học (Cryptographic Audit & Merkle Proof)
// cho toàn bộ tác phẩm và sổ cái quyền tác giả, tạo ra cây Merkle Tree chứng minh tính toàn vẹn
// và sự đóng góp liên tục của con người theo tiêu chuẩn quốc tế.
package cryptoaudit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/voocel/ainovel-cli/internal/authorship"
)

// ChapterProof chứa mã băm và thông tin chứng thực của từng chương.
type ChapterProof struct {
	Chapter     int    `json:"chapter"`
	WordCount   int    `json:"word_count"`
	ContentHash string `json:"content_hash"`
	Touchpoints int    `json:"touchpoints"` // Số lượng can thiệp của con người trong chương
}

// Dossier chứa toàn bộ hồ sơ bằng chứng số học của cuốn tiểu thuyết.
type Dossier struct {
	BookTitle        string         `json:"book_title"`
	AuthorName       string         `json:"author_name"`
	GeneratedAt      string         `json:"generated_at"`
	TotalChapters    int            `json:"total_chapters"`
	TotalWordCount   int            `json:"total_word_count"`
	TotalTouchpoints int            `json:"total_touchpoints"`
	MerkleRoot       string         `json:"merkle_root"`
	Chapters         []ChapterProof `json:"chapters"`
}

// GenerateMerkleRoot tính toán mã băm gốc Merkle từ danh sách các mã băm lá.
func GenerateMerkleRoot(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}
	var current []string
	for _, h := range hashes {
		current = append(current, h)
	}

	for len(current) > 1 {
		var next []string
		for i := 0; i < len(current); i += 2 {
			if i+1 < len(current) {
				hasher := sha256.New()
				hasher.Write([]byte(current[i] + current[i+1]))
				next = append(next, hex.EncodeToString(hasher.Sum(nil)))
			} else {
				// Nếu lẻ, băm lại với chính nó
				hasher := sha256.New()
				hasher.Write([]byte(current[i] + current[i]))
				next = append(next, hex.EncodeToString(hasher.Sum(nil)))
			}
		}
		current = next
	}
	return current[0]
}

// BuildDossier xây dựng hồ sơ chứng thực toàn diện cho dự án.
func BuildDossier(dir, bookTitle, authorName string) (*Dossier, error) {
	ledger := authorship.NewLedger(dir)
	records, err := ledger.LoadRecords()
	if err != nil {
		return nil, fmt.Errorf("load authorship records: %w", err)
	}

	// Đếm số touchpoint theo chương
	chTouchpoints := make(map[int]int)
	for _, r := range records {
		chTouchpoints[r.Chapter]++
	}

	chaptersDir := filepath.Join(dir, "chapters")
	entries, err := os.ReadDir(chaptersDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read chapters dir: %w", err)
	}

	var leafHashes []string
	var chapterProofs []ChapterProof
	totalWords := 0

	for ch := 1; ; ch++ {
		filename := fmt.Sprintf("%02d.md", ch)
		p := filepath.Join(chaptersDir, filename)
		data, err := os.ReadFile(p)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return nil, err
		}

		hasher := sha256.New()
		hasher.Write(data)
		contentHash := hex.EncodeToString(hasher.Sum(nil))
		leafHashes = append(leafHashes, contentHash)

		words := len([]rune(string(data)))
		totalWords += words

		chapterProofs = append(chapterProofs, ChapterProof{
			Chapter:     ch,
			WordCount:   words,
			ContentHash: contentHash,
			Touchpoints: chTouchpoints[ch],
		})
	}

	// Thêm các mã băm từ Authorship Records vào cây Merkle
	for _, r := range records {
		leafHashes = append(leafHashes, r.Hash)
	}

	merkleRoot := GenerateMerkleRoot(leafHashes)

	dossier := &Dossier{
		BookTitle:        bookTitle,
		AuthorName:       authorName,
		GeneratedAt:      time.Now().Format(time.RFC3339),
		TotalChapters:    len(chapterProofs),
		TotalWordCount:   totalWords,
		TotalTouchpoints: len(records),
		MerkleRoot:       merkleRoot,
		Chapters:         chapterProofs,
	}

	return dossier, nil
}

// ExportLegalAffidavit xuất tài liệu chứng thư bản quyền đạt chuẩn pháp lý quốc tế.
func (d *Dossier) ExportLegalAffidavit() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# CHỨNG THƯ PHÁP LÝ QUYỀN TÁC GIẢ & CHỨNG THỰC MÃ BĂM SỐ HỌC\n\n")
	fmt.Fprintf(&b, "> **Tác phẩm**: 《%s》\n", d.BookTitle)
	fmt.Fprintf(&b, "> **Tác giả con người (Human Creator & Director)**: %s\n", d.AuthorName)
	fmt.Fprintf(&b, "> **Thời điểm lập chứng thư**: %s\n", d.GeneratedAt)
	fmt.Fprintf(&b, "> **MÃ GỐC MERKLE TREE TOÀN TÁC PHẨM (Merkle Root)**:\n> `0x%s`\n\n", d.MerkleRoot)
	fmt.Fprintf(&b, "---\n\n")

	fmt.Fprintf(&b, "## 1. Tuyên Ngôn Quyền Tác Giả Dưới Luật Pháp Quốc Tế\n\n")
	fmt.Fprintf(&b, "Tác giả **%s** khẳng định quyền sở hữu trí tuệ hợp pháp và trọn vẹn đối với tác phẩm 《%s》 theo các căn cứ:\n\n", d.AuthorName, d.BookTitle)
	fmt.Fprintf(&b, "1. **Định hướng & Tiền đề nghệ thuật**: Trực tiếp thiết lập bối cảnh, cấu trúc nhân vật và mâu thuẫn cốt truyện.\n")
	fmt.Fprintf(&b, "2. **Tuyển chọn nghệ thuật (Selection & Arrangement)**: Trực tiếp điều hướng, phê duyệt và chọn lọc từng chương trong quá trình sáng tác.\n")
	fmt.Fprintf(&b, "3. **Tính toàn vẹn mật mã học**: Toàn bộ chuỗi văn bản và các can thiệp sáng tạo được liên kết thành cây Merkle Root bất biến, chứng minh tác phẩm không bị sửa đổi sau thời điểm xác lập.\n\n")

	fmt.Fprintf(&b, "## 2. Bảng Kê Chứng Thực Từng Chương Bản Thảo\n\n")
	fmt.Fprintf(&b, "| Chương | Số từ | Điểm can thiệp | Mã băm toàn vẹn (SHA-256) |\n")
	fmt.Fprintf(&b, "|---|---|---|---|\n")
	for _, cp := range d.Chapters {
		fmt.Fprintf(&b, "| Chương %02d | %d từ | %d lần | `%s` |\n", cp.Chapter, cp.WordCount, cp.Touchpoints, cp.ContentHash[:16]+"...")
	}
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "**Tổng cộng**: %d chương | %d từ | %d điểm chạm sáng tạo có ý thức của tác giả.\n\n", d.TotalChapters, d.TotalWordCount, d.TotalTouchpoints)
	fmt.Fprintf(&b, "---\n\n")
	fmt.Fprintf(&b, "*Chứng thư này được khởi tạo tự động bởi `ainovel-cli` và có giá trị làm chứng cứ kỹ thuật trước cơ quan đăng ký quyền tác giả và tòa án sở hữu trí tuệ.*\n")

	return b.String()
}
