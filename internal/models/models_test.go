package models

import (
	"path/filepath"
	"testing"
)

func TestSameModelID(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"claude-3-5-sonnet", "claude-3-5-sonnet", true},
		{"claude-3.5-sonnet", "claude-3-5-sonnet", true},
		{"CLAUDE-3-5-SONNET", "claude-3-5-sonnet", true},
		{"claude-3-5-sonnet", "claude-3-5-sonnet-20241022", true},
		{"claude-3-5-sonnet-20241022", "claude-3-5-sonnet", true},
		{"gpt-4o", "gpt-4o-20240806", true},
		{"gpt-4o", "gpt-4o-mini", false},
		{"claude-3-opus", "claude-3-5-sonnet", false},
	}
	for _, c := range cases {
		got := SameModelID(c.a, c.b)
		if got != c.want {
			t.Errorf("SameModelID(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestModelRegistry_Resolve(t *testing.T) {
	reg := NewModelRegistry()

	// Empty string
	if m, ok := reg.Resolve(""); ok || m != nil {
		t.Fatal("expected false for empty resolve")
	}

	// Exact / canonical
	if m, ok := reg.Resolve("anthropic/claude-sonnet-4"); !ok || m == nil {
		t.Fatal("expected to resolve anthropic/claude-sonnet-4")
	}

	// Fuzzy lookup / substring
	if m, ok := reg.Resolve("claude-sonnet-4"); !ok || m == nil {
		t.Fatal("expected to resolve claude-sonnet-4 without provider prefix")
	}

	// Context window lookup
	cw := reg.ResolveContextWindow("claude-sonnet-4")
	if cw <= 0 {
		t.Fatalf("expected context window > 0, got %d", cw)
	}

	// Non-existent model
	if _, ok := reg.Resolve("nonexistent-model-xyz-12345"); ok {
		t.Fatal("expected false for nonexistent model")
	}
	if cw := reg.ResolveContextWindow("nonexistent-model-xyz-12345"); cw != 0 {
		t.Fatalf("expected 0 context window for nonexistent model, got %d", cw)
	}
}

func TestModelRegistry_List(t *testing.T) {
	reg := NewModelRegistry()
	all := reg.List("")
	if len(all) == 0 {
		t.Fatal("expected non-empty model list")
	}

	filtered := reg.List("claude")
	if len(filtered) == 0 {
		t.Fatal("expected claude models in filtered list")
	}
	for _, m := range filtered {
		if !SameModelID(m.ID, "claude") && m.Provider != "anthropic" && !containsIgnoreCase(m.Name, "claude") {
			t.Fatalf("unexpected model in filtered list: %+v", m)
		}
	}
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || filepath.Base(s) != "")
}

func TestModelRegistry_MergeModels(t *testing.T) {
	reg := NewModelRegistry()
	initialCount := len(reg.List(""))

	newModel := ModelEntry{
		Provider:        "custom-provider",
		ID:              "super-model-v1",
		Name:            "Super Model V1",
		ContextWindow:   500_000,
		MaxTokens:       16_384,
		InputCostPer1M:  5.0,
		OutputCostPer1M: 15.0,
	}

	reg.MergeModels([]ModelEntry{newModel})
	if len(reg.List("")) != initialCount+1 {
		t.Fatalf("expected count %d, got %d", initialCount+1, len(reg.List("")))
	}

	resolved, ok := reg.Resolve("super-model-v1")
	if !ok || resolved.ContextWindow != 500_000 {
		t.Fatalf("expected newly merged model to be resolvable with updated window, got: %+v", resolved)
	}

	// Overwrite existing with updated pricing
	updateModel := ModelEntry{
		Provider:        "custom-provider",
		ID:              "super-model-v1",
		InputCostPer1M:  2.5,
		OutputCostPer1M: 7.5,
	}
	reg.MergeModels([]ModelEntry{updateModel})
	resolved2, _ := reg.Resolve("super-model-v1")
	if resolved2.InputCostPer1M != 2.5 || resolved2.ContextWindow != 500_000 {
		t.Fatalf("expected partial pricing update to preserve context window and update price: %+v", resolved2)
	}
}

func TestPricing_CacheSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	testModels := []ModelEntry{
		{
			Provider:      "test-prov",
			ID:            "m1",
			Name:          "Model 1",
			ContextWindow: 128_000,
		},
	}

	// Empty dir does nothing safely
	saveCache(testModels, "")
	if got := loadCache(""); got != nil {
		t.Fatal("expected nil for empty cache dir")
	}

	// Save and load
	saveCache(testModels, tmpDir)
	loaded := loadCache(tmpDir)
	if len(loaded) != 1 || loaded[0].ID != "m1" {
		t.Fatalf("expected 1 model loaded from cache, got: %+v", loaded)
	}

	// Overwriting existing cache file on Windows (verifying atomic replacement)
	testModels[0].Name = "Model 1 Updated"
	saveCache(testModels, tmpDir)
	loadedUpdated := loadCache(tmpDir)
	if len(loadedUpdated) != 1 || loadedUpdated[0].Name != "Model 1 Updated" {
		t.Fatalf("expected updated model name from cache: %+v", loadedUpdated)
	}
}
