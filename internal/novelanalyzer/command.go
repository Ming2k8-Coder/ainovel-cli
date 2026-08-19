package novelanalyzer

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli analyze` hoặc `ainovel-cli continue`.
func Command(args []string) int {
	fs := flag.NewFlagSet("continue", flag.ExitOnError)
	source := fs.String("source", "", "Đường dẫn tệp/thư mục truyện hoặc file đặc tả có sẵn (ví dụ: ../servers/HUST-sanbaka)")
	targetDir := fs.String("dir", "./novel", "Thư mục tạo dự án viết tiếp")
	_ = fs.Parse(args)

	if *source == "" {
		fmt.Fprintln(os.Stderr, "Lỗi: Vui lòng cung cấp đường dẫn truyện gốc bằng cờ `--source <đường-dẫn>`")
		return 1
	}

	fmt.Printf("🔍 Đang tiến hành phân tích tác phẩm mẫu tại: %s...\n", *source)
	analyzer := NewAnalyzer(*source)
	res, err := analyzer.AnalyzePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi phân tích: %v\n", err)
		return 1
	}

	fmt.Println("============================================================")
	fmt.Printf("📖 Tên tác phẩm trích xuất: %s\n", res.BookTitle)
	fmt.Printf("📊 Số chương phát hiện: %d | Số từ: %d từ\n", res.DetectedChapters, res.TotalWords)
	fmt.Printf("👥 Nhân vật trích xuất (%d): %v\n", len(res.Characters), res.Characters)
	fmt.Println("============================================================")

	if err := SetupContinuationProject(*targetDir, res); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi thiết lập dự án viết tiếp: %v\n", err)
		return 1
	}

	fmt.Printf("✅ DỰ ÁN VIẾT TIẾP ĐÃ ĐƯỢC THIẾT LẬP THÀNH CÔNG TẠI: %s\n", *targetDir)
	fmt.Printf("Tác tử AI sẵn sàng viết tiếp từ Chương %d!\n", res.LastChapterNum+1)
	return 0
}
