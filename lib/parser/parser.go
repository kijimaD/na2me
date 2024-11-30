package parser

import (
	"log"

	"github.com/kijimaD/na2me/lib/ast"
	"github.com/kijimaD/na2me/lib/lexer"
	"github.com/kijimaD/na2me/lib/token"
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
	sentence := &ast.Sentence{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
		return []ast.Node{sentence, &ast.Newpage{Token: p.curToken}}
	}

	return []ast.Node{sentence}
}
