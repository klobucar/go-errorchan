package errorchan

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

// sentinel and a custom error type for Is/As checks.
var errSentinel = errors.New("sentinel boom")

type customError struct{ code int }

func (e *customError) Error() string { return fmt.Sprintf("custom failure code=%d", e.code) }

func TestWrapNil(t *testing.T) {
	if got := Wrap(nil); got != nil {
		t.Errorf("Wrap(nil) = %v, want nil", got)
	}
}

func TestWrapPreservesOriginal(t *testing.T) {
	orig := errors.New("disk on fire")
	styled := Wrap(orig, WithSeed(1))

	if styled.Original != orig { //nolint:errorlint // asserting exact identity of the stored field, not chain membership
		t.Errorf("Original = %v, want %v", styled.Original, orig)
	}
	if styled.OriginalMessage != "disk on fire" {
		t.Errorf("OriginalMessage = %q, want %q", styled.OriginalMessage, "disk on fire")
	}
	if !strings.HasSuffix(styled.Error(), separator+"disk on fire") {
		t.Errorf("Error() = %q, want it to end with the verbatim message", styled.Error())
	}
}

func TestWrapSingleLine(t *testing.T) {
	// Styled errors must stay on one line so they do not wreck logs.
	styled := Wrap(errors.New("line one\nline two"), WithSeed(1))
	framing := styled.Framing()
	if strings.ContainsAny(framing, "\n") {
		t.Errorf("framing contains a newline: %q", framing)
	}
}

func TestWrapUnwrap(t *testing.T) {
	orig := errors.New("boom")
	styled := Wrap(orig, WithSeed(1))
	if got := errors.Unwrap(styled); got != orig { //nolint:errorlint // Unwrap must return the exact wrapped error
		t.Errorf("errors.Unwrap = %v, want %v", got, orig)
	}
}

func TestWrapErrorsIs(t *testing.T) {
	styled := Wrap(fmt.Errorf("context: %w", errSentinel), WithMode(ModeTsun), WithSeed(3))
	if !errors.Is(styled, errSentinel) {
		t.Errorf("errors.Is did not see through the styled error")
	}
	if errors.Is(styled, io.EOF) {
		t.Errorf("errors.Is matched an unrelated sentinel")
	}
}

func TestWrapErrorsAs(t *testing.T) {
	styled := Wrap(fmt.Errorf("wrapped: %w", &customError{code: 7}), WithMode(ModeYan), WithSeed(4))
	var target *customError
	if !errors.As(styled, &target) {
		t.Fatalf("errors.As did not extract the custom error")
	}
	if target.code != 7 {
		t.Errorf("extracted code = %d, want 7", target.code)
	}
}

func TestWrapKeepsWrapChain(t *testing.T) {
	// Re-wrapping a styled error with %w must keep both the sentinel and the
	// styled type reachable.
	styled := Wrap(errSentinel, WithSeed(1))
	rewrapped := fmt.Errorf("outer: %w", styled)
	if !errors.Is(rewrapped, errSentinel) {
		t.Errorf("sentinel lost after re-wrap")
	}
	var target *StyledError
	if !errors.As(rewrapped, &target) {
		t.Errorf("StyledError lost after re-wrap")
	}
}

func TestWrapMode(t *testing.T) {
	styled := Wrap(errors.New("x"), WithMode(ModeYan), WithSeed(1))
	if styled.Mode() != ModeYan {
		t.Errorf("Mode() = %q, want %q", styled.Mode(), ModeYan)
	}
}

func TestWrapUsesGlobalModeByDefault(t *testing.T) {
	restore := Mode()
	t.Cleanup(func() { _ = SetMode(restore) })

	if err := SetMode(ModeTsun); err != nil {
		t.Fatal(err)
	}
	styled := Wrap(errors.New("x"), WithSeed(1))
	if styled.Mode() != ModeTsun {
		t.Errorf("Mode() = %q, want global %q", styled.Mode(), ModeTsun)
	}
}

func TestWrapUnknownModeFallsBackToDere(t *testing.T) {
	styled := Wrap(errors.New("x"), WithMode("uguu"), WithSeed(1))
	if styled.Mode() != ModeDere {
		t.Errorf("unknown mode resolved to %q, want %q", styled.Mode(), ModeDere)
	}
}

func TestWrapEmbedsType(t *testing.T) {
	styled := Wrap(io.EOF, WithSeed(1))
	if !strings.Contains(styled.Framing(), "*errors.errorString*") {
		t.Errorf("framing %q does not embed the error type", styled.Framing())
	}
}

func TestEveryModeProducesNonEmptyFraming(t *testing.T) {
	for _, mode := range []string{ModeDere, ModeTsun, ModeYan} {
		// Sweep seeds so we exercise every intro/kaomoji combination.
		for seed := int64(0); seed < 50; seed++ {
			styled := Wrap(errors.New("base"), WithMode(mode), WithSeed(seed))
			if strings.TrimSpace(styled.Framing()) == "" {
				t.Fatalf("mode %q seed %d produced empty framing", mode, seed)
			}
			if !strings.HasSuffix(styled.Error(), "base") {
				t.Fatalf("mode %q seed %d dropped the original message: %q", mode, seed, styled.Error())
			}
		}
	}
}

func TestTsunSaysBaka(t *testing.T) {
	// About half of the tsun intros call you baka; across a seed sweep it must
	// show up. (The transform leaves "baka" unchanged.)
	found := false
	for seed := int64(0); seed < 50 && !found; seed++ {
		if strings.Contains(Wrap(errors.New("x"), WithMode(ModeTsun), WithSeed(seed)).Error(), "baka") {
			found = true
		}
	}
	if !found {
		t.Error("tsun mode never said baka across 50 seeds")
	}
}

func TestTsunBakaRoughlyHalf(t *testing.T) {
	// Guard the design rule directly against the persona data.
	baka := 0
	for _, intro := range personas[ModeTsun].intros {
		if strings.Contains(intro, "baka") {
			baka++
		}
	}
	total := len(personas[ModeTsun].intros)
	if baka*4 < total || baka*4 > total*3 {
		t.Errorf("tsun baka count = %d/%d, want roughly half", baka, total)
	}
}

func TestUwuifyPublic(t *testing.T) {
	if got := Uwuify("really cool", WithSeed(1)); got != "weawwy coow" {
		t.Errorf("Uwuify = %q, want %q", got, "weawwy coow")
	}
}

func TestUwuifyDeterministicWithRand(t *testing.T) {
	a := Uwuify("rolling rivers", WithRand(newSeededRand(9)))
	b := Uwuify("rolling rivers", WithRand(newSeededRand(9)))
	if a != b {
		t.Errorf("WithRand not reproducible: %q vs %q", a, b)
	}
}

func TestPersonaForFallback(t *testing.T) {
	if got := personaFor("nope").heavy; got != personas[ModeDere].heavy {
		t.Errorf("personaFor unknown did not fall back to dere")
	}
}

func TestWrapWithoutType(t *testing.T) {
	orig := errors.New("connection refused")
	styled := Wrap(orig, WithSeed(42), WithoutType())

	// The Go type must not leak into the framing or the rendered string.
	if strings.Contains(styled.Framing(), "*") {
		t.Errorf("framing still contains a type token: %q", styled.Framing())
	}
	if strings.Contains(styled.Error(), "errorString") {
		t.Errorf("rendered error leaks the Go type: %q", styled.Error())
	}
	// The original message is still preserved verbatim, after the separator.
	if styled.OriginalMessage != "connection refused" {
		t.Errorf("OriginalMessage = %q, want verbatim", styled.OriginalMessage)
	}
	if !strings.HasSuffix(styled.Error(), separator+"connection refused") {
		t.Errorf("rendered error missing verbatim tail: %q", styled.Error())
	}
	// errors.Is still sees through the persona.
	if !errors.Is(styled, orig) {
		t.Error("errors.Is should see through WithoutType styling")
	}
}

func TestWithoutTypeAllModes(t *testing.T) {
	for _, mode := range []string{ModeDere, ModeTsun, ModeYan} {
		styled := Wrap(io.EOF, WithMode(mode), WithoutType(), WithSeed(7))
		if strings.Contains(styled.Framing(), "*") {
			t.Errorf("mode %s: type token leaked: %q", mode, styled.Framing())
		}
	}
}
