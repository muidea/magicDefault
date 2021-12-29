package biz

import (
	"github.com/muidea/magicCommon/event"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
)

func (s *Setting) queryNamespace(sessionInfo *session.SessionInfo, namespace string) (ret *cc.NamespaceView) {
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

	namespaceList, namespaceOK := resultVal.([]*cc.NamespaceView)
	if !namespaceOK || len(namespaceList) <= 0 {
		return
	}

	ret = namespaceList[0]
	return
}
