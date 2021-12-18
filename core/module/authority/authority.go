package authority

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/core/module/authority/biz"
	"github.com/muidea/magicDefault/core/module/authority/service"
)

func init() {
	module.Register(New())
}

type Authority struct {
	service *service.Authority
	biz     *biz.Authority
}

func New() module.Module {
	return &Authority{}
}

func (s *Authority) ID() string {
	return common.AuthorityModule
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
	s.service.RegisterRoute()
}

func (s *Authority) Teardown() {

}
