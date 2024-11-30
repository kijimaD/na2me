package parser

import (
	"fmt"
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

func TestSplitByPeriod(t *testing.T) {
	tests := []struct {
		inputText string
		inputLen  int
		expect    []string
	}{
		{
			inputText: "あああああ。いいいいい。う。",
			inputLen:  5,
			expect: []string{
				"あああああ。",
				"いいいいい。",
				"う。",
			},
		},
		{
			inputText: "あああああ。いいいいい。う。",
			inputLen:  10,
			expect: []string{
				"あああああ。いいいいい。",
				"う。",
			},
		},
		{
			inputText: "あああああ。いいいいい。う。",
			inputLen:  100,
			expect: []string{
				"あああああ。いいいいい。う。",
			},
		},
		{
			inputText: "けれども事実は事実で詐る訳には行かないから、吾輩は「実はとろうとろうと思ってまだ捕らない」と答えた。黒は彼の鼻の先からぴんと突張っている長い髭をびりびりと震わせて非常に笑った。元来黒は自慢をする丈にどこか足りないところがあって、彼の気焔を感心したように咽喉をころころ鳴らして謹聴していればはなはだ御しやすい猫である。",
			inputLen:  10,
			expect: []string{
				"けれども事実は事実で詐る訳には行かないから、吾輩は「実はとろうとろうと思ってまだ捕らない」と答えた。",
				"黒は彼の鼻の先からぴんと突張っている長い髭をびりびりと震わせて非常に笑った。",
				"元来黒は自慢をする丈にどこか足りないところがあって、彼の気焔を感心したように咽喉をころころ鳴らして謹聴していればはなはだ御しやすい猫である。",
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("入力が%sの場合", tt.inputText), func(t *testing.T) {
			got := splitByPeriod(tt.inputText, tt.inputLen)
			assert.Equal(t, tt.expect, got)
		})
	}
}
