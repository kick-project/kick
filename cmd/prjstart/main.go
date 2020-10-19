package main

import (
	"os"
	"path"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/subcmds/installcmd"
	"github.com/kick-project/kick/internal/subcmds/listcmd"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/joho/godotenv"
)

func main() {
	loadDotenv()
	home, err := os.UserHomeDir()
	errutils.Efatalf("error: %w", err)
	s := settings.GetSettings(home)

	if os.Getenv("PRJSTART_DEBUG") == "true" {
		s.LogLevel(log.DebugLevel)
	}

	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		utils.Exit(startcmd.Start(args[1:], s))
	case o.List:
		utils.Exit(listcmd.List(args[1:], s))
	case o.Search:
		utils.Exit(searchcmd.Search(args[1:], s))
	case o.Init:
		utils.Exit(initcmd.InitCmd(args[1:], s))
	case o.Update:
		utils.Exit(updatecmd.Update(args[1:], s))
	case o.Install:
		utils.Exit(installcmd.Install(args[1:], s))
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
