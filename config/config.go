package config

import (
	"fmt"
	fu "github.com/muidea/magicCommon/foundation/util"
)

var localService = "magicdefault"
var localListenPort = "8010"
var batisService = "http://127.0.0.1:8080"
var casService = "http://127.0.0.1:8081"
var fileService = "http://127.0.0.1:8083"
var superNamespace = "super"
var databaseServer = ""
var databaseName = "magicdefault_db"
var databaseUsername = "magicdefault"
var databasePassword = "magicdefault"
var databaseMaxConnection = 10

var configItem *CfgItem

const cfgPath = "/var/app/config/cfg.json"

func init() {
	cfg := &CfgItem{}
	err := fu.LoadConfig(cfgPath, cfg)
	if err != nil {
		return
	}

	configItem = cfg
}

type CfgItem struct {
	BatisService          string `json:"batisService"`
	CasService            string `json:"casService"`
	FileService           string `json:"fileService"`
	DefaultNamespace      string `json:"defaultNamespace"`
	DatabaseServer        string `json:"databaseServer"`
	DatabaseName          string `json:"databaseName"`
	DatabaseUsername      string `json:"databaseUsername"`
	DatabasePassword      string `json:"databasePassword"`
	DatabaseMaxConnection int    `json:"databaseMaxConnection"`
}

func UpdateLocalInfo(service, port string) {
	localService = service
	localListenPort = port
}

func LocalServiceUrl() (ret string) {
	ret = fmt.Sprintf("http://%s:%s", localService, localListenPort)
	return
}

// BatisService baits Service
func BatisService() string {
	if configItem != nil {
		return configItem.BatisService
	}

	return batisService
}

//CasService cas partner addr
func CasService() string {
	if configItem != nil {
		return configItem.CasService
	}

	return casService
}

//FileService file partner addr
func FileService() string {
	if configItem != nil {
		return configItem.FileService
	}

	return fileService
}

//SuperNamespace super namespace
func SuperNamespace() string {
	return superNamespace
}

func DatabaseServer() string {
	if configItem != nil {
		return configItem.DatabaseServer
	}

	return databaseServer
}

func DatabaseName() string {
	if configItem != nil {
		return configItem.DatabaseName
	}
	return databaseName
}

func DatabaseUsername() string {
	if configItem != nil {
		return configItem.DatabaseUsername
	}

	return databaseUsername
}

func DatabaseUserPassword() string {
	if configItem != nil {
		return configItem.DatabasePassword
	}

	return databasePassword
}

func DatabaseMaxConnection() int {
	if configItem != nil {
		return configItem.DatabaseMaxConnection
	}

	return databaseMaxConnection
}
