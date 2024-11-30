package lexer

import (
	"unicode/utf8"

	"github.com/kijimaD/nova/token"
)

type Lexer struct {
	input        string
	position     int // 現在検査中のバイトchの位置
	readPosition int // 入力における次の位置
	ch           rune
}

// ソースコード文字列を引数に取り、初期化する
func NewLexer(input string) *Lexer {
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
	// case '"':
	// 	tok.Type = token.STRING
	// 	tok.Literal = l.readString()
	default:
		// TODO: とりあえず改行まで取る
		tok.Literal = l.readIdentifier()
		return tok
	}

	l.readRune()
	return tok
}

// 半角スペースを読み飛ばす
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readRune()
	}
}

func (l *Lexer) readIdentifier() string {
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
