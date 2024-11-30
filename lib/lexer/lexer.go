package lexer

import (
	"unicode/utf8"

	"github.com/kijimaD/na2me/lib/token"
)

type Lexer struct {
	input        string
	position     int // 現在検査中のバイトchの位置
	readPosition int // 入力における次の位置
	ch           rune
}

// 初期化する
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readRune()

	return l
}

// 次の1文字を読んでinput文字列の現在位置を進める
func (l *Lexer) readRune() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 文字列の終端を示す
		l.position = l.readPosition
		l.readPosition += 1
	} else {
		// 現在の位置から次のruneを読み取る
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size // runeのサイズ分進める
	}
}

// 現在の1文字を読みこんでトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '\n':
		tok.Literal = string(l.ch)
		tok.Type = token.NEWLINE
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		tok.Literal = l.readSentence()
		tok.Type = token.SENTENCE
		return tok
	}
	l.readRune()

	return tok
}

// 半角スペースを読み飛ばす
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readRune()
	}
}

func (l *Lexer) readSentence() string {
	prePosition := l.position
	for l.ch != '\n' {
		l.readRune()
		if l.ch == 0 {
			break
		}
	}

	return l.input[prePosition:l.position]
}

func AddPageTag(str string) {

}
