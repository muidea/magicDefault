package service

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	commonSession "github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	casCommon "github.com/muidea/magicCas/common"
	casToolkit "github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/core/module/remoteHub/biz"
	"github.com/muidea/magicDefault/model"
)

type RemoteHub struct {
	sessionRegistry   commonSession.Registry
	casRouteRegistry  casToolkit.CasRegistry
	roleRouteRegistry casToolkit.RoleRegistry
	validator         fu.Validator

	bizPtr *biz.RemoteHub

	endpointName string
}

func New(
	endpointName string,
	bizPtr *biz.RemoteHub,
) *RemoteHub {
	ptr := &RemoteHub{
		bizPtr:       bizPtr,
		endpointName: endpointName,
		validator:    fu.NewFormValidator(),
	}

	return ptr
}

func (s *RemoteHub) BindRegistry(
	sessionRegistry commonSession.Registry,
	casRouteRegistry casToolkit.CasRegistry,
	roleRouteRegistry casToolkit.RoleRegistry) {

	s.sessionRegistry = sessionRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// Handle middleware handler
func (s *RemoteHub) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonSession.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()
}

// RegisterRoute 注册路由
func (s *RemoteHub) RegisterRoute() {
}

func (s *RemoteHub) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonSession.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *RemoteHub) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *RemoteHub) queryEntity(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *RemoteHub) writeLog(res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(res, req)
	curNamespace := s.getCurrentNamespace(res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
