package biz

import (
	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/base/dao"
	"github.com/muidea/magicDefault/model"
)

type Base struct {
	biz.Base

	baseDao dao.Base

	casService string
}

func New(
	casService string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Base {
	ptr := &Base{
		Base:       biz.New(common.BaseModule, eventHub, backgroundRoutine),
		casService: casService,
		baseDao:    dao.New(batisClient),
	}

	eventHub.Subscribe(common.QueryEntity, ptr)
	eventHub.Subscribe(common.WriteOperateLog, ptr)
	eventHub.Subscribe(common.QueryOperateLog, ptr)

	return ptr
}

func (s *Base) Notify(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	if event.Match(common.QueryEntity) {
		sessionInfo := event.Header().Get("sessionInfo").(*session.SessionInfo)
		id, err := fn.SplitRESTID(event.ID())
		if err != nil {
			return
		}

		entityPtr, entityErr := s.QueryEntity(sessionInfo, id, namespace)
		if result != nil {
			result.Set(entityPtr, entityErr)
		}
		return
	}
	if event.Match(common.WriteOperateLog) {
		ptr := event.Data().(*model.Log)
		s.WriteOperateLog(ptr, namespace)
		return
	}
	if event.Match(common.QueryOperateLog) {
		filterPtr, filterOK := event.Data().(*bc.QueryFilter)
		if !filterOK {
			return
		}

		logsList, _, logErr := s.QueryOperateLog(filterPtr, namespace)
		if result != nil {
			result.Set(logsList, logErr)
		}
		return
	}
}
