package persistence

import (
	"flag"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
)

var batisInitializeOnce sync.Once
var batisUninitializeOnce sync.Once
var batisClient client.Client

var databaseServer = ""
var databaseName = "magicdefault_db"
var databaseUsername = "magicdefault"
var databasePassword = "magicdefault"
var maxConnNum = 10
var batisService = "http://localhost:8080"

func init() {
	flag.StringVar(&databaseServer, "DBServer", databaseServer, "database server address.")
	flag.StringVar(&databaseName, "DBName", databaseName, "database name.")
	flag.StringVar(&databaseUsername, "DBUsername", databaseUsername, "database username.")
	flag.StringVar(&databasePassword, "DBPassword", databasePassword, "database password.")
	flag.IntVar(&maxConnNum, "MaxConnNum", maxConnNum, "max connection number.")
	flag.StringVar(&batisService, "BatisService", batisService, "magicBatis service address.")
}

func Initialize(endpointName string) (err error) {
	batisInitializeOnce.Do(func() {
		clnt := client.NewClient(batisService, endpointName)

		servicePtr := bc.NewService(
			databaseServer,
			databaseName,
			databaseUsername,
			databasePassword,
			maxConnNum,
		)
		err = clnt.RegisterService(servicePtr)
		if err != nil {
			log.Errorf("register instance failed, err:%s", err.Error())
			return
		}

		batisClient = clnt
	})

	return
}

func Uninitialize() {
	batisUninitializeOnce.Do(func() {
		if batisClient != nil {
			batisClient.Release()
			batisClient = nil
		}
	})
}

func GetBatisClient() (ret client.Client) {
	if batisClient == nil {
		err := fmt.Errorf("must initialze persistence first")
		panic(err)
	}

	ret = batisClient
	return
}
