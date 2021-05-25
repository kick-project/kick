package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/sevlyar/go-daemon"
	"github.com/sosedoff/gitkit"
)

var testServerSync sync.Mutex
var testServerUp bool

var usage string = `
Usage:
    testserver [-a <port>] [-s <serverpath>] [-p <pidfile>] [-l <logfile>] 

Options:
    -a <port>           Listen on port number [default: 5000]
    -s <serverpath>     Server path to git repositories [default: tmp/gitserve]
    -p <pidfile>        Path to pidfile [default: tmp/server.pid]
    -l <logfile>        Path to logfile [default: tmp/server.log]
`

func main() {
	args, _ := docopt.ParseDoc(usage)
	port := args["-a"].(string)
	serverpath := args["-s"].(string)
	pidfile := args["-p"].(string)
	logfile := args["-l"].(string)
	home, _ := filepath.Abs("tmp/home")
	srvpath, _ := filepath.Abs(serverpath)
	cntxt := &daemon.Context{
		PidFileName: pidfile,
		PidFilePerm: 0644,
		LogFileName: logfile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"gitkit"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		panic(err)
	}
	if d != nil {
		return
	}
	defer func() {
		err = cntxt.Release()
		if err != nil {
			panic(err)
		}
	}()

	startTestServer(port, srvpath, home)
}

type testServerStruct struct {
	server *http.Server
	home   string
}

func (s *testServerStruct) ServerUp(port, servepath, home string) {
	testServerSync.Lock()
	defer testServerSync.Unlock()
	if testServerUp {
		return
	}
	s.home = home
	addr := getAddr(port)
	s.server = s.gitServer(addr, servepath)
	testServerUp = true
}

func (s *testServerStruct) gitServer(addr string, dir string) *http.Server {
	service := gitkit.New(gitkit.Config{
		Dir:        dir,
		AutoCreate: true,
		AutoHooks:  true,
	})

	// Configure git TestServer. Will create git repos path if it does not exist.
	// If hooks are set, it will also update all repos with new version of hook scripts.
	if err := service.Setup(); err != nil {
		panic(err)
	}

	http.Handle("/", service)
	server := &http.Server{
		Addr: addr,
	}

	err := server.ListenAndServe()
	if err != nil {
		if err.Error() == fmt.Sprintf("listen tcp %s: bind: address already in use", addr) {
			return nil
		}
		panic(err)
	}
	return server
}

func startTestServer(port, servepath, home string) {
	testserver := &testServerStruct{}
	testserver.ServerUp(port, servepath, home)
}

func getAddr(port string) string {
	var addr string
	listen := os.Getenv("KICK_LISTEN") // Set in $PROJECT/.env
	if port != "" {
		addr = "127.0.0.1:" + port
	} else if listen != "" {
		addr = "127.0.0.1:" + listen
	} else {
		panic("Port is not set")
	}
	return addr
}
