package main

import (
	"log"
	"os"

	"github.com/kijimaD/na2me/lib/cmd"
)

func main() {
	app := cmd.NewMainApp()
	err := cmd.RunMainApp(app, os.Args...)
	if err != nil {
		log.Fatal(err)
	}
}
