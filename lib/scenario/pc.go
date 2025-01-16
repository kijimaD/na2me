//go:build linux || windows || darwin

package scenario

import (
	"bytes"
	"os"
)

const statusFileName = "status.json"

func GlobalLoad(smt *ScenarioMasterType) error {
	bs, err := os.ReadFile(statusFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	smt.ImportStatuses(bytes.NewReader(bs))

	return nil
}

func GlobalSave(smt *ScenarioMasterType) error {
	f, err := os.Create(statusFileName)
	if err != nil {
		return err
	}
	smt.ExportStatuses(f)

	return nil
}
