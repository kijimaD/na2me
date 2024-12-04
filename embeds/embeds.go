package embeds

import "embed"

//go:embed scenario/bochan.sce
var Bochan []byte

//go:embed scenario/wagahai_ha_neko_dearu.sce
var WagahaiHaNekoDearu []byte

//go:embed scenario/sanshirou.sce
var Sanshirou []byte

//go:embed scenario/kokoro.sce
var Kokoro []byte

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
	// 本文
	Body []byte
}

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = []Scenario{
		{
			Name:      "Bochan",
			LabelName: "坊っちゃん",
			Body:      Bochan,
		},
		{
			Name:      "WagahaiHaNekoDearu",
			LabelName: "吾輩は猫である",
			Body:      WagahaiHaNekoDearu,
		},
		{
			Name:      "Sanshirou",
			LabelName: "三四郎",
			Body:      Sanshirou,
		},
		{
			Name:      "Kokoro",
			LabelName: "こころ",
			Body:      Kokoro,
		},
	}
	sm.ScenarioIndex = map[string]int{}
	for i, s := range sm.Scenarios {
		sm.ScenarioIndex[s.Name] = i
	}

	ScenarioMaster = sm
}

//go:embed *
var FS embed.FS
