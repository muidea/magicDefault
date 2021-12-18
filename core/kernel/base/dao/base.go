package dao

import (
	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"

	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicDefault/model"
)

type Base interface {
	QueryLog(filter *bc.QueryFilter, namespace string) (ret []*model.Log, total int64, err error)
	WriteLog(log *model.Log, namespace string) (ret *model.Log, err error)
}

// New 新建Base
func New(clnt client.Client) Base {
	return &base{batisClient: clnt}
}

type base struct {
	batisClient client.Client
}

func (s *base) QueryLog(filter *bc.QueryFilter, namespace string) (ret []*model.Log, total int64, err error) {
	ret = []*model.Log{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.SortFilter = &fu.SortFilter{AscSort: false, FieldName: "CreateTime"}
	filter.Equal("Namespace", namespace)

	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *base) WriteLog(log *model.Log, namespace string) (ret *model.Log, err error) {
	log.Namespace = namespace
	err = s.batisClient.InsertEntity(log)
	if err == nil {
		ret = log
	}
	return
}
