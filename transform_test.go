package errorchan

import (
	"math/rand/v2"
	"strings"
	"testing"
)

// newSeededRand mirrors the source WithSeed builds, for direct transform tests.
func newSeededRand(seed int64) *rand.Rand {
	s := uint64(seed)
	return rand.New(rand.NewPCG(s, s))
}

func TestUwuifyWithLight(t *testing.T) {
	// Light mode never stutters, so it is independent of the source.
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"r and l to w", "really cool", "weawwy coow"},
		{"capitals", "Hello World", "Hewwo Wowwd"},
		{"n plus vowel to ny", "no banana", "nyo banyanya"},
		{"capital n plus vowel", "No", "Nyo"},
		{"n without vowel untouched", "knife", "knyife"},
		{"trailing n untouched", "open", "open"},
		{"no letters", ">_< 123", ">_< 123"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uwuifyWith(newSeededRand(1), tt.in, false)
			if got != tt.want {
				t.Errorf("uwuifyWith(light) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUwuifyWithDeterministic(t *testing.T) {
	// Same seed and input must reproduce byte-for-byte, in both intensities.
	for _, heavy := range []bool{false, true} {
		first := uwuifyWith(newSeededRand(42), "really lovely rolling rivers", heavy)
		second := uwuifyWith(newSeededRand(42), "really lovely rolling rivers", heavy)
		if first != second {
			t.Fatalf("heavy=%v not deterministic: %q vs %q", heavy, first, second)
		}
	}
}

func TestUwuifyWithStutterOnlyHeavy(t *testing.T) {
	const in = "lovely little rabbits running rapidly"
	light := uwuifyWith(newSeededRand(2), in, false)
	if strings.Contains(light, "-") {
		t.Errorf("light mode stuttered: %q", light)
	}

	// Heavy mode is expected to stutter for at least one seed.
	stuttered := false
	for seed := int64(0); seed < 20 && !stuttered; seed++ {
		if strings.Contains(uwuifyWith(newSeededRand(seed), in, true), "-") {
			stuttered = true
		}
	}
	if !stuttered {
		t.Error("heavy mode never stuttered across seeds")
	}
}

func TestUwuifyWithStutterShape(t *testing.T) {
	// A stutter duplicates the first letter and a hyphen, leaving the original
	// letters intact afterwards (for example "really" -> "r-weawwy").
	got := uwuifyWith(newSeededRand(2), "really", true)
	if got != "r-weawwy" {
		t.Fatalf("stutter shape = %q, want %q", got, "r-weawwy")
	}
}

func TestIsVowel(t *testing.T) {
	for _, r := range "aeiouAEIOU" {
		if !isVowel(r) {
			t.Errorf("isVowel(%q) = false, want true", r)
		}
	}
	for _, r := range "bcdfgyZ0!" {
		if isVowel(r) {
			t.Errorf("isVowel(%q) = true, want false", r)
		}
	}
}

func TestStyleTextPreservesPlaceholder(t *testing.T) {
	// The %s placeholder is replaced verbatim and never mangled, even though the
	// surrounding text is transformed.
	got := styleText(newSeededRand(1), "a really %s broke", false, "*io.EOF*")
	if !strings.Contains(got, "*io.EOF*") {
		t.Errorf("placeholder not preserved: %q", got)
	}
	if strings.Contains(got, "really") {
		t.Errorf("surrounding text not transformed: %q", got)
	}
}

func TestStyleTextNoPlaceholder(t *testing.T) {
	got := styleText(newSeededRand(1), "really cool", false, "*io.EOF*")
	if got != "weawwy coow" {
		t.Errorf("styleText without placeholder = %q, want %q", got, "weawwy coow")
	}
}
