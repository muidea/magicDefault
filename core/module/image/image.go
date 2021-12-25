package authority

import (
	"github.com/muidea/magicCas/toolkit"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/core/module/image/biz"
	"github.com/muidea/magicDefault/core/module/image/service"
)

func init() {
	module.Register(New())
}

type Image struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	service *service.Image
	biz     *biz.Image
}

func New() module.Module {
	return &Image{}
}

func (s *Image) ID() string {
	return common.ImageModule
}

func (s *Image) BindRegistry(
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

func (s *Image) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(endpointName, config.FileService(),
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(s.biz)
	s.service.RegisterRoute()
}

func (s *Image) Teardown() {
}
