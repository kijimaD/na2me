package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kijimaD/na2me/lib"
	"github.com/kijimaD/na2me/lib/check"
	"github.com/kijimaD/na2me/lib/convert/lexer"
	"github.com/kijimaD/na2me/lib/convert/parser"
	"github.com/kijimaD/na2me/lib/normalize"
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
		CmdCheckNotes,
		CmdPrintChapterTmpl,
		CmdNormalizeBg,
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
	Description: "起動する",
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
	Description: "機械的に改ページタグをつける",
	Action:      cmdConvert,
	Flags:       []cli.Flag{},
}

func cmdConvert(_ *cli.Context) error {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("パイプで加工したいテキストを標準入力に渡す必要がある")
	}

	b, err := io.ReadAll(os.Stdin)
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
	Description: "行の長さが超えてないかチェックする",
	Action:      cmdCheckLen,
	Flags:       []cli.Flag{},
}

func cmdCheckLen(_ *cli.Context) error {
	directory := "./embeds/scenario" // 検索するディレクトリ
	extension := ".sce"              // 対象ファイルの拡張子
	upperThreshold := 240            // 長い行とみなす文字数。より大きい
	lowerThreshold := 3              // 短いとみなす文字数。より小さい

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
			check.WarnLineLen(f, os.Stdout, upperThreshold, lowerThreshold, f.Name())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

var CmdCheckNotes = &cli.Command{
	Name:        "checkNotes",
	Usage:       "checkNotes",
	Description: "おかしくなりがちな脚注部分を表示する",
	Action:      cmdCheckNotes,
	Flags:       []cli.Flag{},
}

func cmdCheckNotes(_ *cli.Context) error {
	directory := "./embeds/scenario" // 検索するディレクトリ
	extension := ".sce"              // 対象ファイルの拡張子

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
			check.WarnNotes(f, os.Stdout, f.Name())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

var CmdPrintChapterTmpl = &cli.Command{
	Name:        "printChapterTmpl",
	Usage:       "printChapterTmpl",
	Description: "章は手動でつけるので、コピペ用のテンプレートを標準出力する",
	Action:      cmdPrintChapterTmpl,
	Flags:       []cli.Flag{},
}

func cmdPrintChapterTmpl(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		return fmt.Errorf("引数が不足している")
	}
	numString := ctx.Args().Get(0)
	num, err := strconv.Atoi(numString)
	if err != nil {
		return err
	}

	{
		str := `*start
[image source="black.png"]
『xxx』xxxx
[p]
`
		fmt.Printf(str)
	}

	for i := 1; i <= num; i++ {
		str := `[jump target="ch%d"]

*ch%d
--------
`
		fmt.Printf(str, i, i)
	}

	{
		str := `[jump target="end"]

*end
終わり
[p]
[jump target="start"]

*footnotes
xxx
[p]
[jump target="start"]
`
		fmt.Printf(str)
	}

	return nil
}

// ================

var CmdNormalizeBg = &cli.Command{
	Name:        "normalizeBg",
	Usage:       "normalizeBg",
	Description: "背景画像を標準化する",
	Action:      runNormalizeBg,
	Flags:       []cli.Flag{},
}

func runNormalizeBg(_ *cli.Context) error {
	const srcDir = "./raw/bg"
	const destDir = "./embeds/bg"

	c, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range c {
		if !strings.HasSuffix(entry.Name(), ".png") && !strings.HasSuffix(entry.Name(), ".jpg") {
			return fmt.Errorf("%sにjpgもしくはpng以外のファイルが含まれている: %s", srcDir, entry.Name())
		}
		src := path.Join(srcDir, entry.Name())
		dest := path.Join(destDir, entry.Name())
		err := normalize.NormalizeBg(src, dest)
		if err != nil {
			return err
		}
	}

	return nil
}
