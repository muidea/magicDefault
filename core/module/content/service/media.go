package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	bc "github.com/muidea/magicBatis/common"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicDefault/common"
)

func (s *Content) filterMediaLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.MediaLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaList, mediaTotal, mediaErr := s.bizPtr.FilterMedia(filter, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		for _, val := range mediaList {
			ptr := &common.MediaLite{}
			ptr.FromMedia(val)
			result.Media = append(result.Media, ptr)
		}

		result.Total = mediaTotal
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

func (s *Content) filterMedia(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.MediaListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaList, mediaTotal, mediaErr := s.bizPtr.FilterMedia(filter, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		for _, val := range mediaList {
			ptr := &common.MediaView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.FromMedia(val, entityPtr)
			result.Media = append(result.Media, ptr)
		}

		result.Total = mediaTotal
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

func (s *Content) FilterMedia(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)
	queryFilter := bc.NewFilter()
	if filter.Pagination != nil {
		queryFilter.Page(filter.Pagination)
	} else {
		queryFilter.Page(fu.NewPagination(20, 1))
	}
	queryFilter.Sort(fu.NewSortFilter("UpdateTime", false))

	modelVal, modelOK := filter.Get("mode")
	if modelOK {
		if modelVal == common.LiteMode {
			s.filterMediaLite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == common.ViewMode {
			s.filterMedia(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.MediaStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.NotifyContentMedia, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		mediaList, mediaTotal, mediaErr := s.bizPtr.FilterMedia(queryFilter, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		for _, val := range mediaList {
			view := &common.MediaView{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.FromMedia(val, entityPtr)
			result.Media = append(result.Media, view)
		}
		result.Total = mediaTotal
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

func (s *Content) QueryMedia(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.MediaResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaPtr, mediaErr := s.bizPtr.QueryMedia(id, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), mediaPtr.Creater, curNamespace)
		result.Media = &common.MediaView{}
		result.Media.FromMedia(mediaPtr, entityPtr)
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

func (s *Content) CreateMedia(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.MediaResult{}
	for {
		param := &common.MediaParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}
		curEntity, curErr := s.getCurrentEntity(ctx)
		if curErr != nil {
			result.ErrorCode = cd.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaPtr, mediaErr := s.bizPtr.CreateMedia(param, curEntity.ID, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		memo := fmt.Sprintf("新增Media:%d", mediaPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Media = &common.MediaView{}
		result.Media.FromMedia(mediaPtr, curEntity)
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

func (s *Content) UpdateMedia(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.MediaResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.MediaParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curEntity, curErr := s.getCurrentEntity(ctx)
		if curErr != nil {
			result.ErrorCode = cd.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaPtr, mediaErr := s.bizPtr.UpdateMedia(id, param, curEntity.ID, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Media:%d", mediaPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Media = &common.MediaView{}
		result.Media.FromMedia(mediaPtr, curEntity)
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

func (s *Content) DeleteMedia(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.MediaResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		mediaPtr, mediaErr := s.bizPtr.DeleteMedia(id, curNamespace)
		if mediaErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = mediaErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Media:%d", mediaPtr.ID)
		s.writeLog(ctx, res, req, memo)

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), mediaPtr.Creater, curNamespace)
		result.Media = &common.MediaView{}
		result.Media.FromMedia(mediaPtr, entityPtr)
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
