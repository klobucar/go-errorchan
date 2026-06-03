package errorchan_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/klobucar/go-errorchan"
)

func ExampleUwuify() {
	// WithSeed makes the phonetic transform reproducible.
	fmt.Println(errorchan.Uwuify("really cool", errorchan.WithSeed(1)))
	// Output: weawwy coow
}

func ExampleWrap() {
	styled := errorchan.Wrap(io.EOF, errorchan.WithMode(errorchan.ModeTsun), errorchan.WithSeed(1))

	// The persona wraps the error, but the original is fully intact.
	fmt.Println(styled)
	fmt.Println("errors.Is(styled, io.EOF):", errors.Is(styled, io.EOF))
	fmt.Println("original:", styled.OriginalMessage)
	// Output:
	// tch, a *errors.errorString*. finye, hewe's the detaiw. don't make me wepeat it >:( — EOF
	// errors.Is(styled, io.EOF): true
	// original: EOF
}

func ExampleWrap_dere() {
	err := errorchan.Wrap(errors.New("connection refused"), errorchan.WithSeed(7))
	fmt.Println(err)
	// Output: evewything was finye u-untiw this *errors.errorString*... g-gomen nyasai >_< — connection refused
}

func ExampleStyled() {
	in := make(chan error, 2)
	in <- errors.New("timeout")
	in <- errors.New("refused")
	close(in)

	for err := range errorchan.Styled(in, errorchan.WithMode(errorchan.ModeDere), errorchan.WithSeed(5)) {
		fmt.Println(err)
	}
	// Output:
	// a *errors.errorString* came o-out... I-I didn't mean to wet it happen (｡•́︿•̀｡) — timeout
	// s-snyiffwe... a *errors.errorString*... I'l-ww twy hawdew n-nyext t-time, p-pwomise (｡•́︿•̀｡) — refused
}

func ExampleRecover() {
	doThing := func() (err error) {
		defer errorchan.Recover(&err, errorchan.WithMode(errorchan.ModeTsun), errorchan.WithSeed(2))
		panic(errors.New("nil pointer somewhere"))
	}

	err := doThing()
	fmt.Println(err)
	// Output: a *errors.errorString*?! sewiouswy? I'm nyot hewping because I wike you ow anything ヽ(`Д´)ﾉ — nil pointer somewhere
}

func ExampleSetMode() {
	// Set the global default; subsequent Wrap calls use it unless overridden.
	previous := errorchan.Mode()
	defer func() { _ = errorchan.SetMode(previous) }()

	if err := errorchan.SetMode(errorchan.ModeYan); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("mode is now:", errorchan.Mode())
	// Output: mode is now: yan
}
