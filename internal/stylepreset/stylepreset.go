// Package stylepreset cung cấp trình quản lý đa phong cách văn học (Multiple Style Presets Manager),
// cho phép chuyển đổi nhanh giữa các phong cách (Chuyển Prompt), tạo phong cách mới từ phản hồi A/B
// và tinh chỉnh (fix/tune) các phong cách hiện có.
package stylepreset

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/abselect"
	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
)

// StyleInfo chứa thông tin hiển thị của một bộ phong cách.
type StyleInfo struct {
	Name        string `json:"name"`        // Ví dụ: "cyberpunk", "xianxia", "mystery"
	Source      string `json:"source"`      // "builtin", "global", "project"
	Active      bool   `json:"active"`      // Có đang kích hoạt cho tác phẩm hiện tại không
	Description string `json:"description"`
}

// Manager quản lý việc chuyển đổi và sửa đổi phong cách văn phong.
type Manager struct {
	dir string
}

// NewManager khởi tạo bộ quản lý phong cách.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

// ListStyles liệt kê toàn bộ các phong cách sẵn có (built-in + global + project).
func (m *Manager) ListStyles(currentActive string) ([]StyleInfo, error) {
	bundle := assets.Load(currentActive, assets.DefaultLoadOptions(m.dir))
	var list []StyleInfo

	for name := range bundle.Styles {
		source := "builtin"
		pPath := filepath.Join(m.dir, "styles", name+".md")
		if _, err := os.Stat(pPath); err == nil {
			source = "project"
		} else {
			homeDir, _ := os.UserHomeDir()
			gPath := filepath.Join(homeDir, ".ainovel", "styles", name+".md")
			if _, err := os.Stat(gPath); err == nil {
				source = "global"
			}
		}

		desc := fmt.Sprintf("Phong cách văn học %s (%s)", name, source)
		list = append(list, StyleInfo{
			Name:        name,
			Source:      source,
			Active:      name == currentActive,
			Description: desc,
		})
	}
	return list, nil
}

// SwitchStyle chuyển đổi nhanh phong cách hoạt động của tác phẩm.
func (m *Manager) SwitchStyle(styleName string) error {
	styleName = strings.TrimSpace(strings.ToLower(styleName))
	if styleName == "" {
		return fmt.Errorf("tên phong cách không được để trống")
	}

	// 1. Cập nhật file config.json của dự án nếu có
	configPath := filepath.Join(m.dir, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var cfg map[string]any
		_ = jsonUnmarshal(data, &cfg)
		cfg["style"] = styleName
		if out, err := jsonMarshalIndent(cfg); err == nil {
			_ = os.WriteFile(configPath, out, 0o644)
		}
	}

	// 2. Tải bundle của style mới và ghi đè trực tiếp vào <project>/style/voice.md
	bundle := assets.Load(styleName, assets.DefaultLoadOptions(m.dir))
	styleContent, exists := bundle.Styles[styleName]
	if !exists {
		styleContent = fmt.Sprintf("# Phong cách %s\n\nƯu tiên nhịp điệu và văn phong đặc trưng của thể loại %s.\n", styleName, styleName)
	}

	projectStyleDir := filepath.Join(m.dir, "style")
	if err := os.MkdirAll(projectStyleDir, 0o755); err != nil {
		return err
	}

	voicePath := filepath.Join(projectStyleDir, "voice.md")
	if err := os.WriteFile(voicePath, []byte(styleContent), 0o644); err != nil {
		return fmt.Errorf("write project voice.md: %w", err)
	}

	return nil
}

// CreateStyleFromAB tạo một bộ phong cách hoàn toàn mới dựa trên kết quả phân tích A/B Selection.
func (m *Manager) CreateStyleFromAB(styleName string) error {
	styleName = strings.TrimSpace(strings.ToLower(styleName))
	if styleName == "" {
		return fmt.Errorf("tên phong cách mới không được để trống")
	}

	ledger := authorship.NewLedger(m.dir)
	abMgr := abselect.NewManager(m.dir, ledger)

	synth, err := abMgr.SynthesizePreferences()
	if err != nil {
		return fmt.Errorf("phân tích phản hồi A/B thất bại: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Phong cách Độc bản: %s (Tạo từ A/B Feedback)\n\n", styleName)
	b.WriteString("> Phong cách này được đúc kết tự động từ các quyết định tuyển chọn nghệ thuật A/B của tác giả.\n\n")

	if len(synth.ProseDirectives) > 0 {
		b.WriteString("## 1. Nguyên Tắc Văn Phong Ưu Tiên\n\n")
		for _, d := range synth.ProseDirectives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}

	if len(synth.Taboos) > 0 {
		b.WriteString("## 2. Ràng Buộc Cấm Kỵ\n\n")
		for _, t := range synth.Taboos {
			fmt.Fprintf(&b, "- %s\n", t)
		}
		b.WriteString("\n")
	}

	// Lưu vào dự án tại <project>/styles/<styleName>.md
	stylesDir := filepath.Join(m.dir, "styles")
	if err := os.MkdirAll(stylesDir, 0o755); err != nil {
		return err
	}

	targetPath := filepath.Join(stylesDir, styleName+".md")
	if err := os.WriteFile(targetPath, []byte(b.String()), 0o644); err != nil {
		return err
	}

	// Tự động chuyển đổi sang phong cách mới vừa tạo
	return m.SwitchStyle(styleName)
}

// TuneStyleWithAB tinh chỉnh (Fix/Tune) một phong cách hiện có bằng dữ liệu A/B mới.
func (m *Manager) TuneStyleWithAB(styleName string) error {
	ledger := authorship.NewLedger(m.dir)
	abMgr := abselect.NewManager(m.dir, ledger)

	synth, err := abMgr.SynthesizePreferences()
	if err != nil {
		return err
	}

	bundle := assets.Load(styleName, assets.DefaultLoadOptions(m.dir))
	existingContent, exists := bundle.Styles[styleName]
	if !exists {
		existingContent = fmt.Sprintf("# Phong cách %s\n\n", styleName)
	}

	var b strings.Builder
	b.WriteString(existingContent)
	b.WriteString("\n\n---\n")
	b.WriteString("## 3. Tinh Chỉnh Bổ Sung Từ Phản Hồi A/B Của Tác Giả (AB Feedback Tune)\n\n")

	for _, d := range synth.ProseDirectives {
		fmt.Fprintf(&b, "- [Bổ sung] %s\n", d)
	}
	for _, t := range synth.Taboos {
		fmt.Fprintf(&b, "- [Hạn chế] %s\n", t)
	}

	stylesDir := filepath.Join(m.dir, "styles")
	if err := os.MkdirAll(stylesDir, 0o755); err != nil {
		return err
	}

	targetPath := filepath.Join(stylesDir, styleName+".md")
	if err := os.WriteFile(targetPath, []byte(b.String()), 0o644); err != nil {
		return err
	}

	return m.SwitchStyle(styleName)
}

func jsonUnmarshal(data []byte, v any) error {
	return bootstrap.JSONUnmarshal(data, v)
}

func jsonMarshalIndent(v any) ([]byte, error) {
	return bootstrap.JSONMarshalIndent(v)
}
