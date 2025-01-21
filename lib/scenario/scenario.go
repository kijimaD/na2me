package scenario

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/kijimaD/na2me/embeds"
)

var ScenarioMaster ScenarioMasterType

func init() {
	scenarios := []Scenario{
		New("フランツカフカ", "変身"),
		New("和辻哲郎", "漱石の人物"),
		New("坂口安吾", "堕落論"),
		New("坂口安吾", "帝銀事件を論ず"),
		New("坂口安吾", "文学のふるさと"),
		New("夏目漱石", "こころ"),
		New("夏目漱石", "それから"),
		New("夏目漱石", "三四郎"),
		New("夏目漱石", "予の描かんと欲する作品"),
		New("夏目漱石", "余と万年筆"),
		New("夏目漱石", "倫敦塔"),
		New("夏目漱石", "写生文"),
		New("夏目漱石", "博士問題の成行"),
		New("夏目漱石", "吾輩は猫である"),
		New("夏目漱石", "坊っちゃん"),
		New("夏目漱石", "夢十夜"),
		New("夏目漱石", "戦争からきた行き違い"),
		New("夏目漱石", "手紙"),
		New("夏目漱石", "文士の生活"),
		New("夏目漱石", "文芸委員は何をするか"),
		New("夏目漱石", "明暗"),
		New("夏目漱石", "模倣と独立"),
		New("夏目漱石", "満韓ところどころ"),
		New("夏目漱石", "無題"),
		New("夏目漱石", "現代日本の開化"),
		New("夏目漱石", "私の個人主義"),
		New("夏目漱石", "私の経過した学生時代"),
		New("夏目漱石", "草枕"),
		New("夏目漱石", "虞美人草"),
		New("夏目漱石", "行人"),
		New("夏目漱石", "趣味の遺伝"),
		New("夏目漱石", "道草"),
		New("夏目漱石", "門"),
		New("太宰治", "きりぎりす"),
		New("太宰治", "トカトントン"),
		New("太宰治", "人間失格"),
		New("太宰治", "女生徒"),
		New("太宰治", "富嶽百景"),
		New("太宰治", "斜陽"),
		New("太宰治", "走れメロス"),
		New("宮本百合子", "漱石の行人について"),
		New("小林多喜二", "蟹工船"),
		New("梶井基次郎", "桜の樹の下には"),
		New("梶井基次郎", "檸檬"),
		New("森鴎外", "舞姫"),
		New("森鴎外", "高瀬舟"),
		New("樋口一葉", "たけくらべ"),
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
	}
	sm := ScenarioMasterType{
		Scenarios:     []Scenario{},
		Statuses:      []Status{},
		ScenarioIndex: map[ScenarioIDType]int{},
	}
	sm.prepare(scenarios)
	sm.loadBody()
	if err := GlobalLoad(&sm); err != nil {
		log.Fatal(err)
	}

	ScenarioMaster = sm
}

type ScenarioMasterType struct {
	Scenarios     []Scenario
	Statuses      []Status
	ScenarioIndex map[ScenarioIDType]int
}

func (master *ScenarioMasterType) prepare(scenarios []Scenario) {
	for i, s := range scenarios {
		id := GenerateScenarioID(s.AuthorName, s.Title)

		master.Scenarios = append(master.Scenarios, Scenario{
			ID:         id,
			Title:      s.Title,
			AuthorName: s.AuthorName,
		})
		master.Statuses = append(master.Statuses, Status{ID: id})
		master.ScenarioIndex[id] = i
	}
}

func (master *ScenarioMasterType) loadBody() error {
	for i, s := range master.Scenarios {
		body, err := embeds.FS.ReadFile(string(s.ID))
		if err != nil {
			return err
		}
		master.Scenarios[i].Body = body
	}

	return nil
}

func (master *ScenarioMasterType) GetScenario(key ScenarioIDType) Scenario {
	idx := master.ScenarioIndex[key]

	return master.Scenarios[idx]
}

func (master *ScenarioMasterType) ExportStatuses(w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(master.Statuses)
	if err != nil {
		return err
	}

	return nil
}

func (master *ScenarioMasterType) ImportStatuses(r io.Reader) error {
	newStatuses := []Status{}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, &newStatuses); err != nil {
		return err
	}
	for _, ns := range newStatuses {
		idx, ok := master.ScenarioIndex[ns.ID]
		if ok {
			master.Statuses[idx].IsRead = ns.IsRead
		}
	}

	return nil
}

type ScenarioIDType string

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

func New(authorName string, title string) Scenario {
	return Scenario{
		Title:      title,
		AuthorName: authorName,
	}
}

type Status struct {
	// 一意の名前
	ID ScenarioIDType
	// 既読
	IsRead bool
}

func GenerateScenarioID(authorName, title string) ScenarioIDType {
	return ScenarioIDType(fmt.Sprintf("scenario/%s/%s.sce", authorName, title))
}
