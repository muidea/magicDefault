package remoteHub

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/remoteHub/biz"
	"github.com/muidea/magicDefault/core/module/remoteHub/service"
)

func init() {
	module.Register(New())
}

type RemoteHub struct {
	service *service.RemoteHub
	biz     *biz.RemoteHub
}

func New() module.Module {
	return &RemoteHub{}
}

func (s *RemoteHub) ID() string {
	return common.RemoteHubModule
}

func (s *RemoteHub) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(endpointName,
		eventHub,
		backgroundRoutine,
	)

	s.service = service.New(endpointName, s.biz)
	s.service.RegisterRoute()
}

func (s *RemoteHub) Teardown() {

}
