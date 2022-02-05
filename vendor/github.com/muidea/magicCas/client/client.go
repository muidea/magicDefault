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

	"github.com/muidea/magicCas/common"
)

// Client client interface
type Client interface {
	RefreshAccessSession() (*common.EntityView, *commonSession.SessionInfo, error)
	LoginAccessAccount(account, password string) (*common.EntityView, *commonSession.SessionInfo, error)
	LogoutAccessAccount() (*commonSession.SessionInfo, error)
	UpdateAccountPassword(ptr *common.UpdatePasswordParam) (ret *common.AccountView, err error)
	VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*common.EntityView, *commonSession.SessionInfo, error)
	VerifyEntityRole(ptr *common.EntityView) (*common.RoleView, error)
	QueryAccessEntity(id int) (*common.EntityView, error)
	FilterAccessLog(ptr *common.EntityView, filter *util.Pagination) ([]*common.LogView, int64, error)

	FilterAccount(filter *util.ContentFilter) (ret []*common.AccountView, total int64, err error)
	FilterAccountLite(filter *util.ContentFilter) (ret []*common.AccountLite, total int64, err error)
	QueryAccount(id int) (ret *common.AccountView, err error)
	CreateAccount(ptr *common.AccountParam) (ret *common.AccountView, err error)
	UpdateAccount(id int, ptr *common.AccountParam) (ret *common.AccountView, err error)
	DeleteAccount(id int) (ret *common.AccountView, err error)
	CheckAccount(account string) (ret []*common.AccountLite, err error)

	FilterEndpoint(filter *util.ContentFilter) (ret []*common.EndpointView, total int64, err error)
	FilterEndpointLite(filter *util.ContentFilter) (ret []*common.EndpointLite, total int64, err error)
	QueryEndpoint(id int) (ret *common.EndpointView, err error)
	CreateEndpoint(ptr *common.EndpointParam) (ret *common.EndpointView, err error)
	UpdateEndpoint(id int, ptr *common.EndpointParam) (ret *common.EndpointView, err error)
	DeleteEndpoint(id int) (ret *common.EndpointView, err error)

	FilterRole(filter *util.ContentFilter) (ret []*common.RoleView, total int64, err error)
	FilterRoleLite(filter *util.ContentFilter) (ret []*common.RoleLite, total int64, err error)
	QueryRole(id int) (ret *common.RoleView, err error)
	CreateRole(ptr *common.RoleParam) (ret *common.RoleView, err error)
	UpdateRole(id int, ptr *common.RoleParam) (ret *common.RoleView, err error)
	DeleteRole(id int) (ret *common.RoleView, err error)

	FilterNamespace(filter *util.ContentFilter) (ret []*common.NamespaceView, total int64, err error)
	FilterNamespaceLite(filter *util.ContentFilter) (ret []*common.NamespaceLite, total int64, err error)
	QueryNamespace(id int) (ret *common.NamespaceView, err error)
	CreateNamespace(ptr *common.NamespaceParam) (ret *common.NamespaceView, err error)
	UpdateNamespace(id int, ptr *common.NamespaceParam) (ret *common.NamespaceView, err error)
	DeleteNamespace(id int) (ret *common.NamespaceView, err error)

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

func (s *client) RefreshAccessSession() (*common.EntityView, *commonSession.SessionInfo, error) {
	result := &common.RefreshResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.RefreshSession}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("刷新会话状态失败，%s", result.Reason)
		return nil, nil, err
	}

	return result.Entity, result.SessionInfo, nil
}

func (s *client) LoginAccessAccount(account, password string) (*common.EntityView, *commonSession.SessionInfo, error) {
	result := &common.LoginResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LoginAccount}, "")
	url.RawQuery = vals.Encode()
	param := &common.LoginParam{Account: account, Password: password}
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
	result := &common.LogoutResult{}

	vals := url.Values{}
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

func (s *client) UpdateAccountPassword(ptr *common.UpdatePasswordParam) (ret *common.AccountView, err error) {
	result := &common.AccountResult{}

	vals := url.Values{}
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

func (s *client) VerifyAccessEndpoint(endpointName, identifyID, authToken string) (*common.EntityView, *commonSession.SessionInfo, error) {
	result := &common.VerifyEndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifyEndpoint}, "")
	url.RawQuery = vals.Encode()
	param := &common.VerifyEndpointParam{Endpoint: endpointName, IdentifyID: identifyID, AuthToken: authToken}
	_, err := net.HTTPPost(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, s.sessionInfo, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("校验终端失败，%s", result.Reason)
	}

	return result.Entity, result.SessionInfo, err
}

func (s *client) VerifyEntityRole(ptr *common.EntityView) (*common.RoleView, error) {
	result := &common.EntityRoleResult{}

	vals := url.Values{}
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

func (s *client) QueryAccessEntity(id int) (ret *common.EntityView, err error) {
	result := &common.QueryEntityResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEntity}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询指定Entity失败，%s", result.Reason)
		return
	}

	ret = result.Entity
	return
}

func (s *client) FilterAccessLog(entityPtr *common.EntityView, filter *util.Pagination) (ret []*common.LogView, total int64, err error) {
	result := &common.AccessLogListResult{}

	vals := url.Values{}
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

func (s *client) FilterAccount(filter *util.ContentFilter) (ret []*common.AccountView, total int64, err error) {
	result := &common.AccountListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "2")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterAccount}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询账号列表失败，%s", result.Reason)
		return
	}

	ret = result.Account
	total = result.Total
	return
}

func (s *client) FilterAccountLite(filter *util.ContentFilter) (ret []*common.AccountLite, total int64, err error) {
	result := &common.AccountLiteListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "1")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterAccount}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询账号列表失败，%s", result.Reason)
		return
	}

	ret = result.Account
	total = result.Total
	return
}

func (s *client) QueryAccount(id int) (ret *common.AccountView, err error) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询指定账号失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) CreateAccount(ptr *common.AccountParam) (ret *common.AccountView, err error) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateAccount}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPPost(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("新建账号失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) UpdateAccount(id int, ptr *common.AccountParam) (ret *common.AccountView, err error) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("修改账号失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) DeleteAccount(id int) (ret *common.AccountView, err error) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("删除账号失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) CheckAccount(account string) (ret []*common.AccountLite, err error) {
	result := &common.AccountLiteListResult{}

	vals := url.Values{}
	vals.Set("account", account)
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CheckAccount}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		return
	}

	ret = result.Account
	return
}

func (s *client) FilterEndpoint(filter *util.ContentFilter) (ret []*common.EndpointView, total int64, err error) {
	result := &common.EndpointListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "2")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterEndpoint}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	total = result.Total
	return
}

func (s *client) FilterEndpointLite(filter *util.ContentFilter) (ret []*common.EndpointLite, total int64, err error) {
	result := &common.EndpointLiteListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "1")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterEndpoint}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	total = result.Total
	return
}

func (s *client) QueryEndpoint(id int) (ret *common.EndpointView, err error) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询指定终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	return
}

func (s *client) CreateEndpoint(ptr *common.EndpointParam) (ret *common.EndpointView, err error) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateEndpoint}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPPost(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("新建终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	return
}

func (s *client) UpdateEndpoint(id int, ptr *common.EndpointParam) (ret *common.EndpointView, err error) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("更新终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	return
}

func (s *client) DeleteEndpoint(id int) (ret *common.EndpointView, err error) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("删除终端失败，%s", result.Reason)
		return
	}

	ret = result.Endpoint
	return
}

func (s *client) FilterRole(filter *util.ContentFilter) (ret []*common.RoleView, total int64, err error) {
	result := &common.RoleListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "2")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterRole}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	total = result.Total
	return
}

func (s *client) FilterRoleLite(filter *util.ContentFilter) (ret []*common.RoleLite, total int64, err error) {
	result := &common.RoleLiteListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "1")

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterRole}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	total = result.Total
	return
}

func (s *client) QueryRole(id int) (ret *common.RoleView, err error) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询指定权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	return
}

func (s *client) CreateRole(ptr *common.RoleParam) (ret *common.RoleView, err error) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateRole}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPPost(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("新建权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	return
}

func (s *client) UpdateRole(id int, ptr *common.RoleParam) (ret *common.RoleView, err error) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("更新权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	return
}

func (s *client) DeleteRole(id int) (ret *common.RoleView, err error) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("删除权限组失败，%s", result.Reason)
		return
	}

	ret = result.Role
	return
}

func (s *client) FilterNamespace(filter *util.ContentFilter) (ret []*common.NamespaceView, total int64, err error) {
	result := &common.NamespaceStatisticResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "2")
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterNamespace}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
	total = result.Total
	return
}

func (s *client) FilterNamespaceLite(filter *util.ContentFilter) (ret []*common.NamespaceLite, total int64, err error) {
	result := &common.NamespaceLiteListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	vals.Set("mode", "1")
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterNamespace}, "")
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
	total = result.Total
	return
}

func (s *client) QueryNamespace(id int) (ret *common.NamespaceView, err error) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("查询指定租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
	return
}

func (s *client) CreateNamespace(ptr *common.NamespaceParam) (ret *common.NamespaceView, err error) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateNamespace}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPPost(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("新建租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
	return
}

func (s *client) UpdateNamespace(id int, ptr *common.NamespaceParam) (ret *common.NamespaceView, err error) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("更新租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
	return
}

func (s *client) DeleteNamespace(id int) (ret *common.NamespaceView, err error) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err = net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("删除租户失败，%s", result.Reason)
		return
	}

	ret = result.Namespace
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
	if s.sessionInfo != nil {
		ret = s.sessionInfo.Encode(ret)
	}
	if s.namespace != "" {
		ret.Set(common.NamespaceID, s.namespace)
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
