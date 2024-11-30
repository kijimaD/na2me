package ast

import (
	"bytes"

	"github.com/kijimaD/na2me/lib/convert/token"
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
		_, isNewline := n.(*Newline)
		output := n.String()
		if !isNewline {
			output = output + "\n"
		}
		out.WriteString(output)
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
	return "\n"
}

// ================
type Newpage struct {
	Token token.Token
}

func (n *Newpage) TokenLiteral() string {
	return n.Token.Literal
}

func (n *Newpage) String() string {
	return "[p]"
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
