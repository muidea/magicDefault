package biz

import (
	"path"
	"sync"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/generator"
	fn "github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

type Base struct {
	id                string
	codeGenerator     sync.Map
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
	ptr := &invokeTask{funcPtr: funcPtr}

	s.backgroundRoutine.Post(ptr)
}

func (s *Base) RootDestination() string {
	return "/#"
}

func (s *Base) InnerDestination() string {
	return s.ID()
}

func (s *Base) GetGenerator(pattern string) generator.Generator {
	val, ok := s.codeGenerator.Load(pattern)
	if !ok {
		generator, _ := generator.New(pattern)
		s.codeGenerator.Store(pattern, generator)
		return generator
	}

	return val.(generator.Generator)
}

func (s *Base) RemoveGenerator(pattern string) {
	s.codeGenerator.Delete(pattern)
}

func (s *Base) GenerateCode(generator generator.Generator) string {
	task := task.NewGeneratorTask(generator)
	s.backgroundRoutine.Invoke(task)
	return task.Result()
}

func (s *Base) GetPatternCode(pattern string) string {
	return s.GenerateCode(s.GetGenerator(pattern))
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

func (s *Base) WriteLog(ptr *model.Log, namespace string) {
	eid := common.WriteOperateLog
	header := event.NewValues()
	header.Set("namespace", namespace)

	writeEvent := event.NewEvent(eid, s.ID(), common.BaseModule, header, ptr)
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
