package biz

import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Content) FilterArticle(filter *bc.QueryFilter, namespace string) (ret []*model.Article, total int64, err error) {
	ret, total, err = s.contentDao.FilterArticle(filter, namespace)
	return
}

func (s *Content) QueryArticle(id int, namespace string) (ret *model.Article, err error) {
	ret, err = s.contentDao.QueryArticle(id, namespace)
	return
}

func (s *Content) CreateArticle(ptr *common.ArticleParam, creater int, namespace string) (ret *model.Article, err error) {
	articlePtr := ptr.ToArticle(nil)
	articlePtr.Creater = creater
	articlePtr.CreateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.CreateArticle(articlePtr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentArticle, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) UpdateArticle(id int, ptr *common.ArticleParam, updater int, namespace string) (ret *model.Article, err error) {
	currentArticle, currentErr := s.contentDao.QueryArticle(id, namespace)
	if currentErr != nil {
		err = currentErr
		return
	}

	currentArticle = ptr.ToArticle(currentArticle)
	currentArticle.Creater = updater
	currentArticle.CreateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.UpdateArticle(currentArticle, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentArticle, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) DeleteArticle(id int, namespace string) (ret *model.Article, err error) {
	ret, err = s.contentDao.DeleteArticle(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentArticle, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
