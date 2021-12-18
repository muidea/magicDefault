package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"

	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	casCommon "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
)

// Client client interface
type Client interface {
	RefreshAccessSession() (*casCommon.EntityView, *commonSession.SessionInfo, error)
	LoginAccessAccount(account, password string) (*casCommon.EntityView, *commonSession.SessionInfo, error)
	LogoutAccessAccount() (*commonSession.SessionInfo, error)
	UpdateAccountPassword(ptr *casCommon.UpdatePasswordParam) (ret *casCommon.AccountView, err error)
	VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*casCommon.EntityView, *commonSession.SessionInfo, error)
	VerifyEntityRole(ptr *casCommon.EntityView) (*casCommon.RoleView, error)
	QueryAccessEntity(id int) (*casCommon.EntityView, error)
	FilterAccessLog(ptr *casCommon.EntityView, filter *util.PageFilter) ([]*casCommon.LogView, int64, error)

	AttachContext(ctx commonSession.ContextInfo)
	DetachContext()

	AttachNameSpace(namespace string)

	BindSession(sessionInfo *commonSession.SessionInfo)
	UnBindSession()

	Release()
}

// NewClient new client
func NewClient(serverURL string) Client {
	clnt := &client{serverURL: serverURL, httpClient: &http.Client{}}

	return clnt
}

type client struct {
	serverURL   string
	sessionInfo *commonSession.SessionInfo
	contextInfo commonSession.ContextInfo
	namespace   string
	httpClient  *http.Client
}

func (s *client) RefreshAccessSession() (*casCommon.EntityView, *commonSession.SessionInfo, error) {
	result := &casCommon.RefreshResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.RefreshSession}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("刷新会话失败，%s", result.Reason)
		return nil, nil, err
	}

	return result.Entity, result.SessionInfo, nil
}

func (s *client) LoginAccessAccount(account, password string) (*casCommon.EntityView, *commonSession.SessionInfo, error) {
	result := &casCommon.LoginResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LoginAccount}, "")
	url.RawQuery = vals.Encode()
	param := &casCommon.LoginParam{Account: account, Password: password}
	_, err := net.HTTPPost(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, s.sessionInfo, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("登录失败，%s", result.Reason)
	}

	return result.Entity, result.SessionInfo, err
}

func (s *client) LogoutAccessAccount() (*commonSession.SessionInfo, error) {
	result := &casCommon.LogoutResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LogoutAccount}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return s.sessionInfo, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("登出失败，%s", result.Reason)
	}

	return result.SessionInfo, err
}

func (s *client) UpdateAccountPassword(ptr *casCommon.UpdatePasswordParam) (ret *casCommon.AccountView, err error) {
	result := &casCommon.AccountResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateAccountPassword}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("修改账号密码失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*casCommon.EntityView, *commonSession.SessionInfo, error) {
	result := &casCommon.VerifyEndpointResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifyEndpoint}, "")
	url.RawQuery = vals.Encode()
	param := &casCommon.VerifyEndpointParam{Endpoint: endpointName, IdentifyID: identifyID, AuthToken: authToken}
	_, err := net.HTTPPost(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, s.sessionInfo, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("校验终端失败，%s", result.Reason)
	}

	return result.Entity, result.SessionInfo, err
}

func (s *client) VerifyEntityRole(ptr *casCommon.EntityView) (*casCommon.RoleView, error) {
	result := &casCommon.EntityRoleResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	vals = ptr.Encode(vals)
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifyEntityRole}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("校验访问对象失败，%s", result.Reason)
		return nil, err
	}

	return result.Role, nil
}

func (s *client) QueryAccessEntity(id int) (ret *casCommon.EntityView, err error) {
	result := &casCommon.QueryEntityResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEntity}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询失败，%s", result.Reason)
		return
	}

	ret = result.Entity
	return
}

func (s *client) FilterAccessLog(entityPtr *casCommon.EntityView, filter *util.PageFilter) (ret []*casCommon.LogView, total int64, err error) {
	result := &casCommon.AccessLogListResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	if entityPtr != nil {
		vals = entityPtr.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryAccessLog}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询日志失败，%s", result.Reason)
		return
	}

	ret = result.AccessLog
	total = result.Total
	return
}

func (s *client) AttachContext(ctx commonSession.ContextInfo) {
	s.contextInfo = ctx
}

func (s *client) DetachContext() {
	s.contextInfo = nil
}

func (s *client) getContextValues() url.Values {
	ret := url.Values{}
	if s.contextInfo != nil {
		ret = s.contextInfo.Encode(ret)
	}
	if s.namespace != "" {
		ret.Set(casCommon.NamespaceID, s.namespace)
	}

	return ret
}

func (s *client) AttachNameSpace(namespace string) {
	s.namespace = namespace
}

func (s *client) BindSession(sessionInfo *commonSession.SessionInfo) {
	s.sessionInfo = sessionInfo
}

func (s *client) UnBindSession() {
	s.sessionInfo = nil
}

func (s *client) Release() {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
		s.httpClient = nil
	}
}
