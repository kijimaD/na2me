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
	ScenarioIndex map[string]int
}

func (master *ScenarioMasterType) GetScenario(key string) Scenario {
	idx := master.ScenarioIndex[key]

	return master.Scenarios[idx]
}

type Scenario struct {
	// 一意の名前
	ID string
	// タイトル
	Title string
	// 著者名
	AuthorName string
	// 本文
	Body []byte
}

func NewScenario(authorName string, title string) Scenario {
	return Scenario{
		Title:      title,
		AuthorName: authorName,
	}
}

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = append(sm.Scenarios,
		NewScenario("フランツカフカ", "変身"),
		NewScenario("和辻哲郎", "漱石の人物"),
		NewScenario("坂口安吾", "堕落論"),
		NewScenario("夏目漱石", "こころ"),
		NewScenario("夏目漱石", "それから"),
		NewScenario("夏目漱石", "三四郎"),
		NewScenario("夏目漱石", "吾輩は猫である"),
		NewScenario("夏目漱石", "坊っちゃん"),
		NewScenario("夏目漱石", "明暗"),
		NewScenario("夏目漱石", "模倣と独立"),
		NewScenario("夏目漱石", "満韓ところどころ"),
		NewScenario("夏目漱石", "私の個人主義"),
		NewScenario("夏目漱石", "虞美人草"),
		NewScenario("夏目漱石", "道草"),
		NewScenario("夏目漱石", "門"),
		NewScenario("太宰治", "人間失格"),
		NewScenario("太宰治", "斜陽"),
		NewScenario("太宰治", "走れメロス"),
		NewScenario("梶井基次郎", "檸檬"),
		NewScenario("梶井基次郎", "桜の樹の下には"),
		NewScenario("森鴎外", "舞姫"),
		NewScenario("石原莞爾", "最終戦争論"),
		NewScenario("福沢諭吉", "学問のすすめ"),
		NewScenario("芥川龍之介", "或阿呆の一生"),
		NewScenario("芥川龍之介", "羅生門"),
		NewScenario("魯迅", "故郷"),
		// NewScenario("", ""),
	)

	sm.ScenarioIndex = map[string]int{}
	for i, s := range sm.Scenarios {
		fname := fmt.Sprintf("scenario/%s/%s.sce", s.AuthorName, s.Title)
		sm.Scenarios[i].ID = fname
		sm.ScenarioIndex[fname] = i

		body, err := FS.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		sm.Scenarios[i].Body = body
	}

	ScenarioMaster = sm
}
