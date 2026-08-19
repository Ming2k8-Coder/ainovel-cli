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

	// Branching & Checkout
	createBranch := fs.String("create-branch", "", "Tạo nhánh Git mới")
	checkoutBranch := fs.String("checkout", "", "Chuyển sang nhánh Git chỉ định")
	listBranches := fs.Bool("list-branches", false, "Liệt kê danh sách các nhánh Git")

	// Issues Management
	createIssue := fs.String("create-issue", "", "Tạo một Issue mới")
	issueDesc := fs.String("desc", "", "Mô tả chi tiết cho Issue / PR")
	startIssue := fs.String("start-issue", "", "Bắt đầu thực hiện Issue theo mã ID")
	makeFeat := fs.Bool("make-feat", true, "Tạo nhánh feat/ mới khi thực hiện Issue (nếu false sẽ thực hiện trực tiếp trên main)")
	listIssues := fs.Bool("list-issues", false, "Xem danh sách các Issues")

	// Pull Request & Review
	createPR := fs.String("create-pr", "", "Mở một Pull Request mới với tiêu đề")
	sourceBranch := fs.String("source", "", "Nhánh nguồn của PR (mặc định là nhánh hiện tại)")
	targetBranch := fs.String("target", "main", "Nhánh đích của PR (mặc định là main)")
	reviewPR := fs.String("review-pr", "", "Thẩm định/Đánh giá Pull Request theo mã ID")
	approve := fs.Bool("approve", true, "Phê duyệt PR (nếu false là từ chối)")
	comment := fs.String("comment", "", "Lời bình luận thẩm định cho PR")
	mergePR := fs.String("merge-pr", "", "Gộp nhánh Pull Request đã duyệt vào nhánh đích")
	listPRs := fs.Bool("list-prs", false, "Xem danh sách các Pull Requests")

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

	// 1. Quản Lý Nhánh (Branching & Checkout)
	if *createBranch != "" {
		if err := mgr.CreateBranch(*createBranch); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tạo nhánh: %v\n", err)
			return 1
		}
		fmt.Printf("🌿 Đã tạo thành công nhánh Git mới: %s\n", *createBranch)
		return 0
	}

	if *checkoutBranch != "" {
		if err := mgr.Checkout(*checkoutBranch, false); err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi chuyển nhánh: %v\n", err)
			return 1
		}
		fmt.Printf("🔀 Đã chuyển thành công sang nhánh Git: %s\n", *checkoutBranch)
		return 0
	}

	if *listBranches {
		branches, current, err := mgr.ListBranches()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi danh sách nhánh: %v\n", err)
			return 1
		}
		fmt.Println("🌿 DANH SÁCH CÁC NHÁNH GIT (GIT BRANCHES):")
		for _, b := range branches {
			if b == current {
				fmt.Printf("  * %s (Đang làm việc)\n", b)
			} else {
				fmt.Printf("    %s\n", b)
			}
		}
		return 0
	}

	// 2. Quản Lý Issues
	if *createIssue != "" {
		issue, err := mgr.CreateIssue(*createIssue, *issueDesc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tạo Issue: %v\n", err)
			return 1
		}
		fmt.Printf("📌 Đã tạo Issue mới [%s]: %s\n", issue.ID, issue.Title)
		return 0
	}

	if *startIssue != "" {
		issue, err := mgr.StartIssue(*startIssue, *makeFeat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi bắt đầu Issue: %v\n", err)
			return 1
		}
		fmt.Printf("🚀 Đã chuyển Issue [%s] sang IN_PROGRESS. Nhánh hoạt động: %s\n", issue.ID, issue.BranchName)
		return 0
	}

	if *listIssues {
		issues, err := mgr.LoadIssues()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi danh sách Issues: %v\n", err)
			return 1
		}
		fmt.Println("📌 DANH SÁCH ISSUES NỘI DUNG TÁC PHẨM:")
		if len(issues) == 0 {
			fmt.Println("Chưa có Issue nào. Dùng `--create-issue \"Tiêu đề\"` để tạo.")
			return 0
		}
		for _, is := range issues {
			fmt.Printf("• [%s] %s | Trạng thái: %s | Nhánh: %s\n  ➔ Mô tả: %s\n",
				is.ID, is.Title, is.Status, is.BranchName, is.Description)
		}
		return 0
	}

	// 3. Quản Lý Pull Requests (PR Review)
	if *createPR != "" {
		pr, err := mgr.CreatePR(*createPR, *sourceBranch, *targetBranch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi tạo Pull Request: %v\n", err)
			return 1
		}
		fmt.Printf("🔍 Đã mở Pull Request [%s]: %s (%s -> %s)\n", pr.ID, pr.Title, pr.SourceBranch, pr.TargetBranch)
		return 0
	}

	if *reviewPR != "" {
		pr, err := mgr.ReviewPR(*reviewPR, *approve, "Editor / Author", *comment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi thẩm định PR: %v\n", err)
			return 1
		}
		fmt.Printf("✍️ Đã thẩm định PR [%s] | Trạng thái: %s | Nhận xét: %s\n", pr.ID, pr.Status, pr.ReviewComment)
		return 0
	}

	if *mergePR != "" {
		pr, err := mgr.MergePR(*mergePR)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi gộp nhánh PR: %v\n", err)
			return 1
		}
		fmt.Printf("🔀 Đã gộp thành công PR [%s] (%s -> %s) vào %s!\n", pr.ID, pr.SourceBranch, pr.TargetBranch, pr.TargetBranch)
		return 0
	}

	if *listPRs {
		prs, err := mgr.LoadPRs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lỗi danh sách PRs: %v\n", err)
			return 1
		}
		fmt.Println("🔍 DANH SÁCH PULL REQUESTS & THẨM ĐỊNH NỘI DUNG:")
		if len(prs) == 0 {
			fmt.Println("Chưa có PR nào. Dùng `--create-pr \"Tiêu đề\"` để tạo.")
			return 0
		}
		for _, pr := range prs {
			fmt.Printf("• [%s] %s (%s -> %s) | Trạng thái: %s | Người duyệt: %s\n  ➔ Nhận xét: %s\n",
				pr.ID, pr.Title, pr.SourceBranch, pr.TargetBranch, pr.Status, pr.Reviewer, pr.ReviewComment)
		}
		return 0
	}

	// Mặc định: Commit hoặc Status
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

	fmt.Println("🔧 QUẢN LÝ PHIÊN BẢN GIT (GIT BRANCHING, ISSUES & PR REVIEW)")
	fmt.Println("============================================================")
	fmt.Println(status)
	fmt.Println("\n📜 Lịch sử Commit gần đây:")
	hist, _ := mgr.History()
	fmt.Println(hist)
	return 0
}
