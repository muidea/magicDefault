package authority

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/core/module/authority/biz"
	"github.com/muidea/magicDefault/core/module/authority/service"
)

func init() {
	module.Register(New())
}

type Authority struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	service *service.Authority
	biz     *biz.Authority
}

func New() module.Module {
	return &Authority{}
}

func (s *Authority) ID() string {
	return common.AuthorityModule
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

func (s *Authority) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(config.CasService(),
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(s.biz)
	s.service.BindRegistry(s.routeRegistry, s.casRouteRegistry, s.roleRouteRegistry)
	s.service.RegisterRoute()
}

func (s *Authority) Teardown() {

}
