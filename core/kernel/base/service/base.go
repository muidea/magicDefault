package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/cihub/seelog"

	bc "github.com/muidea/magicBatis/common"
	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	commonSession "github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	casCommon "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	fileCommon "github.com/muidea/magicFile/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
	"github.com/muidea/magicDefault/core/kernel/base/biz"
	"github.com/muidea/magicDefault/model"
)

// Base BaseService
type Base struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

	bizPtr    *biz.Base
	systemBiz biz.System

	endpointName string
	fileService  string
}

// New create base
func New(endpointName, fileService string, bizPtr *biz.Base) *Base {
	ptr := &Base{
		endpointName: endpointName,
		fileService:  fileService,
		bizPtr:       bizPtr,
		validator:    fu.NewFormValidator(),
	}

	return ptr
}

func (s *Base) BindRegistry(
	routeRegistry toolkit.RouteRegistry,
	casRouteRegistry toolkit.CasRegistry,
	roleRouteRegistry toolkit.RoleRegistry) {

	s.routeRegistry = routeRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// LoginAccount login account
func (s *Base) LoginAccount(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)
	result := &casCommon.LoginResult{}
	for {
		param := &casCommon.LoginParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Account == "" || param.Password == "" {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		loginEntity, loginSession, loginErr := s.bizPtr.LoginAccount(curSession.GetSessionInfo(), param.Account, param.Password, curNamespace)
		if loginErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = loginErr.Error()
			break
		}

		entityRole, entityErr := s.bizPtr.VerifyEntityRole(loginSession, loginEntity, curNamespace)
		if entityErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = entityErr.Error()
			break
		}

		curSession.SetSessionInfo(loginSession)
		curSession.SetOption(commonSession.AuthAccount, loginEntity)
		curSession.SetOption(commonSession.AuthRole, entityRole)
		curSession.SetOption(commonSession.AuthNamespace, curNamespace)
		curSession.SetOption(commonSession.AuthRemoteAddress, fn.GetHTTPRemoteAddress(req))
		curSession.Flush(res, req)

		result.Entity = loginEntity
		result.SessionInfo = loginSession
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

// LogoutAccount logout account
func (s *Base) LogoutAccount(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	result := &casCommon.LogoutResult{}
	for {
		namespace := s.getCurrentNamespace(ctx, res, req)
		logoutSession, logoutErr := s.bizPtr.LogoutAccount(curSession.GetSessionInfo(), namespace)
		if logoutErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = logoutErr.Error()
			break
		}

		curSession.SetSessionInfo(logoutSession)
		curSession.RemoveOption(commonSession.AuthAccount)
		curSession.RemoveOption(commonSession.AuthRole)
		curSession.RemoveOption(commonSession.AuthRemoteAddress)
		curSession.Flush(res, req)

		result.SessionInfo = logoutSession
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

func (s *Base) UpdateAccountPassword(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	result := &casCommon.AccountResult{}
	for {
		param := &casCommon.UpdatePasswordParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		namespace := s.getCurrentNamespace(ctx, res, req)
		accountPtr, accountErr := s.bizPtr.UpdateAccountPassword(curSession.GetSessionInfo(), param, namespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		result.Account = accountPtr
		result.ErrorCode = commonDef.Success
		break
	}
	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

// VerifyEndpoint verify endpoint
func (s *Base) VerifyEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	result := &casCommon.VerifyEndpointResult{}
	for {
		ptr := &casCommon.VerifyEndpointParam{}
		err := fn.ParseJSONBody(req, s.validator, ptr)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		verifyEntity, verifySession, verifyErr := s.bizPtr.VerifyEndpoint(curSession.GetSessionInfo(), ptr.Endpoint, ptr.IdentifyID, ptr.AuthToken, curNamespace)
		if verifyErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = verifyErr.Error()
			break
		}

		entityRole, entityErr := s.bizPtr.VerifyEntityRole(verifySession, verifyEntity, curNamespace)
		if entityErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = entityErr.Error()
			break
		}

		curSession.SetSessionInfo(verifySession)
		curSession.SetOption(commonSession.AuthAccount, verifyEntity)
		curSession.SetOption(commonSession.AuthRole, entityRole)
		curSession.SetOption(commonSession.AuthNamespace, curNamespace)
		curSession.SetOption(commonSession.AuthRemoteAddress, fn.GetHTTPRemoteAddress(req))
		curSession.Flush(res, req)

		result.Entity = verifyEntity
		result.SessionInfo = verifySession
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

// RefreshSession refresh commonSession status
func (s *Base) RefreshSession(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	result := &casCommon.RefreshResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		refreshEntity, refreshSession, refreshErr := s.bizPtr.RefreshSession(curSession.GetSessionInfo(), curNamespace)
		if refreshErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = refreshErr.Error()
			break
		}

		entityRole, entityErr := s.bizPtr.VerifyEntityRole(refreshSession, refreshEntity, curNamespace)
		if entityErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = entityErr.Error()
			break
		}

		curSession.SetSessionInfo(refreshSession)
		curSession.SetOption(commonSession.AuthAccount, refreshEntity)
		curSession.SetOption(commonSession.AuthRole, entityRole)
		curSession.SetOption(commonSession.AuthNamespace, curNamespace)
		curSession.SetOption(commonSession.AuthRemoteAddress, fn.GetHTTPRemoteAddress(req))
		curSession.Flush(res, req)

		result.SessionInfo = refreshSession
		result.Entity = refreshEntity
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

// QueryAccessLog query access log
func (s *Base) QueryAccessLog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	filter := fu.NewPagination()
	filter.Decode(req)

	result := &casCommon.AccessLogListResult{}
	for {
		curEntity, curErr := s.getCurrentEntity(ctx, res, req)
		if curErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "查询访问日志失败"
			break
		}

		namespace := s.getCurrentNamespace(ctx, res, req)
		logList, logCount, logErr := s.bizPtr.QueryAccessLog(curSession.GetSessionInfo(), curEntity, filter, namespace)
		if logErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = logErr.Error()
			break
		}

		result.AccessLog = logList
		result.Total = logCount
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
	return
}

// QueryOperateLog query operate log
func (s *Base) QueryOperateLog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)

	queryFilter := bc.NewFilter()
	queryFilter.Page(filter.Pagination)
	result := &common.OperateLogListResult{}
	for {
		namespace := s.getCurrentNamespace(ctx, res, req)
		curEntity, curErr := s.getCurrentEntity(ctx, res, req)
		if curErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = curErr.Error()
			break
		}

		queryFilter.Equal("Creater", curEntity.ID)
		logList, logCount, logErr := s.bizPtr.QueryOperateLog(queryFilter, namespace)
		if logErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = logErr.Error()
			break
		}

		for _, val := range logList {
			view := &common.LogView{}
			view.FromLog(val, curEntity)
			result.OperateLog = append(result.OperateLog, view)
		}

		result.Total = logCount
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

	return
}

func (s *Base) EnumPrivate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	namespace := s.getCurrentNamespace(ctx, res, req)
	result := &common.EnumPrivateItemResult{Private: []*casCommon.PrivateItem{}}
	for {
		items := s.roleRouteRegistry.GetAllPrivateItem()
		for _, val := range items {
			// if special route must be super namespace
			if s.systemBiz.IsSpecialRoute(val.Path) && !s.isSuperNamespace(namespace) {
				continue
			}

			// if super namespace, only setting route
			if s.isSuperNamespace(namespace) && !s.systemBiz.IsSettingRoute(val.Path) {
				continue
			}

			result.Private = append(result.Private, val)
		}
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusInternalServerError)
}

// QueryBaseInfo get system info
func (s *Base) QueryBaseInfo(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)

	type getResult struct {
		commonDef.Result
		Route   []*biz.Route   `json:"route"`
		Content []*biz.Content `json:"content"`
	}

	result := &getResult{}
	for {
		curEntity, curErr := s.getCurrentEntity(ctx, res, req)
		if curErr != nil {
			result.ErrorCode = commonDef.InvalidAuthority
			result.Reason = "无效账号"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		entityRole, entityErr := s.bizPtr.VerifyEntityRole(curSession.GetSessionInfo(), curEntity, curNamespace)
		if entityErr != nil {
			result.ErrorCode = commonDef.InvalidAuthority
			result.Reason = "无效权限"
			break
		}

		isSuper := s.isSuperNamespace(curNamespace)
		result.Route, result.Content = s.systemBiz.SystemInfo(entityRole, isSuper)
		break
	}
	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// Handle middleware handler
func (s *Base) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Context().Value(commonSession.AuthSession).(commonSession.Session)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonSession.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)

	urlPath := req.URL.Path
	if len(urlPath) > len(common.ApiVersion) {
		urlPath = urlPath[len(common.ApiVersion):]
	}

	curNamespace := s.getCurrentNamespace(ctx.Context(), res, req)
	switch urlPath {
	case common.UploadFile, common.ViewFile:
		req.Header.Set("source", fmt.Sprintf("%s_%s", s.endpointName, curNamespace))
	}

	req.URL.RawQuery = values.Encode()

	ctx.Next()
}

// RegisterRoute 注册路由
func (s *Base) RegisterRoute() {
	// account login,logout, confirm, refresh
	//---------------------------------------------------------------------------------------
	loginRoute := engine.CreateRoute(common.LoginAccount, "POST", s.LoginAccount)
	s.routeRegistry.AddRoute(loginRoute, s)

	logoutRoute := engine.CreateRoute(common.LogoutAccount, "DELETE", s.LogoutAccount)
	s.casRouteRegistry.AddRoute(logoutRoute, s)

	changePasswordRoute := engine.CreateRoute(common.UpdateAccountPassword, "PUT", s.UpdateAccountPassword)
	s.casRouteRegistry.AddRoute(changePasswordRoute, s)

	verifyEndpointRoute := engine.CreateRoute(common.VerifyEndpoint, "POST", s.VerifyEndpoint)
	s.routeRegistry.AddRoute(verifyEndpointRoute, s)

	s.casRouteRegistry.AddHandler(common.RefreshSession, "GET", s.RefreshSession)

	s.casRouteRegistry.AddHandler(common.QueryAccessLog, "GET", s.QueryAccessLog)

	s.casRouteRegistry.AddHandler(common.QueryOperateLog, "GET", s.QueryOperateLog)

	// upload file
	//---------------------------------------------------------------------------------------
	uploadFileURL := fn.JoinSuffix(s.fileService, strings.Join([]string{fileCommon.ApiVersion, fileCommon.UploadFileURL}, ""))
	uploadFileRoute := engine.CreateProxyRoute(common.UploadFile, "POST", uploadFileURL, true)
	s.routeRegistry.AddRoute(uploadFileRoute, s)

	viewFileURL := fn.JoinSuffix(s.fileService, strings.Join([]string{fileCommon.ApiVersion, fileCommon.DownloadFileURL}, ""))
	viewFileRoute := engine.CreateProxyRoute(common.ViewFile, "GET", viewFileURL, true)
	s.routeRegistry.AddRoute(viewFileRoute, s)

	s.casRouteRegistry.AddHandler(common.EnumPrivate, "GET", s.EnumPrivate)
	s.casRouteRegistry.AddHandler(common.QueryBaseInfo, "GET", s.QueryBaseInfo)
}

func (s *Base) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(ctx, res, req)
	curNamespace := s.getCurrentNamespace(ctx, res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	_, logErr := s.bizPtr.WriteOperateLog(logPtr, curNamespace)
	if logErr != nil {
		log.Errorf("Write log failed, err:%s", logErr.Error())
	}
}

func (s *Base) getCurrentEntity(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := ctx.Value(commonSession.AuthSession).(commonSession.Session)
	authVal, ok := curSession.GetOption(commonSession.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *Base) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(casCommon.NamespaceID)
	if namespace != "" {
		ret = namespace
		return
	}

	items := strings.Split(req.Host, ":")
	if nil != net.ParseIP(items[0]) {
		ret = "default"
		return
	}

	items = strings.Split(items[0], ".")
	if len(items) <= 1 {
		ret = "default"
		return
	}

	ret = items[0]
	return
}

func (s *Base) isSuperNamespace(namespace string) bool {
	return namespace == config.SuperNamespace()
}
