package main

import (
	"os"
	"path"

	"github.com/apex/log"
	"github.com/joho/godotenv"
	"github.com/kick-project/kick/internal"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/subcmds/installcmd"
	"github.com/kick-project/kick/internal/subcmds/listcmd"
	"github.com/kick-project/kick/internal/subcmds/removecmd"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
)

func main() {
	loadDotenv()
	home, err := os.UserHomeDir()
	errutils.Efatalf("error: %w", err)
	inject := di.Setup(home)

	if os.Getenv("KICK_DEBUG") == "true" {
		inject.LogLevel(log.DebugLevel)
	}

	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		utils.Exit(startcmd.Start(args[1:], inject))
	case o.List:
		utils.Exit(listcmd.List(args[1:], inject))
	case o.Search:
		utils.Exit(searchcmd.Search(args[1:], inject))
	case o.Init:
		utils.Exit(initcmd.InitCmd(args[1:], inject))
	case o.Update:
		utils.Exit(updatecmd.Update(args[1:], inject))
	case o.Install:
		utils.Exit(installcmd.Install(args[1:], inject))
	case o.Remove:
		utils.Exit(removecmd.Remove(args[1:], inject))
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
