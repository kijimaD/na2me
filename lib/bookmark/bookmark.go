package bookmark

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/kijimaD/na2me/lib/scenario"
)

var Bookmarks BookmarksType

func init() {
	Bookmarks = BookmarksType{
		Bookmarks:     []Bookmark{},
		BookmarkIndex: map[scenario.ScenarioIDType]int{},
	}

	if err := GlobalLoad(); err != nil {
		log.Fatal(err)
	}
}

type Bookmark struct {
	// 一意の名前。現在タイトルごとになっているので、タイトル名を入れる
	ID scenario.ScenarioIDType
	// タイトル名
	ScenarioName string
	// 章ラベル
	Label string
	// 保存日付
	SavedAt time.Time
}

func NewBookmark(id scenario.ScenarioIDType, scenarioName string, label string) Bookmark {
	return Bookmark{
		ID:           id,
		ScenarioName: scenarioName,
		Label:        label,
		SavedAt:      time.Now(),
	}
}

type BookmarksType struct {
	Bookmarks     []Bookmark
	BookmarkIndex map[scenario.ScenarioIDType]int
}

func (master *BookmarksType) Get(key scenario.ScenarioIDType) (Bookmark, bool) {
	idx, ok := master.BookmarkIndex[key]
	if !ok {
		return Bookmark{}, false
	}

	return master.Bookmarks[idx], true
}

func (master *BookmarksType) Add(bookmark Bookmark) {
	idx, ok := master.BookmarkIndex[bookmark.ID]
	if ok {
		master.Bookmarks[idx] = bookmark

		return
	}

	master.Bookmarks = append(master.Bookmarks, bookmark)
	for i, bm := range master.Bookmarks {
		master.BookmarkIndex[bm.ID] = i
	}
}

func (master *BookmarksType) Delete(key scenario.ScenarioIDType) {
	idx, ok := master.BookmarkIndex[key]
	if !ok {
		return
	}

	master.Bookmarks = append(master.Bookmarks[:idx], master.Bookmarks[idx+1:]...)
	delete(master.BookmarkIndex, key)

	// 削除でindexがずれるので再計算
	for i, bm := range master.Bookmarks {
		master.BookmarkIndex[bm.ID] = i
	}
}

func (master *BookmarksType) Export(w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(master.Bookmarks)
	if err != nil {
		return err
	}

	return nil
}

func (master *BookmarksType) Import(r io.Reader) error {
	newBM := BookmarksType{
		Bookmarks:     []Bookmark{},
		BookmarkIndex: map[scenario.ScenarioIDType]int{},
	}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, &newBM.Bookmarks); err != nil {
		return err
	}
	for i, bm := range newBM.Bookmarks {
		newBM.BookmarkIndex[bm.ID] = i
	}
	master.Bookmarks = newBM.Bookmarks
	master.BookmarkIndex = newBM.BookmarkIndex

	return nil
}
