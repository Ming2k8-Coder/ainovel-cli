// Package nlcommander cung cấp Bộ Điều Khiển Ngôn Ngữ Tự Nhiên (Natural Language Commander).
// Cho phép tác giả gõ bất kỳ câu lệnh bằng tiếng Việt hoặc tiếng Anh bình thường,
// hệ thống sẽ tự động phân tích ý định (Intent Recognition) và thực thi đúng chuỗi công việc tương ứng trong ainovel-cli.
package nlcommander

import (
	"fmt"
	"strings"

	"github.com/voocel/ainovel-cli/internal/abselect"
	"github.com/voocel/ainovel-cli/internal/analytics"
	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/cryptoaudit"
	"github.com/voocel/ainovel-cli/internal/gitmgr"
	"github.com/voocel/ainovel-cli/internal/novelanalyzer"
	"github.com/voocel/ainovel-cli/internal/publisher"
	"github.com/voocel/ainovel-cli/internal/server"
	"github.com/voocel/ainovel-cli/internal/stylepreset"
)

// Intent đại diện cho ý định được trích xuất từ câu lệnh ngôn ngữ tự nhiên.
type Intent struct {
	Action string            `json:"action"` // "webui", "serve", "analyze", "copyright", "style", "git", "publish", "stats", "ab"
	Params map[string]string `json:"params"`
	Raw    string            `json:"raw"`
}

// Commander quản lý việc phân tích và điều phối lệnh.
type Commander struct {
	dir string
}

// NewCommander khởi tạo Commander.
func NewCommander(dir string) *Commander {
	return &Commander{dir: dir}
}

// ParseIntent phân tích câu lệnh tự nhiên để xác định Action và Parameters.
func (c *Commander) ParseIntent(input string) Intent {
	inputLower := strings.ToLower(input)
	intent := Intent{
		Params: make(map[string]string),
		Raw:    input,
	}

	// 1. Phán đoán intent liên quan tới WebUI / Viewer
	if strings.Contains(inputLower, "webui") || strings.Contains(inputLower, "giao diện web") || strings.Contains(inputLower, "dashboard") {
		intent.Action = "webui"
		return intent
	}

	if strings.Contains(inputLower, "xem truyện") || strings.Contains(inputLower, "serve") || strings.Contains(inputLower, "reader") {
		intent.Action = "serve"
		return intent
	}

	// 2. Phán đoán intent liên quan tới Phân tích / Viết tiếp truyện có sẵn
	if strings.Contains(inputLower, "viết tiếp") || strings.Contains(inputLower, "phân tích") || strings.Contains(inputLower, "analyze") || strings.Contains(inputLower, "continue") {
		intent.Action = "analyze"
		// Trích xuất đường dẫn nếu có trong ngoặc hoặc đằng sau từ "từ"
		if idx := strings.Index(inputLower, "từ "); idx != -1 {
			pathPart := strings.TrimSpace(input[idx+len("từ "):])
			intent.Params["source"] = pathPart
		} else {
			intent.Params["source"] = "./source"
		}
		return intent
	}

	// 3. Phán đoán intent liên quan tới Bản quyền / Merkle Proof
	if strings.Contains(inputLower, "bản quyền") || strings.Contains(inputLower, "copyright") || strings.Contains(inputLower, "merkle") || strings.Contains(inputLower, "chứng thư") {
		intent.Action = "copyright"
		return intent
	}

	// 4. Phán đoán intent liên quan tới Chuyển Style / Prompt
	if strings.Contains(inputLower, "style") || strings.Contains(inputLower, "văn phong") || strings.Contains(inputLower, "phong cách") || strings.Contains(inputLower, "chuyển prompt") {
		intent.Action = "style"
		if strings.Contains(inputLower, "cyberpunk") {
			intent.Params["switch"] = "cyberpunk"
		} else if strings.Contains(inputLower, "xianxia") || strings.Contains(inputLower, "tiên hiệp") {
			intent.Params["switch"] = "xianxia"
		} else if strings.Contains(inputLower, "mystery") || strings.Contains(inputLower, "trinh thám") {
			intent.Params["switch"] = "mystery"
		} else {
			intent.Params["list"] = "true"
		}
		return intent
	}

	// 5. Phán đoán intent liên quan tới Git Workflow
	if strings.Contains(inputLower, "git") || strings.Contains(inputLower, "commit") || strings.Contains(inputLower, "nhánh") || strings.Contains(inputLower, "pr") || strings.Contains(inputLower, "issue") {
		intent.Action = "git"
		return intent
	}

	// 6. Phán đoán intent liên quan tới Xuất bản MDX / FullBook
	if strings.Contains(inputLower, "xuất bản") || strings.Contains(inputLower, "publish") || strings.Contains(inputLower, "mdx") {
		intent.Action = "publish"
		return intent
	}

	// 7. Phán đoán intent liên quan tới Thống kê Stats / Analytics
	if strings.Contains(inputLower, "thống kê") || strings.Contains(inputLower, "chi phí") || strings.Contains(inputLower, "stats") || strings.Contains(inputLower, "token") {
		intent.Action = "stats"
		return intent
	}

	// Mặc định: Phân tích A/B Feedback
	intent.Action = "ab"
	return intent
}

// Execute thực thi Intent tương ứng.
func (c *Commander) Execute(intent Intent) error {
	fmt.Printf("🤖 NATURAL LANGUAGE COMMANDER: Đã nhận diện ý định [%s]\n", strings.ToUpper(intent.Action))
	fmt.Println("------------------------------------------------------------")

	switch intent.Action {
	case "webui":
		fmt.Println("🚀 Khởi chạy WebUI Dashboard...")
		return server.NewServer(c.dir, 3000).Start()

	case "serve":
		fmt.Println("📖 Khởi chạy Web Reader Server...")
		return server.NewServer(c.dir, 8080).Start()

	case "analyze":
		source := intent.Params["source"]
		if source == "" {
			source = "./"
		}
		fmt.Printf("🔍 Phân tích tác phẩm có sẵn tại %s...\n", source)
		res, err := novelanalyzer.NewAnalyzer(source).AnalyzePath()
		if err != nil {
			return err
		}
		return novelanalyzer.SetupContinuationProject(c.dir, res)

	case "copyright":
		fmt.Println("🛡️ Lập chứng thư bản quyền số học Merkle Tree...")
		dossier, err := cryptoaudit.BuildDossier(c.dir, "Tác Phẩm", "Human Author")
		if err != nil {
			return err
		}
		fmt.Println(dossier.ExportLegalAffidavit())
		return nil

	case "style":
		styleMgr := stylepreset.NewManager(c.dir)
		if targetStyle, ok := intent.Params["switch"]; ok {
			fmt.Printf("🎨 Chuyển phong cách prompt sang [%s]...\n", targetStyle)
			return styleMgr.SwitchStyle(targetStyle)
		}
		styles, _ := styleMgr.ListStyles("cyberpunk")
		fmt.Printf("🎨 Các phong cách sẵn có (%d): %v\n", len(styles), styles)
		return nil

	case "git":
		statusStr, err := gitmgr.NewManager(c.dir).Status()
		if err != nil {
			return err
		}
		fmt.Println(statusStr)
		return nil

	case "publish":
		fmt.Println("📚 Xuất bản tác phẩm sang MDX tương thích sanbaka-web...")
		return publisher.CompileMDX(nil, publisher.ExportOptions{OutputDir: "./dist/mdx", BookTitle: "Tác Phẩm"})

	case "stats":
		return analytics.PrintReport(c.dir)

	case "ab":
		ledger := authorship.NewLedger(c.dir)
		synth, err := abselect.NewManager(c.dir, ledger).SynthesizePreferences()
		if err != nil {
			return err
		}
		fmt.Printf("📊 Sở thích A/B đã học được: %v\n", synth.PreferredTraits)
		return nil

	default:
		return fmt.Errorf("không thể nhận diện lệnh tự nhiên: %s", intent.Raw)
	}
}
