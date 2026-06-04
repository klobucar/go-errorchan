package errorchan

import (
	"fmt"
	"strings"
)

// hiddenTypeNoun replaces the Go type slot in the framing when [WithoutType] is
// set. It is a plain English noun so the phonetic transform can mangle it to fit
// the active persona.
const hiddenTypeNoun = "error"

// separator divides the persona framing from the original error message in
// [StyledError.Error]. It is a single inline glyph so styled errors stay on one
// line and remain greppable in logs.
const separator = " — "

// StyledError wraps an error with personality framing while keeping the
// original error fully intact and reachable. It implements error and Unwrap, so
// [errors.Is], [errors.As], and [errors.Unwrap] all see through to the wrapped
// error.
type StyledError struct {
	// Original is the untouched error that was wrapped.
	Original error
	// OriginalMessage is Original.Error() captured at wrap time, preserved
	// verbatim.
	OriginalMessage string

	mode    string // the mode used to render framing
	framing string // the styled intro plus kaomoji, without the original message
}

// Error returns the persona framing followed by the original message, separated
// by [separator] on a single line. The original message is preserved verbatim.
func (e *StyledError) Error() string {
	return e.framing + separator + e.OriginalMessage
}

// Unwrap returns the original wrapped error, keeping %w-style chains and
// errors.Is/errors.As working through the persona.
func (e *StyledError) Unwrap() error {
	return e.Original
}

// Mode reports the personality mode used to render this error.
func (e *StyledError) Mode() string {
	return e.mode
}

// Framing returns just the persona framing (intro plus kaomoji) without the
// original message, which is occasionally handy for custom formatting.
func (e *StyledError) Framing() string {
	return e.framing
}

// Wrap returns err styled in the configured mode, or nil if err is nil. The
// returned *StyledError unwraps to err, so sentinel and type checks against the
// result keep working.
func Wrap(err error, opts ...Option) *StyledError {
	if err == nil {
		return nil
	}
	return resolve(opts...).style(err)
}

// style renders err into a StyledError using an already-resolved config. It is
// the single core that every public surface funnels through.
func (c config) style(err error) *StyledError {
	p := personaFor(c.mode)
	intro := p.intros[c.rng.IntN(len(p.intros))]
	kaomoji := p.kaomoji[c.rng.IntN(len(p.kaomoji))]

	var framing string
	if c.hideType {
		// Suppress the Go type: swap the type verb for a neutral noun and run the
		// whole intro through the transform, so the noun is uwuified like the rest
		// (for example "error" becomes "ewwow" in heavy modes).
		framing = uwuifyWith(c.rng, strings.ReplaceAll(intro, "%s", hiddenTypeNoun), p.heavy)
	} else {
		framing = styleText(c.rng, intro, p.heavy, typeName(err))
	}
	framing = strings.TrimSpace(framing) + " " + kaomoji

	return &StyledError{
		Original:        err,
		OriginalMessage: err.Error(),
		mode:            c.mode,
		framing:         framing,
	}
}

// typeName renders the Go type of err wrapped in asterisks for emphasis, with
// any leading pointer marker folded in so pointer types do not double up:
// io.EOF becomes "*errors.errorString*" and *fs.PathError becomes
// "*fs.PathError*".
func typeName(err error) string {
	return "*" + strings.TrimPrefix(fmt.Sprintf("%T", err), "*") + "*"
}

// Uwuify applies the raw phonetic transform to s and returns the result, for
// example "really cool" becomes "weawwy coow". The active mode selects the
// transform intensity (heavy for [ModeDere] and [ModeYan], light for
// [ModeTsun]); pass [WithSeed] or [WithRand] for deterministic output.
func Uwuify(s string, opts ...Option) string {
	c := resolve(opts...)
	return uwuifyWith(c.rng, s, personaFor(c.mode).heavy)
}
