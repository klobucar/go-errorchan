package errorchan

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

// newTestLogger returns a logger whose error attributes are styled, writing to
// the returned buffer in text form.
func newTestLogger(buf *bytes.Buffer, opts ...Option) *slog.Logger {
	base := slog.NewTextHandler(buf, &slog.HandlerOptions{
		// Drop time so output is stable across runs.
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(NewSlogHandler(base, opts...))
}

func TestSlogHandlerStylesErrorAttr(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf, WithMode(ModeTsun), WithSeed(1))

	log.Error("request failed", slog.Any("err", errSentinel))

	out := buf.String()
	if !strings.Contains(out, "sentinel boom") {
		t.Errorf("original message missing from log output: %q", out)
	}
	if !strings.Contains(out, separator) {
		t.Errorf("styled framing missing from log output: %q", out)
	}
}

func TestSlogHandlerLeavesNonErrorAttrs(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf, WithSeed(1))

	log.Info("hello", slog.String("user", "really"), slog.Int("count", 3))

	out := buf.String()
	if !strings.Contains(out, "user=really") {
		t.Errorf("non-error string attr was altered: %q", out)
	}
	if !strings.Contains(out, "count=3") {
		t.Errorf("non-error int attr was altered: %q", out)
	}
}

func TestSlogHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf, WithSeed(1)).With(slog.Any("err", errSentinel))

	log.Info("carried")

	if !strings.Contains(buf.String(), "sentinel boom") {
		t.Errorf("error in WithAttrs was not styled/preserved: %q", buf.String())
	}
}

func TestSlogHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf, WithSeed(1)).WithGroup("op")

	log.Error("grouped", slog.Any("err", errSentinel))

	out := buf.String()
	if !strings.Contains(out, "op.err=") {
		t.Errorf("group prefix missing: %q", out)
	}
	if !strings.Contains(out, "sentinel boom") {
		t.Errorf("error inside group not preserved: %q", out)
	}
}

func TestSlogHandlerNestedGroupAttr(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf, WithSeed(1))

	log.Error("nested", slog.Group("ctx", slog.Any("err", errSentinel), slog.String("k", "v")))

	out := buf.String()
	if !strings.Contains(out, "sentinel boom") {
		t.Errorf("error nested in group not preserved: %q", out)
	}
	if !strings.Contains(out, "ctx.k=v") {
		t.Errorf("sibling attr in group lost: %q", out)
	}
}

func TestSlogHandlerEnabled(t *testing.T) {
	base := slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{Level: slog.LevelWarn})
	h := NewSlogHandler(base)
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Enabled(Info) = true, want false for a Warn base handler")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("Enabled(Error) = false, want true")
	}
}

func TestSlogHandlerConcurrent(t *testing.T) {
	// Seeded option builds a fresh source per record, so concurrent logging is
	// race-free under go test -race.
	var mu sync.Mutex
	var buf bytes.Buffer
	base := slog.NewTextHandler(&lockedWriter{mu: &mu, buf: &buf}, nil)
	log := slog.New(NewSlogHandler(base, WithSeed(5)))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Error("boom", slog.Any("err", errors.New("parallel")))
		}()
	}
	wg.Wait()
}

// lockedWriter serializes writes so the text handler's buffer is not itself the
// thing being raced on in the concurrency test.
type lockedWriter struct {
	mu  *sync.Mutex
	buf *bytes.Buffer
}

func (w *lockedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}
