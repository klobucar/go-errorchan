package errorchan

import (
	"context"
	"log/slog"
)

// slogHandler is a slog.Handler middleware that styles error-valued attributes
// as they pass through to a base handler.
type slogHandler struct {
	base slog.Handler
	opts []Option
}

// NewSlogHandler wraps base so that any error-valued attribute in a log record
// is restyled with the configured mode before being passed on. Non-error
// attributes are left untouched, and attributes nested in groups are styled
// recursively.
//
// Because each error is styled independently, pass [WithSeed] for reproducible
// output. Avoid [WithRand] here: a handler may be invoked from many goroutines
// at once, and a shared *rand.Rand is not safe for concurrent use.
func NewSlogHandler(base slog.Handler, opts ...Option) slog.Handler {
	return &slogHandler{base: base, opts: opts}
}

// Enabled defers to the base handler.
func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

// Handle rebuilds the record with error attributes styled, then forwards it to
// the base handler.
func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	styled := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.Attrs(func(a slog.Attr) bool {
		styled.AddAttrs(h.styleAttr(a))
		return true
	})
	return h.base.Handle(ctx, styled)
}

// WithAttrs styles the supplied attributes and returns a handler that carries
// them on the base handler.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	styled := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		styled[i] = h.styleAttr(a)
	}
	return &slogHandler{base: h.base.WithAttrs(styled), opts: h.opts}
}

// WithGroup returns a handler that opens a group on the base handler.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{base: h.base.WithGroup(name), opts: h.opts}
}

// styleAttr returns a, with any error value styled and any group recursed into.
func (h *slogHandler) styleAttr(a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		if err, ok := a.Value.Any().(error); ok && err != nil {
			return slog.Any(a.Key, Wrap(err, h.opts...))
		}
	case slog.KindGroup:
		group := a.Value.Group()
		styled := make([]slog.Attr, len(group))
		for i, ga := range group {
			styled[i] = h.styleAttr(ga)
		}
		return slog.Attr{Key: a.Key, Value: slog.GroupValue(styled...)}
	}
	return a
}
