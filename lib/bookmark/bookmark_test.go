package bookmark

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kijimaD/na2me/embeds"
	"github.com/stretchr/testify/assert"
)

func TestExport_JSONに変換できる(t *testing.T) {
	bmt := BookmarksType{
		Bookmarks:     []Bookmark{},
		BookmarkIndex: map[embeds.ScenarioIDType]int{},
	}
	testDate := time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)
	{
		bm := NewBookmark(embeds.ScenarioIDType("id/こころ"), "こころ", "ch1")
		bm.SavedAt = testDate
		bmt.Add(bm)
	}
	{
		bm := NewBookmark(embeds.ScenarioIDType("id/道草"), "道草", "ch2")
		bm.SavedAt = testDate
		bmt.Add(bm)
	}
	buf := bytes.Buffer{}
	assert.NoError(t, bmt.Export(&buf))

	expect := `[{"ID":"id/こころ","ScenarioName":"こころ","Label":"ch1","SavedAt":"2024-12-21T00:00:00Z"},{"ID":"id/道草","ScenarioName":"道草","Label":"ch2","SavedAt":"2024-12-21T00:00:00Z"}]
`
	assert.Equal(t, expect, buf.String())
}

func TestImport_JSONから読み込める(t *testing.T) {
	bmt := BookmarksType{}
	input := `[{"ID":"id/こころ","ScenarioName":"こころ","Label":"ch1","SavedAt":"2024-12-21T00:00:00Z"},{"ID":"id/道草","ScenarioName":"道草","Label":"ch2","SavedAt":"2024-12-21T00:00:00Z"}]
`
	r := strings.NewReader(input)
	assert.NoError(t, bmt.Import(r))

	expect := BookmarksType{
		Bookmarks: []Bookmark{
			Bookmark{ID: "id/こころ", ScenarioName: "こころ", Label: "ch1", SavedAt: time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)},
			Bookmark{ID: "id/道草", ScenarioName: "道草", Label: "ch2", SavedAt: time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)}},
		BookmarkIndex: map[embeds.ScenarioIDType]int{"id/こころ": 0, "id/道草": 1}}
	assert.Equal(t, expect, bmt)
}
