package biz

import (
	log "github.com/cihub/seelog"
	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	casCommon "github.com/muidea/magicCas/common"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/totalizer/dao"
	"github.com/muidea/magicDefault/model"
)

type Totalizer struct {
	biz.Base

	totalizerDao dao.Totalizer

	namespace2Totalizer Namespace2Totalizer
}

func New(
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Totalizer {
	totalizer := &Totalizer{
		Base:                biz.New(common.TotalizerModule, eventHub, backgroundRoutine),
		totalizerDao:        dao.New(batisClient),
		namespace2Totalizer: Namespace2Totalizer{},
	}

	eventHub.Subscribe(common.InitializeAuthorityNamespace, totalizer)
	eventHub.Subscribe(common.CreateTotalizer, totalizer)
	eventHub.Subscribe(common.DeleteTotalizer, totalizer)
	eventHub.Subscribe(common.UpdateTotalizer, totalizer)
	eventHub.Subscribe(common.QueryTotalizer, totalizer)
	eventHub.Subscribe(common.NotifyTimer, totalizer)

	eventHub.Subscribe(common.CreateEventMask, totalizer)
	eventHub.Subscribe(common.DeleteEventMask, totalizer)
	eventHub.Subscribe(common.UpdateEventMask, totalizer)

	return totalizer
}

func (s *Totalizer) Notify(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	owner := event.Header().GetString("owner")
	log.Infof("notify event, id:%s,source:%s,destination:%s, namespace:%s, owner:%s", event.ID(), event.Source(), event.Destination(), namespace, owner)
	if event.Match(common.InitializeAuthorityNamespace) {
		namespaceLite, namespaceOK := event.Data().(*casCommon.NamespaceView)
		if !namespaceOK {
			return
		}

		s.onInitializeNamespace(namespaceLite)
		return
	}
	if event.Match(common.CreateTotalizer) {
		totalizerPtr, totalizerOK := event.Data().(*model.Totalizer)
		if !totalizerOK {
			return
		}

		s.onCreateTotalizer(totalizerPtr)
		return
	}
	if event.Match(common.DeleteTotalizer) {
		totalizerPtr, totalizerOK := event.Data().(*model.Totalizer)
		if !totalizerOK {
			return
		}

		s.onDeleteTotalizer(totalizerPtr)
		return
	}
	if event.Match(common.UpdateTotalizer) {
		return
	}
	if event.Match(common.QueryTotalizer) {
		filterPtr, filterOK := event.Data().(*bc.QueryFilter)
		if !filterOK {
			return
		}
		totalizerList, totalizerErr := s.filterTotalizer(filterPtr, namespace)
		if result != nil {
			result.Set(totalizerList, totalizerErr)
		}
		return
	}
	if event.Match(common.NotifyTimer) {
		eventPtr, eventOK := event.Data().(*common.TimerNotify)
		if !eventOK {
			return
		}
		s.onTimerNotify(eventPtr)
		return
	}
	if event.Match(common.CreateEventMask) {
		s.onIncreaseTotalizer(owner, namespace)
		return
	}
	if event.Match(common.UpdateEventMask) {
		return
	}
	if event.Match(common.DeleteEventMask) {
		s.onDecreaseTotalizer(owner, namespace)
		return
	}
}
