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
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/core/module/remoteHub/biz"
	"github.com/muidea/magicDefault/model"
)

type RemoteHub struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

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
	routeRegistry toolkit.RouteRegistry,
	casRouteRegistry toolkit.CasRegistry,
	roleRouteRegistry toolkit.RoleRegistry) {

	s.routeRegistry = routeRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// RegisterRoute 注册路由
func (s *RemoteHub) RegisterRoute() {
}

func (s *RemoteHub) getCurrentEntity(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *RemoteHub) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *RemoteHub) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *RemoteHub) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(ctx, res, req)
	curNamespace := s.getCurrentNamespace(ctx, res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
