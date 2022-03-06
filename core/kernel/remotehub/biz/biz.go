package biz

import (
	"github.com/muidea/magicBatis/client"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/remotehub/dao"
)

type RemoteHub struct {
	biz.Base

	remotehubDao dao.RemoteHub

	endpointName string
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *RemoteHub {

	ptr := &RemoteHub{
		Base:         biz.New(common.RemoteHubModule, eventHub, backgroundRoutine),
		remotehubDao: dao.New(batisClient),
		endpointName: endpointName,
	}

	return ptr
}

func (s *RemoteHub) Notify(event event.Event, result event.Result) {
	return
}
