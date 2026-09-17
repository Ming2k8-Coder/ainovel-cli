// Package universe cung cấp Động Cơ Quản Lý Vũ Trụ Truyện Quy Mô Lớn (MCU-Scale Shared Universe & Tracing Engine).
// Công cụ giải quyết triệt để bài toán theo vết (tracing), nhất quán quy luật vũ trụ (cross-series canon),
// và quản lý liên minh/bảo vật dùng chung giữa hàng chục bộ truyện song song (tương tự MCU, Star Wars, Warhammer).
package universe

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// EntityEvent lưu lại vết sự kiện của nhân vật/bảo vật trong vũ trụ.
type EntityEvent struct {
	UniverseYear int    `json:"universe_year"` // Năm vũ trụ (ví dụ 1956, 2026, 2099)
	SeriesName   string `json:"series_name"`   // Tên bộ truyện (ví dụ "HUST-sanbaka", "VNU-hivemind")
	Chapter      int    `json:"chapter"`
	Description  string `json:"description"`   // Mô tả biến cố / cột mốc
	Location     string `json:"location"`      // Địa điểm xảy ra
}

// EntityRecord vết lịch sử đầy đủ của 1 nhân vật / bảo vật trong toàn vũ trụ.
type EntityRecord struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"` // "character", "artifact", "faction", "cosmic_law"
	OriginSeries string       `json:"origin_series"`
	CurrentOwner string       `json:"current_owner,omitempty"`
	Status      string        `json:"status"` // "active", "deceased", "sealed", "destroyed"
	History     []EntityEvent `json:"history"`
}

// CanonRule quy tắc định luật tối cao của vũ trụ (không bộ truyện nào được vi phạm).
type CanonRule struct {
	ID          string `json:"id"`
	RuleName    string `json:"rule_name"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "fatal", "strict", "guideline"
}

// UniverseManifest hồ sơ tổng quan của toàn bộ Vũ Trụ.
type UniverseManifest struct {
	UniverseName string         `json:"universe_name"`
	SeriesList   []string       `json:"series_list"`
	CanonRules   []CanonRule    `json:"canon_rules"`
	Entities     []EntityRecord `json:"entities"`
	UpdatedAt    string         `json:"updated_at"`
}

// Manager quản lý vũ trụ truyện và bộ truy vết MCU.
type Manager struct {
	dir string
	mu  sync.RWMutex
}

// NewManager khởi tạo bộ quản lý vũ trụ.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) manifestFile() string {
	return filepath.Join(m.dir, "meta", "universe_manifest.json")
}

// LoadManifest nạp hồ sơ vũ trụ.
func (m *Manager) LoadManifest() (*UniverseManifest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.manifestFile()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &UniverseManifest{
				UniverseName: "SANBAKA Shared Universe",
				SeriesList:   []string{"HUST-sanbaka", "VNU-hivemind", "HMU-biokernel"},
				CanonRules: []CanonRule{
					{
						ID:          "rule-01",
						RuleName:    "Quy tắc bảo toàn năng lượng Cyber-Soul",
						Description: "Năng lượng linh hồn số không tự nhiên sinh ra hay mất đi, chỉ chuyển đổi giữa các chip bán dẫn.",
						Severity:    "fatal",
					},
				},
				Entities:  []EntityRecord{},
				UpdatedAt: time.Now().Format(time.RFC3339),
			}, nil
		}
		return nil, err
	}

	var manifest UniverseManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

// SaveManifest lưu hồ sơ vũ trụ.
func (m *Manager) SaveManifest(manifest *UniverseManifest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	manifest.UpdatedAt = time.Now().Format(time.RFC3339)
	p := m.manifestFile()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// RegisterEvent ghi vết một sự kiện mới của nhân vật/bảo vật trong vũ trụ.
func (m *Manager) RegisterEvent(entityName, entityType, seriesName string, chapter, year int, desc, loc string) error {
	manifest, err := m.LoadManifest()
	if err != nil {
		return err
	}

	found := false
	for i := range manifest.Entities {
		if strings.EqualFold(manifest.Entities[i].Name, entityName) {
			manifest.Entities[i].History = append(manifest.Entities[i].History, EntityEvent{
				UniverseYear: year,
				SeriesName:   seriesName,
				Chapter:      chapter,
				Description:  desc,
				Location:     loc,
			})
			found = true
			break
		}
	}

	if !found {
		manifest.Entities = append(manifest.Entities, EntityRecord{
			Name:         entityName,
			Type:         entityType,
			OriginSeries: seriesName,
			Status:       "active",
			History: []EntityEvent{
				{
					UniverseYear: year,
					SeriesName:   seriesName,
					Chapter:      chapter,
					Description:  desc,
					Location:     loc,
				},
			},
		})
	}

	return m.SaveManifest(manifest)
}

// TraceEntity truy vết toàn bộ lịch sử xuất hiện và biến cố của 1 thực thể qua mọi bộ truyện.
func (m *Manager) TraceEntity(entityName string) (*EntityRecord, error) {
	manifest, err := m.LoadManifest()
	if err != nil {
		return nil, err
	}

	for _, e := range manifest.Entities {
		if strings.EqualFold(e.Name, entityName) || strings.Contains(strings.ToLower(e.Name), strings.ToLower(entityName)) {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("không tìm thấy thực thể '%s' trong cơ sở dữ liệu vũ trụ", entityName)
}

// AuditContinuity kiểm tra xung đột tính nhất quán (Canon Paradoxes) toàn vũ trụ.
func (m *Manager) AuditContinuity() ([]string, error) {
	manifest, err := m.LoadManifest()
	if err != nil {
		return nil, err
	}

	var issues []string

	for _, e := range manifest.Entities {
		if len(e.History) > 1 {
			// Kiểm tra thứ tự thời gian vũ trụ
			for i := 0; i < len(e.History)-1; i++ {
				if e.History[i].UniverseYear > e.History[i+1].UniverseYear {
					issues = append(issues, fmt.Sprintf("⚠️ Nghịch lý thời gian của [%s]: Xuất hiện năm %d (%s Ch.%d) sau khi đã ở năm %d (%s Ch.%d)",
						e.Name, e.History[i].UniverseYear, e.History[i].SeriesName, e.History[i].Chapter,
						e.History[i+1].UniverseYear, e.History[i+1].SeriesName, e.History[i+1].Chapter))
				}
			}
		}
	}

	return issues, nil
}
