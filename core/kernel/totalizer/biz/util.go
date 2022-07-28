package biz

import (
	"fmt"
	"time"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func NewTotalizer(owner, trigger string, period []int, namespace string, serializer common.Serializer) common.Totalizer {
	normalPtr := &normalTotalizer{
		owner:      owner,
		trigger:    trigger,
		serializer: serializer,
	}

	log.Infof("new totalizer, owner:%s, trigger:%s, period:%v, namespace:%v", owner, trigger, period, namespace)

	normalPtr.periodTotalizer = map[int]*model.Totalizer{}
	for _, val := range period {
		ptr, err := serializer.Load(owner, val, model.TotalizeCurrent, namespace)
		if err != nil {
			ptr = model.NewTotalizer(owner, val, namespace)
			ptr, err = serializer.Save(ptr)
			if err != nil {
				log.Errorf("save totalizer failed, owner:%s, trigger:%s, period:%d, err:%s", owner, trigger, val, err.Error())
				continue
			}
		}
		normalPtr.periodTotalizer[val] = ptr
	}

	return normalPtr
}

type normalTotalizer struct {
	owner           string
	trigger         string
	periodTotalizer map[int]*model.Totalizer
	serializer      common.Serializer
}

func (s *normalTotalizer) Owner() string {
	return s.owner
}

func (s *normalTotalizer) Trigger(event event.Event) (err error) {
	if !event.Match(s.trigger) {
		return
	}

	action := event.Header().GetInt("action")
	for idx, val := range s.periodTotalizer {
		switch action {
		case common.Create:
			val.Value += 1
		case common.Delete:
			val.Value -= 1
		}
		if val.Value < 0 {
			val.Value = 0
		}

		val.TimeStamp = time.Now().UTC().Unix()
		if s.serializer != nil {
			val, err = s.serializer.Save(val)
			if err != nil {
				log.Errorf("save totalizer failed, owner:%s, trigger:%s, period:%d, err:%s", s.owner, s.trigger, idx, err.Error())
			}
		}
	}
	return
}

func (s *normalTotalizer) Period(typeVal int, timeStamp time.Time) (err error) {
	log.Infof("totalizer period, typeVal:%d, time:%s", typeVal, timeStamp.Local().Format(time.RFC3339))
	valPtr, valOK := s.periodTotalizer[typeVal]
	if !valOK {
		err = fmt.Errorf("illegal period type, type value:%d", typeVal)
		log.Errorf("period totalizer failed, err:%s", err.Error())
		return
	}

	historyPtr := valPtr.DuplicateHistory()
	if s.serializer != nil {
		_, historyErr := s.serializer.Save(historyPtr)
		if historyErr != nil {
			log.Errorf("save history totalizer failed, owner:%s, trigger:%s, period:%d, err:%s", s.owner, s.trigger, typeVal, historyErr.Error())
		}
	}

	valPtr.Reset(timeStamp)
	_, currentErr := s.serializer.Save(valPtr)
	if currentErr != nil {
		log.Errorf("save current totalizer failed, owner:%s, trigger:%s, period:%d, err:%s", s.owner, s.trigger, typeVal, currentErr.Error())
	}

	s.periodTotalizer[typeVal] = valPtr
	return
}

func (s *normalTotalizer) Get(typeVal int) (ret *model.Totalizer, err error) {
	valPtr, valOK := s.periodTotalizer[typeVal]
	if !valOK {
		err = fmt.Errorf("illegal period type, type value:%d", typeVal)
		return
	}

	ret = valPtr
	return
}

func (s *normalTotalizer) Same(owner, trigger string) bool {
	return s.owner == owner && s.trigger == trigger
}
