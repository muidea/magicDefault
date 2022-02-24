package biz

import (
	log "github.com/cihub/seelog"
	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/totalizer/dao"
	"github.com/muidea/magicDefault/model"
)

type Totalizer struct {
	biz.Base

	totalizerDao dao.Totalizer

	endpointName        string
	namespace2Totalizer Namespace2Totalizer
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Totalizer {
	totalizer := &Totalizer{
		Base:                biz.New(common.TotalizerModule, eventHub, backgroundRoutine),
		totalizerDao:        dao.New(batisClient),
		endpointName:        endpointName,
		namespace2Totalizer: Namespace2Totalizer{},
	}

	eventHub.Subscribe(common.NotifyAuthorityNamespace, totalizer)

	eventHub.Subscribe(common.QueryTotalizer, totalizer)
	eventHub.Subscribe(common.CreateTotalizer, totalizer)
	eventHub.Subscribe(common.DeleteTotalizer, totalizer)

	eventHub.Subscribe(common.NotifyTimer, totalizer)
	eventHub.Subscribe(common.NotifyEventMask, totalizer)

	return totalizer
}

func (s *Totalizer) Notify(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	log.Infof("notify event, id:%s,source:%s,destination:%s, namespace:%s", event.ID(), event.Source(), event.Destination(), namespace)
	if event.Match(common.NotifyAuthorityNamespace) {
		namespaceLite, namespaceOK := event.Data().(*cc.NamespaceView)
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
	if event.Match(common.NotifyEventMask) {
		owner := event.Header().GetString("owner")
		action := event.Header().GetInt("action")
		s.onNotifyTotalizer(action, owner, namespace)
		return
	}
}
