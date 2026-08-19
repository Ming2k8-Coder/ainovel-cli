// Package gitmgr cung cấp công cụ tự động hóa và quản lý phiên bản Git (Git Version Control & Auto-Commit)
// cho toàn bộ tác phẩm, hỗ trợ phân nhánh (Branching), chuyển nhánh (Checkout), quản lý Issue
// và quy trình kiểm duyệt Pull Request (PR Review).
package gitmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Issue đại diện cho một nhiệm vụ/lỗi/yêu cầu nội dung cần xử lý trong tác phẩm.
type Issue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "open", "in_progress", "closed"
	BranchName  string `json:"branch_name,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// PullRequest đại diện cho một yêu cầu gộp nhánh và thẩm định nội dung.
type PullRequest struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	SourceBranch  string `json:"source_branch"`
	TargetBranch  string `json:"target_branch"`
	Status        string `json:"status"` // "open", "approved", "merged", "rejected"
	Reviewer      string `json:"reviewer,omitempty"`
	ReviewComment string `json:"review_comment,omitempty"`
	CreatedAt     string `json:"created_at"`
}

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

// CreateBranch tạo một nhánh Git mới.
func (m *Manager) CreateBranch(branchName string) error {
	if !m.IsGitRepo() {
		if err := m.InitRepo(); err != nil {
			return err
		}
	}
	_, err := m.runGit("branch", branchName)
	return err
}

// Checkout chuyển sang một nhánh Git chỉ định.
func (m *Manager) Checkout(branchName string, createIfMissing bool) error {
	if !m.IsGitRepo() {
		if err := m.InitRepo(); err != nil {
			return err
		}
	}
	if createIfMissing {
		_, err := m.runGit("checkout", "-b", branchName)
		return err
	}
	_, err := m.runGit("checkout", branchName)
	return err
}

// ListBranches liệt kê danh sách các nhánh Git và trả về nhánh hiện tại.
func (m *Manager) ListBranches() ([]string, string, error) {
	if !m.IsGitRepo() {
		return nil, "", fmt.Errorf("chưa khởi tạo Git Repository")
	}
	out, err := m.runGit("branch")
	if err != nil {
		return nil, "", err
	}

	var branches []string
	current := ""
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "* ") {
			b := strings.TrimPrefix(l, "* ")
			current = b
			branches = append(branches, b)
		} else {
			branches = append(branches, l)
		}
	}
	return branches, current, nil
}

// AutoCommitChapter tự động commit chương vừa hoàn thành vào Git.
func (m *Manager) AutoCommitChapter(chapter int, chapterTitle string, wordCount int) (string, error) {
	if !m.IsGitRepo() {
		if err := m.InitRepo(); err != nil {
			return "", fmt.Errorf("init git repo: %w", err)
		}
	}

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

// CreateTag tạo nhãn đánh dấu mốc tác phẩm.
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

// --- QUẢN LÝ ISSUES ---

func (m *Manager) issuesFile() string {
	return filepath.Join(m.dir, "meta", "git_issues.json")
}

// LoadIssues nạp danh sách Issues.
func (m *Manager) LoadIssues() ([]Issue, error) {
	p := m.issuesFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return []Issue{}, nil
		}
		return nil, err
	}
	var issues []Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (m *Manager) saveIssues(issues []Issue) error {
	p := m.issuesFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// CreateIssue tạo một Issue mới.
func (m *Manager) CreateIssue(title, description string) (*Issue, error) {
	issues, err := m.LoadIssues()
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("ISSUE-%d", len(issues)+1)
	issue := Issue{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      "open",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	issues = append(issues, issue)
	if err := m.saveIssues(issues); err != nil {
		return nil, err
	}
	return &issue, nil
}

// StartIssue chuyển Issue sang trạng thái in_progress và tạo nhánh feat/ (nếu yêu cầu).
func (m *Manager) StartIssue(issueID string, makeFeatBranch bool) (*Issue, error) {
	issues, err := m.LoadIssues()
	if err != nil {
		return nil, err
	}
	var target *Issue
	for i := range issues {
		if issues[i].ID == issueID {
			target = &issues[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("không tìm thấy issue %s", issueID)
	}

	target.Status = "in_progress"

	if makeFeatBranch {
		slug := slugify(target.Title)
		branchName := fmt.Sprintf("feat/%s-%s", strings.ToLower(target.ID), slug)
		if err := m.Checkout(branchName, true); err != nil {
			return nil, fmt.Errorf("create feat branch: %w", err)
		}
		target.BranchName = branchName
	} else {
		// Ở lại main
		_ = m.Checkout("main", false)
		target.BranchName = "main"
	}

	if err := m.saveIssues(issues); err != nil {
		return nil, err
	}
	return target, nil
}

// --- QUẢN LÝ PULL REQUESTS (PR REVIEW) ---

func (m *Manager) prsFile() string {
	return filepath.Join(m.dir, "meta", "git_pull_requests.json")
}

// LoadPRs nạp danh sách PRs.
func (m *Manager) LoadPRs() ([]PullRequest, error) {
	p := m.prsFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return []PullRequest{}, nil
		}
		return nil, err
	}
	var prs []PullRequest
	if err := json.Unmarshal(data, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}

func (m *Manager) savePRs(prs []PullRequest) error {
	p := m.prsFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(prs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// CreatePR mở một Pull Request để kiểm duyệt.
func (m *Manager) CreatePR(title, sourceBranch, targetBranch string) (*PullRequest, error) {
	prs, err := m.LoadPRs()
	if err != nil {
		return nil, err
	}
	if targetBranch == "" {
		targetBranch = "main"
	}
	if sourceBranch == "" {
		_, current, _ := m.ListBranches()
		sourceBranch = current
	}
	id := fmt.Sprintf("PR-%d", len(prs)+1)
	pr := PullRequest{
		ID:           id,
		Title:        title,
		SourceBranch: sourceBranch,
		TargetBranch: targetBranch,
		Status:       "open",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
	prs = append(prs, pr)
	if err := m.savePRs(prs); err != nil {
		return nil, err
	}
	return &pr, nil
}

// ReviewPR thực hiện thẩm định/đánh giá Pull Request (Approve hoặc Reject kèm bình luận).
func (m *Manager) ReviewPR(prID string, approved bool, reviewer, comment string) (*PullRequest, error) {
	prs, err := m.LoadPRs()
	if err != nil {
		return nil, err
	}
	var target *PullRequest
	for i := range prs {
		if prs[i].ID == prID {
			target = &prs[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("không tìm thấy PR %s", prID)
	}

	target.Reviewer = reviewer
	target.ReviewComment = comment
	if approved {
		target.Status = "approved"
	} else {
		target.Status = "rejected"
	}

	if err := m.savePRs(prs); err != nil {
		return nil, err
	}
	return target, nil
}

// MergePR thực hiện gộp nhánh Pull Request đã duyệt vào nhánh đích (chẳng hạn main).
func (m *Manager) MergePR(prID string) (*PullRequest, error) {
	prs, err := m.LoadPRs()
	if err != nil {
		return nil, err
	}
	var target *PullRequest
	for i := range prs {
		if prs[i].ID == prID {
			target = &prs[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("không tìm thấy PR %s", prID)
	}

	if target.Status != "approved" && target.Status != "open" {
		return nil, fmt.Errorf("PR %s đang ở trạng thái %s, không thể gộp", prID, target.Status)
	}

	// Thực hiện git checkout targetBranch và git merge sourceBranch
	if err := m.Checkout(target.TargetBranch, false); err != nil {
		return nil, fmt.Errorf("checkout target branch %s: %w", target.TargetBranch, err)
	}

	msg := fmt.Sprintf("merge: Pull Request #%s (%s -> %s): %s", target.ID, target.SourceBranch, target.TargetBranch, target.Title)
	if _, err := m.runGit("merge", "--no-ff", target.SourceBranch, "-m", msg); err != nil {
		return nil, fmt.Errorf("git merge: %w", err)
	}

	target.Status = "merged"
	if err := m.savePRs(prs); err != nil {
		return nil, err
	}
	return target, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
