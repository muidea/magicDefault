package biz

import (
	log "github.com/cihub/seelog"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/totalizer/dao"
)

type TotalizerList []common.Totalizer

type Trigger2Totalizer map[string]TotalizerList

type Namespace2Totalizer map[string]Trigger2Totalizer

type Totalizer struct {
	biz.Base

	totalizerDao dao.Totalizer

	endpointName string

	namespace2Totalizer Namespace2Totalizer
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Totalizer {
	ptr := &Totalizer{
		Base:                biz.New(common.TotalizerModule, eventHub, backgroundRoutine),
		totalizerDao:        dao.New(batisClient),
		endpointName:        endpointName,
		namespace2Totalizer: Namespace2Totalizer{},
	}

	ptr.SubscribeFunc(common.QueryTotalizer, ptr.queryTotalizer)
	ptr.SubscribeFunc(common.CreateTotalizer, ptr.createTotalizer)
	ptr.SubscribeFunc(common.DeleteTotalizer, ptr.deleteTotalizer)
	ptr.SubscribeFunc(common.NotifyTimer, ptr.notifyTimer)

	return ptr
}

func (s *Totalizer) queryTotalizer(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
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

func (s *Totalizer) createTotalizer(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	paramPtr, paramOK := event.Data().(*common.TotalizeParam)
	if !paramOK {
		return
	}

	s.onCreateTotalizer(paramPtr, namespace)
	return
}

func (s *Totalizer) deleteTotalizer(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	paramPtr, paramOK := event.Data().(*common.TotalizeParam)
	if !paramOK {
		return
	}

	s.onDeleteTotalizer(paramPtr, namespace)
	return
}

func (s *Totalizer) notifyTimer(event event.Event, result event.Result) {
	eventPtr, eventOK := event.Data().(*common.TimerNotify)
	if !eventOK {
		return
	}
	s.onTimerNotify(eventPtr)
	return
}

func (s *Totalizer) Notify(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	log.Infof("notify event, id:%s,source:%s,destination:%s, namespace:%s", event.ID(), event.Source(), event.Destination(), namespace)
	if event.Match(common.NotifyEventMask) {
		s.onNotifyTotalizer(event, namespace)
		return
	}
}
