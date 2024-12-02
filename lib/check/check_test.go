package check

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarnLongLine(t *testing.T) {
	buf := bytes.Buffer{}
	r := strings.NewReader(`hello
worldxxxxx
safe
12345678
あああああああ
安全
警告を受ける文字長
受けない文字長
`)
	WarnLongLine(r, &buf, 8, "sample.txt")

	expect := `sample.txt
Line: 2, Length: 10
  worldxxxxx
Line: 7, Length: 9
  警告を受ける文字長
--------------------------------------------------------------------------------
`
	assert.Equal(t, expect, buf.String())
}

func TestWarnNotes(t *testing.T) {
	buf := bytes.Buffer{}
	r := strings.NewReader(`sample
＃警告1
※警告2
hello
world
`)
	WarnNotes(r, &buf, "sample.txt")

	expect := `sample.txt
Line: 2
  ＃警告1
Line: 3
  ※警告2
--------------------------------------------------------------------------------
`
	assert.Equal(t, expect, buf.String())
}
