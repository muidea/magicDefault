package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	fc "github.com/muidea/magicFile/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/image/biz"
	"github.com/muidea/magicDefault/model"
)

type Image struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

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
	routeRegistry toolkit.RouteRegistry,
	casRouteRegistry toolkit.CasRegistry,
	roleRouteRegistry toolkit.RoleRegistry) {

	s.routeRegistry = routeRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry
}

// RegisterRoute 注册路由
func (s *Image) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterImage, "GET", cc.ReadPrivate, s.FilterImage)
	s.roleRouteRegistry.AddHandler(common.QueryImage, "GET", cc.ReadPrivate, s.QueryImage)
	s.roleRouteRegistry.AddHandler(common.UpdateImage, "PUT", cc.WritePrivate, s.UpdateImage)
	s.roleRouteRegistry.AddHandler(common.DeleteImage, "DELETE", cc.DeletePrivate, s.DeleteImage)
}

func (s *Image) FilterImage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)

	result := &common.ImageStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		imageList, imageTotal, imageErr := s.bizPtr.FilterImage(filter, curNamespace)
		if imageErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = imageErr.Error()
			break
		}
		for _, val := range imageList {
			view := &fc.FileView{}
			view.FromFileDetail(val)
			result.Image = append(result.Image, view)
		}

		result.Total = imageTotal
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) QueryImage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		imagePtr, imageErr := s.bizPtr.QueryImage(id, curNamespace)
		if imageErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fc.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) UpdateImage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &fc.FileParam{}
		err = fn.ParseJSONBody(req, nil, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		imagePtr, imageErr := s.bizPtr.UpdateImage(id, param, curNamespace)
		if imageErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fc.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) DeleteImage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.ImageResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		imagePtr, imageErr := s.bizPtr.DeleteImage(id, curNamespace)
		if imageErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = imageErr.Error()
			break
		}

		result.Image = &fc.FileView{}
		result.Image.FromFileDetail(imagePtr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Image) getCurrentEntity(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *Image) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(cc.NamespaceID)
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

func (s *Image) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *Image) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	curEntity, _ := s.getCurrentEntity(ctx, res, req)
	curNamespace := s.getCurrentNamespace(ctx, res, req)
	logPtr := &model.Log{Address: address, Memo: memo, Creater: curEntity.ID, CreateTime: time.Now().UTC().Unix()}
	s.bizPtr.WriteLog(logPtr, curNamespace)
}
