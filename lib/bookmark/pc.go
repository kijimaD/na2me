//go:build linux || windows || darwin

package bookmark

import (
	"bytes"
	"os"
)

const bookmarkFileName = "bm.json"

func GlobalLoad() error {
	bs, err := os.ReadFile(bookmarkFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	newBM := BookmarksType{}
	newBM.Import(bytes.NewReader(bs))
	Bookmarks = newBM

	return nil
}

func GlobalSave() error {
	f, err := os.Create(bookmarkFileName)
	if err != nil {
		return err
	}

	Bookmarks.Export(f)

	return nil
}
