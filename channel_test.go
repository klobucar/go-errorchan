package errorchan

import (
	"errors"
	"strings"
	"testing"
)

func TestStyledChannel(t *testing.T) {
	in := make(chan error, 3)
	in <- errors.New("first")
	in <- errors.New("second")
	in <- errors.New("third")
	close(in)

	out := Styled(in, WithMode(ModeTsun), WithSeed(1))

	var got []error
	for err := range out {
		got = append(got, err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d errors, want 3", len(got))
	}
	for i, want := range []string{"first", "second", "third"} {
		var styled *StyledError
		if !errors.As(got[i], &styled) {
			t.Fatalf("element %d is not a *StyledError: %v", i, got[i])
		}
		if styled.OriginalMessage != want {
			t.Errorf("element %d original = %q, want %q", i, styled.OriginalMessage, want)
		}
	}
}

func TestStyledChannelForwardsNil(t *testing.T) {
	in := make(chan error, 2)
	in <- nil
	in <- errors.New("real")
	close(in)

	out := Styled(in, WithSeed(1))
	first := <-out
	if first != nil {
		t.Errorf("nil was not forwarded as nil: %#v", first)
	}
	second := <-out
	if second == nil || !strings.HasSuffix(second.Error(), "real") {
		t.Errorf("second error mangled: %v", second)
	}
}

func TestStyledChannelClosesOutput(t *testing.T) {
	in := make(chan error)
	close(in)
	out := Styled(in, WithSeed(1))
	if _, ok := <-out; ok {
		t.Error("output channel not closed after input closed")
	}
}

func TestStyledChannelDeterministicSequence(t *testing.T) {
	// A seeded stream reproduces its sequence of framings exactly.
	run := func() []string {
		in := make(chan error, 4)
		for _, m := range []string{"a", "b", "c", "d"} {
			in <- errors.New(m)
		}
		close(in)
		var framings []string
		for err := range Styled(in, WithMode(ModeDere), WithSeed(77)) {
			var styled *StyledError
			if errors.As(err, &styled) {
				framings = append(framings, styled.Framing())
			}
		}
		return framings
	}
	first, second := run(), run()
	if strings.Join(first, "|") != strings.Join(second, "|") {
		t.Errorf("seeded stream not reproducible:\n%v\n%v", first, second)
	}
}
