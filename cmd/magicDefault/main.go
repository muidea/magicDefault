package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicCommon/application"

	"github.com/muidea/magicDefault/core"
)

var listenPort = "8080"
var endpointName = "magicdefault"

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

	core, err := core.New(endpointName, listenPort)
	if err != nil {
		log.Errorf("create core service failed, err:%s", err.Error())
		return
	}

	application.Startup(core)
	application.Run()
	application.Shutdown()
}
