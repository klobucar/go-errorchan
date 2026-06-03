package errorchan

// Styled is the pun: it restyles every error flowing through a channel. It
// reads errors from in, wraps each non-nil error in the configured mode, and
// forwards the result on the returned channel. The output channel is closed
// once in is drained and closed.
//
// A single random source is resolved once and shared across the stream, so with
// [WithSeed] or [WithRand] the sequence of styled errors is reproducible while
// still varying error-to-error. nil values are forwarded unchanged (a nil error
// stays a nil error, never a non-nil wrapper around nil).
func Styled(in <-chan error, opts ...Option) <-chan error {
	cfg := resolve(opts...)
	out := make(chan error)
	go func() {
		defer close(out)
		for err := range in {
			if err == nil {
				out <- nil
				continue
			}
			out <- cfg.style(err)
		}
	}()
	return out
}
