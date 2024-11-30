package ast

import (
	"bytes"

	"github.com/kijimaD/na2me/lib/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// ================
type Scenario struct {
	Nodes []Node
}

func (p *Scenario) TokenLiteral() string {
	if len(p.Nodes) > 0 {
		return p.Nodes[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Scenario) String() string {
	var out bytes.Buffer

	for _, n := range p.Nodes {
		out.WriteString(n.String())
	}

	return out.String()
}

// ================
type Newline struct {
	Token token.Token
}

func (n *Newline) TokenLiteral() string {
	return n.Token.Literal
}

func (n *Newline) String() string {
	return n.Token.Literal
}

// ================
type Sentence struct {
	Token token.Token
	Value string
}

func (s *Sentence) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Sentence) String() string {
	return s.Value
}
