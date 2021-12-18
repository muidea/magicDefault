package setting

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/kernel/setting/biz"
	"github.com/muidea/magicDefault/core/kernel/setting/service"
)

func init() {
	module.Register(New())
}

type Setting struct {
	service *service.Setting
	biz     *biz.Setting
}

func New() module.Module {
	return &Setting{}
}

func (s *Setting) ID() string {
	return common.SettingModule
}

func (s *Setting) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(s.biz)
	s.service.RegisterRoute()
}

func (s *Setting) Teardown() {

}
