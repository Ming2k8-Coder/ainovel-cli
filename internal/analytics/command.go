package analytics

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli stats` hoặc `ainovel-cli analytics`.
func Command(args []string) int {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	_ = fs.Parse(args)

	if err := PrintReport(*dir); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi tạo bảng thống kê analytics: %v\n", err)
		return 1
	}
	return 0
}
