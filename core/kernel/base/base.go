package base

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/core/kernel/base/biz"
	"github.com/muidea/magicDefault/core/kernel/base/service"
)

func init() {
	module.Register(New())
}

type Base struct {
	service *service.Base
	biz     *biz.Base
}

func New() module.Module {
	return &Base{}
}

func (s *Base) ID() string {
	return common.BaseModule
}

func (s *Base) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		config.CasService(),
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(endpointName, config.FileService(), s.biz)
	s.service.RegisterRoute()
}

func (s *Base) Teardown() {

}
