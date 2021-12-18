package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	commonDef "github.com/muidea/magicCommon/def"
	fu "github.com/muidea/magicCommon/foundation/util"
	commonSession "github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	casCommon "github.com/muidea/magicCas/common"
	casToolkit "github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/kernel/setting/biz"
)

type Setting struct {
	sessionRegistry   commonSession.Registry
	casRouteRegistry  casToolkit.CasRegistry
	roleRouteRegistry casToolkit.RoleRegistry
	validator         fu.Validator

	settingBiz *biz.Setting
}

func New(
	settingBiz *biz.Setting,
) *Setting {
	ptr := &Setting{
		settingBiz: settingBiz,
		validator:  fu.NewFormValidator(),
	}

	return ptr
}

// Handle middleware handler
func (s *Setting) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonSession.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()
}

// RegisterRoute 注册路由
func (s *Setting) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.ViewSetting, "GET", casCommon.ReadPrivate, s.ViewSetting)

	s.roleRouteRegistry.AddHandler(common.ViewSettingProfile, "GET", casCommon.ReadPrivate, s.ViewSettingProfile)
}

func (s *Setting) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonSession.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *Setting) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *Setting) ViewSetting(res http.ResponseWriter, req *http.Request) {
	result := &common.SettingResult{Item: []*common.Content{}}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		curEntity, curErr := s.getCurrentEntity(res, req)
		if curErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "查询用户失败"
			break
		}

		itemList, itemErr := s.settingBiz.QuerySetting(curSession.GetSessionInfo(), curNamespace)
		if itemErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "查询数据失败"
			break
		}

		if curEntity.IsAdmin() {
			onlineItem := &common.Content{Name: "在线人数", Content: fmt.Sprintf("%d", s.sessionRegistry.CountSession(&namespaceFilter{namespace: curNamespace}))}
			itemList = append(itemList, onlineItem)
		}

		result.Item = itemList
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

type namespaceFilter struct {
	namespace string
}

func (s *namespaceFilter) Enable(val interface{}) bool {
	session := val.(commonSession.Session)
	if session == nil {
		return false
	}
	nVal, nOK := session.GetOption(commonSession.AuthNamespace)
	if !nOK {
		return false
	}
	return nVal.(string) == s.namespace
}
