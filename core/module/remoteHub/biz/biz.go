package biz

import (
	log "github.com/cihub/seelog"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/module/remoteHub/dao"
)

type RemoteHub struct {
	biz.Base

	remoteHubDao dao.RemoteHub

	endpointName string
}

func New(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *RemoteHub {
	batisClnt, batisErr := persistence.GetBatisClient()
	if batisErr != nil {
		log.Criticalf("get batis client failed, err:%s", batisErr.Error())
	}

	ptr := &RemoteHub{
		Base:         biz.New(common.RemoteHubModule, eventHub, backgroundRoutine),
		remoteHubDao: dao.New(batisClnt),
		endpointName: endpointName,
	}

	return ptr
}

func (s *RemoteHub) Notify(event event.Event, result event.Result) {
	return
}
