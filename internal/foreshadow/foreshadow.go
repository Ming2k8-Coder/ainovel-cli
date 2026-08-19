// Package foreshadow cung cấp ma trận quản lý phục bút (Foreshadowing Matrix & Mystery Tracker),
// cho phép theo dõi các manh mối đã gài, các câu hỏi bí ẩn chưa giải đáp và đảm bảo các nút thắt cốt truyện
// được tháo gỡ hoàn chỉnh.
package foreshadow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Item đại diện cho một chi tiết phục bút hoặc bí ẩn được gài vào tác phẩm.
type Item struct {
	ID             string   `json:"id"`
	ChapterPlanted int      `json:"chapter_planted"` // Chương bắt đầu gài phục bút
	Description    string   `json:"description"`     // Chi tiết phục bút / manh mối
	Category       string   `json:"category"`        // "mystery", "item_secret", "character_motive", "world_lore"
	Resolved       bool     `json:"resolved"`
	ChapterResolved int     `json:"chapter_resolved,omitempty"`
	ResolutionNote string   `json:"resolution_note,omitempty"`
	CreatedAt      string   `json:"created_at"`
}

// Matrix chứa toàn bộ danh sách phục bút.
type Matrix struct {
	Items []Item `json:"items"`
}

// Manager quản lý mạng lưới phục bút.
type Manager struct {
	dir string
	mu  sync.RWMutex
}

// NewManager khởi tạo bộ quản lý phục bút.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) matrixFile() string {
	return filepath.Join(m.dir, "meta", "foreshadow_matrix.json")
}

// LoadMatrix nạp ma trận phục bút từ đĩa.
func (m *Manager) LoadMatrix() (*Matrix, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.matrixFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Matrix{Items: []Item{}}, nil
		}
		return nil, err
	}

	var mat Matrix
	if err := json.Unmarshal(data, &mat); err != nil {
		return nil, err
	}
	return &mat, nil
}

// SaveMatrix lưu ma trận phục bút.
func (m *Manager) SaveMatrix(mat *Matrix) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.matrixFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(mat, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// AddItem gài thêm một phục bút mới.
func (m *Manager) AddItem(chapter int, desc, category string) (*Item, error) {
	mat, err := m.LoadMatrix()
	if err != nil {
		return nil, err
	}

	item := Item{
		ID:             fmt.Sprintf("fs-%d-%d", time.Now().UnixNano(), chapter),
		ChapterPlanted: chapter,
		Description:    desc,
		Category:       category,
		Resolved:       false,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	mat.Items = append(mat.Items, item)
	if err := m.SaveMatrix(mat); err != nil {
		return nil, err
	}
	return &item, nil
}

// ResolveItem đánh dấu một phục bút đã được giải đáp.
func (m *Manager) ResolveItem(id string, chapter int, note string) error {
	mat, err := m.LoadMatrix()
	if err != nil {
		return err
	}

	found := false
	for i := range mat.Items {
		if mat.Items[i].ID == id || strings.HasPrefix(mat.Items[i].ID, id) {
			mat.Items[i].Resolved = true
			mat.Items[i].ChapterResolved = chapter
			mat.Items[i].ResolutionNote = note
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("không tìm thấy phục bút có mã %s", id)
	}

	return m.SaveMatrix(mat)
}
