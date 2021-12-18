package core

import (
	"fmt"
	"github.com/muidea/magicCommon/module"
	"net"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	casCommon "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/assist/registry"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/model"

	_ "github.com/muidea/magicDefault/core/kernel/base"
	_ "github.com/muidea/magicDefault/core/kernel/setting"
	_ "github.com/muidea/magicDefault/core/kernel/totalizer"
	_ "github.com/muidea/magicDefault/core/module/authority"
	_ "github.com/muidea/magicDefault/core/module/image"
	_ "github.com/muidea/magicDefault/core/module/remoteHub"
)

type loadNamespaceTask struct {
	eventHub event.Hub
}

func (s *loadNamespaceTask) Run() {
	log.Info("load namespace task running!")
	eid := common.LoadAuthorityNamespace
	header := event.NewValues()
	header.Set("namespace", config.SuperNamespace())

	filter := def.NewFilter()
	eventPtr := event.NewEvent(eid, "/", common.AuthorityModule, header, filter)
	result := s.eventHub.Call(eventPtr)
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		return
	}
	namespaceList, namespaceOK := resultVal.([]*casCommon.NamespaceView)
	if !namespaceOK {
		return
	}

	eid = common.InitializeAuthorityNamespace
	for _, val := range namespaceList {
		eventPtr = event.NewEvent(eid, "/", "/#", header, val)
		s.eventHub.Post(eventPtr)
	}
}

// New 新建Core
func New(endpointName, listenPort string) (ret *Core, err error) {
	batisClnt, batisErr := persistence.GetBatisClient()
	if batisErr != nil {
		err = batisErr
		return
	}

	err = model.InitializeModel(batisClnt)
	if err != nil {
		log.Errorf("initializeModel failed, err:%s", err.Error())
		return
	}

	core := &Core{
		endpointName: endpointName,
		listenPort:   listenPort,
	}

	ret = core
	return
}

// Core Core对象
type Core struct {
	endpointName string
	listenPort   string

	httpServer engine.HTTPServer
}

func (s *Core) Name() string {
	return s.endpointName
}

// Startup 启动
func (s *Core) Startup() {
	router := engine.NewRouter()
	registry.New(s, s, s, router)

	s.httpServer = engine.NewHTTPServer(s.listenPort)
	s.httpServer.Bind(router)

	modules := module.GetModules()
	for _, val := range modules {
		val.Setup(s.endpointName)
	}
}

func (s *Core) Run() {
	s.httpServer.Run()
}

// Shutdown 销毁
func (s *Core) Shutdown() {
	modules := module.GetModules()
	for _, val := range modules {
		val.Teardown()
	}
}

func (s *Core) OnTimeOut(session session.Session) {

}

// Verify verify current session
func (s *Core) Verify(res http.ResponseWriter, req *http.Request) (err error) {
	_, err = s.getCurrentEntity(res, req)
	return
}

// VerifyRole verify current role
func (s *Core) VerifyRole(res http.ResponseWriter, req *http.Request) (ret *casCommon.RoleView, err error) {
	curRole, curErr := s.getCurrentRole(res, req)
	if curErr != nil {
		err = curErr
		return
	}

	ret = curRole
	return
}

func (s *Core) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(casCommon.NamespaceID)
	if namespace != "" {
		ret = namespace
		return
	}

	items := strings.Split(req.Host, ":")
	if nil != net.ParseIP(items[0]) {
		ret = "default"
		return
	}

	items = strings.Split(items[0], ".")
	if len(items) <= 1 {
		ret = "default"
		return
	}

	ret = items[0]
	return
}

func (s *Core) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := registry.GetSession(res, req)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}

	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *Core) getCurrentRole(res http.ResponseWriter, req *http.Request) (ret *casCommon.RoleView, err error) {
	curSession := registry.GetSession(res, req)
	authVal, ok := curSession.GetOption(session.AuthRole)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}

	ret = authVal.(*casCommon.RoleView)
	return
}
