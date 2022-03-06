package persistence

import (
	"fmt"
	"sync"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/model"
)

var batisInitializeOnce sync.Once
var batisUninitializeOnce sync.Once
var batisClient client.Client

func Initialize(endpointName string) (err error) {
	batisInitializeOnce.Do(func() {
		clnt := client.NewClient(config.BatisService(), endpointName)

		servicePtr := bc.NewService(
			config.DatabaseServer(),
			config.DatabaseName(),
			config.DatabaseUsername(),
			config.DatabaseUserPassword(),
			config.DatabaseMaxConnection(),
		)
		err = clnt.RegisterService(servicePtr)
		if err != nil {
			log.Errorf("register instance failed, err:%s", err.Error())
			return
		}

		err = model.InitializeModel(clnt)
		if err != nil {
			log.Errorf("initialize model failed, err:%s", err.Error())
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
