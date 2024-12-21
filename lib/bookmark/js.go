//go:build js

package bookmark

import (
	"bytes"
	"strings"
	"syscall/js"
)

const bookmarkKey = "bm.json"

func GlobalLoad() error {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	jsVal := localStorage.Get(bookmarkKey)
	if !jsVal.Truthy() {
		return nil
	}

	bmt := BookmarksType{}
	bmt.Import(strings.NewReader(jsVal.String()))
	Bookmarks = bmt

	return nil
}

func GlobalSave() error {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")

	buf := bytes.Buffer{}
	Bookmarks.Export(&buf)
	localStorage.Set(bookmarkKey, buf.String())

	return nil
}
