package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/abselect"
	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/entry/headless"
	"github.com/voocel/ainovel-cli/internal/entry/startup"
	"github.com/voocel/ainovel-cli/internal/entry/tui"
	"github.com/voocel/ainovel-cli/internal/eval"
	"github.com/voocel/ainovel-cli/internal/rules"
	buildversion "github.com/voocel/ainovel-cli/internal/version"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// headlessMode Ghi lại xem lần này có khởi động chế độ headless hay không, giúp die quyết định dừng màn hình khi thoát lỗi.
var headlessMode bool

func main() {
	// Lệnh con (subcommand) được chặn trước khi phân tích cờ (flag) thông thường:
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "eval":
			os.Exit(eval.Command(os.Args[2:]))
		case "copyright", "authorship":
			os.Exit(authorship.Command(os.Args[2:]))
		case "ab-feedback", "ab-select":
			os.Exit(abselect.Command(os.Args[2:]))
		}
	}

	opts, args, err := parseCLIOptions(os.Args[1:])
	if err != nil {
		die("cờ (flags): %v", err)
	}
	if opts.Version {
		buildversion.Print(os.Stdout, versionInfo())
		return
	}
	if opts.Update {
		if err := runSelfUpdate(opts.UpdateVersion); err != nil {
			fmt.Fprintf(os.Stderr, "cập nhật: %v\n", err)
			os.Exit(1)
		}
		return
	}
	headlessMode = opts.Headless

	// Hướng dẫn cài đặt lần đầu
	if bootstrap.NeedsSetup() {
		if opts.Headless {
			die("lỗi: chế độ headless không hỗ trợ hướng dẫn cài đặt lần đầu, vui lòng chạy giao diện TUI một lần để hoàn tất cấu hình")
		}
		setupCfg, err := bootstrap.RunSetup()
		if err != nil {
			die("cài đặt (setup): %v", err)
		}
		// Sau khi hoàn tất cài đặt, tiếp tục với cấu hình được tạo
		runWithConfig(setupCfg, opts, args)
		return
	}

	// Tải cấu hình
	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		die("cấu hình (config): %v", err)
	}

	runWithConfig(cfg, opts, args)
}

// die Xử lý tập trung thoát lỗi nghiêm trọng: in ra stderr, ghi vào đĩa tại ~/.ainovel/last-error.log,
// và tạm dừng chờ Enter trong terminal tương tác (không phải headless) -- khi nhấp kép khởi động trên console,
// nó sẽ đóng ngay cùng tiến trình; nếu không tạm dừng, lỗi sẽ biến mất quá nhanh khiến người dùng không thể kiểm tra.
func die(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, msg)
	if path := bootstrap.WriteStartupError(msg); path != "" {
		fmt.Fprintf(os.Stderr, "（Lỗi chi tiết đã được ghi vào %s）\n", path)
	}
	if !headlessMode && stdinIsTerminal() {
		fmt.Fprint(os.Stderr, "\nNhấn phím Enter để thoát...")
		fmt.Fscanln(os.Stdin)
	}
	os.Exit(1)
}

// stdinIsTerminal Kiểm tra xem đầu vào chuẩn (stdin) có kết nối với terminal (thiết bị ký tự) hay không.
// Nhấp kép khởi động / terminal tương tác là true; ống dẫn (pipe), chuyển hướng, CI là false.
func stdinIsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func runWithConfig(cfg bootstrap.Config, opts cliOptions, args []string) {
	rules.EnsureHomeRulesDir()

	if len(args) > 0 {
		die("lỗi: không còn hỗ trợ truyền trực tiếp yêu cầu tiểu thuyết từ dòng lệnh, vui lòng nhập trong ô nhập liệu TUI sau khi khởi động")
	}

	// FillDefaults phải thực hiện trước khi tải tài nguyên: OutputDir là trường thời gian chạy, giá trị mặc định được quy chuẩn ở đây --
	// nếu không, trong cấu hình mặc định, việc ghi đè văn phong cấp sách tại <Thư mục sách>/style/ sẽ không bao giờ được tải.
	cfg.FillDefaults()
	bundle := assets.Load(cfg.Style, assets.DefaultLoadOptions(cfg.OutputDir))
	if opts.Headless {
		prompt, err := loadPrompt(opts)
		if err != nil {
			die("lỗi: %v", err)
		}
		if err := headless.Run(cfg, bundle, headless.Options{Prompt: prompt}); err != nil {
			die("lỗi: %v", err)
		}
		return
	}
	if opts.Prompt != "" || opts.PromptFile != "" {
		die("lỗi: --prompt/--prompt-file chỉ có thể sử dụng trong chế độ --headless")
	}
	if err := tui.Run(cfg, bundle, versionInfo()); err != nil {
		die("lỗi: %v", err)
	}
}

type cliOptions struct {
	Headless      bool
	Prompt        string
	PromptFile    string
	Version       bool
	Update        bool
	UpdateVersion string
}

// parseCLIOptions Trích xuất các cờ CLI, trả về các tùy chọn và tham số còn lại.
func parseCLIOptions(argv []string) (cliOptions, []string, error) {
	var opts cliOptions
	var args []string
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "--version", "-v":
			opts.Version = true
		case "version":
			if i+1 < len(argv) {
				return opts, nil, fmt.Errorf("version không nhận tham số")
			}
			opts.Version = true
		case "update":
			if opts.Update {
				return opts, nil, fmt.Errorf("update chỉ được chỉ định một lần")
			}
			opts.Update = true
			if i+1 < len(argv) {
				if strings.HasPrefix(argv[i+1], "-") {
					return opts, nil, fmt.Errorf("update chỉ nhận một tham số phiên bản tùy chọn")
				}
				opts.UpdateVersion = argv[i+1]
				i++
			}
			if i+1 < len(argv) {
				return opts, nil, fmt.Errorf("update chỉ nhận một tham số phiên bản tùy chọn")
			}
		case "--headless":
			opts.Headless = true
		case "--prompt":
			if i+1 >= len(argv) {
				return opts, nil, fmt.Errorf("--prompt thiếu giá trị")
			}
			opts.Prompt = argv[i+1]
			i++
		case "--prompt-file":
			if i+1 >= len(argv) {
				return opts, nil, fmt.Errorf("--prompt-file thiếu giá trị")
			}
			opts.PromptFile = argv[i+1]
			i++
		default:
			args = append(args, argv[i])
		}
	}
	if opts.Prompt != "" && opts.PromptFile != "" {
		return opts, nil, fmt.Errorf("--prompt và --prompt-file không thể dùng cùng lúc")
	}
	if opts.Version && (opts.Update || opts.Headless || opts.Prompt != "" || opts.PromptFile != "" || len(args) > 0) {
		return opts, nil, fmt.Errorf("version không thể kết hợp với các tham số khởi động khác")
	}
	if opts.Update && (opts.Headless || opts.Prompt != "" || opts.PromptFile != "" || len(args) > 0) {
		return opts, nil, fmt.Errorf("update không thể kết hợp với các tham số khởi động khác")
	}
	return opts, args, nil
}

func versionInfo() buildversion.Info {
	return buildversion.Resolve(buildversion.Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
}

func runSelfUpdate(target string) error {
	info := versionInfo()
	result, err := buildversion.Update(context.Background(), buildversion.UpdateOptions{
		Repo:           "voocel/ainovel-cli",
		BinaryName:     "ainovel-cli",
		TargetVersion:  target,
		CurrentVersion: info.Version,
	})
	if err != nil {
		return err
	}
	if !result.Updated {
		fmt.Printf("ainovel-cli đã là phiên bản mới nhất %s\n", result.Version)
		return nil
	}
	fmt.Printf("ainovel-cli đã được cập nhật lên %s\n", result.Version)
	fmt.Printf("Vị trí cài đặt: %s\n", result.Path)
	return nil
}

func loadPrompt(opts cliOptions) (string, error) {
	return loadPromptFrom(opts, os.Stdin)
}

func loadPromptFrom(opts cliOptions, stdin io.Reader) (string, error) {
	if opts.PromptFile == "" {
		return strings.TrimSpace(opts.Prompt), nil
	}

	if opts.PromptFile == "-" {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("Đọc prompt thất bại: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return startup.LoadPromptFile(opts.PromptFile)
}
