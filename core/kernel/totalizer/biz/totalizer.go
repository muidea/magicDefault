package biz

import (
	"time"

	log "github.com/cihub/seelog"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func EventID(ptr *common.TotalizeParam) string {
	eventPath, _ := fn.SplitRESTURL(ptr.Trigger)
	return fn.FormatRoutePattern(eventPath, nil)
}

func (s *Totalizer) Load(owner string, typeVal, catalogVal int, namespace string) (ret *model.Totalizer, err error) {
	curTotalizer, curErr := s.totalizerDao.QueryTotalizer(owner, typeVal, catalogVal, namespace)
	if curErr != nil {
		log.Errorf("query totalizer failed, totalizer owner:%v, type:%v, catalog:%v, namespace:%v, err:%v", owner, typeVal, catalogVal, namespace, curErr)
		ptr := &model.Totalizer{
			Owner:     owner,
			Type:      typeVal,
			TimeStamp: time.Now().UTC().Unix(),
			Value:     0,
			Catalog:   catalogVal,
		}

		curTotalizer, curErr = s.totalizerDao.CreateTotalizer(ptr, namespace)
		if curErr != nil {
			log.Errorf("create totalizer failed, totalizer owner:%v, type:%v, catalog:%v, namespace:%v, err:%v", owner, typeVal, catalogVal, namespace, curErr)
			err = curErr
			return
		}
	}

	ret = curTotalizer
	return
}

func (s *Totalizer) Save(ptr *model.Totalizer) (ret *model.Totalizer, err error) {
	curTotalizer, curErr := s.totalizerDao.QueryTotalizer(ptr.Owner, ptr.Type, ptr.Catalog, ptr.Namespace)
	if curErr != nil {
		curTotalizer, curErr = s.totalizerDao.CreateTotalizer(ptr, ptr.Namespace)
		if curErr != nil {
			log.Errorf("create totalizer failed, totalizer owner:%v, type:%v, catalog:%v, namespace:%v, err:%v",
				ptr.Owner, ptr.Type, ptr.Catalog, ptr.Namespace, curErr)
			return
		}

		ret = curTotalizer
		return
	}

	curTotalizer.Value = ptr.Value
	curTotalizer.TimeStamp = ptr.TimeStamp
	curTotalizer, curErr = s.totalizerDao.UpdateTotalizer(curTotalizer, curTotalizer.Namespace)
	if curErr != nil {
		log.Errorf("update totalizer failed, totalizer owner:%v, type:%v, catalog:%v, namespace:%v, err:%v",
			ptr.Owner, ptr.Type, ptr.Catalog, ptr.Namespace, curErr)
		err = curErr
		return
	}

	ret = curTotalizer
	return
}

func (s *Totalizer) initializerTotalizer(paramPtr *common.TotalizeParam, namespace string) {
	trigger2Totalizer, ok := s.namespace2Totalizer[namespace]
	if !ok {
		trigger2Totalizer = Trigger2Totalizer{}
	}

	triggerEvent := EventID(paramPtr)
	totalizerList, ok := trigger2Totalizer[triggerEvent]
	if !ok {
		totalizerList = TotalizerList{}
	}

	existFlag := false
	for _, val := range totalizerList {
		if val.Same(paramPtr.Owner, paramPtr.Trigger) {
			existFlag = true
			break
		}
	}
	if existFlag {
		log.Errorf("duplicate totalizer, owner:%s,trigger:%s", paramPtr.Owner, paramPtr.Trigger)
		return
	}

	totalizerPtr := NewTotalizer(paramPtr.Owner, paramPtr.Trigger, paramPtr.Period, namespace, s)
	totalizerList = append(totalizerList, totalizerPtr)
	trigger2Totalizer[triggerEvent] = totalizerList
	s.namespace2Totalizer[namespace] = trigger2Totalizer
	s.Subscribe(triggerEvent, s)
}

func (s *Totalizer) uninitializedTotalizer(paramPtr *common.TotalizeParam, namespace string) {
	trigger2Totalizer, ok := s.namespace2Totalizer[namespace]
	if !ok {
		return
	}

	triggerEvent := EventID(paramPtr)
	totalizerList, ok := trigger2Totalizer[triggerEvent]
	if !ok {
		return
	}

	newList := TotalizerList{}
	for _, val := range totalizerList {
		if val.Same(paramPtr.Owner, paramPtr.Trigger) {
			continue
		}

		newList = append(newList, val)
	}

	if len(newList) == 0 {
		delete(trigger2Totalizer, triggerEvent)
		s.Unsubscribe(triggerEvent, s)
	}
	if len(trigger2Totalizer) == 0 {
		delete(s.namespace2Totalizer, namespace)
	}
}

func (s *Totalizer) onCreateTotalizer(paramPtr *common.TotalizeParam, namespace string) {
	s.Invoke(func() {
		s.initializerTotalizer(paramPtr, namespace)
	})
}

func (s *Totalizer) onDeleteTotalizer(paramPtr *common.TotalizeParam, namespace string) {
	s.Invoke(func() {
		s.uninitializedTotalizer(paramPtr, namespace)
	})
}

func (s *Totalizer) onNotifyTotalizer(event event.Event, namespace string) {
	s.Invoke(func() {
		trigger2Totalizer, ok := s.namespace2Totalizer[namespace]
		if !ok {
			return
		}

		for k, v := range trigger2Totalizer {
			if !event.Match(k) {
				continue
			}

			for _, sv := range v {
				sv.Trigger(event)
			}
		}
	})
}

func (s *Totalizer) onTimerNotify(eventPtr *common.TimerNotify) {
	s.Invoke(func() {
		_, preWeek := eventPtr.PreTime.ISOWeek()
		_, curWeek := eventPtr.CurTime.ISOWeek()
		refreshWeek := curWeek != preWeek
		refreshDaily := eventPtr.PreTime.Day() != eventPtr.CurTime.Day()
		refreshMonth := eventPtr.PreTime.Month() != eventPtr.CurTime.Month()
		log.Infof("onTimerNotify, refreshDaily:%v, refreshWeek:%v, refreshMonth:%v,curWeek:%v,curMonth:%v", refreshDaily, refreshWeek, refreshMonth, eventPtr.CurTime.Weekday(), eventPtr.CurTime.Month())
		for _, nsv := range s.namespace2Totalizer {
			for _, tv := range nsv {
				if refreshDaily {
					// 每日更新
					for _, nv := range tv {
						nv.Period(common.TotalizeDaily, eventPtr.CurTime)
					}
				}
				if refreshWeek {
					// 每周一 更新周统计
					if eventPtr.CurTime.Weekday() == time.Monday {
						for _, nv := range tv {
							nv.Period(common.TotalizeWeek, eventPtr.CurTime)
						}
					}
				}
				if refreshMonth {
					// 每月第一天 更新月统计
					if eventPtr.CurTime.Day() == 1 {
						for _, nv := range tv {
							nv.Period(common.TotalizeMonth, eventPtr.CurTime)
						}
					}
				}
			}
		}
	})
}

func (s *Totalizer) saveHistoryTotalizer(totalizerPtr *model.Totalizer) (ret *model.Totalizer) {
	if totalizerPtr.Catalog != model.TotalizeHistory {
		log.Errorf("illegal totalizer catalog, owner:%v, type:%v, catalog:%v, namespace:%v",
			totalizerPtr.Owner,
			totalizerPtr.Type,
			totalizerPtr.Catalog,
			totalizerPtr.Namespace)
		return
	}

	curTotalizer, curErr := s.totalizerDao.CreateTotalizer(totalizerPtr, totalizerPtr.Namespace)
	if curErr != nil {
		log.Errorf("save history totalizer failed, totalizer:%v, err:%v", curTotalizer, curErr)
		return
	}

	ret = curTotalizer
	return
}

func (s *Totalizer) filterTotalizer(filter *bc.QueryFilter, namespace string) (ret []*model.Totalizer, err error) {
	totalizerList, _, totalizerErr := s.totalizerDao.FilterTotalizer(filter, namespace)
	if totalizerErr != nil {
		err = totalizerErr
		return
	}

	ret = totalizerList
	return
}
