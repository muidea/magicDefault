package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/authority/biz"
)

type Authority struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

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
	routeRegistry toolkit.RouteRegistry,
	casRouteRegistry toolkit.CasRegistry,
	roleRouteRegistry toolkit.RoleRegistry) {

	s.routeRegistry = routeRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry

	s.routeRegistry.SetApiVersion(common.ApiVersion)
	s.casRouteRegistry.SetApiVersion(common.ApiVersion)
	s.roleRouteRegistry.SetApiVersion(common.ApiVersion)
}

// RegisterRoute 注册路由
func (s *Authority) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterAuthorityAccount, "GET", cc.ReadPermission, s.FilterAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityAccount, "GET", cc.ReadPermission, s.QueryAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityAccount, "POST", cc.WritePermission, s.CreateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityAccount, "PUT", cc.WritePermission, s.UpdateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityAccount, "DELETE", cc.DeletePermission, s.DeleteAuthorityAccount)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityEndpoint, "GET", cc.ReadPermission, s.FilterAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityEndpoint, "GET", cc.ReadPermission, s.QueryAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityEndpoint, "POST", cc.WritePermission, s.CreateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityEndpoint, "PUT", cc.WritePermission, s.UpdateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityEndpoint, "DELETE", cc.DeletePermission, s.DeleteAuthorityEndpoint)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityRole, "GET", cc.ReadPermission, s.FilterAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityRole, "GET", cc.ReadPermission, s.QueryAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityRole, "POST", cc.WritePermission, s.CreateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityRole, "PUT", cc.WritePermission, s.UpdateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityRole, "DELETE", cc.DeletePermission, s.DeleteAuthorityRole)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityNamespace, "GET", cc.ReadPermission, s.FilterAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityNamespace, "GET", cc.ReadPermission, s.QueryAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityNamespace, "POST", cc.WritePermission, s.CreateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityNamespace, "PUT", cc.WritePermission, s.UpdateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityNamespace, "DELETE", cc.DeletePermission, s.DeleteAuthorityNamespace)
}

func (s *Authority) getCurrentEntity(ctx context.Context) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *Authority) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *Authority) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	entityPtr, _ := s.getCurrentEntity(ctx)
	namespace := s.getCurrentNamespace(ctx, res, req)
	s.bizPtr.WriteLog(memo, address, entityPtr, namespace)
}
