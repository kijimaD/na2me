package embeds

import (
	"embed"
	"fmt"
	"log"
)

//go:embed *
var FS embed.FS

// ================

type ScenarioMasterType struct {
	Scenarios     []Scenario
	ScenarioIndex map[ScenarioIDType]int
}

func (master *ScenarioMasterType) GetScenario(key ScenarioIDType) Scenario {
	idx := master.ScenarioIndex[key]

	return master.Scenarios[idx]
}

type Scenario struct {
	// 一意の名前
	ID ScenarioIDType
	// タイトル
	Title string
	// 著者名
	AuthorName string
	// 本文
	Body []byte
}

type ScenarioIDType string

func NewScenario(authorName string, title string) Scenario {
	return Scenario{
		Title:      title,
		AuthorName: authorName,
	}
}

func generateScenarioID(authorName, title string) ScenarioIDType {
	return ScenarioIDType(fmt.Sprintf("scenario/%s/%s.sce", authorName, title))
}

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = append(sm.Scenarios,
		NewScenario("フランツカフカ", "変身"),
		NewScenario("和辻哲郎", "漱石の人物"),
		NewScenario("坂口安吾", "堕落論"),
		NewScenario("坂口安吾", "帝銀事件を論ず"),
		NewScenario("夏目漱石", "こころ"),
		NewScenario("夏目漱石", "それから"),
		NewScenario("夏目漱石", "三四郎"),
		NewScenario("夏目漱石", "吾輩は猫である"),
		NewScenario("夏目漱石", "坊っちゃん"),
		NewScenario("夏目漱石", "明暗"),
		NewScenario("夏目漱石", "模倣と独立"),
		NewScenario("夏目漱石", "満韓ところどころ"),
		NewScenario("夏目漱石", "現代日本の開化"),
		NewScenario("夏目漱石", "私の個人主義"),
		NewScenario("夏目漱石", "草枕"),
		NewScenario("夏目漱石", "虞美人草"),
		NewScenario("夏目漱石", "道草"),
		NewScenario("夏目漱石", "門"),
		NewScenario("太宰治", "人間失格"),
		NewScenario("太宰治", "斜陽"),
		NewScenario("太宰治", "走れメロス"),
		NewScenario("太宰治", "女生徒"),
		NewScenario("梶井基次郎", "桜の樹の下には"),
		NewScenario("梶井基次郎", "檸檬"),
		NewScenario("森鴎外", "舞姫"),
		NewScenario("泉鏡花", "高野聖"),
		NewScenario("田山花袋", "蒲団"),
		NewScenario("石原莞爾", "最終戦争論"),
		NewScenario("福沢諭吉", "学問のすすめ"),
		NewScenario("芥川龍之介", "トロッコ"),
		NewScenario("芥川龍之介", "地獄変"),
		NewScenario("芥川龍之介", "或阿呆の一生"),
		NewScenario("芥川龍之介", "羅生門"),
		NewScenario("芥川龍之介", "葉"),
		NewScenario("芥川龍之介", "蜜柑"),
		NewScenario("芥川龍之介", "鼻"),
		NewScenario("魯迅", "故郷"),
		// NewScenario("", ""),
	)

	sm.ScenarioIndex = map[ScenarioIDType]int{}
	for i, s := range sm.Scenarios {
		id := generateScenarioID(s.AuthorName, s.Title)
		sm.Scenarios[i].ID = id
		sm.ScenarioIndex[id] = i

		body, err := FS.ReadFile(string(id))
		if err != nil {
			log.Fatal(err)
		}
		sm.Scenarios[i].Body = body
	}

	ScenarioMaster = sm
}
