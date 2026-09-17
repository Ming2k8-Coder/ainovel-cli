package server

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli serve`.
func Command(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	port := fs.Int("port", 8080, "Cổng HTTP của Web Reader")
	_ = fs.Parse(args)

	srv := NewServer(*dir, *port)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi khởi chạy Web Reader: %v\n", err)
		return 1
	}
	return 0
}
