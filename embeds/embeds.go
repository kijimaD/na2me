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

//go:embed scenario/ningen_shikkaku.sce
var NingenShikkaku []byte

//go:embed scenario/hashire_merosu.sce
var HashireMerosu []byte

//go:embed scenario/mohouto_dokuritu.sce
var MohoutoDokuritu []byte

//go:embed scenario/meian.sce
var Meian []byte

//go:embed scenario/souseki_no_jinbutu.sce
var SousekiNoJinbutu []byte

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
			Name:       "Bochan",
			LabelName:  "坊っちゃん",
			AuthorName: "夏目漱石",
			Body:       Bochan,
		},
		{
			Name:       "WagahaiHaNekoDearu",
			LabelName:  "吾輩は猫である",
			AuthorName: "夏目漱石",
			Body:       WagahaiHaNekoDearu,
		},
		{
			Name:       "Sanshirou",
			LabelName:  "三四郎",
			AuthorName: "夏目漱石",
			Body:       Sanshirou,
		},
		{
			Name:       "Kokoro",
			LabelName:  "こころ",
			AuthorName: "夏目漱石",
			Body:       Kokoro,
		},
		{
			Name:       "MohoutoDokuritu",
			LabelName:  "模倣と独立",
			AuthorName: "夏目漱石",
			Body:       MohoutoDokuritu,
		},
		{
			Name:       "Meian",
			LabelName:  "明暗",
			AuthorName: "夏目漱石",
			Body:       Meian,
		},
		{
			Name:       "NingenShikkaku",
			LabelName:  "人間失格",
			AuthorName: "太宰治",
			Body:       NingenShikkaku,
		},
		{
			Name:       "HashireMerosu",
			LabelName:  "走れメロス",
			AuthorName: "太宰治",
			Body:       HashireMerosu,
		},
		{
			Name:       "SousekiNoJinbutu",
			LabelName:  "漱石の人物",
			AuthorName: "和辻哲郎",
			Body:       SousekiNoJinbutu,
		},
		// {
		// 	Name:       "",
		// 	LabelName:  "",
		// 	AuthorName: "",
		// 	Body:       ,
		// },
	}
	sm.ScenarioIndex = map[string]int{}
	for i, s := range sm.Scenarios {
		sm.ScenarioIndex[s.Name] = i
	}

	ScenarioMaster = sm
}

//go:embed *
var FS embed.FS
