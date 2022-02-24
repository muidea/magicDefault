package biz

import (
	log "github.com/cihub/seelog"
	bc "github.com/muidea/magicBatis/common"
	cClnt "github.com/muidea/magicCas/client"
	cc "github.com/muidea/magicCas/common"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicDefault/model"
)

func (s *Base) LoginAccount(curSessionInfo *session.SessionInfo, account, password, namespace string) (entityPtr *cc.EntityView, sessionInfo *session.SessionInfo, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	loginEntity, loginSession, loginErr := clnt.LoginAccessAccount(account, password)
	if loginErr != nil {
		err = loginErr
		log.Errorf("login account failed, err:%s", err.Error())
		return
	}

	entityPtr = loginEntity
	sessionInfo = loginSession
	return
}

func (s *Base) LogoutAccount(curSessionInfo *session.SessionInfo, namespace string) (ret *session.SessionInfo, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	logoutSession, logoutErr := clnt.LogoutAccessAccount()
	if logoutErr != nil {
		err = logoutErr
		log.Errorf("logout account failed, err:%s", err.Error())
		return
	}

	ret = logoutSession

	return
}

func (s *Base) UpdateAccountPassword(curSessionInfo *session.SessionInfo, ptr *cc.UpdatePasswordParam, namespace string) (ret *cc.AccountView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateAccountPassword(ptr)
	return
}

func (s *Base) VerifyEndpoint(curSessionInfo *session.SessionInfo, endpointName, identifyID, authToken, namespace string) (entityPtr *cc.EntityView, sessionInfo *session.SessionInfo, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	confirmEntity, confirmSession, confirmErr := clnt.VerifyAccessEndpoint(endpointName, identifyID, authToken)
	if confirmErr != nil {
		err = confirmErr
		log.Errorf("confirm endpoint failed, err:%s", err.Error())
		return
	}

	entityPtr = confirmEntity
	sessionInfo = confirmSession

	return
}

func (s *Base) RefreshSession(curSessionInfo *session.SessionInfo, namespace string) (entityPtr *cc.EntityView, sessionInfoPtr *session.SessionInfo, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	refreshEntity, refreshSessionInfo, refreshErr := clnt.RefreshAccessSession()
	if refreshErr != nil {
		err = refreshErr
		return
	}

	entityPtr = refreshEntity
	sessionInfoPtr = refreshSessionInfo
	return
}

func (s *Base) VerifyEntityRole(curSessionInfo *session.SessionInfo, ptr *cc.EntityView, namespace string) (ret *cc.RoleView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	verifyRole, verifyErr := clnt.VerifyEntityRole(ptr)
	if verifyErr != nil {
		err = verifyErr
		log.Errorf("verify account failed, err:%s", err.Error())
		return
	}

	ret = verifyRole
	return
}

func (s *Base) QueryEntity(curSessionInfo *session.SessionInfo, id int, namespace string) (*cc.EntityView, error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	return clnt.QueryAccessEntity(id)
}

func (s *Base) QueryAccessLog(curSessionInfo *session.SessionInfo, entityPtr *cc.EntityView, filter *fu.Pagination, namespace string) (ret []*cc.LogView, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(curSessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterAccessLog(entityPtr, filter)
	return
}

func (s *Base) WriteOperateLog(log *model.Log, namespace string) (ret *model.Log, err error) {
	ret, err = s.baseDao.WriteLog(log, namespace)

	return
}

func (s *Base) QueryOperateLog(filter *bc.QueryFilter, namespace string) (ret []*model.Log, total int64, err error) {
	ret, total, err = s.baseDao.QueryLog(filter, namespace)

	return
}
