package authorship

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Command thực thi lệnh dòng lệnh `ainovel-cli copyright`.
func Command(args []string) int {
	fs := flag.NewFlagSet("copyright", flag.ExitOnError)
	author := fs.String("author", "Human Author", "Tên tác giả / chủ sở hữu bản quyền")
	outputDir := fs.String("dir", "./novel", "Thư mục chứa dữ liệu tác phẩm")
	outputFile := fs.String("output", "", "Đường dẫn xuất tệp báo cáo (mặc định in ra màn hình hoặc lưu vào meta/COPYRIGHT_AUTHORSHIP_REPORT.md)")
	_ = fs.Parse(args)

	ledger := NewLedger(*outputDir)
	report, err := ledger.GenerateCopyrightReport(filepath.Base(*outputDir), *author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi tạo hồ sơ bản quyền: %v\n", err)
		return 1
	}

	targetPath := *outputFile
	if targetPath == "" {
		targetPath = filepath.Join(*outputDir, "meta", "COPYRIGHT_AUTHORSHIP_REPORT.md")
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi tạo thư mục đích: %v\n", err)
		return 1
	}

	if err := os.WriteFile(targetPath, []byte(report), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi ghi tệp báo cáo: %v\n", err)
		return 1
	}

	fmt.Printf("✅ Đã xuất thành công Hồ sơ Chứng nhận Quyền tác giả tại: %s\n", targetPath)
	fmt.Println("------------------------------------------------------------")
	fmt.Println(report)
	return 0
}
