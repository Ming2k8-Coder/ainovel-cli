package charvector

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli char-matrix` hoặc `ainovel-cli char-vector`.
func Command(args []string) int {
	fs := flag.NewFlagSet("char-matrix", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	charA := fs.String("a", "", "Tên nhân vật A")
	charB := fs.String("b", "", "Tên nhân vật B")
	trust := fs.Int("trust", 0, "Chỉ số tin tưởng (-100 đến +100)")
	tension := fs.Int("tension", 0, "Chỉ số căng thẳng (0 đến 100)")
	alliance := fs.String("alliance", "neutral", "Mối quan hệ: 'ally', 'neutral', 'rival', 'enemy'")
	note := fs.String("note", "", "Ghi chú diễn biến quan hệ")
	chapter := fs.Int("chapter", 1, "Chương cập nhật")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *charA != "" && *charB != "" {
		if err := mgr.UpdateVector(*charA, *charB, *trust, *tension, *alliance, *note, *chapter); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi cập nhật quan hệ nhân vật: %v\n", err)
			return 1
		}
		fmt.Printf("✅ Đã cập nhật thành công quan hệ giữa [%s] và [%s] tại Chương %d!\n", *charA, *charB, *chapter)
		return 0
	}

	mat, err := mgr.LoadMatrix()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi nạp ma trận nhân vật: %v\n", err)
		return 1
	}

	fmt.Println("👥 MA TRẬN PHÂN TÍCH QUAN HỆ & CHỈ SỐ TÂM LÝ NHÂN VẬT")
	fmt.Println("============================================================")
	if len(mat.Vectors) == 0 {
		fmt.Println("Chưa có mối quan hệ nào được thiết lập. Dùng `-a \"Nhân vật A\" -b \"Nhân vật B\" --trust 50` để tạo.")
		return 0
	}

	for _, v := range mat.Vectors {
		fmt.Printf("• %s ⚡ %s [Thân thiết: %d | Căng thẳng: %d | Loại: %s] (Ch.%d)\n  ➔ Ghi chú: %s\n",
			v.CharA, v.CharB, v.Trust, v.Tension, v.Alliance, v.LastUpdateCh, v.Note)
	}
	return 0
}
