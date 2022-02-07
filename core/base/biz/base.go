package biz

import (
	"path"
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

type Base struct {
	id                string
	eventHub          event.Hub
	backgroundRoutine task.BackgroundRoutine
}

type invokeTask struct {
	funcPtr func()
}

func (s *invokeTask) Run() {
	s.funcPtr()
}

func New(
	id string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) Base {
	return Base{
		id:                id,
		eventHub:          eventHub,
		backgroundRoutine: backgroundRoutine,
	}
}

func (s *Base) ID() string {
	return s.id
}

func (s *Base) PostEvent(event event.Event) {
	s.eventHub.Post(event)
}

func (s *Base) SendEvent(event event.Event) event.Result {
	return s.eventHub.Send(event)
}

func (s *Base) CallEvent(event event.Event) event.Result {
	return s.eventHub.Call(event)
}

func (s *Base) Invoke(funcPtr func()) {
	taskPtr := &invokeTask{funcPtr: funcPtr}
	s.backgroundRoutine.Post(taskPtr)
}

func (s *Base) Timer(intervalValue time.Duration, offsetValue time.Duration, funcPtr func()) {
	taskPtr := &invokeTask{funcPtr: funcPtr}
	s.backgroundRoutine.Timer(taskPtr, intervalValue, offsetValue)
}

func (s *Base) BroadCast(eid string, header event.Values, val interface{}) {
	event := event.NewEvent(eid, s.ID(), s.RootDestination(), header, val)
	s.eventHub.Post(event)
}

func (s *Base) RootDestination() string {
	return "/#"
}

func (s *Base) InnerDestination() string {
	return s.ID()
}

func (s *Base) QueryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	eid := fn.FormatID(common.QueryEntity, id)
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("sessionInfo", sessionInfo)

	queryEvent := event.NewEvent(eid, s.ID(), common.BaseModule, header, nil)
	queryResult := s.CallEvent(queryEvent)
	resultVal, resultErr := queryResult.Get()
	if resultErr != nil {
		return
	}

	ret = resultVal.(*cc.EntityView)
	return
}

func (s *Base) WriteLog(memo, address string, entityPtr *cc.EntityView, namespace string) {
	eid := common.WriteOperateLog
	header := event.NewValues()
	header.Set("namespace", namespace)

	logPtr := &model.Log{Address: address, Memo: memo, Creater: entityPtr.ID, CreateTime: time.Now().UTC().Unix()}
	writeEvent := event.NewEvent(eid, s.ID(), common.BaseModule, header, logPtr)
	s.PostEvent(writeEvent)
	return
}

func (s *Base) CheckIn(feature, namespace string) {
	owner := path.Join(s.ID(), feature)

	eid := common.CreateTotalizer
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("owner", owner)

	rtdTotalizer := model.NewTotalizer(owner, model.TotalizeRealtime, namespace)
	rtdEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, rtdTotalizer)
	s.PostEvent(rtdEvent)

	weekTotalizer := model.NewTotalizer(owner, model.TotalizeWeek, namespace)
	weekEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, weekTotalizer)
	s.PostEvent(weekEvent)

	monthTotalizer := model.NewTotalizer(owner, model.TotalizeMonth, namespace)
	monthEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, monthTotalizer)
	s.PostEvent(monthEvent)
}

func (s *Base) CheckOut(feature, namespace string) {
	owner := path.Join(s.ID(), feature)

	eid := common.DeleteTotalizer
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("owner", owner)

	rtdTotalizer := model.NewTotalizer(owner, model.TotalizeRealtime, namespace)
	rtdEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, rtdTotalizer)
	s.PostEvent(rtdEvent)

	weekTotalizer := model.NewTotalizer(owner, model.TotalizeWeek, namespace)
	weekEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, weekTotalizer)
	s.PostEvent(weekEvent)

	monthTotalizer := model.NewTotalizer(owner, model.TotalizeMonth, namespace)
	monthEvent := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, monthTotalizer)
	s.PostEvent(monthEvent)
}

func (s *Base) QuerySummary(feature, namespace string) (ret []*model.Totalizer) {
	owner := path.Join(s.ID(), feature)

	eid := common.QueryTotalizer
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("owner", owner)

	filter := bc.NewFilter()
	filter.Equal("Owner", owner)
	filter.Equal("Catalog", model.TotalizeCurrent)

	eventPtr := event.NewEvent(eid, s.ID(), common.TotalizerModule, header, filter)
	resultPtr := s.CallEvent(eventPtr)
	resultVal, resultErr := resultPtr.Get()
	if resultErr != nil {
		return
	}

	totalizerList, totalizerOK := resultVal.([]*model.Totalizer)
	if !totalizerOK {
		return
	}

	ret = totalizerList
	return
}
