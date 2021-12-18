package main

import (
	"flag"
	"fmt"
	"github.com/muidea/magicCommon/application"
	"github.com/muidea/magicDefault/assist/persistence"
	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicDefault/core"
)

var listenPort = "8890"
var endpointName = "magicDefault"

func initPprofMonitor(listenPort string) error {
	var err error
	addr := ":1" + listenPort

	go func() {
		err = http.ListenAndServe(addr, nil)
		if err != nil {
			log.Critical("funcRetErr=http.ListenAndServe||err=%s", err.Error())
		}
	}()

	return err
}

func main() {
	flag.StringVar(&listenPort, "ListenPort", listenPort, "magicDefault listen address")
	flag.StringVar(&endpointName, "EndpointName", endpointName, "application endpoint name.")
	flag.Parse()

	initPprofMonitor(listenPort)

	fmt.Printf("magicDefault V1.0\n")
	err := persistence.Initialize(endpointName)
	if err != nil {
		log.Errorf("initialize persistence failed, err:%s", err.Error())
		return
	}
	defer persistence.Uninitialize()

	core, err := core.New(endpointName, listenPort)
	if err != nil {
		log.Errorf("create core service failed, err:%s", err.Error())
		return
	}

	app := application.GetApp()
	app.Startup(core)
	app.Run()
	app.Shutdown()
}
