package cryptoaudit

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Command thực thi lệnh `ainovel-cli audit-proof` hoặc `ainovel-cli merkle-proof`.
func Command(args []string) int {
	fs := flag.NewFlagSet("audit-proof", flag.ExitOnError)
	author := fs.String("author", "Human Author", "Tên tác giả / người định hướng nghệ thuật")
	bookTitle := fs.String("title", "Tác phẩm", "Tên tác phẩm")
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	outputFile := fs.String("output", "", "Đường dẫn lưu chứng thư (mặc định meta/LEGAL_AFFIDAVIT.md)")
	_ = fs.Parse(args)

	dossier, err := BuildDossier(*dir, *bookTitle, *author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi xây dựng hồ sơ chứng thực số học: %v\n", err)
		return 1
	}

	affidavit := dossier.ExportLegalAffidavit()

	targetPath := *outputFile
	if targetPath == "" {
		targetPath = filepath.Join(*dir, "meta", "LEGAL_AFFIDAVIT.md")
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi tạo thư mục: %v\n", err)
		return 1
	}

	if err := os.WriteFile(targetPath, []byte(affidavit), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi ghi chứng thư: %v\n", err)
		return 1
	}

	fmt.Printf("🛡️ ĐÃ TẠO THÀNH CÔNG CHỨNG THƯ BẢN QUYỀN MERKLE ROOT TẠI: %s\n", targetPath)
	fmt.Printf("• Merkle Root: 0x%s\n", dossier.MerkleRoot)
	fmt.Printf("• Tổng số chương: %d | Tổng số từ: %d | Điểm chạm: %d\n", dossier.TotalChapters, dossier.TotalWordCount, dossier.TotalTouchpoints)
	return 0
}
