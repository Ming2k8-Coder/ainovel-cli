// Package webui cung cấp giao diện đồ họa người dùng trên trình duyệt (WebUI Portal),
// hoạt động song song cùng với giao diện TUI (Terminal UI) và chế độ Headless (CLI),
// mang lại trải nghiệm tương tác trực quan cao cấp nhất cho tác giả.
package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/cryptoaudit"
	"github.com/voocel/ainovel-cli/internal/gitmgr"
	"github.com/voocel/ainovel-cli/internal/host"
	"github.com/voocel/ainovel-cli/internal/store"
)

// Options tùy chọn khởi chạy WebUI.
type Options struct {
	Port int
	Dir  string
}

// Run khởi chạy máy chủ WebUI đầy đủ tính năng.
func Run(cfg bootstrap.Config, bundle assets.Bundle, opts Options) error {
	if opts.Port <= 0 {
		opts.Port = 3000
	}
	if opts.Dir == "" {
		opts.Dir = cfg.OutputDir
	}

	st := store.NewStore(opts.Dir)
	h := host.New(cfg, bundle, st)
	git := gitmgr.NewManager(opts.Dir)

	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		progress, _ := st.Progress.Load()
		runMeta, _ := st.RunMeta.Load()
		book, _ := st.Book.Load()

		bookTitle := filepath.Base(opts.Dir)
		if book != nil && book.Title != "" {
			bookTitle = book.Title
		}

		chCount := 0
		words := 0
		if progress != nil {
			chCount = len(progress.CompletedChapters)
			words = progress.TotalWordCount
		}

		dossier, _ := cryptoaudit.BuildDossier(opts.Dir, bookTitle, "Human Author")
		merkle := ""
		if dossier != nil {
			merkle = dossier.MerkleRoot
		}

		phase := "Khởi tạo"
		if progress != nil {
			phase = string(progress.Phase)
		}

		res := map[string]any{
			"title":       bookTitle,
			"phase":       phase,
			"chapters":    chCount,
			"words":       words,
			"merkle_root": merkle,
			"advance_mode": func() string {
				if runMeta != nil { return string(runMeta.AdvanceMode) }
				return "manual"
			}(),
		}
		_ = json.NewEncoder(w).Encode(res)
	})

	mux.HandleFunc("/api/next", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := h.GrantAdvancePermit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "message": "Đã cấp phép phát hành chương tiếp theo"})
	})

	mux.HandleFunc("/api/steer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Text) == "" {
			http.Error(w, "văn bản can thiệp không hợp lệ", http.StatusBadRequest)
			return
		}
		if err := h.Steer(req.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "message": "Đã gửi định hướng can thiệp đến Host/Arbiter"})
	})

	mux.HandleFunc("/api/git/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		statusStr, _ := git.Status()
		historyStr, _ := git.History()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  statusStr,
			"history": historyStr,
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ainovel-cli — WebUI Interactive Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700&family=Outfit:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #090d16;
            --card: rgba(22, 30, 49, 0.75);
            --border: rgba(255, 255, 255, 0.1);
            --accent: #38bdf8;
            --accent-glow: rgba(56, 189, 248, 0.25);
            --success: #34d399;
            --text: #f8fafc;
            --text-sub: #94a3b8;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Inter', sans-serif;
            background: var(--bg);
            color: var(--text);
            padding: 32px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        .top-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 32px;
            padding-bottom: 16px;
            border-bottom: 1px solid var(--border);
        }
        .logo { font-family: 'Outfit', sans-serif; font-size: 1.75rem; color: var(--accent); }
        .grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 32px; }
        .card {
            background: var(--card);
            backdrop-filter: blur(12px);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 20px;
        }
        .card h3 { font-size: 0.85rem; color: var(--text-sub); text-transform: uppercase; margin-bottom: 8px; }
        .card .val { font-size: 1.5rem; font-weight: 700; color: var(--accent); }
        .controls { display: grid; grid-template-columns: 2fr 1fr; gap: 20px; }
        .panel {
            background: var(--card);
            backdrop-filter: blur(12px);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 24px;
        }
        .panel h2 { font-size: 1.25rem; font-family: 'Outfit', sans-serif; margin-bottom: 16px; }
        textarea {
            width: 100%;
            height: 100px;
            background: rgba(0,0,0,0.3);
            border: 1px solid var(--border);
            border-radius: 8px;
            color: var(--text);
            padding: 12px;
            font-family: inherit;
            resize: none;
            margin-bottom: 12px;
        }
        button {
            background: var(--accent);
            color: #0f172a;
            border: none;
            padding: 12px 24px;
            font-weight: 700;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s;
        }
        button:hover { opacity: 0.9; transform: translateY(-1px); }
        .btn-success { background: var(--success); }
    </style>
</head>
<body>
    <div class="container">
        <div class="top-bar">
            <div class="logo">ainovel-cli WebUI</div>
            <div style="font-size: 0.9rem; color: var(--text-sub);">Interactive Master Dashboard</div>
        </div>

        <div class="grid">
            <div class="card">
                <h3>Tác phẩm</h3>
                <div class="val" id="title">...</div>
            </div>
            <div class="card">
                <h3>Giai đoạn</h3>
                <div class="val" id="phase">...</div>
            </div>
            <div class="card">
                <h3>Tiến độ</h3>
                <div class="val" id="chapters">...</div>
            </div>
            <div class="card">
                <h3>Bản quyền Merkle Root</h3>
                <div class="val" style="font-size: 0.9rem; word-break: break-all;" id="merkle">...</div>
            </div>
        </div>

        <div class="controls">
            <div class="panel">
                <h2>Can Thiệp Sáng Tạo Thời Gian Thực (Creative Steer)</h2>
                <textarea id="steer-text" placeholder="Nhập định hướng tình tiết, hành động nhân vật hoặc chuyển hướng cốt truyện..."></textarea>
                <button onclick="sendSteer()">Gửi Định Hướng Can Thiệp</button>
            </div>

            <div class="panel" style="display: flex; flex-direction: column; justify-content: center; align-items: center; text-align: center;">
                <h2>Cấp Phép Phát Hành</h2>
                <p style="color: var(--text-sub); font-size: 0.85rem; margin-bottom: 16px;">Nghiệm thu chương hiện tại và cho phép AI tiếp tục sáng tác chương sau</p>
                <button class="btn-success" style="width: 100%; padding: 16px;" onclick="sendNext()">Duyệt & Cấp Phép (/next)</button>
            </div>
        </div>
    </div>

    <script>
        async function updateStatus() {
            const data = await fetch('/api/status').then(r => r.json());
            document.getElementById('title').innerText = data.title;
            document.getElementById('phase').innerText = data.phase;
            document.getElementById('chapters').innerText = data.chapters + ' chương (' + data.words + ' từ)';
            document.getElementById('merkle').innerText = data.merkle_root ? ('0x' + data.merkle_root.substring(0, 12) + '...') : 'Chưa đóng dấu';
        }

        async function sendSteer() {
            const text = document.getElementById('steer-text').value;
            if (!text) return alert('Vui lòng nhập nội dung can thiệp');
            const res = await fetch('/api/steer', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({text})
            }).then(r => r.json());
            alert(res.message);
            document.getElementById('steer-text').value = '';
        }

        async function sendNext() {
            const res = await fetch('/api/next', {method: 'POST'}).then(r => r.json());
            alert(res.message);
            updateStatus();
        }

        updateStatus();
        setInterval(updateStatus, 5000);
    </script>
</body>
</html>`
		_, _ = w.Write([]byte(html))
	})

	addr := fmt.Sprintf(":%d", opts.Port)
	fmt.Printf("🌐 ĐÃ KHỞI CHẠY WEBUI DASHBOARD TẠI: http://localhost:%d\n", opts.Port)
	fmt.Println("Truy cập trình duyệt để theo dõi và điều khiển sáng tác trực quan.")
	return http.ListenAndServe(addr, mux)
}
