package biz

import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Content) FilterLink(filter *bc.QueryFilter, namespace string) (ret []*model.Link, total int64, err error) {
	ret, total, err = s.contentDao.FilterLink(filter, namespace)
	return
}

func (s *Content) QueryLink(id int, namespace string) (ret *model.Link, err error) {
	ret, err = s.contentDao.QueryLink(id, namespace)
	return
}

func (s *Content) CreateLink(ptr *common.LinkParam, creater int, namespace string) (ret *model.Link, err error) {
	linkPtr := ptr.ToLink(nil)
	linkPtr.Creater = creater
	linkPtr.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.CreateLink(linkPtr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentLink, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) UpdateLink(id int, ptr *common.LinkParam, updater int, namespace string) (ret *model.Link, err error) {
	currentLink, currentErr := s.contentDao.QueryLink(id, namespace)
	if currentErr != nil {
		err = currentErr
		return
	}

	currentLink = ptr.ToLink(currentLink)
	currentLink.Creater = updater
	currentLink.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.UpdateLink(currentLink, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentLink, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) DeleteLink(id int, namespace string) (ret *model.Link, err error) {
	ret, err = s.contentDao.DeleteLink(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentLink, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
