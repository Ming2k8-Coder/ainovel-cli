// Package publisher cung cấp công cụ xuất bản tiểu thuyết sang đa định dạng:
// MDX (cho web tĩnh Astro/Next.js/sanbaka-web), EPUB (cho thiết bị đọc sách),
// và Sách điện tử Markdown hoàn chỉnh.
package publisher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/voocel/ainovel-cli/internal/store"
)

// ExportOptions tùy chọn xuất bản tác phẩm.
type ExportOptions struct {
	OutputDir   string // Thư mục xuất ra (mặc định ./dist)
	Format      string // "mdx", "md", "epub"
	BookTitle   string
	Author      string
	Description string
}

// CompileMDX xuất bản toàn bộ tác phẩm thành các tệp `.mdx` tương thích với sanbaka-web / Next.js / Astro.
func CompileMDX(st *store.Store, opts ExportOptions) error {
	if opts.OutputDir == "" {
		opts.OutputDir = "./dist/mdx"
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return err
	}

	progress, err := st.Progress.Load()
	if err != nil {
		return err
	}
	if progress == nil {
		return fmt.Errorf("không tìm thấy dữ liệu tiến độ tác phẩm")
	}

	for _, ch := range progress.CompletedChapters {
		content, _, err := st.Drafts.LoadChapterContent(ch)
		if err != nil {
			continue
		}
		summary, _ := st.Summaries.LoadSummary(ch)
		title := fmt.Sprintf("Chương %d", ch)
		if summary != nil && summary.Title != "" {
			title = fmt.Sprintf("Chương %d: %s", ch, summary.Title)
		}

		var b strings.Builder
		b.WriteString("---\n")
		fmt.Fprintf(&b, "title: %q\n", title)
		fmt.Fprintf(&b, "chapter: %d\n", ch)
		fmt.Fprintf(&b, "author: %q\n", opts.Author)
		fmt.Fprintf(&b, "wordCount: %d\n", len([]rune(content)))
		fmt.Fprintf(&b, "publishedAt: %q\n", time.Now().Format("2006-01-02"))
		if summary != nil {
			fmt.Fprintf(&b, "summary: %q\n", summary.Summary)
		}
		b.WriteString("---\n\n")
		fmt.Fprintf(&b, "# %s\n\n", title)
		b.WriteString(content)
		b.WriteString("\n")

		targetFile := filepath.Join(opts.OutputDir, fmt.Sprintf("chapter-%03d.mdx", ch))
		if err := os.WriteFile(targetFile, []byte(b.String()), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// CompileFullBook xuất bản toàn bộ tác phẩm thành một tệp Markdown duy nhất hoàn chỉnh.
func CompileFullBook(st *store.Store, targetPath string, bookTitle, author string) error {
	progress, err := st.Progress.Load()
	if err != nil {
		return err
	}
	if progress == nil {
		return fmt.Errorf("không tìm thấy dữ liệu tiến độ tác phẩm")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", bookTitle)
	fmt.Fprintf(&b, "> **Tác giả**: %s\n", author)
	fmt.Fprintf(&b, "> **Tổng số chương**: %d chương\n", len(progress.CompletedChapters))
	fmt.Fprintf(&b, "> **Tổng số từ**: %d từ\n", progress.TotalWordCount)
	fmt.Fprintf(&b, "> **Ngày hoàn thiện**: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "---\n\n")

	for _, ch := range progress.CompletedChapters {
		content, _, err := st.Drafts.LoadChapterContent(ch)
		if err != nil {
			continue
		}
		summary, _ := st.Summaries.LoadSummary(ch)
		title := fmt.Sprintf("Chương %d", ch)
		if summary != nil && summary.Title != "" {
			title = fmt.Sprintf("Chương %d: %s", ch, summary.Title)
		}

		fmt.Fprintf(&b, "## %s\n\n", title)
		b.WriteString(content)
		b.WriteString("\n\n---\n\n")
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(targetPath, []byte(b.String()), 0o644)
}
