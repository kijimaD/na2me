package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF      = "EOF"
	SENTENCE = "SENTENCE"
	NEWLINE  = "NEWLINE"
	NEWPAGE  = "NEWPAGE"
)
