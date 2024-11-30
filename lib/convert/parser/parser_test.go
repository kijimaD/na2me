package parser

import (
	"fmt"
	"testing"

	"github.com/kijimaD/na2me/lib/convert/ast"
	"github.com/kijimaD/na2me/lib/convert/lexer"
	"github.com/kijimaD/na2me/lib/convert/token"
	"github.com/stretchr/testify/assert"
)

func TestScenarioPrint1(t *testing.T) {
	l := lexer.New("あああ\n\nいいい\nううう\n")
	p := New(l)
	scenario := p.ParseScenario()
	expect := `あああ
[p]

いいい
[p]
ううう
[p]
`
	assert.Equal(t, expect, scenario.String())
}

func TestScenarioPrint2(t *testing.T) {
	l := lexer.New("あああ\nいいい\nううう")
	p := New(l)
	scenario := p.ParseScenario()
	expect := `あああ
[p]
いいい
[p]
ううう
[p]
`
	assert.Equal(t, expect, scenario.String())
}

func TestParseScenario(t *testing.T) {
	l := lexer.New("あああ\n\nいいい\n\nううう\n")
	p := New(l)
	scenario := p.ParseScenario()

	expect := []ast.Node{
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "あああ"}, Value: "あああ"},
		&ast.Newpage{Token: token.Token{Type: "NEWPAGE", Literal: "NEWPAGE"}},
		&ast.Newline{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "いいい"}, Value: "いいい"},
		&ast.Newpage{Token: token.Token{Type: "NEWPAGE", Literal: "NEWPAGE"}},
		&ast.Newline{Token: token.Token{Type: "NEWLINE", Literal: "\n"}},
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "ううう"}, Value: "ううう"},
		&ast.Newpage{Token: token.Token{Type: "NEWPAGE", Literal: "NEWPAGE"}},
	}
	assert.Equal(t, expect, scenario.Nodes)
}

func TestParseScenario_長い文章は分割する(t *testing.T) {
	l := lexer.New(`ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ。いいい`)
	p := New(l)
	scenario := p.ParseScenario()

	expect := []ast.Node{
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ。いいい"}, Value: "ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ。"},
		&ast.Newpage{Token: token.Token{Type: "NEWPAGE", Literal: "NEWPAGE"}},
		&ast.Sentence{Token: token.Token{Type: "SENTENCE", Literal: "ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ,ああああああああああ。いいい"}, Value: "いいい"},
		&ast.Newpage{Token: token.Token{Type: "NEWPAGE", Literal: "NEWPAGE"}},
	}
	assert.Equal(t, expect, scenario.Nodes)
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
				"けれども事実は事実で詐る訳には行かないから、吾輩は「実はとろうとろうと思ってまだ捕らない」",
				"と答えた。黒は彼の鼻の先からぴんと突張っている長い髭をびりびりと震わせて非常に笑った。",
				"元来黒は自慢をする丈にどこか足りないところがあって、彼の気焔を感心したように咽喉をころころ鳴らして謹聴していればはなはだ御しやすい猫である。",
			},
		},
		{
			inputText: "「君歯をどうかしたかね」と主人は問題を転じた。「ええ実はある所で椎茸を食いましてね」「何を食ったって？」「その、少し椎茸を食ったんで。椎茸の傘を前歯で噛み切ろうとしたらぼろりと歯が欠けましたよ」「椎茸で前歯がかけるなんざ、何だか爺々臭いね。俳句にはなるかも知れないが、恋にはならんようだな」と平手で吾輩の頭を軽く叩く。「ああその猫が例のですか、なかなか肥ってるじゃありませんか、それなら車屋の黒にだって負けそうもありませんね、立派なものだ」と寒月君は大に吾輩を賞める。",
			inputLen:  10,
			expect: []string{
				"「君歯をどうかしたかね」",
				"と主人は問題を転じた。",
				"「ええ実はある所で椎茸を食いましてね」",
				"「何を食ったって？」「その、少し椎茸を食ったんで。椎茸の傘を前歯で噛み切ろうとしたらぼろりと歯が欠けましたよ」",
				"「椎茸で前歯がかけるなんざ、何だか爺々臭いね。俳句にはなるかも知れないが、恋にはならんようだな」",
				"と平手で吾輩の頭を軽く叩く。",
				"「ああその猫が例のですか、なかなか肥ってるじゃありませんか、それなら車屋の黒にだって負けそうもありませんね、立派なものだ」",
				"と寒月君は大に吾輩を賞める。",
			},
		},
		{
			inputText: "「何おめでてえ？　正月でおめでたけりゃ、御めえなんざあ年が年中おめでてえ方だろう。気をつけろい、この吹い子の向う面め」吹い子の向うづらという句は罵詈の言語であるようだが、吾輩には了解が出来なかった。「ちょっと伺がうが吹い子の向うづらと云うのはどう云う意味かね」「へん、手めえが悪体をつかれてる癖に、その訳を聞きゃ世話あねえ、だから正月野郎だって事よ」正月野郎は詩的であるが、その意味に至ると吹い子の何とかよりも一層不明瞭な文句である。参考のためちょっと聞いておきたいが、聞いたって明瞭な答弁は得られぬに極まっているから、面と対ったまま無言で立っておった。いささか手持無沙汰の体である。すると突然黒のうちの神さんが大きな声を張り揚げて「おや棚へ上げて置いた鮭がない。大変だ。またあの黒の畜生が取ったんだよ。ほんとに憎らしい猫だっちゃありゃあしない。今に帰って来たら、どうするか見ていやがれ」と怒鳴る。",
			inputLen:  10,
			expect: []string{
				"「何おめでてえ？　正月でおめでたけりゃ、御めえなんざあ年が年中おめでてえ方だろう。気をつけろい、この吹い子の向う面め」",
				"吹い子の向うづらという句は罵詈の言語であるようだが、吾輩には了解が出来なかった。",
				"「ちょっと伺がうが吹い子の向うづらと云うのはどう云う意味かね」",
				"「へん、手めえが悪体をつかれてる癖に、その訳を聞きゃ世話あねえ、だから正月野郎だって事よ」",
				"正月野郎は詩的であるが、その意味に至ると吹い子の何とかよりも一層不明瞭な文句である。",
				"参考のためちょっと聞いておきたいが、聞いたって明瞭な答弁は得られぬに極まっているから、面と対ったまま無言で立っておった。",
				"いささか手持無沙汰の体である。",
				"すると突然黒のうちの神さんが大きな声を張り揚げて「おや棚へ上げて置いた鮭がない。大変だ。またあの黒の畜生が取ったんだよ。ほんとに憎らしい猫だっちゃありゃあしない。今に帰って来たら、どうするか見ていやがれ」",
				"と怒鳴る。",
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
