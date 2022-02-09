package biz

import (
	"time"

	log "github.com/cihub/seelog"
	bc "github.com/muidea/magicBatis/common"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

/*
namespace  -> type -> owner -> totalizer
*/
func (s *Totalizer) onInitializeNamespace(namespacePtr *cc.NamespaceView) {
	s.Invoke(func() {
		totalizerList, _, totalizerErr := s.totalizerDao.FilterTotalizer(nil, namespacePtr.Name)
		if totalizerErr != nil {
			return
		}

		type2Totalizer, ok := s.namespace2Totalizer[namespacePtr.Name]
		if !ok {
			type2Totalizer = Type2Totalizer{}
		}

		for _, val := range totalizerList {
			n2Totalizer, ok := type2Totalizer[val.Type]
			if !ok {
				n2Totalizer = Owner2Totalizer{}
			}

			n2Totalizer[val.Owner] = val
			type2Totalizer[val.Type] = n2Totalizer
		}
		s.namespace2Totalizer[namespacePtr.Name] = type2Totalizer
	})
}

func (s *Totalizer) onCreateTotalizer(totalizerPtr *model.Totalizer) {
	s.Invoke(func() {
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
	})
}

func (s *Totalizer) onDeleteTotalizer(totalizerPtr *model.Totalizer) {
	s.Invoke(func() {
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
	})
}

func (s *Totalizer) onNotifyTotalizer(action int, owner, namespace string) {
	switch action {
	case common.Create:
		s.onIncreaseTotalizer(owner, namespace)
	case common.Delete:
		s.onDecreaseTotalizer(owner, namespace)
	}
}

func (s *Totalizer) onIncreaseTotalizer(owner, namespace string) {
	s.Invoke(func() {
		type2Totalizer, ok := s.namespace2Totalizer[namespace]
		if !ok {
			return
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
			return
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
				if refreshWeek && tk == model.TotalizeWeek {
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
				if refreshMonth && tk == model.TotalizeMonth {
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
