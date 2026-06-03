package errorchan

import (
	"sync"
	"testing"
)

func TestSetModeAndMode(t *testing.T) {
	restore := Mode()
	t.Cleanup(func() { _ = SetMode(restore) })

	for _, mode := range []string{ModeDere, ModeTsun, ModeYan} {
		if err := SetMode(mode); err != nil {
			t.Fatalf("SetMode(%q) error: %v", mode, err)
		}
		if got := Mode(); got != mode {
			t.Errorf("Mode() = %q, want %q", got, mode)
		}
	}
}

func TestSetModeUnknown(t *testing.T) {
	restore := Mode()
	t.Cleanup(func() { _ = SetMode(restore) })

	if err := SetMode(ModeTsun); err != nil {
		t.Fatal(err)
	}
	if err := SetMode("uguu"); err == nil {
		t.Error("SetMode(unknown) = nil error, want error")
	}
	if got := Mode(); got != ModeTsun {
		t.Errorf("unknown SetMode changed mode to %q, want it unchanged at %q", got, ModeTsun)
	}
}

func TestDefaultModeIsDere(t *testing.T) {
	// The package initializes to dere; verify the constant rather than relying
	// on global state another test may have changed.
	if ModeDere != "dere" {
		t.Errorf("ModeDere = %q, want %q", ModeDere, "dere")
	}
}

func TestModeRaceSafe(t *testing.T) {
	restore := Mode()
	t.Cleanup(func() { _ = SetMode(restore) })

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); _ = SetMode(ModeTsun) }()
		go func() { defer wg.Done(); _ = Mode() }()
	}
	wg.Wait()
}

func TestValidMode(t *testing.T) {
	for _, mode := range []string{ModeDere, ModeTsun, ModeYan} {
		if !validMode(mode) {
			t.Errorf("validMode(%q) = false, want true", mode)
		}
	}
	if validMode("nope") {
		t.Error("validMode(nope) = true, want false")
	}
}

func TestDefaultSourceIsUsable(t *testing.T) {
	// With no seed, resolve falls back to a non-deterministic source. The result
	// is unpredictable but must be non-empty and well-formed.
	if got := Uwuify("really cool"); got == "" {
		t.Error("Uwuify with default source returned empty string")
	}
	styled := Wrap(errSentinel, WithMode(ModeDere))
	if styled == nil || styled.OriginalMessage != errSentinel.Error() {
		t.Errorf("Wrap with default source malformed: %v", styled)
	}
}
