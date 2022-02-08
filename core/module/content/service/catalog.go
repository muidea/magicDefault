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

func (s *Content) filterCatalogLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.CatalogLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		catalogList, catalogTotal, catalogErr := s.bizPtr.FilterCatalog(filter, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		for _, val := range catalogList {
			ptr := &common.CatalogLite{}
			ptr.FromCatalog(val)
			result.Catalog = append(result.Catalog, ptr)
		}

		result.Total = catalogTotal
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

func (s *Content) filterCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CatalogListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		catalogList, catalogTotal, catalogErr := s.bizPtr.FilterCatalog(filter, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		for _, val := range catalogList {
			ptr := &common.CatalogView{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.FromCatalog(val, entityPtr)
			result.Catalog = append(result.Catalog, ptr)
		}

		result.Total = catalogTotal
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

func (s *Content) FilterCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
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
		if modelVal == "1" {
			s.filterCatalogLite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == "2" {
			s.filterCatalog(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CatalogStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.NotifyContentCatalog, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		catalogList, catalogTotal, catalogErr := s.bizPtr.FilterCatalog(queryFilter, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		for _, val := range catalogList {
			view := &common.CatalogView{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.FromCatalog(val, entityPtr)
			result.Catalog = append(result.Catalog, view)
		}
		result.Total = catalogTotal
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

func (s *Content) QueryCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CatalogResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		catalogPtr, catalogErr := s.bizPtr.QueryCatalog(id, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), catalogPtr.Creater, curNamespace)
		result.Catalog = &common.CatalogView{}
		result.Catalog.FromCatalog(catalogPtr, entityPtr)
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

func (s *Content) CreateCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.CatalogResult{}
	for {
		param := &common.CatalogParam{}
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
		catalogPtr, catalogErr := s.bizPtr.CreateCatalog(param, curEntity.ID, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		memo := fmt.Sprintf("新增Catalog:%d", catalogPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Catalog = &common.CatalogView{}
		result.Catalog.FromCatalog(catalogPtr, curEntity)
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

func (s *Content) UpdateCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.CatalogResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.CatalogParam{}
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
		catalogPtr, catalogErr := s.bizPtr.UpdateCatalog(id, param, curEntity.ID, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Catalog:%d", catalogPtr.ID)
		s.writeLog(ctx, res, req, memo)

		result.Catalog = &common.CatalogView{}
		result.Catalog.FromCatalog(catalogPtr, curEntity)
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

func (s *Content) DeleteCatalog(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.CatalogResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		catalogPtr, catalogErr := s.bizPtr.DeleteCatalog(id, curNamespace)
		if catalogErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = catalogErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Catalog:%d", catalogPtr.ID)
		s.writeLog(ctx, res, req, memo)

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), catalogPtr.Creater, curNamespace)
		result.Catalog = &common.CatalogView{}
		result.Catalog.FromCatalog(catalogPtr, entityPtr)
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
