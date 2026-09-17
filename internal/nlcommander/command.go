package nlcommander

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command thực thi lệnh `ainovel-cli cmd` hoặc `ainovel-cli do`.
func Command(args []string) int {
	fs := flag.NewFlagSet("cmd", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	_ = fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if strings.TrimSpace(prompt) == "" {
		fmt.Fprintln(os.Stderr, "Lỗi: Vui lòng nhập câu lệnh ngôn ngữ tự nhiên (Ví dụ: `ainovel-cli cmd \"Phân tích truyện HUST-sanbaka và mở WebUI\"`)")
		return 1
	}

	cmder := NewCommander(*dir)
	intent := cmder.ParseIntent(prompt)

	if err := cmder.Execute(intent); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi thực thi lệnh tự nhiên: %v\n", err)
		return 1
	}
	return 0
}
