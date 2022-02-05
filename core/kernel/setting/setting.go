package setting

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/kernel/setting/biz"
	"github.com/muidea/magicDefault/core/kernel/setting/service"
)

func init() {
	module.Register(New())
}

type Setting struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	service *service.Setting
	biz     *biz.Setting
}

func New() *Setting {
	return &Setting{}
}

func (s *Setting) ID() string {
	return common.SettingModule
}

func (s *Setting) BindRegistry(
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

func (s *Setting) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		persistence.GetBatisClient(),
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(s.biz)
	s.service.BindRegistry(s.routeRegistry, s.casRouteRegistry, s.roleRouteRegistry)
	s.service.RegisterRoute()
}

func (s *Setting) Teardown() {

}
