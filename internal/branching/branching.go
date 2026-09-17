// Package branching cung cấp tính năng rẽ nhánh cốt truyện (Multiverse / Parallel Timelines),
// cho phép tác giả "Fork" tác phẩm tại bất kỳ điểm kiểm soát hoặc chương nào để khám phá các hướng đi khác nhau
// mà không làm ảnh hưởng đến dòng thời gian chính (Main Timeline).
package branching

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Branch đại diện cho một nhánh cốt truyện song song.
type Branch struct {
	Name        string `json:"name"`
	ForkChapter int    `json:"fork_chapter"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Active      bool   `json:"active"`
}

// Meta lưu trữ cấu hình các nhánh trong tác phẩm.
type Meta struct {
	ActiveBranch string   `json:"active_branch"`
	Branches     []Branch `json:"branches"`
}

// Manager quản lý việc rẽ nhánh và chuyển đổi dòng thời gian cốt truyện.
type Manager struct {
	dir string
	mu  sync.RWMutex
}

// NewManager khởi tạo bộ quản lý nhánh.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) metaPath() string {
	return filepath.Join(m.dir, "meta", "branches.json")
}

// LoadMeta nạp thông tin danh sách nhánh.
func (m *Manager) LoadMeta() (*Meta, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.metaPath()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Meta{
				ActiveBranch: "main",
				Branches: []Branch{
					{
						Name:        "main",
						ForkChapter: 0,
						Description: "Dòng thời gian chính thống",
						CreatedAt:   time.Now().Format(time.RFC3339),
						Active:      true,
					},
				},
			}, nil
		}
		return nil, err
	}

	var meta Meta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// SaveMeta lưu thông tin cấu hình nhánh.
func (m *Manager) SaveMeta(meta *Meta) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	p := m.metaPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// CreateBranch tạo một nhánh cốt truyện mới xuất phát từ chương forkChapter.
func (m *Manager) CreateBranch(name, description string, forkChapter int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("tên nhánh không được để trống")
	}

	meta, err := m.LoadMeta()
	if err != nil {
		return err
	}

	for _, b := range meta.Branches {
		if b.Name == name {
			return fmt.Errorf("nhánh '%s' đã tồn tại", name)
		}
	}

	newBranch := Branch{
		Name:        name,
		ForkChapter: forkChapter,
		Description: description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Active:      false,
	}

	meta.Branches = append(meta.Branches, newBranch)

	// Tạo thư mục snapshot lưu trạng thái nhánh
	branchDir := filepath.Join(m.dir, "meta", "branches_data", name)
	if err := os.MkdirAll(branchDir, 0o755); err != nil {
		return fmt.Errorf("create branch dir: %w", err)
	}

	return m.SaveMeta(meta)
}

// ListBranches liệt kê danh sách các nhánh cốt truyện hiện có.
func (m *Manager) ListBranches() ([]Branch, string, error) {
	meta, err := m.LoadMeta()
	if err != nil {
		return nil, "", err
	}
	return meta.Branches, meta.ActiveBranch, nil
}
