package parser

import (
	"log"
	"strings"

	"github.com/kijimaD/na2me/lib/convert/ast"
	"github.com/kijimaD/na2me/lib/convert/lexer"
	"github.com/kijimaD/na2me/lib/convert/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token // 現在のトークン
	peekToken token.Token // 次のトークン
	parseFns  map[token.TokenType]parseFn
}

type parseFn func() []ast.Node

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 前置トークン
	p.parseFns = make(map[token.TokenType]parseFn)
	p.registerFunc(token.NEWLINE, p.parseNewline)
	p.registerFunc(token.SENTENCE, p.parseSentence)

	// 2つトークンを読み込む。curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

// 次のトークンに進む
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 次のトークンと引数の型を比較する
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// パースを開始する。トークンを1つずつ辿る
func (p *Parser) ParseScenario() *ast.Scenario {
	scenario := &ast.Scenario{}
	scenario.Nodes = []ast.Node{}

	for p.curToken.Type != token.EOF {
		node := p.parseNode()
		if node != nil {
			scenario.Nodes = append(scenario.Nodes, node...)
		}
		p.nextToken()
	}

	return scenario
}

func (p *Parser) parseNode() []ast.Node {
	f := p.parseFns[p.curToken.Type]
	if f == nil {
		// p.noPrefixParseFnError(p.curToken.Type)
		// return nil
		log.Fatalf("no parse function %s", string(p.curToken.Type))
	}
	result := f()

	return result
}

func (p *Parser) registerFunc(t token.TokenType, f parseFn) {
	p.parseFns[t] = f
}

func (p *Parser) parseNewline() []ast.Node {
	return []ast.Node{&ast.Newline{Token: p.curToken}}
}

func (p *Parser) parseSentence() []ast.Node {
	const sentenceLength = 50
	const forceLength = 100
	nodes := []ast.Node{}
	splits := splitByPeriod(p.curToken.Literal, sentenceLength, forceLength)
	for _, s := range splits {
		nodes = append(nodes, &ast.Sentence{Token: p.curToken, Value: s})
		nodes = append(nodes, &ast.Newpage{Token: token.Token{Type: token.NEWPAGE, Literal: token.NEWPAGE}})
	}

	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	return nodes
}

// TODO: 先読みをするようにしたほうがよさそう
// 指定文字以上経過後に"。"で分割する
// minLengthを超えると、"。"や"」"で改行する。括弧の中では改行しない
// forceLengthを超えると、括弧の中でも読点で分割を行う
func splitByPeriod(text string, minLength int, forceLength int) []string {
	const jpPeriodChar = '。'
	const jpCommaChar = '、'
	const jpBracketStartChar1 = '「'
	const jpBracketEndChar1 = '」'
	const parentStartChar = '('
	const parentEndChar = ')'
	const jpBracketStartChar2 = '［' // 半角の "[" とは違うことに注意!!
	const jpBracketEndChar2 = '］'   // 半角の "]" とは違うことに注意!!
	result := []string{}
	current := ""

	// 会話
	var insideJPBrackets1 bool
	// 大体脚注用なので中では改行しない
	var insideJPBrackets2 bool
	// 大体脚注用なので中では改行しない
	var insideParentheses bool

	for _, char := range text {
		current += string(char)

		if char == jpBracketStartChar1 {
			insideJPBrackets1 = true
		} else if char == jpBracketEndChar1 {
			insideJPBrackets1 = false
		}
		if char == parentStartChar {
			insideParentheses = true
		} else if char == parentEndChar {
			insideParentheses = false
		}
		if char == jpBracketStartChar2 {
			insideJPBrackets2 = true
		} else if char == jpBracketEndChar2 {
			insideJPBrackets2 = false
		}

		if insideParentheses || insideJPBrackets2 {
			continue
		}

		if (char == jpPeriodChar || char == jpBracketEndChar1) && !insideJPBrackets1 && len([]rune(current)) > minLength {
			result = append(result, strings.TrimSpace(current))
			current = ""
		} else if char == jpPeriodChar && len([]rune(current)) > forceLength {
			result = append(result, strings.TrimSpace(current))
			current = ""
		} else if char == jpCommaChar && len([]rune(current)) > forceLength {
			result = append(result, strings.TrimSpace(current))
			current = ""
		}
	}

	if len(strings.TrimSpace(current)) > 0 {
		result = append(result, strings.TrimSpace(current))
	}

	return result
}
