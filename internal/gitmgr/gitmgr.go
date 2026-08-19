// Package gitmgr cung cấp công cụ tự động hóa và quản lý phiên bản Git (Git Version Control & Auto-Commit)
// cho toàn bộ tác phẩm, giúp lưu trữ tiến trình sáng tác theo từng chương và từng mốc lịch sử.
package gitmgr

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Manager quản lý các thao tác Git trong dự án.
type Manager struct {
	dir string
}

// NewManager khởi tạo bộ quản lý Git.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = m.dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// IsGitRepo kiểm tra xem thư mục có phải là repository Git hay không.
func (m *Manager) IsGitRepo() bool {
	_, err := m.runGit("rev-parse", "--is-inside-work-tree")
	return err == nil
}

// InitRepo khởi tạo repository Git nếu chưa có.
func (m *Manager) InitRepo() error {
	if m.IsGitRepo() {
		return nil
	}
	_, err := m.runGit("init")
	return err
}

// AutoCommitChapter tự động commit chương vừa hoàn thành vào Git.
func (m *Manager) AutoCommitChapter(chapter int, chapterTitle string, wordCount int) (string, error) {
	if !m.IsGitRepo() {
		if err := m.InitRepo(); err != nil {
			return "", fmt.Errorf("init git repo: %w", err)
		}
	}

	// Stage toàn bộ thay đổi trong thư mục tác phẩm
	if _, err := m.runGit("add", "."); err != nil {
		return "", fmt.Errorf("git add: %w", err)
	}

	msg := fmt.Sprintf("feat(ch%02d): %s (%d từ)", chapter, chapterTitle, wordCount)
	out, err := m.runGit("commit", "-m", msg)
	if err != nil {
		if strings.Contains(out, "nothing to commit") {
			return "no changes", nil
		}
		return "", fmt.Errorf("git commit: %w", err)
	}
	return msg, nil
}

// CreateTag tạo nhãn đánh dấu mốc tác phẩm (ví dụ: vol1-complete, arc2-done).
func (m *Manager) CreateTag(tagName, message string) error {
	if !m.IsGitRepo() {
		return fmt.Errorf("chưa khởi tạo git repository")
	}
	_, err := m.runGit("tag", "-a", tagName, "-m", message)
	return err
}

// Status kiểm tra trạng thái Git hiện tại.
func (m *Manager) Status() (string, error) {
	if !m.IsGitRepo() {
		return "Chưa khởi tạo Git Repository", nil
	}
	branch, _ := m.runGit("branch", "--show-current")
	status, err := m.runGit("status", "--short")
	if err != nil {
		return "", err
	}
	if status == "" {
		status = "Working tree clean (Không có thay đổi chưa commit)"
	}
	return fmt.Sprintf("Nhánh hiện tại: %s\n\nTrạng thái:\n%s", branch, status), nil
}

// History lấy lịch sử 10 commit gần nhất.
func (m *Manager) History() (string, error) {
	if !m.IsGitRepo() {
		return "Chưa có lịch sử Git", nil
	}
	return m.runGit("log", "--oneline", "-n", "10")
}
