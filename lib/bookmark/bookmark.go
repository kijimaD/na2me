package bookmark

import "time"

var Bookmarks BookmarksType

func init() {
	Bookmarks = BookmarksType{
		Bookmarks:     []Bookmark{},
		BookmarkIndex: map[string]int{},
	}
}

type Bookmark struct {
	// 一意の名前。現在タイトルごとになっているので、タイトル名を入れる
	ID string
	// タイトル名
	ScenarioName string
	// 章ラベル
	Label string
	// イベント位置
	EventIdx int
	// 保存日付
	SavedAt time.Time
}

func NewBookmark(scenarioName string, label string, eventIdx int) Bookmark {
	return Bookmark{
		ID:           scenarioName,
		ScenarioName: scenarioName,
		Label:        label,
		EventIdx:     eventIdx,
		SavedAt:      time.Now(),
	}
}

type BookmarksType struct {
	Bookmarks     []Bookmark
	BookmarkIndex map[string]int
}

func (master BookmarksType) Get(key string) Bookmark {
	idx := master.BookmarkIndex[key]

	return master.Bookmarks[idx]
}

func (master BookmarksType) Add(bookmark Bookmark) {
	master.Bookmarks = append(master.Bookmarks, bookmark)
	for i, bm := range master.Bookmarks {
		master.BookmarkIndex[bm.ID] = i
	}
}
