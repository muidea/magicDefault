package biz

import (
	casCommon "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCommon/event"
	fu "github.com/muidea/magicCommon/foundation/util"
	commonSession "github.com/muidea/magicCommon/session"
	"github.com/muidea/magicDefault/common"
)

func (s *Setting) queryNamespace(sessionInfo *commonSession.SessionInfo, namespace string) (ret *casCommon.NamespaceView) {
	eid := common.FilterAuthorityNamespace
	header := event.NewValues()
	header.Set("namespace", namespace)
	header.Set("sessionInfo", sessionInfo)

	filter := fu.NewFilter()
	filter.Set("name", namespace)
	eventPtr := event.NewEvent(eid, s.ID(), common.AuthorityModule, header, filter)
	result := s.CallEvent(eventPtr)
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		return
	}

	namespaceList, namespaceOK := resultVal.([]*casCommon.NamespaceView)
	if !namespaceOK || len(namespaceList) <= 0 {
		return
	}

	ret = namespaceList[0]
	return
}
