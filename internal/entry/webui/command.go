package webui

import (
	"flag"
	"fmt"
	"os"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
)

// Command thực thi lệnh `ainovel-cli webui`.
func Command(args []string) int {
	fs := flag.NewFlagSet("webui", flag.ExitOnError)
	dir := fs.String("dir", "./novel", "Thư mục tác phẩm")
	port := fs.Int("port", 3000, "Cổng HTTP của WebUI Portal")
	_ = fs.Parse(args)

	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		cfg = bootstrap.Config{}
	}
	if *dir != "" {
		cfg.OutputDir = *dir
	}
	cfg.FillDefaults()

	bundle := assets.Load(cfg.Style, assets.DefaultLoadOptions(cfg.OutputDir))

	opts := Options{
		Port: *port,
		Dir:  cfg.OutputDir,
	}

	if err := Run(cfg, bundle, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi khởi chạy WebUI: %v\n", err)
		return 1
	}
	return 0
}
