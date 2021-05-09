package main

import (
	"os"
	"path"

	"github.com/apex/log"
	"github.com/joho/godotenv"
	"github.com/kick-project/kick/internal"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/subcmds/installcmd"
	"github.com/kick-project/kick/internal/subcmds/listcmd"
	"github.com/kick-project/kick/internal/subcmds/removecmd"
	"github.com/kick-project/kick/internal/subcmds/repocmd"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/kick-project/kick/internal/subcmds/setupcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
)

//gocyclo:ignore
func main() {
	loadDotenv()
	home, err := os.UserHomeDir()
	errs.FatalF("error: %w", err)
	inject := di.Setup(home)

	if os.Getenv("KICK_DEBUG") == "true" {
		inject.LogLevel(log.DebugLevel)
	}

	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		exit.Exit(startcmd.Start(args[1:], inject))
	case o.List:
		exit.Exit(listcmd.List(args[1:], inject))
	case o.Search:
		exit.Exit(searchcmd.Search(args[1:], inject))
	case o.Setup:
		exit.Exit(setupcmd.SetupCmd(args[1:], inject))
	case o.Update:
		exit.Exit(updatecmd.Update(args[1:], inject))
	case o.Install:
		exit.Exit(installcmd.Install(args[1:], inject))
	case o.Remove:
		exit.Exit(removecmd.Remove(args[1:], inject))
	case o.Init:
		exit.Exit(initcmd.Init(args[1:], inject))
	case o.Repo:
		exit.Exit(repocmd.Repo(args[1:], inject))
	}
	exit.Exit(255)
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
