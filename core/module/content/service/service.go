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
	"github.com/muidea/magicDefault/core/module/content/biz"
)

type Content struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

	bizPtr *biz.Content
}

func New(
	contentBiz *biz.Content,
) *Content {
	ptr := &Content{
		bizPtr:    contentBiz,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *Content) BindRegistry(
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
func (s *Content) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterContentArticle, "GET", cc.ReadPrivate, s.FilterArticle)
	s.roleRouteRegistry.AddHandler(common.QueryContentArticle, "GET", cc.ReadPrivate, s.QueryArticle)
	s.roleRouteRegistry.AddHandler(common.CreateContentArticle, "POST", cc.WritePrivate, s.CreateArticle)
	s.roleRouteRegistry.AddHandler(common.UpdateContentArticle, "PUT", cc.WritePrivate, s.UpdateArticle)
	s.roleRouteRegistry.AddHandler(common.DeleteContentArticle, "DELETE", cc.DeletePrivate, s.DeleteArticle)

	s.roleRouteRegistry.AddHandler(common.FilterContentCatalog, "GET", cc.ReadPrivate, s.FilterCatalog)
	s.roleRouteRegistry.AddHandler(common.QueryContentCatalog, "GET", cc.ReadPrivate, s.QueryCatalog)
	s.roleRouteRegistry.AddHandler(common.CreateContentCatalog, "POST", cc.WritePrivate, s.CreateCatalog)
	s.roleRouteRegistry.AddHandler(common.UpdateContentCatalog, "PUT", cc.WritePrivate, s.UpdateCatalog)
	s.roleRouteRegistry.AddHandler(common.DeleteContentCatalog, "DELETE", cc.DeletePrivate, s.DeleteCatalog)

	s.roleRouteRegistry.AddHandler(common.FilterContentLink, "GET", cc.ReadPrivate, s.FilterLink)
	s.roleRouteRegistry.AddHandler(common.QueryContentLink, "GET", cc.ReadPrivate, s.QueryLink)
	s.roleRouteRegistry.AddHandler(common.CreateContentLink, "POST", cc.WritePrivate, s.CreateLink)
	s.roleRouteRegistry.AddHandler(common.UpdateContentLink, "PUT", cc.WritePrivate, s.UpdateLink)
	s.roleRouteRegistry.AddHandler(common.DeleteContentLink, "DELETE", cc.DeletePrivate, s.DeleteLink)

	s.roleRouteRegistry.AddHandler(common.FilterContentMedia, "GET", cc.ReadPrivate, s.FilterMedia)
	s.roleRouteRegistry.AddHandler(common.QueryContentMedia, "GET", cc.ReadPrivate, s.QueryMedia)
	s.roleRouteRegistry.AddHandler(common.CreateContentMedia, "POST", cc.WritePrivate, s.CreateMedia)
	s.roleRouteRegistry.AddHandler(common.UpdateContentMedia, "PUT", cc.WritePrivate, s.UpdateMedia)
	s.roleRouteRegistry.AddHandler(common.DeleteContentMedia, "DELETE", cc.DeletePrivate, s.DeleteMedia)

	s.roleRouteRegistry.AddHandler(common.FilterContentComment, "GET", cc.ReadPrivate, s.FilterComment)
	s.roleRouteRegistry.AddHandler(common.QueryContentComment, "GET", cc.ReadPrivate, s.QueryComment)
	s.roleRouteRegistry.AddHandler(common.CreateContentComment, "POST", cc.WritePrivate, s.CreateComment)
	s.roleRouteRegistry.AddHandler(common.UpdateContentComment, "PUT", cc.WritePrivate, s.UpdateComment)
	s.roleRouteRegistry.AddHandler(common.DeleteContentComment, "DELETE", cc.DeletePrivate, s.DeleteComment)

}

func (s *Content) getCurrentEntity(ctx context.Context) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *Content) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *Content) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *Content) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	entityPtr, _ := s.getCurrentEntity(ctx)
	namespace := s.getCurrentNamespace(ctx, res, req)
	s.bizPtr.WriteLog(memo, address, entityPtr, namespace)
}
