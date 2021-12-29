package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/common"
)

// Client client interface
type Client interface {
	RefreshAccessSession() (*cc.EntityView, *session.SessionInfo, error)
	LoginAccessAccount(account, password string) (*cc.EntityView, *session.SessionInfo, error)
	LogoutAccessAccount() (*session.SessionInfo, error)
	UpdateAccountPassword(ptr *cc.UpdatePasswordParam) (ret *cc.AccountView, err error)
	VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*cc.EntityView, *session.SessionInfo, error)
	VerifyEntityRole(ptr *cc.EntityView) (*cc.RoleView, error)
	QueryAccessEntity(id int) (*cc.EntityView, error)
	FilterAccessLog(ptr *cc.EntityView, filter *util.Pagination) ([]*cc.LogView, int64, error)

	AttachContext(ctx session.ContextInfo)
	DetachContext()

	AttachNameSpace(namespace string)

	BindSession(sessionInfo *session.SessionInfo)
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
	sessionInfo *session.SessionInfo
	contextInfo session.ContextInfo
	namespace   string
	httpClient  *http.Client
}

func (s *client) RefreshAccessSession() (*cc.EntityView, *session.SessionInfo, error) {
	result := &cc.RefreshResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("刷新会话失败，%s", result.Reason)
		return nil, nil, err
	}

	return result.Entity, result.SessionInfo, nil
}

func (s *client) LoginAccessAccount(account, password string) (*cc.EntityView, *session.SessionInfo, error) {
	result := &cc.LoginResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LoginAccount}, "")
	url.RawQuery = vals.Encode()
	param := &cc.LoginParam{Account: account, Password: password}
	_, err := net.HTTPPost(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, s.sessionInfo, err
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("登录失败，%s", result.Reason)
	}

	return result.Entity, result.SessionInfo, err
}

func (s *client) LogoutAccessAccount() (*session.SessionInfo, error) {
	result := &cc.LogoutResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("登出失败，%s", result.Reason)
	}

	return result.SessionInfo, err
}

func (s *client) UpdateAccountPassword(ptr *cc.UpdatePasswordParam) (ret *cc.AccountView, err error) {
	result := &cc.AccountResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("修改账号密码失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*cc.EntityView, *session.SessionInfo, error) {
	result := &cc.VerifyEndpointResult{}

	vals := url.Values{}
	if s.sessionInfo != nil {
		vals = s.sessionInfo.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifyEndpoint}, "")
	url.RawQuery = vals.Encode()
	param := &cc.VerifyEndpointParam{Endpoint: endpointName, IdentifyID: identifyID, AuthToken: authToken}
	_, err := net.HTTPPost(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, s.sessionInfo, err
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("校验终端失败，%s", result.Reason)
	}

	return result.Entity, result.SessionInfo, err
}

func (s *client) VerifyEntityRole(ptr *cc.EntityView) (*cc.RoleView, error) {
	result := &cc.EntityRoleResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("校验访问对象失败，%s", result.Reason)
		return nil, err
	}

	return result.Role, nil
}

func (s *client) QueryAccessEntity(id int) (ret *cc.EntityView, err error) {
	result := &cc.QueryEntityResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询失败，%s", result.Reason)
		return
	}

	ret = result.Entity
	return
}

func (s *client) FilterAccessLog(entityPtr *cc.EntityView, filter *util.Pagination) (ret []*cc.LogView, total int64, err error) {
	result := &cc.AccessLogListResult{}

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

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询日志失败，%s", result.Reason)
		return
	}

	ret = result.AccessLog
	total = result.Total
	return
}

func (s *client) AttachContext(ctx session.ContextInfo) {
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
		ret.Set(cc.NamespaceID, s.namespace)
	}

	return ret
}

func (s *client) AttachNameSpace(namespace string) {
	s.namespace = namespace
}

func (s *client) BindSession(sessionInfo *session.SessionInfo) {
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
