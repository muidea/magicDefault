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

func (s *Content) filterCommentLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.CommentLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		commentList, commentTotal, commentErr := s.bizPtr.FilterComment(filter, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		for _, val := range commentList {
			ptr := &common.CommentLite{}
			ptr.FromComment(val)
			result.Comment = append(result.Comment, ptr)
		}

		result.Total = commentTotal
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

func (s *Content) filterComment(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CommentListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		commentList, commentTotal, commentErr := s.bizPtr.FilterComment(filter, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		for _, val := range commentList {
			ptr := &common.CommentView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.FromComment(val, entityPtr)
			result.Comment = append(result.Comment, ptr)
		}

		result.Total = commentTotal
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

func (s *Content) FilterComment(ctx context.Context, res http.ResponseWriter, req *http.Request) {
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
			s.filterCommentLite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == common.ViewMode {
			s.filterComment(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CommentStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.NotifyContentComment, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		commentList, commentTotal, commentErr := s.bizPtr.FilterComment(queryFilter, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		for _, val := range commentList {
			view := &common.CommentView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.FromComment(val, entityPtr)
			result.Comment = append(result.Comment, view)
		}
		result.Total = commentTotal
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

func (s *Content) QueryComment(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CommentResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		commentPtr, commentErr := s.bizPtr.QueryComment(id, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		result.Comment = &common.CommentView{}
		entityPtr := s.queryEntity(curSession.GetSessionInfo(), commentPtr.Creater, curNamespace)
		result.Comment.FromComment(commentPtr, entityPtr)
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

func (s *Content) CreateComment(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.CommentResult{}
	for {
		param := &common.CommentParam{}
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
		commentPtr, commentErr := s.bizPtr.CreateComment(param, curEntity.ID, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		memo := fmt.Sprintf("新增Comment:%d", commentPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Comment = &common.CommentView{}
		result.Comment.FromComment(commentPtr, curEntity)
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

func (s *Content) UpdateComment(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.CommentResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.CommentParam{}
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
		commentPtr, commentErr := s.bizPtr.UpdateComment(id, param, curEntity.ID, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Comment:%d", commentPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Comment = &common.CommentView{}
		result.Comment.FromComment(commentPtr, curEntity)
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

func (s *Content) DeleteComment(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CommentResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		commentPtr, commentErr := s.bizPtr.DeleteComment(id, curNamespace)
		if commentErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = commentErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Comment:%d", commentPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Comment = &common.CommentView{}
		entityPtr := s.queryEntity(curSession.GetSessionInfo(), commentPtr.Creater, curNamespace)
		result.Comment.FromComment(commentPtr, entityPtr)
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
