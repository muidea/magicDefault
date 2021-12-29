package biz

import (
	"fmt"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Setting) QuerySettingProfile(sessionInfo *session.SessionInfo, entityPtr *cc.EntityView, namespace string) (ret *cc.AccountView, err error) {
	eid := fn.FormatID(common.QueryAuthorityAccount, entityPtr.EID)
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("sessionInfo", sessionInfo)

	event := event.NewEvent(eid, s.ID(), common.AuthorityModule, header, nil)
	result := s.CallEvent(event)
	if result == nil {
		err = fmt.Errorf("can't query account")
		return
	}
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		err = resultErr
		return
	}
	accountPtr := resultVal.(*cc.AccountView)
	if accountPtr == nil {
		err = fmt.Errorf("illegal account value")
		return
	}

	ret = accountPtr
	return
}

func (s *Setting) QueryOperateLog(sessionInfo *session.SessionInfo, filter *bc.QueryFilter, namespace string) (ret []*model.Log, err error) {
	eid := common.QueryOperateLog
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("owner", s.ID())

	event := event.NewEvent(eid, s.ID(), common.BaseModule, header, filter)
	result := s.CallEvent(event)
	if result == nil {
		err = fmt.Errorf("can't query operate log")
		return
	}
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		err = resultErr
		return
	}
	logList, logOK := resultVal.([]*model.Log)
	if !logOK {
		err = fmt.Errorf("query operate log failed")
		return
	}

	ret = logList
	return
}
