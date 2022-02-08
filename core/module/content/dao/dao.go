package dao

import (
	"fmt"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	"github.com/muidea/magicDefault/model"
)

type Content interface {
	FilterArticle(filter *bc.QueryFilter, namespace string) (ret []*model.Article, total int64, err error)
	QueryArticle(id int, namespace string) (ret *model.Article, err error)
	CreateArticle(ptr *model.Article, namespace string) (ret *model.Article, err error)
	UpdateArticle(ptr *model.Article, namespace string) (ret *model.Article, err error)
	DeleteArticle(id int, namespace string) (ret *model.Article, err error)

	FilterCatalog(filter *bc.QueryFilter, namespace string) (ret []*model.Catalog, total int64, err error)
	QueryCatalog(id int, namespace string) (ret *model.Catalog, err error)
	CreateCatalog(ptr *model.Catalog, namespace string) (ret *model.Catalog, err error)
	UpdateCatalog(ptr *model.Catalog, namespace string) (ret *model.Catalog, err error)
	DeleteCatalog(id int, namespace string) (ret *model.Catalog, err error)

	FilterLink(filter *bc.QueryFilter, namespace string) (ret []*model.Link, total int64, err error)
	QueryLink(id int, namespace string) (ret *model.Link, err error)
	CreateLink(ptr *model.Link, namespace string) (ret *model.Link, err error)
	UpdateLink(ptr *model.Link, namespace string) (ret *model.Link, err error)
	DeleteLink(id int, namespace string) (ret *model.Link, err error)

	FilterMedia(filter *bc.QueryFilter, namespace string) (ret []*model.Media, total int64, err error)
	QueryMedia(id int, namespace string) (ret *model.Media, err error)
	CreateMedia(ptr *model.Media, namespace string) (ret *model.Media, err error)
	UpdateMedia(ptr *model.Media, namespace string) (ret *model.Media, err error)
	DeleteMedia(id int, namespace string) (ret *model.Media, err error)

	FilterComment(filter *bc.QueryFilter, namespace string) (ret []*model.Comment, total int64, err error)
	QueryComment(id int, namespace string) (ret *model.Comment, err error)
	CreateComment(ptr *model.Comment, namespace string) (ret *model.Comment, err error)
	UpdateComment(ptr *model.Comment, namespace string) (ret *model.Comment, err error)
	DeleteComment(id int, namespace string) (ret *model.Comment, err error)
}

func New(clnt client.Client) Content {
	return &content{batisClient: clnt}
}

type content struct {
	batisClient client.Client
}

func (s *content) FilterArticle(filter *bc.QueryFilter, namespace string) (ret []*model.Article, total int64, err error) {
	ret = []*model.Article{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&model.Article{Catalog: []*model.Catalog{}})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *content) QueryArticle(id int, namespace string) (ret *model.Article, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.Article{Catalog: []*model.Catalog{}}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) CreateArticle(ptr *model.Article, namespace string) (ret *model.Article, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) UpdateArticle(ptr *model.Article, namespace string) (ret *model.Article, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) DeleteArticle(id int, namespace string) (ret *model.Article, err error) {
	valPtr, valErr := s.QueryArticle(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}

func (s *content) FilterCatalog(filter *bc.QueryFilter, namespace string) (ret []*model.Catalog, total int64, err error) {
	ret = []*model.Catalog{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&model.Catalog{Catalog: &model.Catalog{}})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *content) QueryCatalog(id int, namespace string) (ret *model.Catalog, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.Catalog{Catalog: &model.Catalog{}}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) CreateCatalog(ptr *model.Catalog, namespace string) (ret *model.Catalog, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) UpdateCatalog(ptr *model.Catalog, namespace string) (ret *model.Catalog, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) DeleteCatalog(id int, namespace string) (ret *model.Catalog, err error) {
	valPtr, valErr := s.QueryCatalog(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}

func (s *content) FilterLink(filter *bc.QueryFilter, namespace string) (ret []*model.Link, total int64, err error) {
	ret = []*model.Link{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&model.Link{Catalog: []*model.Catalog{}})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *content) QueryLink(id int, namespace string) (ret *model.Link, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.Link{Catalog: []*model.Catalog{}}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) CreateLink(ptr *model.Link, namespace string) (ret *model.Link, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) UpdateLink(ptr *model.Link, namespace string) (ret *model.Link, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) DeleteLink(id int, namespace string) (ret *model.Link, err error) {
	valPtr, valErr := s.QueryLink(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}

func (s *content) FilterMedia(filter *bc.QueryFilter, namespace string) (ret []*model.Media, total int64, err error) {
	ret = []*model.Media{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&model.Media{Catalog: []*model.Catalog{}})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *content) QueryMedia(id int, namespace string) (ret *model.Media, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.Media{Catalog: []*model.Catalog{}}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) CreateMedia(ptr *model.Media, namespace string) (ret *model.Media, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) UpdateMedia(ptr *model.Media, namespace string) (ret *model.Media, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) DeleteMedia(id int, namespace string) (ret *model.Media, err error) {
	valPtr, valErr := s.QueryMedia(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}

func (s *content) FilterComment(filter *bc.QueryFilter, namespace string) (ret []*model.Comment, total int64, err error) {
	ret = []*model.Comment{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&model.Comment{})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *content) QueryComment(id int, namespace string) (ret *model.Comment, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.Comment{}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) CreateComment(ptr *model.Comment, namespace string) (ret *model.Comment, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) UpdateComment(ptr *model.Comment, namespace string) (ret *model.Comment, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *content) DeleteComment(id int, namespace string) (ret *model.Comment, err error) {
	valPtr, valErr := s.QueryComment(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}
