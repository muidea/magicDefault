package biz

import (
	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
)

type Authority struct {
	biz.Base

	casService string
}

func New(
	casService string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Authority {
	ptr := &Authority{
		Base:       biz.New(common.AuthorityModule, eventHub, backgroundRoutine),
		casService: casService,
	}

	eventHub.Subscribe(common.LoadAuthorityNamespace, ptr)
	eventHub.Subscribe(common.FilterAuthorityNamespace, ptr)

	eventHub.Subscribe(common.QueryAuthorityAccount, ptr)
	eventHub.Subscribe(common.CreateAuthorityAccount, ptr)

	return ptr
}

func (s *Authority) Notify(event event.Event, result event.Result) {
	namespace := event.Header().GetString("namespace")
	var sessionInfo *session.SessionInfo
	val := event.Header().Get("sessionInfo")
	if val != nil {
		sessionInfo = val.(*session.SessionInfo)
	}
	if event.Match(common.LoadAuthorityNamespace) {
		filter, ok := event.Data().(*fu.ContentFilter)
		if !ok {
			return
		}

		namespaceList, _, namespaceErr := s.FilterAuthorityNamespace(sessionInfo, filter, namespace)
		if namespaceErr != nil {
			return
		}
		if result != nil {
			result.Set(namespaceList, namespaceErr)
		}
		return
	}
	if event.Match(common.FilterAuthorityNamespace) {
		filter, ok := event.Data().(*fu.ContentFilter)
		if !ok {
			return
		}

		namespaceList, _, namespaceErr := s.FilterAuthorityNamespace(sessionInfo, filter, namespace)
		if namespaceErr != nil {
			return
		}
		if result != nil {
			result.Set(namespaceList, namespaceErr)
		}
	}

	if event.Match(common.QueryAuthorityAccount) {
		idVal, idErr := fn.SplitRESTID(event.ID())
		if idErr != nil {
			return
		}

		accountPtr, accountErr := s.QueryAuthorityAccount(sessionInfo, idVal, namespace)
		if accountErr != nil {
			return
		}
		if result != nil {
			result.Set(accountPtr, accountErr)
		}
		return
	}

	if event.Match(common.CreateAuthorityAccount) {
		paramPtr := event.Data().(*cc.AccountParam)
		if paramPtr == nil {
			return
		}

		accountPtr, accountErr := s.CreateAuthorityAccount(sessionInfo, paramPtr, namespace)
		if accountErr != nil {
			return
		}
		if result != nil {
			result.Set(accountPtr, accountErr)
		}
		return
	}
}
