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
	Name string
	// 表示名
	LabelName string
	// 著者名
	AuthorName string
	// 本文
	Body []byte
}

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = []Scenario{
		{
			LabelName:  "こころ",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "坊っちゃん",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "吾輩は猫である",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "三四郎",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "模倣と独立",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "明暗",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "それから",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "門",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "私の個人主義",
			AuthorName: "夏目漱石",
		},
		{
			LabelName:  "人間失格",
			AuthorName: "太宰治",
		},
		{
			LabelName:  "走れメロス",
			AuthorName: "太宰治",
		},
		{
			LabelName:  "漱石の人物",
			AuthorName: "和辻哲郎",
		},
		{
			LabelName:  "羅生門",
			AuthorName: "芥川龍之介",
		},
		{
			LabelName:  "学問のすすめ",
			AuthorName: "福沢諭吉",
		},
		{
			LabelName:  "故郷",
			AuthorName: "魯迅",
		},
	}

	sm.ScenarioIndex = map[string]int{}
	for i, s := range sm.Scenarios {
		fname := fmt.Sprintf("scenario/%s/%s.sce", s.AuthorName, s.LabelName)
		sm.Scenarios[i].Name = fname
		sm.ScenarioIndex[fname] = i

		body, err := FS.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		sm.Scenarios[i].Body = body
	}

	ScenarioMaster = sm
}
