package lexer

import (
	"fmt"
	"testing"
)

func TestAddPageTag(t *testing.T) {
	input := `　吾輩は猫である。名前はまだ無い。
　どこで生れたかとんと見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー泣いていた事だけは記憶している。吾輩はここで始めて人間というものを見た。しかもあとで聞くとそれは書生という人間中で一番獰悪な種族であったそうだ。この書生というのは時々我々を捕えて煮て食うという話である。`
	l := NewLexer(input)

	result := []string{}

	for _, _ = range input {
		result = append(result, string(l.ch))
		l.readRune()
	}

	// fmt.Printf("%#v\n", result)
}

func TestReadRune(t *testing.T) {
	tests := []struct {
		input          string
		expectCh       string
		expectPosition int
	}{
		{
			input:          "あいう",
			expectCh:       "い",
			expectPosition: 3,
		},
		{
			input:          "123",
			expectCh:       "2",
			expectPosition: 1,
		},
		{
			input:          "abc",
			expectCh:       "b",
			expectPosition: 1,
		},
		{
			input:          "ABC",
			expectCh:       "B",
			expectPosition: 1,
		},
		{
			input:          "123",
			expectCh:       "2",
			expectPosition: 1,
		},
		{
			input:          "(x)",
			expectCh:       "x",
			expectPosition: 1,
		},
		{
			input:          "🔥emoji",
			expectCh:       "e",
			expectPosition: 4,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("入力が%sの場合", tt.input), func(t *testing.T) {
			l := NewLexer(tt.input)
			l.readRune()
			{
				got := l.ch
				if string(got) != tt.expectCh {
					t.Errorf("got %s want %s", string(got), tt.expectCh)
				}
			}
			{
				got := l.position
				if got != tt.expectPosition {
					t.Errorf("got %d want %d", got, tt.expectPosition)
				}
			}
		})
	}
}

func TestSkipShiteSpace(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "hello",
			expect: "h",
		},
		{
			input:  " hello",
			expect: "h",
		},
		{
			input:  "     hello",
			expect: "h",
		},
		{
			input:  "     1",
			expect: "1",
		},
		{
			input:  "     あ",
			expect: "あ",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("入力が%sの場合", tt.input), func(t *testing.T) {
			l := NewLexer(tt.input)
			l.skipWhitespace()
			got := l.ch
			if string(got) != tt.expect {
				t.Errorf("got %s want %s", string(got), tt.expect)
			}
		})
	}
}

func TestReadIdentifier(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "あいう",
			expect: "あいう",
		},
		{
			input:  "あいう\nかきく",
			expect: "あいう",
		},
		{
			input:  "\nあいう",
			expect: "",
		},
		// {
		// 	name:   "",
		// 	input:  "",
		// 	expect: "",
		// },
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("入力が%sの場合", tt.input), func(t *testing.T) {
			l := NewLexer(tt.input)
			got := l.readIdentifier()
			if got != tt.expect {
				t.Errorf("got %s want %s", got, tt.expect)
			}
		})
	}
}
