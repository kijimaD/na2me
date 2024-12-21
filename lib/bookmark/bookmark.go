package bookmark

import (
	"time"

	"github.com/kijimaD/na2me/embeds"
)

var Bookmarks BookmarksType

func init() {
	Bookmarks = BookmarksType{
		Bookmarks:     []Bookmark{},
		BookmarkIndex: map[embeds.ScenarioIDType]int{},
	}
}

type Bookmark struct {
	// 一意の名前。現在タイトルごとになっているので、タイトル名を入れる
	ID embeds.ScenarioIDType
	// タイトル名
	ScenarioName string
	// 章ラベル
	Label string
	// 保存日付
	SavedAt time.Time
}

func NewBookmark(id embeds.ScenarioIDType, scenarioName string, label string) Bookmark {
	return Bookmark{
		ID:           id,
		ScenarioName: scenarioName,
		Label:        label,
		SavedAt:      time.Now(),
	}
}

type BookmarksType struct {
	Bookmarks     []Bookmark
	BookmarkIndex map[embeds.ScenarioIDType]int
}

func (master *BookmarksType) Get(key embeds.ScenarioIDType) (Bookmark, bool) {
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
