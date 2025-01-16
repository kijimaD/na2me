package scenario

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportStatuses_JSONに変換できる(t *testing.T) {
	smt := ScenarioMasterType{
		Statuses: []Status{
			{
				ID:     "aaa",
				IsRead: true,
			},
			{
				ID:     "bbb",
				IsRead: false,
			},
		},
	}

	buf := bytes.Buffer{}
	assert.NoError(t, smt.ExportStatuses(&buf))

	expect := `[{"ID":"aaa","IsRead":true},{"ID":"bbb","IsRead":false}]
`
	assert.Equal(t, expect, buf.String())
}

func TestImportStatuses_JSONから読み込める(t *testing.T) {
	smt := ScenarioMasterType{
		Scenarios:     []Scenario{},
		Statuses:      []Status{},
		ScenarioIndex: map[ScenarioIDType]int{},
	}
	scenarios := []Scenario{
		New("夏目漱石", "こころ"),
		New("夏目漱石", "道草"),
		New("フランツカフカ", "変身"),
	}
	for i, s := range scenarios {
		id := GenerateScenarioID(s.AuthorName, s.Title)
		smt.Scenarios = append(smt.Scenarios, Scenario{
			ID:         id,
			Title:      s.Title,
			AuthorName: s.AuthorName,
		})
		smt.Statuses = append(smt.Statuses, Status{ID: id})
		smt.ScenarioIndex[id] = i
	}

	input := `[{"ID":"scenario/夏目漱石/こころ.sce","IsRead":true},{"ID":"scenario/フランツカフカ/変身.sce","IsRead":true}]
`
	r := strings.NewReader(input)
	assert.NoError(t, smt.ImportStatuses(r))
	expect := ScenarioMasterType{
		Scenarios: []Scenario{
			Scenario{
				ID:         "scenario/夏目漱石/こころ.sce",
				Title:      "こころ",
				AuthorName: "夏目漱石",
				Body:       []uint8(nil),
			},
			Scenario{
				ID:         "scenario/夏目漱石/道草.sce",
				Title:      "道草",
				AuthorName: "夏目漱石",
				Body:       []uint8(nil),
			},
			Scenario{
				ID:         "scenario/フランツカフカ/変身.sce",
				Title:      "変身",
				AuthorName: "フランツカフカ",
				Body:       []uint8(nil),
			},
		},
		Statuses: []Status{
			Status{
				ID:     "scenario/夏目漱石/こころ.sce",
				IsRead: true,
			},
			Status{
				ID:     "scenario/夏目漱石/道草.sce",
				IsRead: false,
			},
			Status{
				ID:     "scenario/フランツカフカ/変身.sce",
				IsRead: true,
			},
		},
		ScenarioIndex: map[ScenarioIDType]int{
			"scenario/夏目漱石/こころ.sce":   0,
			"scenario/夏目漱石/道草.sce":    1,
			"scenario/フランツカフカ/変身.sce": 2,
		},
	}
	assert.Equal(t, expect, smt)
}
