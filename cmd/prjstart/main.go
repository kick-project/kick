package main

import (
	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/subcmds"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"
)

func main() {
	loadDotenv()

	internal.Config = internal.LoadConfig("", "")
	args := os.Args
	o := internal.GetOptMain(args)
	switch {
	case o.Start:
		utils.Exit(subcmds.Start(args[1:]))
	case o.List:
		utils.Exit(subcmds.List(args[1:]))
	case o.Install:
		utils.Exit(subcmds.Install(args[1:]))
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
