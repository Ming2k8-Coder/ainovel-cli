package llmretry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/voocel/agentcore"
)

type mockGenerator struct {
	calls     int
	responses []*agentcore.LLMResponse
	errors    []error
}

func (m *mockGenerator) Generate(ctx context.Context, msgs []agentcore.Message, tools []agentcore.ToolSpec, opts ...agentcore.CallOption) (*agentcore.LLMResponse, error) {
	idx := m.calls
	m.calls++
	if idx < len(m.errors) && m.errors[idx] != nil {
		return nil, m.errors[idx]
	}
	if idx < len(m.responses) {
		return m.responses[idx], nil
	}
	return &agentcore.LLMResponse{
		Message: agentcore.Message{
			Role:    agentcore.RoleAssistant,
			Content: []agentcore.ContentBlock{agentcore.TextBlock("default")},
		},
	}, nil
}

type mockRetryableErr struct {
	msg string
}

func (e *mockRetryableErr) Error() string   { return e.msg }
func (e *mockRetryableErr) Retryable() bool { return true }

type mockHintErr struct {
	mockRetryableErr
	after time.Duration
}

func (e *mockHintErr) RetryAfter() time.Duration { return e.after }

func TestGenerate_SuccessImmediate(t *testing.T) {
	gen := &mockGenerator{
		responses: []*agentcore.LLMResponse{
			{
				Message: agentcore.Message{
					Role:    agentcore.RoleAssistant,
					Content: []agentcore.ContentBlock{agentcore.TextBlock("hello")},
				},
			},
		},
	}
	resp, err := Generate(context.Background(), gen, Config{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Message.Content) == 0 || resp.Message.Content[0].Text != "hello" {
		t.Fatalf("resp text = %v, want hello", resp.Message.Content)
	}
	if gen.calls != 1 {
		t.Fatalf("expected 1 call, got %d", gen.calls)
	}
}

func TestGenerate_NonRetryableError(t *testing.T) {
	gen := &mockGenerator{
		errors: []error{errors.New("unauthorized 401")},
	}
	_, err := Generate(context.Background(), gen, Config{}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if gen.calls != 1 {
		t.Fatalf("expected 1 call (no retries for non-retryable error), got %d", gen.calls)
	}
}

func TestGenerate_RetryableRecovers(t *testing.T) {
	retryEvents := 0
	cfg := Config{
		Agent: "test-agent",
		OnRetry: func(e Event) {
			retryEvents++
		},
	}

	gen := &mockGenerator{
		errors: []error{
			&mockHintErr{mockRetryableErr: mockRetryableErr{msg: "rate limit 429"}, after: 10 * time.Millisecond},
		},
		responses: []*agentcore.LLMResponse{
			nil,
			{
				Message: agentcore.Message{
					Role:    agentcore.RoleAssistant,
					Content: []agentcore.ContentBlock{agentcore.TextBlock("recovered")},
				},
			},
		},
	}

	resp, err := Generate(context.Background(), gen, cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Message.Content) == 0 || resp.Message.Content[0].Text != "recovered" {
		t.Fatalf("resp text = %v, want recovered", resp.Message.Content)
	}
	if gen.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", gen.calls)
	}
	if retryEvents != 1 {
		t.Fatalf("expected 1 OnRetry event, got %d", retryEvents)
	}
}

func TestGenerate_ContextCanceledDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	gen := &mockGenerator{
		errors: []error{
			&mockHintErr{mockRetryableErr: mockRetryableErr{msg: "server error 503"}, after: 500 * time.Millisecond},
		},
	}

	_, err := Generate(ctx, gen, Config{}, nil)
	if err == nil {
		t.Fatal("expected error when context expires during backoff")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", err)
	}
}
