package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kijimaD/na2me/lib"
	"github.com/kijimaD/na2me/lib/check"
	"github.com/kijimaD/na2me/lib/convert/lexer"
	"github.com/kijimaD/na2me/lib/convert/parser"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

const splash = `───────────────────────────────────────────────────────
███    ██  █████  ██████  ███    ███ ███████
████   ██ ██   ██      ██ ████  ████ ██
██ ██  ██ ███████  █████  ██ ████ ██ █████
██  ██ ██ ██   ██ ██      ██  ██  ██ ██
██   ████ ██   ██ ███████ ██      ██ ███████
───────────────────────────────────────────────────────
`

func NewMainApp() *cli.App {
	app := cli.NewApp()
	app.Name = "na2me"
	app.Usage = "na2me [subcommand] [args]"
	app.Description = "na2me novel file converter"
	app.DefaultCommand = CmdLaunch.Name
	app.Version = "v0.0.0"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		CmdLaunch,
		CmdConvert,
		CmdCheckLen,
	}
	cli.AppHelpTemplate = fmt.Sprintf(`%s
%s
`, splash, cli.AppHelpTemplate)

	return app
}

func RunMainApp(app *cli.App, args ...string) error {
	err := app.Run(args)
	if err != nil {
		return fmt.Errorf("コマンド実行が失敗した: %w", err)
	}

	return nil
}

// ================

var CmdLaunch = &cli.Command{
	Name:        "launch",
	Usage:       "launch",
	Description: "launch",
	Action:      cmdLaunch,
	Flags:       []cli.Flag{},
}

func cmdLaunch(_ *cli.Context) error {
	game := &lib.Game{}
	game.Start()

	return nil
}

// ================

var CmdConvert = &cli.Command{
	Name:        "convert",
	Usage:       "convert",
	Description: "convert",
	Action:      cmdConvert,
	Flags:       []cli.Flag{},
}

func cmdConvert(_ *cli.Context) error {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("パイプで加工したいテキストを標準入力に渡す必要がある")
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	if len(b) <= 0 {
		return fmt.Errorf("パイプで加工したいテキストを標準入力に渡す必要がある")
	}

	l := lexer.New(string(b))
	p := parser.New(l)
	scenario := p.ParseScenario()

	fmt.Println(scenario)

	return nil
}

// ================

var CmdCheckLen = &cli.Command{
	Name:        "checkLen",
	Usage:       "checkLen",
	Description: "checkLen",
	Action:      cmdCheckLen,
	Flags:       []cli.Flag{},
}

func cmdCheckLen(_ *cli.Context) error {
	directory := "./embeds/scenario" // 検索するディレクトリ
	extension := ".sce"              // 対象ファイルの拡張子
	threshold := 240                 // 長い行とみなす文字数の閾値

	// ディレクトリを再帰的に検索
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		// ファイルであり、指定の拡張子を持つ場合のみ処理
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			check.WarnLongLine(f, os.Stdout, threshold, f.Name())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
