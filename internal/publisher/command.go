package publisher

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/voocel/ainovel-cli/internal/store"
)

// Command thực thi lệnh `ainovel-cli publish` hoặc `ainovel-cli export-mdx`.
func Command(args []string) int {
	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	format := fs.String("format", "mdx", "Định dạng xuất bản: 'mdx' (từng chương có frontmatter) hoặc 'full-md' (toàn bộ tác phẩm thành 1 file)")
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	outputDir := fs.String("out", "./dist", "Thư mục xuất bản đầu ra")
	author := fs.String("author", "Human Author", "Tên tác giả")
	_ = fs.Parse(args)

	st := store.NewStore(*dir)

	bookTitle := filepath.Base(*dir)
	if book, err := st.Book.Load(); err == nil && book != nil && book.Title != "" {
		bookTitle = book.Title
	}

	if *format == "full-md" {
		targetFile := filepath.Join(*outputDir, fmt.Sprintf("%s_FullBook.md", bookTitle))
		if err := CompileFullBook(st, targetFile, bookTitle, *author); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi xuất bản FullBook: %v\n", err)
			return 1
		}
		fmt.Printf("📚 Đã xuất bản toàn bộ tác phẩm thành công tại: %s\n", targetFile)
		return 0
	}

	// Mặc định xuất bản MDX cho Astro / Next.js / sanbaka-web
	mdxDir := filepath.Join(*outputDir, "mdx")
	opts := ExportOptions{
		OutputDir: mdxDir,
		BookTitle: bookTitle,
		Author:    *author,
	}
	if err := CompileMDX(st, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi xuất bản MDX: %v\n", err)
		return 1
	}

	fmt.Printf("✨ Đã biên dịch toàn bộ tác phẩm sang MDX tương thích sanbaka-web tại: %s\n", mdxDir)
	return 0
}
