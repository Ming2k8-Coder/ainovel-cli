package universe

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli universe` hoặc `ainovel-cli mcu`.
func Command(args []string) int {
	fs := flag.NewFlagSet("universe", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	trace := fs.String("trace", "", "Tên nhân vật / bảo vật cần truy vết dòng thời gian toàn vũ trụ")
	record := fs.String("record", "", "Tên thực thể cần ghi nhận vết sự kiện mới")
	series := fs.String("series", "HUST-sanbaka", "Tên bộ truyện trong vũ trụ")
	chapter := fs.Int("chapter", 1, "Số chương")
	year := fs.Int("year", 2026, "Năm vũ trụ (Universe Year)")
	desc := fs.String("desc", "", "Mô tả sự kiện / cột mốc")
	loc := fs.String("loc", "ĐHBK Hà Nội", "Địa điểm diễn ra")
	checkCanon := fs.Bool("check-canon", false, "Kiểm tra xung đột tính nhất quán và nghịch lý thời gian toàn vũ trụ")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *trace != "" {
		rec, err := mgr.TraceEntity(*trace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi truy vết: %v\n", err)
			return 1
		}
		fmt.Printf("🌌 DÒNG THỜI GIAN VŨ TRỤ (MCU-SCALE TRACE) CỦA: [%s]\n", rec.Name)
		fmt.Println("============================================================")
		fmt.Printf("• Loại thực thể: %s | Bộ truyện xuất phát: %s | Trạng thái: %s\n\n", rec.Type, rec.OriginSeries, rec.Status)
		fmt.Println("📜 Lịch sử xuất hiện & Biến cố xuyên suốt các bộ truyện:")
		for i, h := range rec.History {
			fmt.Printf("  %d. [Năm %d] Bộ truyện: %s (Ch.%d) tại %s\n     ➔ %s\n",
				i+1, h.UniverseYear, h.SeriesName, h.Chapter, h.Location, h.Description)
		}
		return 0
	}

	if *record != "" {
		if err := mgr.RegisterEvent(*record, "character", *series, *chapter, *year, *desc, *loc); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi ghi nhận sự kiện vũ trụ: %v\n", err)
			return 1
		}
		fmt.Printf("✅ Đã ghi nhận thành công vết sự kiện của [%s] trong bộ truyện %s (Năm %d, Ch.%d)!\n", *record, *series, *year, *chapter)
		return 0
	}

	if *checkCanon {
		issues, err := mgr.AuditContinuity()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi kiểm tra Canon: %v\n", err)
			return 1
		}
		fmt.Println("🛡️ KẾT QUẢ KIỂM TRA TÍNH NHẤT QUÁN QUY LUẬT VŨ TRỤ (CANON CONTINUITY AUDIT)")
		fmt.Println("============================================================")
		if len(issues) == 0 {
			fmt.Println("✅ Không phát hiện nghịch lý thời gian hay xung đột Canon nào! Vũ trụ hoàn toàn nhất quán.")
			return 0
		}
		for _, issue := range issues {
			fmt.Println(issue)
		}
		return 0
	}

	manifest, err := mgr.LoadManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi nạp vũ trụ: %v\n", err)
		return 1
	}

	fmt.Printf("🌌 TỔNG QUAN VŨ TRỤ TRUYỆN: [%s]\n", manifest.UniverseName)
	fmt.Println("============================================================")
	fmt.Printf("• Danh sách bộ truyện liên kết (%d): %v\n", len(manifest.SeriesList), manifest.SeriesList)
	fmt.Printf("• Số lượng thực thể đã theo vết: %d thực thể\n", len(manifest.Entities))
	fmt.Println("\n📜 Định luật Vũ trụ Tối cao (Canon Rules):")
	for _, r := range manifest.CanonRules {
		fmt.Printf("  - [%s] %s: %s\n", r.Severity, r.RuleName, r.Description)
	}
	fmt.Println("\n💡 Mẹo: Dùng `ainovel-cli universe --trace \"Tên Nhân Vật\"` để truy vết dòng thời gian kiểu MCU!")
	return 0
}
