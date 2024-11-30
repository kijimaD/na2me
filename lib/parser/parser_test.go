package parser

import (
	"testing"

	"github.com/kijimaD/na2me/lib/ast"
	"github.com/kijimaD/na2me/lib/lexer"
	"github.com/kijimaD/na2me/lib/token"
	"github.com/stretchr/testify/assert"
)

func TestParseScenario(t *testing.T) {
	l := lexer.New("あああ\n\nいいい\nううう\n")
	p := New(l)
	scenario := p.ParseScenario()
	nodes := []ast.Node{}
	for _, n := range scenario.Nodes {
		nodes = append(nodes, n)
	}

	expect := []ast.Node{
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "あああ"}, Value: "あああ"},
		&ast.Newpage{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
		&ast.Newline{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "いいい"}, Value: "いいい"},
		&ast.Newpage{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "ううう"}, Value: "ううう"},
		&ast.Newpage{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
	}

	assert.Equal(t, expect, nodes)
}
