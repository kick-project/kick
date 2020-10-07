package main

import (
	"log"
	"os"
	"path"

	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/subcmds/initcmd"
	"github.com/crosseyed/prjstart/internal/subcmds/listcmd"
	"github.com/crosseyed/prjstart/internal/subcmds/start"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	loadDotenv()
	s := settings.GetSettings("")

	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		utils.Exit(start.Start(args[1:], s))
	case o.List:
		utils.Exit(listcmd.List(args[1:], s))
	case o.Init:
		utils.Exit(initcmd.InitCmd(args[1:], s))
	}
	utils.Exit(255)
}

func loadDotenv() {
	// TODO - optional .env support
	for _, envfile := range []string{path.Join(os.Getenv("HOME"), ".env"), ".env"} {
		if _, err := os.Stat(envfile); err != nil {
			continue
		}
		err := godotenv.Load(envfile)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
