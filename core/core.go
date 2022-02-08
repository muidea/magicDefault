package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicCommon/event"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"
	engine "github.com/muidea/magicEngine"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"

	_ "github.com/muidea/magicDefault/core/kernel/base"
	_ "github.com/muidea/magicDefault/core/kernel/setting"
	_ "github.com/muidea/magicDefault/core/kernel/totalizer"
	_ "github.com/muidea/magicDefault/core/module/authority"
	_ "github.com/muidea/magicDefault/core/module/content"
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

	filter := fu.NewFilter()
	eventPtr := event.NewEvent(eid, "/", common.AuthorityModule, header, filter)
	result := s.eventHub.Call(eventPtr)
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		return
	}
	namespaceList, namespaceOK := resultVal.([]*cc.NamespaceView)
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

	sessionRegistry session.Registry
	httpServer      engine.HTTPServer
}

// Startup 启动
func (s *Core) Startup(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	router := engine.NewRouter()
	routeRegistry := toolkit.NewRouteRegistry(router)
	casRegistry := toolkit.NewCasRegistry(s, router)
	roleRegistry := toolkit.NewRoleRegistry(s, router)

	s.sessionRegistry = session.CreateRegistry(s)

	persistence.Initialize(s.endpointName)

	s.httpServer = engine.NewHTTPServer(s.listenPort)
	s.httpServer.Bind(router)
	s.httpServer.Use(s)

	modules := module.GetModules()
	for _, val := range modules {

		module.BindBatisClient(val, persistence.GetBatisClient())

		module.BindRegistry(val, routeRegistry, casRegistry, roleRegistry)

		module.Setup(val, s.endpointName, eventHub, backgroundRoutine)
	}
}

func (s *Core) Run() {
	s.httpServer.Run()
}

// Shutdown 销毁
func (s *Core) Shutdown() {
	modules := module.GetModules()
	for _, val := range modules {
		module.Teardown(val)
	}

	persistence.Uninitialize()
}

func (s *Core) OnTimeOut(session session.Session) {

}

// Verify verify current session
func (s *Core) Verify(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	namespaceVal, ok := curSession.GetOption(session.AuthNamespace)
	if !ok {
		err = fmt.Errorf("无效Namespace,请先登录")
		return
	}

	requestNamespace := s.getRequestNamespace(res, req)
	if namespaceVal.(string) != requestNamespace {
		log.Errorf("namespace mismatch, request namespace:%s, current namespace:%s", requestNamespace, namespaceVal.(string))
		err = fmt.Errorf("无效Namespace,未通过验证")
		return
	}

	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}

	ret = authVal.(*cc.EntityView)
	return
}

func (s *Core) VerifyRole(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *cc.RoleView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	namespaceVal, ok := curSession.GetOption(session.AuthNamespace)
	if !ok {
		err = fmt.Errorf("无效Namespace,请先登录")
		return
	}

	requestNamespace := s.getRequestNamespace(res, req)
	if namespaceVal.(string) != requestNamespace {
		log.Errorf("namespace mismatch, request namespace:%s, current namespace:%s", requestNamespace, namespaceVal.(string))
		err = fmt.Errorf("无效Namespace,未通过验证")
		return
	}

	roleVal, ok := curSession.GetOption(session.AuthRole)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}

	ret = roleVal.(*cc.RoleView)
	return
}

func (s *Core) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = session.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Update(context.WithValue(ctx.Context(), session.AuthSession, curSession))
	ctx.Next()
}

func (s *Core) getRequestNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(cc.NamespaceID)
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
