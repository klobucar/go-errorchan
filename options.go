package errorchan

import "math/rand/v2"

// Option configures a styling operation. Options are applied with the
// functional-options pattern and are accepted by every public surface
// ([Uwuify], [Wrap], [Styled], [Recover], [NewSlogHandler]).
type Option func(*config)

// config is the resolved configuration for a single styling operation.
type config struct {
	rng      *rand.Rand
	mode     string
	hideType bool
}

// WithMode selects the personality mode for this operation, overriding the
// global default from [Mode]. An unrecognized mode falls back to [ModeDere].
func WithMode(mode string) Option {
	return func(c *config) {
		c.mode = mode
	}
}

// WithoutType suppresses the error's Go type in the persona framing. Normally
// the framing names the wrapped type (for example *errors.errorString*) to aid
// debugging; with this option the type slot is replaced by a neutral noun that
// is run through the same phonetic transform as the rest of the framing. The
// wrapped error and its verbatim message are unaffected. Useful when the framing
// is shown to end users who should not see Go internals.
func WithoutType() Option {
	return func(c *config) {
		c.hideType = true
	}
}

// WithSeed makes the output deterministic by deriving a fresh random source
// from seed. Each operation that resolves this option gets its own source, so
// it is safe to reuse the same seeded option across goroutines.
func WithSeed(seed int64) Option {
	return func(c *config) {
		s := uint64(seed)
		c.rng = rand.New(rand.NewPCG(s, s))
	}
}

// WithRand makes the output deterministic by using the supplied source. Unlike
// [WithSeed], the same *rand.Rand is shared by every operation that resolves
// this option; callers are responsible for not using it concurrently from
// multiple goroutines.
func WithRand(r *rand.Rand) Option {
	return func(c *config) {
		c.rng = r
	}
}

// resolve builds a config from the global defaults and the supplied options. An
// unknown mode is coerced to [ModeDere] and a missing source is replaced with a
// non-deterministic one, so the returned config is always usable.
func resolve(opts ...Option) config {
	c := config{mode: Mode()}
	for _, opt := range opts {
		opt(&c)
	}
	if !validMode(c.mode) {
		c.mode = ModeDere
	}
	if c.rng == nil {
		c.rng = newDefaultRand()
	}
	return c
}

// newDefaultRand returns a non-deterministic source seeded from the
// auto-seeded top-level generator, which is itself safe for concurrent use.
func newDefaultRand() *rand.Rand {
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}
