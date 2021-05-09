package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/sevlyar/go-daemon"
	"github.com/sosedoff/gitkit"
)

var testServerSync sync.Mutex
var testServerUp bool

func main() {
	home, _ := filepath.Abs("tmp/home")
	srvpath, _ := filepath.Abs("tmp/gitserve")
	cntxt := &daemon.Context{
		PidFileName: "tmp/server.pid",
		PidFilePerm: 0644,
		LogFileName: "tmp/server.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"gitkit"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer func() {
		err = cntxt.Release()
		errs.Panic(err)
	}()

	startTestServer(home, srvpath)
}

type testServerStruct struct {
	server *http.Server
	config *config.File
	home   string
}

func (s *testServerStruct) ServerUp(home string, servpath string) {
	testServerSync.Lock()
	defer testServerSync.Unlock()
	if testServerUp {
		return
	}
	s.home = home
	s.server = s.gitServer(servpath)
	s.loadConfig()
	spathexists := true
	homeexists := true
	_, err := os.Stat(servpath)
	if os.IsNotExist(err) {
		spathexists = false
	}

	_, err = os.Stat(home)
	if os.IsNotExist(err) {
		homeexists = false
	}
	testServerUp = true
	log.Printf("HOMEDIR: %s EXISTS=%t", home, homeexists)
	log.Printf("SERVDIR: %s EXISTS=%t", servpath, spathexists)
}

func (s *testServerStruct) GetHome() string {
	return s.home
}

func (s *testServerStruct) gitServer(dir string) *http.Server {
	addr := getAddr()

	service := gitkit.New(gitkit.Config{
		Dir:        dir,
		AutoCreate: true,
		AutoHooks:  true,
	})

	// Configure git TestServer. Will create git repos path if it does not exist.
	// If hooks are set, it will also update all repos with new version of hook scripts.
	if err := service.Setup(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", service)
	server := &http.Server{
		Addr: addr,
	}

	m := `Starting Git server %s, path %s`
	log.Printf(m, addr, dir)
	err := server.ListenAndServe()
	if err != nil {
		if err.Error() == fmt.Sprintf("listen tcp %s: bind: address already in use", addr) {
			return nil
		}
		log.Fatalf("%v", err)
	}
	return server
}

func (s *testServerStruct) loadConfig() {
	conf := &config.File{
		PathUserConf:     filepath.Join(s.home, ".kick", "config.yml"),
		PathTemplateConf: filepath.Join(s.home, ".kick", "templates.yml"),
	}
	err := conf.Load()
	errs.Panic(err)

	s.config = conf
}

func startTestServer(home string, servpath string) {
	testserver := &testServerStruct{}
	testserver.ServerUp(home, servpath)
}

func getAddr() string {
	addr := os.Getenv("TESTGITADDR") // Set in $PROJECT/.env
	if addr == "" {
		addr = "127.0.0.1:5000"
	} else {
		addr = "127.0.0.1" + addr
	}
	return addr
}
