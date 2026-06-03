package errorchan

import (
	"math/rand/v2"
	"strings"
	"unicode"
)

// stutterProbability is the per-word chance that a "heavy" transform prepends a
// cutesy stutter (for example "hello" -> "h-hello"). It is consulted only in
// heavy mode and only at the start of a word, so the number of random draws is
// a deterministic function of the input given a fixed source.
const stutterProbability = 0.3

// uwuifyWith applies the phonetic transform to s using rng as its only source
// of randomness. When heavy is true the transform may also insert stutters at
// word boundaries; the light variant performs the same letter substitutions but
// is "too irritated to stutter cutely".
//
// This layer has no knowledge of errors or personas: given the same source and
// input it is fully deterministic, which is what makes seeded output testable.
func uwuifyWith(rng *rand.Rand, s string, heavy bool) string {
	var b strings.Builder
	b.Grow(len(s) + len(s)/4)

	runes := []rune(s)
	atWordStart := true
	for i, r := range runes {
		letter := unicode.IsLetter(r)

		if atWordStart && letter && heavy && rng.Float64() < stutterProbability {
			b.WriteRune(r)
			b.WriteRune('-')
		}

		switch {
		case r == 'r' || r == 'l':
			b.WriteRune('w')
		case r == 'R' || r == 'L':
			b.WriteRune('W')
		case (r == 'n' || r == 'N') && i+1 < len(runes) && isVowel(runes[i+1]):
			// "n" + vowel becomes "ny" + vowel: no -> nyo, na -> nya.
			b.WriteRune(r)
			b.WriteRune('y')
		default:
			b.WriteRune(r)
		}

		atWordStart = !letter
	}

	return b.String()
}

// isVowel reports whether r is an ASCII vowel (either case).
func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}

// styleText runs a persona intro template through the transform while keeping a
// single "%s" placeholder pristine. The placeholder is reserved for the error's
// Go type, which must stay readable, so the segments around it are transformed
// independently and rejoined with the untouched typeName.
func styleText(rng *rand.Rand, template string, heavy bool, typeName string) string {
	if i := strings.Index(template, "%s"); i >= 0 {
		head := uwuifyWith(rng, template[:i], heavy)
		tail := uwuifyWith(rng, template[i+2:], heavy)
		return head + typeName + tail
	}
	return uwuifyWith(rng, template, heavy)
}
