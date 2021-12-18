package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"

	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	casCommon "github.com/muidea/magicCas/common"
	casToolkit "github.com/muidea/magicCas/toolkit"

	fileCommon "github.com/muidea/magicFile/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/image/biz"
	"github.com/muidea/magicDefault/model"
)

type Image struct {
	sessionRegistry   commonSession.Registry
	casRouteRegistry  casToolkit.CasRegistry
	roleRouteRegistry casToolkit.RoleRegistry
	validator         fu.Validator

	bizPtr *biz.Image
}

func New(
	fileBiz *biz.Image,
) *Image {
	ptr := &Image{
		bizPtr:    fileBiz,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *Image) BindRegistry(
	sessionRegistry commonSession.Registry,
	casRouteRegistry casToolkit.CasRegistry,
	roleRouteRegistry casToolkit.RoleRegistry) {

	s.sessionRegistry = sessionRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// Handle middleware handler
func (s *Image) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonSession.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()
}

// RegisterRoute 注册路由
func (s *Image) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterImage, "GET", casCommon.ReadPrivate, s.FilterImage)
	s.roleRouteRegistry.AddHandler(common.QueryImage, "GET", casCommon.ReadPrivate, s.QueryImage)
	s.roleRouteRegistry.AddHandler(common.UpdateImage, "PUT", casCommon.WritePrivate, s.UpdateImage)
	s.roleRouteRegistry.AddHandler(common.DeleteImage, "DELETE", casCommon.DeletePrivate, s.DeleteImage)
}

func (s *Image) FilterImage(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	result := &common.ImageStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(res, req)
		imageList, imageTotal, imageErr := s.bizPtr.FilterImage(filter, curNamespace)
		if imageErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = imageErr.Error()
			break
		}
		for _, val := range imageList {
			view := &fileCommon.FileView{}
			view.FromFileDetail(val)
			result.Image = append(result.Image, view)
		}

		result.Total = imageTotal
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) QueryImage(res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(res, req)
		imagePtr, imageErr := s.bizPtr.QueryImage(id, curNamespace)
		if imageErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fileCommon.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) UpdateImage(res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		param := &fileCommon.FileParam{}
		err = fn.ParseJSONBody(req, nil, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(res, req)
		imagePtr, imageErr := s.bizPtr.UpdateImage(id, param, curNamespace)
		if imageErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fileCommon.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) DeleteImage(res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(res, req)
		imagePtr, imageErr := s.bizPtr.DeleteImage(id, curNamespace)
		if imageErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fileCommon.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casCommon.EntityView, err error) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonSession.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casCommon.EntityView)
	return
}

func (s *Image) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *Image) queryEntity(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *Image) writeLog(res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(res, req)
	curNamespace := s.getCurrentNamespace(res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
