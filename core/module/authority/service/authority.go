package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/authority/biz"
	"github.com/muidea/magicDefault/model"
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
	s.roleRouteRegistry.AddHandler(common.FilterAuthorityAccount, "GET", cc.ReadPrivate, s.FilterAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityAccount, "GET", cc.ReadPrivate, s.QueryAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityAccount, "POST", cc.WritePrivate, s.CreateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityAccount, "PUT", cc.WritePrivate, s.UpdateAuthorityAccount)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityAccount, "DELETE", cc.DeletePrivate, s.DeleteAuthorityAccount)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityEndpoint, "GET", cc.ReadPrivate, s.FilterAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityEndpoint, "GET", cc.ReadPrivate, s.QueryAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityEndpoint, "POST", cc.WritePrivate, s.CreateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityEndpoint, "PUT", cc.WritePrivate, s.UpdateAuthorityEndpoint)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityEndpoint, "DELETE", cc.DeletePrivate, s.DeleteAuthorityEndpoint)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityRole, "GET", cc.ReadPrivate, s.FilterAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityRole, "GET", cc.ReadPrivate, s.QueryAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityRole, "POST", cc.WritePrivate, s.CreateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityRole, "PUT", cc.WritePrivate, s.UpdateAuthorityRole)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityRole, "DELETE", cc.DeletePrivate, s.DeleteAuthorityRole)

	s.roleRouteRegistry.AddHandler(common.FilterAuthorityNamespace, "GET", cc.ReadPrivate, s.FilterAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.QueryAuthorityNamespace, "GET", cc.ReadPrivate, s.QueryAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.CreateAuthorityNamespace, "POST", cc.WritePrivate, s.CreateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.UpdateAuthorityNamespace, "PUT", cc.WritePrivate, s.UpdateAuthorityNamespace)
	s.roleRouteRegistry.AddHandler(common.DeleteAuthorityNamespace, "DELETE", cc.DeletePrivate, s.DeleteAuthorityNamespace)
}

func (s *Authority) getCurrentEntity(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *cc.EntityView, err error) {
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

func (s *Authority) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *Authority) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(ctx, res, req)
	curNamespace := s.getCurrentNamespace(ctx, res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
