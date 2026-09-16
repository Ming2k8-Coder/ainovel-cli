package eval

import (
	"path/filepath"
	"testing"
)

func TestRunDir(t *testing.T) {
	out := "/tmp/eval"
	cases := []struct {
		arm    string
		repeat int
		total  int
		want   string
	}{
		{ArmSingle, 1, 1, filepath.Join(out, "artifacts", "case1")},
		{ArmBaseline, 1, 1, filepath.Join(out, "artifacts", "case1", "baseline")},
		{ArmVariant, 2, 3, filepath.Join(out, "artifacts", "case1", "r2", "variant")},
		{ArmSingle, 2, 3, filepath.Join(out, "artifacts", "case1", "r2")},
	}
	for _, tc := range cases {
		if got := runDir(out, "case1", tc.arm, tc.repeat, tc.total); got != tc.want {
			t.Fatalf("runDir(%s,%d,%d) = %s want %s", tc.arm, tc.repeat, tc.total, got, tc.want)
		}
	}
}

func TestCommandRejectsInvalidRepeat(t *testing.T) {
	code := Command([]string{"--cases", "evals/cases/smoke", "--repeat", "0"})
	if code != 2 {
		t.Fatalf("repeat=0 应返回用法错误 2，得到 %d", code)
	}
}

// TestIssue10_MultiAgentEvaluationDelta verifies Issue #10:
// A/B delta comparison computes metrics delta between baseline and variant runs.
func TestIssue10_MultiAgentEvaluationDelta(t *testing.T) {
	c := writerSmokeCase()
	base := cleanResult()
	base.Metrics.TotalWords = 5000
	base.Metrics.Usage.CostUSD = 0.10
	base.Metrics.Usage.UsageRecorded = true

	variant := cleanResult()
	variant.Metrics.TotalWords = 6000
	variant.Metrics.Usage.CostUSD = 0.08 // More words, lower cost
	variant.Metrics.Usage.UsageRecorded = true

	delta := GradeDelta(c, base, variant)
	if delta.Metrics.CostDeltaRatio >= 0 {
		t.Fatalf("expected negative cost delta ratio, got %f", delta.Metrics.CostDeltaRatio)
	}
	if delta.Metrics.TotalWordsRatio <= 1.0 {
		t.Fatalf("expected words ratio > 1.0, got %f", delta.Metrics.TotalWordsRatio)
	}
}

