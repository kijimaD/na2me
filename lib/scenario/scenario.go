package scenario

import (
	"fmt"
	"log"

	"github.com/kijimaD/na2me/embeds"
)

var ScenarioMaster ScenarioMasterType

func init() {
	sm := ScenarioMasterType{}
	sm.Scenarios = append(sm.Scenarios,
		New("フランツカフカ", "変身"),
		New("和辻哲郎", "漱石の人物"),
		New("坂口安吾", "堕落論"),
		New("坂口安吾", "帝銀事件を論ず"),
		New("坂口安吾", "文学のふるさと"),
		New("夏目漱石", "こころ"),
		New("夏目漱石", "それから"),
		New("夏目漱石", "三四郎"),
		New("夏目漱石", "倫敦塔"),
		New("夏目漱石", "吾輩は猫である"),
		New("夏目漱石", "坊っちゃん"),
		New("夏目漱石", "文士の生活"),
		New("夏目漱石", "明暗"),
		New("夏目漱石", "模倣と独立"),
		New("夏目漱石", "満韓ところどころ"),
		New("夏目漱石", "現代日本の開化"),
		New("夏目漱石", "私の個人主義"),
		New("夏目漱石", "草枕"),
		New("夏目漱石", "虞美人草"),
		New("夏目漱石", "行人"),
		New("夏目漱石", "道草"),
		New("夏目漱石", "門"),
		New("太宰治", "きりぎりす"),
		New("太宰治", "トカトントン"),
		New("太宰治", "人間失格"),
		New("太宰治", "女生徒"),
		New("太宰治", "富嶽百景"),
		New("太宰治", "斜陽"),
		New("太宰治", "走れメロス"),
		New("小林多喜二", "蟹工船"),
		New("梶井基次郎", "桜の樹の下には"),
		New("梶井基次郎", "檸檬"),
		New("森鴎外", "舞姫"),
		New("森鴎外", "高瀬舟"),
		New("泉鏡花", "高野聖"),
		New("田山花袋", "蒲団"),
		New("石原莞爾", "最終戦争論"),
		New("福沢諭吉", "学問のすすめ"),
		New("芥川龍之介", "トロッコ"),
		New("芥川龍之介", "地獄変"),
		New("芥川龍之介", "或阿呆の一生"),
		New("芥川龍之介", "歯車"),
		New("芥川龍之介", "羅生門"),
		New("芥川龍之介", "葉"),
		New("芥川龍之介", "藪の中"),
		New("芥川龍之介", "蜜柑"),
		New("芥川龍之介", "鼻"),
		New("高村光太郎", "道程"),
		New("魯迅", "故郷"),
		New("魯迅", "狂人日記"),
		// New("", ""),
	)

	sm.ScenarioIndex = map[ScenarioIDType]int{}
	for i, s := range sm.Scenarios {
		id := GenerateScenarioID(s.AuthorName, s.Title)
		sm.Scenarios[i].ID = id
		sm.ScenarioIndex[id] = i

		body, err := embeds.FS.ReadFile(string(id))
		if err != nil {
			log.Fatal(err)
		}
		sm.Scenarios[i].Body = body
	}

	ScenarioMaster = sm
}

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

func New(authorName string, title string) Scenario {
	return Scenario{
		Title:      title,
		AuthorName: authorName,
	}
}

func GenerateScenarioID(authorName, title string) ScenarioIDType {
	return ScenarioIDType(fmt.Sprintf("scenario/%s/%s.sce", authorName, title))
}
