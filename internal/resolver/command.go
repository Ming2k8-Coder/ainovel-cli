package resolver

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli resolve` hoặc `ainovel-cli fix-logic`.
func Command(args []string) int {
	fs := flag.NewFlagSet("resolve", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	scan := fs.Bool("scan", true, "Quét phát hiện mâu thuẫn logic và lỗ hổng cốt truyện")
	fix := fs.String("fix", "", "Mã ID của mâu thuẫn logic cần tự động khắc phục")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *fix != "" {
		if err := mgr.ApplyFixPlan(*fix); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi thực thi khắc phục logic: %v\n", err)
			return 1
		}
		fmt.Printf("🛠️ Đã áp dụng thành công phương án sửa lỗi logic [%s] vào hàng đợi PendingRewrites!\n", *fix)
		return 0
	}

	if *scan {
		issues, err := mgr.ScanParadoxes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi quét mâu thuẫn logic: %v\n", err)
			return 1
		}

		fmt.Println("🔍 BỘ PHÁT HIỆN MÂU THUẪN LOGIC CỐT TRUYỆN (PLOT PARADOX RESOLVER)")
		fmt.Println("============================================================")
		if len(issues) == 0 {
			fmt.Println("✅ Không phát hiện mâu thuẫn logic hay lỗ hổng cốt truyện nào! Tác phẩm đạt độ nhất quán cao.")
			return 0
		}

		for _, is := range issues {
			fmt.Printf("• [%s] Chương %d (%s) | Mức độ: %s\n  ➔ Mô tả: %s\n  ➔ Phương án sửa: %s\n\n",
				is.ID, is.Chapter, is.Type, is.Severity, is.Description, is.FixPlan)
		}
		fmt.Printf("💡 Mẹo: Chạy `ainovel-cli resolve --fix <ID>` để tự động áp dụng phương án sửa đổi!\n")
	}

	return 0
}
