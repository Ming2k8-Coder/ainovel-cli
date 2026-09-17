package foreshadow

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli foreshadow`.
func Command(args []string) int {
	fs := flag.NewFlagSet("foreshadow", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	add := fs.String("add", "", "Nội dung phục bút mới cần gài vào")
	chapter := fs.Int("chapter", 1, "Chương liên quan")
	category := fs.String("category", "mystery", "Loại phục bút: 'mystery', 'character_motive', 'item_secret', 'world_lore'")
	resolve := fs.String("resolve", "", "Mã ID của phục bút cần đánh dấu đã giải quyết")
	note := fs.String("note", "", "Ghi chú giải đáp")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *add != "" {
		item, err := mgr.AddItem(*chapter, *add, *category)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi thêm phục bút: %v\n", err)
			return 1
		}
		fmt.Printf("✅ Đã gài phục bút thành công tại Chương %d [ID: %s]: %s\n", item.ChapterPlanted, item.ID, item.Description)
		return 0
	}

	if *resolve != "" {
		if err := mgr.ResolveItem(*resolve, *chapter, *note); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tháo gỡ phục bút: %v\n", err)
			return 1
		}
		fmt.Printf("🎉 Đã tháo gỡ thành công phục bút %s tại Chương %d!\n", *resolve, *chapter)
		return 0
	}

	mat, err := mgr.LoadMatrix()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi nạp ma trận phục bút: %v\n", err)
		return 1
	}

	fmt.Println("📜 MẠNG LƯỚI PHỤC BÚT & BÍ ẨN CỐT TRUYỆN")
	fmt.Println("============================================================")
	if len(mat.Items) == 0 {
		fmt.Println("Chưa có phục bút nào được ghi nhận. Dùng `--add \"Nội dung\" --chapter 1` để gài thêm.")
		return 0
	}

	unresolvedCount := 0
	for _, it := range mat.Items {
		status := "⏳ Chưa giải đáp"
		if it.Resolved {
			status = fmt.Sprintf("✅ Đã tháo gỡ ở Ch.%d (%s)", it.ChapterResolved, it.ResolutionNote)
		} else {
			unresolvedCount++
		}
		fmt.Printf("[%s] Ch.%d | %s | %s\n    ➔ %s\n", it.ID[:12], it.ChapterPlanted, stringsToUpper(it.Category), status, it.Description)
	}

	fmt.Printf("\n📊 Tổng số phục bút: %d | Chưa tháo gỡ: %d\n", len(mat.Items), unresolvedCount)
	return 0
}

func stringsToUpper(s string) string {
	if s == "mystery" { return "BÍ ẨN" }
	if s == "character_motive" { return "ĐỘNG CƠ NHÂN VẬT" }
	if s == "item_secret" { return "BẢO VẬT/CÔNG NGHỆ" }
	if s == "world_lore" { return "BỐI CẢNH THẾ GIỚI" }
	return s
}
