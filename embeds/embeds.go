package embeds

import "embed"

//go:embed scenario/bochan.sce
var Bochan []byte

//go:embed scenario/wagahai_ha_neko_dearu.sce
var WagahaiHaNekoDearu []byte

//go:embed scenario/sanshirou.sce
var Sanshirou []byte

// ================

type ScenarioMasterType struct {
	Scenarios     []Scenario
	ScenarioIndex map[string]int
}

type Scenario struct {
	LabelName string
	Name      string
	Body      []byte
}

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = []Scenario{
		{
			LabelName: "坊っちゃん",
			Name:      "Bochan",
			Body:      Bochan,
		},
		{
			LabelName: "吾輩は猫である",
			Name:      "WagahaiHaNekoDearu",
			Body:      WagahaiHaNekoDearu,
		},
		{
			LabelName: "三四郎",
			Name:      "Sanshirou",
			Body:      Sanshirou,
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
