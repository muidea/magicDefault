package biz

import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Content) FilterMedia(filter *bc.QueryFilter, namespace string) (ret []*model.Media, total int64, err error) {
	ret, total, err = s.contentDao.FilterMedia(filter, namespace)
	return
}

func (s *Content) QueryMedia(id int, namespace string) (ret *model.Media, err error) {
	ret, err = s.contentDao.QueryMedia(id, namespace)
	return
}

func (s *Content) CreateMedia(ptr *common.MediaParam, creater int, namespace string) (ret *model.Media, err error) {
	mediaPtr := ptr.ToMedia(nil)
	mediaPtr.Creater = creater
	mediaPtr.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.CreateMedia(mediaPtr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentMedia, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) UpdateMedia(id int, ptr *common.MediaParam, updater int, namespace string) (ret *model.Media, err error) {
	currentMedia, currentErr := s.contentDao.QueryMedia(id, namespace)
	if currentErr != nil {
		err = currentErr
		return
	}

	currentMedia = ptr.ToMedia(currentMedia)
	currentMedia.Creater = updater
	currentMedia.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.UpdateMedia(currentMedia, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentMedia, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) DeleteMedia(id int, namespace string) (ret *model.Media, err error) {
	ret, err = s.contentDao.DeleteMedia(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentMedia, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
