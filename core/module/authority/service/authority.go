package service

import (
	"fmt"
	"github.com/muidea/magicDefault/core/module/authority/biz"
	"github.com/muidea/magicDefault/model"
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

	"github.com/muidea/magicDefault/common"
)

type Authority struct {
	sessionRegistry   commonSession.Registry
	casRouteRegistry  casToolkit.CasRegistry
	roleRouteRegistry casToolkit.RoleRegistry
	validator         fu.Validator

	bizPtr *biz.Authority
}

func New(
	authorityBiz *biz.Authority,
) *Authority {
	ptr := &Authority{
		bizPtr:    authorityBiz,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *Authority) BindRegistry(
	sessionRegistry commonSession.Registry,
	casRouteRegistry casToolkit.CasRegistry,
	roleRouteRegistry casToolkit.RoleRegistry) {

	s.sessionRegistry = sessionRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// Handle middleware handler
func (s *Authority) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonSession.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()
}

// RegisterRoute 注册路由
func (s *Authority) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterAuthorityAccount, "GET", casCommon.ReadPrivate, s.FilterAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityAccount, "GET", casCommon.ReadPrivate, s.QueryAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityAccount, "POST", casCommon.WritePrivate, s.CreateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityAccount, "PUT", casCommon.WritePrivate, s.UpdateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityAccount, "DELETE", casCommon.DeletePrivate, s.DeleteAuthorityAccount)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityEndpoint, "GET", casCommon.ReadPrivate, s.FilterAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityEndpoint, "GET", casCommon.ReadPrivate, s.QueryAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityEndpoint, "POST", casCommon.WritePrivate, s.CreateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityEndpoint, "PUT", casCommon.WritePrivate, s.UpdateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityEndpoint, "DELETE", casCommon.DeletePrivate, s.DeleteAuthorityEndpoint)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityRole, "GET", casCommon.ReadPrivate, s.FilterAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityRole, "GET", casCommon.ReadPrivate, s.QueryAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityRole, "POST", casCommon.WritePrivate, s.CreateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityRole, "PUT", casCommon.WritePrivate, s.UpdateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityRole, "DELETE", casCommon.DeletePrivate, s.DeleteAuthorityRole)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityNamespace, "GET", casCommon.ReadPrivate, s.FilterAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityNamespace, "GET", casCommon.ReadPrivate, s.QueryAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityNamespace, "POST", casCommon.WritePrivate, s.CreateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityNamespace, "PUT", casCommon.WritePrivate, s.UpdateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityNamespace, "DELETE", casCommon.DeletePrivate, s.DeleteAuthorityNamespace)
}

func (s *Authority) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonSession.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *Authority) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *Authority) queryEntity(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *Authority) writeLog(res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(res, req)
	curNamespace := s.getCurrentNamespace(res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
