//go:build js

package scenario

import (
	"bytes"
	"strings"
	"syscall/js"
)

const statusKey = "status.json"

func GlobalLoad(smt *ScenarioMasterType) error {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	jsVal := localStorage.Get(statusKey)
	if !jsVal.Truthy() {
		return nil
	}
	smt.ImportStatuses(strings.NewReader(jsVal.String()))

	return nil
}

func GlobalSave(smt *ScenarioMasterType) error {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")

	buf := bytes.Buffer{}
	smt.ExportStatuses(&buf)
	localStorage.Set(statusKey, buf.String())

	return nil
}
