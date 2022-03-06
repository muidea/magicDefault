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

func (s *Content) filterLinkLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.LinkLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		linkList, linkTotal, linkErr := s.bizPtr.FilterLink(filter, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		for _, val := range linkList {
			ptr := &common.LinkLite{}
			ptr.FromLink(val)
			result.Link = append(result.Link, ptr)
		}

		result.Total = linkTotal
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

func (s *Content) filterLink(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.LinkListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		linkList, linkTotal, linkErr := s.bizPtr.FilterLink(filter, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		for _, val := range linkList {
			ptr := &common.LinkView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.FromLink(val, entityPtr)
			result.Link = append(result.Link, ptr)
		}

		result.Total = linkTotal
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

func (s *Content) FilterLink(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)
	queryFilter := bc.NewFilter()
	if filter.Pagination != nil {
		queryFilter.Page(filter.Pagination)
	} else {
		queryFilter.Page(fu.NewPagination(20, 1))
	}

	modelVal, modelOK := filter.Get("mode")
	if modelOK {
		if modelVal == common.LiteMode {
			s.filterLinkLite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == common.ViewMode {
			s.filterLink(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.LinkStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.NotifyContentLink, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		linkList, linkTotal, linkErr := s.bizPtr.FilterLink(queryFilter, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		for _, val := range linkList {
			view := &common.LinkView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.FromLink(val, entityPtr)
			result.Link = append(result.Link, view)
		}
		result.Total = linkTotal
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

func (s *Content) QueryLink(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.LinkResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		linkPtr, linkErr := s.bizPtr.QueryLink(id, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		result.Link = &common.LinkView{}
		entityPtr := s.queryEntity(curSession.GetSessionInfo(), linkPtr.Creater, curNamespace)
		result.Link.FromLink(linkPtr, entityPtr)
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

func (s *Content) CreateLink(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.LinkResult{}
	for {
		param := &common.LinkParam{}
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
		linkPtr, linkErr := s.bizPtr.CreateLink(param, curEntity.ID, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		memo := fmt.Sprintf("新增Link:%d", linkPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Link = &common.LinkView{}
		result.Link.FromLink(linkPtr, curEntity)
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

func (s *Content) UpdateLink(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.LinkResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.LinkParam{}
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
		linkPtr, linkErr := s.bizPtr.UpdateLink(id, param, curEntity.ID, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Link:%d", linkPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Link = &common.LinkView{}
		result.Link.FromLink(linkPtr, curEntity)
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

func (s *Content) DeleteLink(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.LinkResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		linkPtr, linkErr := s.bizPtr.DeleteLink(id, curNamespace)
		if linkErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = linkErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Link:%d", linkPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Link = &common.LinkView{}
		entityPtr := s.queryEntity(curSession.GetSessionInfo(), linkPtr.Creater, curNamespace)
		result.Link.FromLink(linkPtr, entityPtr)
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
