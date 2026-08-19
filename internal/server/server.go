// Package server cung cấp máy chủ web xem truyện cục bộ (Local Web Reader & Previewer),
// cho phép tác giả đọc và duyệt tác phẩm trực quan trên trình duyệt web với giao diện Dark Mode & Glassmorphism.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/voocel/ainovel-cli/internal/authorship"
	"github.com/voocel/ainovel-cli/internal/cryptoaudit"
	"github.com/voocel/ainovel-cli/internal/store"
)

// Server quản lý dịch vụ xem truyện web.
type Server struct {
	dir   string
	port  int
	store *store.Store
}

// NewServer khởi tạo máy chủ web reader.
func NewServer(dir string, port int) *Server {
	return &Server{
		dir:   dir,
		port:  port,
		store: store.NewStore(dir),
	}
}

// Start khởi chạy HTTP server.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/info", s.handleAPIInfo)
	mux.HandleFunc("/api/chapters", s.handleAPIChapters)
	mux.HandleFunc("/api/chapter/", s.handleAPIChapterDetail)
	mux.HandleFunc("/api/authorship", s.handleAPIAuthorship)
	mux.HandleFunc("/", s.handleIndex)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 Đang khởi chạy Web Reader tại http://localhost:%d\n", s.port)
	fmt.Println("Nhấn Ctrl+C để dừng máy chủ web.")
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookTitle := filepath.Base(s.dir)
	if book, err := s.store.Book.Load(); err == nil && book != nil && book.Title != "" {
		bookTitle = book.Title
	}

	progress, _ := s.store.Progress.Load()
	totalCh := 0
	totalWords := 0
	if progress != nil {
		totalCh = len(progress.CompletedChapters)
		totalWords = progress.TotalWordCount
	}

	dossier, _ := cryptoaudit.BuildDossier(s.dir, bookTitle, "Human Author")
	merkleRoot := ""
	if dossier != nil {
		merkleRoot = dossier.MerkleRoot
	}

	res := map[string]any{
		"title":       bookTitle,
		"chapters":    totalCh,
		"word_count":  totalWords,
		"merkle_root": merkleRoot,
		"dir":         s.dir,
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleAPIChapters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	progress, err := s.store.Progress.Load()
	if err != nil || progress == nil {
		_ = json.NewEncoder(w).Encode([]any{})
		return
	}

	type ChapterItem struct {
		Chapter   int    `json:"chapter"`
		Title     string `json:"title"`
		WordCount int    `json:"word_count"`
		Summary   string `json:"summary"`
	}

	var items []ChapterItem
	for _, ch := range progress.CompletedChapters {
		content, _, err := s.store.Drafts.LoadChapterContent(ch)
		if err != nil {
			continue
		}
		summary, _ := s.store.Summaries.LoadSummary(ch)
		title := fmt.Sprintf("Chương %d", ch)
		sumText := ""
		if summary != nil {
			if summary.Title != "" {
				title = fmt.Sprintf("Chương %d: %s", ch, summary.Title)
			}
			sumText = summary.Summary
		}

		items = append(items, ChapterItem{
			Chapter:   ch,
			Title:     title,
			WordCount: len([]rune(content)),
			Summary:   sumText,
		})
	}

	_ = json.NewEncoder(w).Encode(items)
}

func (s *Server) handleAPIChapterDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid chapter path", http.StatusBadRequest)
		return
	}
	chNum, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid chapter number", http.StatusBadRequest)
		return
	}

	content, _, err := s.store.Drafts.LoadChapterContent(chNum)
	if err != nil {
		http.Error(w, "chapter not found", http.StatusNotFound)
		return
	}

	summary, _ := s.store.Summaries.LoadSummary(chNum)
	title := fmt.Sprintf("Chương %d", chNum)
	if summary != nil && summary.Title != "" {
		title = fmt.Sprintf("Chương %d: %s", chNum, summary.Title)
	}

	res := map[string]any{
		"chapter":    chNum,
		"title":      title,
		"content":    content,
		"word_count": len([]rune(content)),
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleAPIAuthorship(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ledger := authorship.NewLedger(s.dir)
	records, err := ledger.LoadRecords()
	if err != nil {
		_ = json.NewEncoder(w).Encode([]any{})
		return
	}
	_ = json.NewEncoder(w).Encode(records)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SANBAKA Web Reader — ainovel-cli</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700&family=Outfit:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #0b0f19;
            --panel: rgba(22, 30, 49, 0.7);
            --border: rgba(255, 255, 255, 0.08);
            --accent: #38bdf8;
            --accent-glow: rgba(56, 189, 248, 0.2);
            --text: #f1f5f9;
            --text-sub: #94a3b8;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg);
            color: var(--text);
            display: flex;
            height: 100vh;
            overflow: hidden;
        }
        sidebar {
            width: 320px;
            background: var(--panel);
            backdrop-filter: blur(16px);
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
        }
        .header {
            padding: 24px;
            border-bottom: 1px solid var(--border);
        }
        .header h1 {
            font-family: 'Outfit', sans-serif;
            font-size: 1.25rem;
            color: var(--accent);
            margin-bottom: 6px;
        }
        .meta-badge {
            font-size: 0.75rem;
            color: var(--text-sub);
            background: rgba(255, 255, 255, 0.05);
            padding: 4px 8px;
            border-radius: 4px;
            display: inline-block;
        }
        .chapter-list {
            flex: 1;
            overflow-y: auto;
            padding: 12px;
        }
        .chapter-item {
            padding: 12px 16px;
            border-radius: 8px;
            margin-bottom: 6px;
            cursor: pointer;
            transition: all 0.2s;
            border: 1px solid transparent;
        }
        .chapter-item:hover {
            background: rgba(255, 255, 255, 0.04);
        }
        .chapter-item.active {
            background: var(--accent-glow);
            border-color: var(--accent);
            color: var(--accent);
        }
        .chapter-item h3 { font-size: 0.9rem; font-weight: 600; }
        .chapter-item p { font-size: 0.75rem; color: var(--text-sub); margin-top: 4px; }
        main {
            flex: 1;
            padding: 40px 80px;
            overflow-y: auto;
            max-width: 900px;
            margin: 0 auto;
        }
        .reader-title {
            font-family: 'Outfit', sans-serif;
            font-size: 2rem;
            margin-bottom: 24px;
            padding-bottom: 16px;
            border-bottom: 1px solid var(--border);
        }
        .reader-content {
            font-size: 1.1rem;
            line-height: 1.8;
            color: #e2e8f0;
            white-space: pre-wrap;
        }
    </style>
</head>
<body>
    <sidebar>
        <div class="header">
            <h1 id="book-title">Đang tải...</h1>
            <div class="meta-badge" id="book-stats">...</div>
        </div>
        <div class="chapter-list" id="chapters"></div>
    </sidebar>
    <main>
        <h2 class="reader-title" id="chapter-title">Chọn chương để đọc</h2>
        <div class="reader-content" id="chapter-content">Vui lòng chọn một chương từ danh sách bên trái.</div>
    </main>

    <script>
        async function init() {
            const info = await fetch('/api/info').then(r => r.json());
            document.getElementById('book-title').innerText = info.title;
            document.getElementById('book-stats').innerText = info.chapters + ' chương • ' + info.word_count + ' từ';

            const chapters = await fetch('/api/chapters').then(r => r.json());
            const list = document.getElementById('chapters');
            list.innerHTML = '';
            chapters.forEach(ch => {
                const div = document.createElement('div');
                div.className = 'chapter-item';
                div.innerHTML = '<h3>' + ch.title + '</h3><p>' + ch.word_count + ' từ</p>';
                div.onclick = () => loadChapter(ch.chapter, div);
                list.appendChild(div);
            });

            if (chapters.length > 0) {
                loadChapter(chapters[0].chapter, list.children[0]);
            }
        }

        async function loadChapter(chNum, el) {
            document.querySelectorAll('.chapter-item').forEach(e => e.classList.remove('active'));
            if (el) el.classList.add('active');

            const data = await fetch('/api/chapter/' + chNum).then(r => r.json());
            document.getElementById('chapter-title').innerText = data.title;
            document.getElementById('chapter-content').innerText = data.content;
        }

        init();
    </script>
</body>
</html>`
	_, _ = w.Write([]byte(html))
}
