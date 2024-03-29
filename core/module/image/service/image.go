package service

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	fc "github.com/muidea/magicFile/common"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/module/image/biz"
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

	s.routeRegistry.SetApiVersion(common.ApiVersion)
	s.casRouteRegistry.SetApiVersion(common.ApiVersion)
	s.roleRouteRegistry.SetApiVersion(common.ApiVersion)
}

// RegisterRoute 注册路由
func (s *Image) RegisterRoute() {
	s.roleRouteRegistry.AddHandler(common.FilterImage, "GET", cc.ReadPermission, s.FilterImage)
	s.roleRouteRegistry.AddHandler(common.QueryImage, "GET", cc.ReadPermission, s.QueryImage)
	s.roleRouteRegistry.AddHandler(common.UpdateImage, "PUT", cc.WritePermission, s.UpdateImage)
	s.roleRouteRegistry.AddHandler(common.DeleteImage, "DELETE", cc.DeletePermission, s.DeleteImage)
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
