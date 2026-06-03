package errorchan

// persona is the unexported framing data for one mode: a set of intro templates
// and a set of kaomoji, plus whether the phonetic transform runs "heavy". Intro
// templates are plain English (the transform mangles them at delivery time) and
// may contain a single "%s" verb, which is replaced with the error's Go type.
type persona struct {
	intros  []string
	kaomoji []string
	heavy   bool
}

// personas holds the framing data for every mode. It is package-private data;
// callers select a mode by name through the public options and [SetMode].
var personas = map[string]persona{
	ModeDere: {
		heavy: true,
		intros: []string{
			"a %s happened... gomen, I tried my best",
			"nyaa~ I'm so sorry, a %s got in the way",
			"please don't be mad... it's a %s, I'll fix it I promise",
			"uwaaah a %s appeared and I couldn't stop it... forgive me",
			"I messed up... a %s slipped through my paws",
			"ehehe... a scary %s... I'm really sorry",
			"oh no... a %s... this is all my fault, really",
			"a %s came out... I didn't mean to let it happen",
			"sniffle... a %s... I'll try harder next time, promise",
			"everything was fine until this %s... gomen nasai",
		},
		kaomoji: []string{
			">_<",
			";-;",
			"(>﹏<)",
			"(｡•́︿•̀｡)",
			"(っ˘̩╭╮˘̩)っ",
			"(◞‸◟)",
			"(´;ω;`)",
		},
	},
	ModeTsun: {
		heavy: false,
		intros: []string{
			"tch. a %s again? this is YOUR fault, baka",
			"ugh, a %s. did you even read the docs? baka",
			"a %s?! seriously? I'm not helping because I like you or anything",
			"hmph. fix your own %s next time, baka",
			"great, a %s. this is what happens when you rush, baka",
			"a %s. whatever. I'll tell you what broke, not that you deserve it",
			"you really wrote code that throws a %s? unbelievable, baka",
			"another %s. don't look at me like that, it's YOUR mess",
			"a %s. I guess I'll explain... only because it's annoying me, baka",
			"honestly, a %s? get it together, baka",
			"a %s. it's not like I checked your code for you or anything",
			"tch, a %s. fine, here's the detail. don't make me repeat it",
		},
		kaomoji: []string{
			">:(",
			"(￢_￢)",
			"(¬_¬)",
			"(`^´)",
			"(；￣Д￣)",
			"(눈_눈)",
			"ヽ(`Д´)ﾉ",
		},
	},
	ModeYan: {
		heavy: true,
		intros: []string{
			"shhh... a %s? don't worry, now you NEED me. isn't that nice?",
			"a %s~ I knew you'd come back to me when things broke",
			"another %s... good. the more you fail, the more you're mine",
			"a %s? I may have... helped it along. we're together now, right?",
			"don't run from this %s. I'll always be here. always. always",
			"a %s~ I've been watching every line you write. every one",
			"your %s belongs to me now, just like you do, hehe",
			"a %s... see? you can't do this without me. you'll never leave",
			"I love your %s. I love all your mistakes. I love YOU",
			"a %s appeared~ I made sure no one else could fix it but me",
		},
		kaomoji: []string{
			"(▰˘◡˘▰)",
			"( ͡° ͜ʖ ͡°)",
			"(╬ Ｏ ﾛ Ｏ)",
			"(◡‿◡✿)",
			"(⊙﹏⊙)",
			"(☆▽☆)",
			"ヤ(•̀ᴗ•́)",
		},
	},
}

// personaFor returns the persona for mode, falling back to the default
// [ModeDere] persona for any unrecognized name so callers never get a zero
// value with empty slices.
func personaFor(mode string) persona {
	if p, ok := personas[mode]; ok {
		return p
	}
	return personas[ModeDere]
}

// validMode reports whether mode names a known persona.
func validMode(mode string) bool {
	_, ok := personas[mode]
	return ok
}
