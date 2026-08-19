// Package charvector cung cấp ma trận quản lý quan hệ và chỉ số tâm lý nhân vật (Character Vector & Dynamics),
// giúp theo dõi độ tin cậy (Trust), xung đột (Tension), liên minh (Alliance) và nợ nần (Debt) giữa các nhân vật.
package charvector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// RelationshipVector chứa các chỉ số quan hệ giữa hai nhân vật.
type RelationshipVector struct {
	CharA     string `json:"char_a"`
	CharB     string `json:"char_b"`
	Trust     int    `json:"trust"`    // -100 (Thù hằn) -> +100 (Tuyệt đối tin tưởng)
	Tension   int    `json:"tension"`  // 0 (Bình yên) -> 100 (Căng thẳng cực độ)
	Alliance  string `json:"alliance"` // "ally", "neutral", "rival", "enemy"
	LastUpdateCh int `json:"last_update_ch"`
	Note      string `json:"note"`
}

// Matrix chứa toàn bộ ma trận quan hệ nhân vật.
type Matrix struct {
	Vectors []RelationshipVector `json:"vectors"`
}

// Manager quản lý ma trận chỉ số tâm lý & quan hệ nhân vật.
type Manager struct {
	dir string
	mu  sync.RWMutex
}

// NewManager khởi tạo bộ quản lý vector quan hệ nhân vật.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) vectorFile() string {
	return filepath.Join(m.dir, "meta", "character_vectors.json")
}

// LoadMatrix nạp ma trận vector quan hệ từ đĩa.
func (m *Manager) LoadMatrix() (*Matrix, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.vectorFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Matrix{Vectors: []RelationshipVector{}}, nil
		}
		return nil, err
	}

	var mat Matrix
	if err := json.Unmarshal(data, &mat); err != nil {
		return nil, err
	}
	return &mat, nil
}

// SaveMatrix lưu ma trận quan hệ.
func (m *Manager) SaveMatrix(mat *Matrix) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.vectorFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(mat, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// UpdateVector cập nhật chỉ số quan hệ giữa hai nhân vật.
func (m *Manager) UpdateVector(charA, charB string, trust, tension int, alliance, note string, ch int) error {
	mat, err := m.LoadMatrix()
	if err != nil {
		return err
	}

	found := false
	for i := range mat.Vectors {
		if (mat.Vectors[i].CharA == charA && mat.Vectors[i].CharB == charB) ||
			(mat.Vectors[i].CharA == charB && mat.Vectors[i].CharB == charA) {
			mat.Vectors[i].Trust = trust
			mat.Vectors[i].Tension = tension
			mat.Vectors[i].Alliance = alliance
			mat.Vectors[i].Note = note
			mat.Vectors[i].LastUpdateCh = ch
			found = true
			break
		}
	}

	if !found {
		mat.Vectors = append(mat.Vectors, RelationshipVector{
			CharA:        charA,
			CharB:        charB,
			Trust:        trust,
			Tension:      tension,
			Alliance:     alliance,
			LastUpdateCh: ch,
			Note:         note,
		})
	}

	return m.SaveMatrix(mat)
}
