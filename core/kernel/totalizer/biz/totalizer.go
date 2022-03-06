package biz

import (
	"time"

	log "github.com/cihub/seelog"

	bc "github.com/muidea/magicBatis/common"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Totalizer) initializerTotalizer(totalizerPtr *model.Totalizer) {
	type2Totalizer, ok := s.namespace2Totalizer[totalizerPtr.Namespace]
	if !ok {
		type2Totalizer = Type2Totalizer{}
	}

	n2Totalizer, ok := type2Totalizer[totalizerPtr.Type]
	if !ok {
		n2Totalizer = Owner2Totalizer{}
	}

	curTotalizer := s.loadTotalizer(totalizerPtr)
	if curTotalizer != nil {
		totalizerPtr = curTotalizer
	} else {
		totalizerPtr = s.saveTotalizer(totalizerPtr)
	}

	n2Totalizer[totalizerPtr.Owner] = totalizerPtr
	type2Totalizer[totalizerPtr.Type] = n2Totalizer
	s.namespace2Totalizer[totalizerPtr.Namespace] = type2Totalizer
}

func (s *Totalizer) uninitializedTotalizer(totalizerPtr *model.Totalizer) {
	type2Totalizer, ok := s.namespace2Totalizer[totalizerPtr.Namespace]
	if !ok {
		return
	}

	n2Totalizer, ok := type2Totalizer[totalizerPtr.Type]
	if !ok {
		return
	}

	delete(n2Totalizer, totalizerPtr.Owner)
	if len(n2Totalizer) == 0 {
		delete(type2Totalizer, totalizerPtr.Type)
	}
	if len(type2Totalizer) == 0 {
		delete(s.namespace2Totalizer, totalizerPtr.Namespace)
	}
}

func (s *Totalizer) createInternal(owner, namespace string) {
	log.Infof("create internal totalizer, owner:%v, namespace:%v", owner, namespace)
	ptr := model.NewTotalizer(owner, common.TotalizeRealtime, namespace)
	ptr.Value = 1
	s.initializerTotalizer(ptr)

	ptr = model.NewTotalizer(owner, common.TotalizeWeek, namespace)
	ptr.Value = 1
	s.initializerTotalizer(ptr)

	ptr = model.NewTotalizer(owner, common.TotalizeMonth, namespace)
	ptr.Value = 1
	s.initializerTotalizer(ptr)
}

func (s *Totalizer) onCreateTotalizer(totalizerPtr *model.Totalizer) {
	s.Invoke(func() {
		s.initializerTotalizer(totalizerPtr)
	})
}

func (s *Totalizer) onDeleteTotalizer(totalizerPtr *model.Totalizer) {
	s.Invoke(func() {
		s.uninitializedTotalizer(totalizerPtr)
	})
}

func (s *Totalizer) onNotifyTotalizer(eventID string, action int, namespace string) {
	eventPath, _ := fn.SplitRESTURL(eventID)
	feature := fn.FormatRoutePattern(eventPath, nil)
	switch action {
	case common.Create:
		s.onIncreaseTotalizer(feature, namespace)
	case common.Delete:
		s.onDecreaseTotalizer(feature, namespace)
	}
}

func (s *Totalizer) onIncreaseTotalizer(owner, namespace string) {
	s.Invoke(func() {
		type2Totalizer, ok := s.namespace2Totalizer[namespace]
		if !ok {
			s.createInternal(owner, namespace)
			type2Totalizer, ok = s.namespace2Totalizer[namespace]
		}

		//ntv := Type2Totalizer{}
		for _, tv := range type2Totalizer {
			//nnv := Owner2Totalizer{}
			for k, v := range tv {
				if k != owner {
					//nnv[k] = v
					continue
				}

				v.Value += 1
				v.TimeStamp = time.Now().UTC().Unix()

				v = s.saveTotalizer(v)
				if v == nil {
					continue
				}

				//nnv[k] = v
				tv[k] = v
			}
			//ntv[tk] = nnv
		}
	})
}

func (s *Totalizer) onDecreaseTotalizer(owner, namespace string) {
	s.Invoke(func() {
		type2Totalizer, ok := s.namespace2Totalizer[namespace]
		if !ok {
			s.createInternal(owner, namespace)
			type2Totalizer, ok = s.namespace2Totalizer[namespace]
		}

		//ntv := Type2Totalizer{}
		for _, tv := range type2Totalizer {
			//nnv := Owner2Totalizer{}
			for k, v := range tv {
				if k != owner {
					//nnv[k] = v
					continue
				}

				v.Value -= 1
				if v.Value <= 0 {
					v.Value = 0
				}

				v.TimeStamp = time.Now().UTC().Unix()

				v = s.saveTotalizer(v)
				if v == nil {
					continue
				}

				//nnv[k] = v
				tv[k] = v
			}
			//ntv[tk] = nnv
		}
	})
}

func (s *Totalizer) onTimerNotify(eventPtr *common.TimerNotify) {
	s.Invoke(func() {
		_, preWeek := eventPtr.PreTime.ISOWeek()
		_, curWeek := eventPtr.CurTime.ISOWeek()
		refreshWeek := curWeek != preWeek
		refreshMonth := eventPtr.PreTime.Month() != eventPtr.CurTime.Month()
		log.Infof("onTimerNotify, refreshWeek:%v, refreshMonth:%v,curWeek:%v,curMonth:%v", refreshWeek, refreshMonth, eventPtr.CurTime.Weekday(), eventPtr.CurTime.Month())
		for _, nsv := range s.namespace2Totalizer {
			for tk, tv := range nsv {
				if refreshWeek && tk == common.TotalizeWeek {
					// 每周一 更新周统计
					if eventPtr.CurTime.Weekday() == time.Monday {
						for _, nv := range tv {
							nv.Catalog = model.TotalizeHistory
							nv.TimeStamp = eventPtr.CurTime.UTC().Unix()
							s.saveTotalizer(nv)

							nv.Reset()
							s.saveTotalizer(nv)
						}
					}
				}
				if refreshMonth && tk == common.TotalizeMonth {
					// 每月第一天 更新月统计
					if eventPtr.CurTime.Day() == 1 {
						for _, nv := range tv {
							nv.Catalog = model.TotalizeHistory
							nv.TimeStamp = eventPtr.CurTime.UTC().Unix()
							s.saveTotalizer(nv)

							nv.Reset()
							s.saveTotalizer(nv)
						}
					}
				}
			}
		}
	})
}

func (s *Totalizer) saveTotalizer(totalizerPtr *model.Totalizer) (ret *model.Totalizer) {
	curTotalizer, curErr := s.totalizerDao.QueryTotalizer(totalizerPtr.Owner, totalizerPtr.Type, totalizerPtr.Catalog, totalizerPtr.Namespace)
	if curErr != nil {
		curTotalizer, curErr = s.totalizerDao.CreateTotalizer(totalizerPtr, totalizerPtr.Namespace)
		if curErr != nil {
			log.Errorf("create totalizer failed, totalizer:%v , err:%v", totalizerPtr, curErr)
			return
		}

		ret = curTotalizer
		return
	}

	curTotalizer.Value = totalizerPtr.Value
	curTotalizer.TimeStamp = totalizerPtr.TimeStamp
	curTotalizer, curErr = s.totalizerDao.UpdateTotalizer(curTotalizer, curTotalizer.Namespace)
	if curErr != nil {
		log.Errorf("update totalizer failed, totalizer:%v, err:%v", curTotalizer, curErr)
		return
	}

	ret = curTotalizer
	return
}

func (s *Totalizer) loadTotalizer(totalizerPtr *model.Totalizer) (ret *model.Totalizer) {
	curTotalizer, curErr := s.totalizerDao.QueryTotalizer(totalizerPtr.Owner, totalizerPtr.Type, totalizerPtr.Catalog, totalizerPtr.Namespace)
	if curErr != nil {
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
