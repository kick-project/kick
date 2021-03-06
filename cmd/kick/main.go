package main

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/kick-project/kick/internal"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/subcmds/installcmd"
	"github.com/kick-project/kick/internal/subcmds/removecmd"
	"github.com/kick-project/kick/internal/subcmds/repocmd"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/kick-project/kick/internal/subcmds/setupcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
)

//nolint
//gocyclo:ignore
func main() {
	loadDotenv()
	home, err := os.UserHomeDir()
	errs.FatalF("error: %w", err)
	inject := di.New(&di.Options{Home: home})
	exitHdlr := inject.MakeExitHandler()

	// open log file and close on exit
	logfile := os.Getenv("KICK_LOG")
	if logfile != "" {
		lf := inject.MakeLogFile(logfile)
		defer func() {
			lf.Close()
			info, err := os.Stat(lf.Name())
			if err != nil && info.Size() == 0 {
				os.Remove(lf.Name())
			}
		}()
	}

	if os.Getenv("KICK_DEBUG") == "true" {
		inject.LogLevel(logger.DebugLevel)
	}

	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		startcmd.Start(args[1:], inject)
	case o.Search:
		exitHdlr.Exit(searchcmd.Search(args[1:], inject))
	case o.Setup:
		exitHdlr.Exit(setupcmd.SetupCmd(args[1:], inject))
	case o.Update:
		exitHdlr.Exit(updatecmd.Update(args[1:], inject))
	case o.Install:
		exitHdlr.Exit(installcmd.Install(args[1:], inject))
	case o.Remove:
		exitHdlr.Exit(removecmd.Remove(args[1:], inject))
	case o.Init:
		exitHdlr.Exit(initcmd.Init(args[1:], inject))
	case o.Repo:
		exitHdlr.Exit(repocmd.Repo(args[1:], inject))
	}
	exitHdlr.Exit(255)
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
