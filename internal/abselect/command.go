package abselect

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/voocel/ainovel-cli/internal/authorship"
)

// Command thực thi lệnh dòng lệnh `ainovel-cli ab-feedback` hoặc `ainovel-cli ab-select`.
func Command(args []string) int {
	fs := flag.NewFlagSet("ab-feedback", flag.ExitOnError)
	outputDir := fs.String("dir", "./novel", "Thư mục chứa dữ liệu tác phẩm")
	apply := fs.Bool("apply", true, "Tự động áp dụng sở thích đã học vào style/voice.md của dự án")
	_ = fs.Parse(args)

	ledger := authorship.NewLedger(*outputDir)
	mgr := NewManager(*outputDir, ledger)

	synth, err := mgr.SynthesizePreferences()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi phân tích sở thích A/B: %v\n", err)
		return 1
	}

	fmt.Println("📊 KẾT QUẢ PHÂN TÍCH PHẢN HỒI LỰA CHỌN A/B TỪ TÁC GIẢ")
	fmt.Println("============================================================")
	fmt.Printf("• Đặc trưng ưa chuộng (%d): %v\n", len(synth.PreferredTraits), synth.PreferredTraits)
	fmt.Printf("• Đặc trưng cần tránh (%d): %v\n\n", len(synth.AvoidedTraits), synth.AvoidedTraits)

	if len(synth.ProseDirectives) > 0 {
		fmt.Println("📝 Chỉ thị văn phong ưu tiên (Được bổ sung vào System Prompt):")
		for _, d := range synth.ProseDirectives {
			fmt.Printf("  + %s\n", d)
		}
		fmt.Println()
	}

	if len(synth.Taboos) > 0 {
		fmt.Println("🚫 Ràng buộc cấm kỵ (Được bổ sung vào System Prompt):")
		for _, t := range synth.Taboos {
			fmt.Printf("  - %s\n", t)
		}
		fmt.Println()
	}

	if *apply {
		if err := mgr.ApplyFeedbackToProject(); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi áp dụng vào style/voice.md: %v\n", err)
			return 1
		}
		voicePath := filepath.Join(*outputDir, "style", "voice.md")
		fmt.Printf("✅ Đã tự động cập nhật System Prompts & Voice Rules tại: %s\n", voicePath)
	}

	return 0
}
