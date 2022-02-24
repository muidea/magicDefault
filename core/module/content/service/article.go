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

func (s *Content) filterArticleLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.ArticleLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		articleList, articleTotal, articleErr := s.bizPtr.FilterArticle(filter, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		for _, val := range articleList {
			ptr := &common.ArticleLite{}
			ptr.FromArticle(val)
			result.Article = append(result.Article, ptr)
		}

		result.Total = articleTotal
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

func (s *Content) filterArticle(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.ArticleListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		articleList, articleTotal, articleErr := s.bizPtr.FilterArticle(filter, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		for _, val := range articleList {
			ptr := &common.ArticleView{}
			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.FromArticle(val, entityPtr)
			result.Article = append(result.Article, ptr)
		}

		result.Total = articleTotal
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

func (s *Content) FilterArticle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
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
			s.filterArticleLite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == common.ViewMode {
			s.filterArticle(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.ArticleStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.NotifyContentArticle, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		articleList, articleTotal, articleErr := s.bizPtr.FilterArticle(queryFilter, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		for _, val := range articleList {
			view := &common.ArticleView{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.FromArticle(val, entityPtr)
			result.Article = append(result.Article, view)
		}
		result.Total = articleTotal
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

func (s *Content) QueryArticle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.ArticleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		articlePtr, articleErr := s.bizPtr.QueryArticle(id, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), articlePtr.Creater, curNamespace)
		result.Article = &common.ArticleView{}
		result.Article.FromArticle(articlePtr, entityPtr)
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

func (s *Content) CreateArticle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.ArticleResult{}
	for {
		param := &common.ArticleParam{}
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
		articlePtr, articleErr := s.bizPtr.CreateArticle(param, curEntity.ID, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		memo := fmt.Sprintf("新增Article:%d", articlePtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Article = &common.ArticleView{}
		result.Article.FromArticle(articlePtr, curEntity)
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

func (s *Content) UpdateArticle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.ArticleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.ArticleParam{}
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
		articlePtr, articleErr := s.bizPtr.UpdateArticle(id, param, curEntity.ID, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Article:%d", articlePtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Article = &common.ArticleView{}
		result.Article.FromArticle(articlePtr, curEntity)
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

func (s *Content) DeleteArticle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.ArticleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		articlePtr, articleErr := s.bizPtr.DeleteArticle(id, curNamespace)
		if articleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = articleErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Article:%d", articlePtr.ID)
		s.writeLog(ctx, res, req, memo)

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), articlePtr.Creater, curNamespace)
		result.Article = &common.ArticleView{}
		result.Article.FromArticle(articlePtr, entityPtr)
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
