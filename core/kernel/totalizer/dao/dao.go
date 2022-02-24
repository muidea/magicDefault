package dao

import (
	"fmt"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	"github.com/muidea/magicDefault/model"
)

type Totalizer interface {
	FilterTotalizer(filter *bc.QueryFilter, namespace string) (ret []*model.Totalizer, total int64, err error)
	QueryTotalizer(name string, typeVal, catalog int, namespace string) (ret *model.Totalizer, err error)
	CreateTotalizer(totalizer *model.Totalizer, namespace string) (ret *model.Totalizer, err error)
	DestroyTotalizer(id int, namespace string) (ret *model.Totalizer, err error)
	UpdateTotalizer(totalizer *model.Totalizer, namespace string) (ret *model.Totalizer, err error)
}

func New(clnt client.Client) Totalizer {
	return &totalizer{batisClient: clnt}
}

type totalizer struct {
	batisClient client.Client
}

func (s *totalizer) FilterTotalizer(filter *bc.QueryFilter, namespace string) (ret []*model.Totalizer, total int64, err error) {
	ret = []*model.Totalizer{}
	if filter == nil {
		filter = bc.NewFilter()
	}

	filter.Equal("Namespace", namespace)
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)

	return
}

func (s *totalizer) QueryTotalizer(name string, typeVal, catalog int, namespace string) (ret *model.Totalizer, err error) {
	val := &model.Totalizer{Owner: name, Type: typeVal, Catalog: catalog, Namespace: namespace}
	err = s.batisClient.QueryEntity(val)
	if err == nil {
		ret = val
	}

	return
}

func (s *totalizer) CreateTotalizer(totalizer *model.Totalizer, namespace string) (ret *model.Totalizer, err error) {
	if totalizer == nil {
		err = fmt.Errorf("illegal totalizer")
		return
	}

	totalizer.Namespace = namespace
	err = s.batisClient.InsertEntity(totalizer)
	if err == nil {
		ret = totalizer
	}

	return
}

func (s *totalizer) queryByID(id int, namespace string) (ret *model.Totalizer, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal totalizer id, id:%d", id)
		return
	}

	val := &model.Totalizer{ID: id, Namespace: namespace}
	err = s.batisClient.QueryEntity(val)
	if err == nil {
		ret = val
	}

	return
}

func (s *totalizer) DestroyTotalizer(id int, namespace string) (ret *model.Totalizer, err error) {
	ret, err = s.queryByID(id, namespace)
	if err == nil {
		err = s.batisClient.DeleteEntity(ret)
	}

	return
}

func (s *totalizer) UpdateTotalizer(totalizer *model.Totalizer, namespace string) (ret *model.Totalizer, err error) {
	if totalizer == nil || totalizer.ID <= 0 {
		err = fmt.Errorf("illegal totalizer, val:%v", totalizer)
		return
	}

	totalizer.Namespace = namespace
	err = s.batisClient.UpdateEntity(totalizer)
	if err == nil {
		ret = totalizer
	}

	return
}
