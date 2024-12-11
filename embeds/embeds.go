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

//go:embed scenario/sorekara.sce
var Sorekara []byte

//go:embed scenario/mon.sce
var Mon []byte

//go:embed scenario/watashi_no_kojinsyugi.sce
var WatashiNoKojinsyugi []byte

//go:embed scenario/rasyomon.sce
var Rasyomon []byte

//go:embed scenario/gakumon_no_susume.sce
var GakumonNoSusume []byte

//go:embed scenario/kokyo.sce
var Kokyo []byte

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
			Name:       "Sorekara",
			LabelName:  "それから",
			AuthorName: "夏目漱石",
			Body:       Sorekara,
		},
		{
			Name:       "Mon",
			LabelName:  "門",
			AuthorName: "夏目漱石",
			Body:       Mon,
		},
		{
			Name:       "WatashiNoKojinsyugi",
			LabelName:  "私の個人主義",
			AuthorName: "夏目漱石",
			Body:       WatashiNoKojinsyugi,
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
		{
			Name:       "Rasyomon",
			LabelName:  "羅生門",
			AuthorName: "芥川龍之介",
			Body:       Rasyomon,
		},
		{
			Name:       "GakumonNoSusume",
			LabelName:  "学問のすすめ",
			AuthorName: "福沢諭吉",
			Body:       GakumonNoSusume,
		},
		{
			Name:       "Kokyo",
			LabelName:  "故郷",
			AuthorName: "魯迅",
			Body:       Kokyo,
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
