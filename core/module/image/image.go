package authority

import (
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
	service *service.Image
	biz     *biz.Image
}

func New() module.Module {
	return &Image{}
}

func (s *Image) ID() string {
	return common.ImageModule
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
