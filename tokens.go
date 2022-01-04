package main

type (
	Token interface {
	}

	EOF struct {
		Token
	}

	Number struct {
		Token
		x int
	}

	Ident struct {
		Token
		x string
	}

	Keyword struct {
		Token
		x string
	}

	Delim struct {
		Token
		x rune
	}
)
