package remoteHub

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/remoteHub/biz"
	"github.com/muidea/magicDefault/core/module/remoteHub/service"
)

func init() {
	module.Register(New())
}

type RemoteHub struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	service *service.RemoteHub
	biz     *biz.RemoteHub
}

func New() module.Module {
	return &RemoteHub{}
}

func (s *RemoteHub) ID() string {
	return common.RemoteHubModule
}

func (s *RemoteHub) BindRegistry(
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

func (s *RemoteHub) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(endpointName,
		persistence.GetBatisClient(),
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(endpointName, s.biz)
	s.service.BindRegistry(s.routeRegistry, s.casRouteRegistry, s.roleRouteRegistry)
	s.service.RegisterRoute()
}

func (s *RemoteHub) Teardown() {

}
