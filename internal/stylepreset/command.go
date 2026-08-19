package stylepreset

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli style`.
func Command(args []string) int {
	fs := flag.NewFlagSet("style", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	list := fs.Bool("list", false, "Xem danh sách các bộ phong cách Prompt (Style Presets)")
	switchStyle := fs.String("switch", "", "Chuyển nhanh sang bộ phong cách chỉ định (ví dụ: 'cyberpunk', 'xianxia', 'mystery')")
	createFromAB := fs.String("create-from-ab", "", "Tạo phong cách mới từ sở thích phản hồi A/B với tên đặt trước")
	tuneStyle := fs.String("tune", "", "Tinh chỉnh/Sửa lỗi phong cách hiện có bằng phản hồi A/B")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *switchStyle != "" {
		if err := mgr.SwitchStyle(*switchStyle); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi chuyển phong cách: %v\n", err)
			return 1
		}
		fmt.Printf("🎨 ĐÃ CHUYỂN PHONG CÁCH PROMPT THÀNH CÔNG SANG: [%s]\n", *switchStyle)
		fmt.Println("Văn phong mới đã được nạp tự động vào <dự án>/style/voice.md!")
		return 0
	}

	if *createFromAB != "" {
		if err := mgr.CreateStyleFromAB(*createFromAB); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tạo phong cách mới từ A/B: %v\n", err)
			return 1
		}
		fmt.Printf("✨ ĐÃ TẠO VÀ KÍCH HOẠT PHONG CÁCH MỚI TỪ PHẢN HỒI A/B: [%s]\n", *createFromAB)
		return 0
	}

	if *tuneStyle != "" {
		if err := mgr.TuneStyleWithAB(*tuneStyle); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tinh chỉnh phong cách: %v\n", err)
			return 1
		}
		fmt.Printf("🛠️ ĐÃ TINH CHỈNH THÀNH CÔNG PHONG CÁCH: [%s] KẾT HỢP DỮ LIỆU A/B!\n", *tuneStyle)
		return 0
	}

	if *list || len(args) == 0 {
		styles, err := mgr.ListStyles("cyberpunk")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi danh sách phong cách: %v\n", err)
			return 1
		}
		fmt.Println("🎨 DANH SÁCH BỘ PHONG CÁCH PROMPT (STYLE PRESETS) SẴN CÓ:")
		fmt.Println("============================================================")
		for _, s := range styles {
			activeFlag := "  "
			if s.Active {
				activeFlag = "► "
			}
			fmt.Printf("%s[%s] (Nguồn: %s) — %s\n", activeFlag, s.Name, s.Source, s.Description)
		}
		fmt.Println("\n💡 Mẹo: Dùng `ainovel-cli style --switch <tên-phong-cách>` để chuyển nhanh phong cách!")
		return 0
	}

	return 0
}
