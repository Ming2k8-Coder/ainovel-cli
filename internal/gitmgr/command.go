package gitmgr

import (
	"flag"
	"fmt"
	"os"
)

// Command thực thi lệnh `ainovel-cli git`.
func Command(args []string) int {
	fs := flag.NewFlagSet("git", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	init := fs.Bool("init", false, "Khởi tạo Git Repository")
	commitMsg := fs.String("commit", "", "Tự động commit với thông điệp")
	tag := fs.String("tag", "", "Tạo nhãn Git Tag")
	_ = fs.Parse(args)

	mgr := NewManager(*dir)

	if *init {
		if err := mgr.InitRepo(); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi khởi tạo Git: %v\n", err)
			return 1
		}
		fmt.Println("✅ Đã khởi tạo Git Repository thành công!")
		return 0
	}

	if *commitMsg != "" {
		if _, err := mgr.runGit("add", "."); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi git add: %v\n", err)
			return 1
		}
		out, err := mgr.runGit("commit", "-m", *commitMsg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi git commit: %v\n", err)
			return 1
		}
		fmt.Printf("✅ Đã commit thành công: %s\n", out)
		return 0
	}

	if *tag != "" {
		if err := mgr.CreateTag(*tag, "Milestone tag "+*tag); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tạo Tag: %v\n", err)
			return 1
		}
		fmt.Printf("🏷️ Đã tạo nhãn Git Tag thành công: %s\n", *tag)
		return 0
	}

	status, err := mgr.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi kiểm tra trạng thái Git: %v\n", err)
		return 1
	}

	fmt.Println("🔧 QUẢN LÝ PHIÊN BẢN GIT (GIT VERSION CONTROL)")
	fmt.Println("============================================================")
	fmt.Println(status)
	fmt.Println("\n📜 Lịch sử Commit gần đây:")
	hist, _ := mgr.History()
	fmt.Println(hist)
	return 0
}
