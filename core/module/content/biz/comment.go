package biz

import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Content) FilterComment(filter *bc.QueryFilter, namespace string) (ret []*model.Comment, total int64, err error) {
	ret, total, err = s.contentDao.FilterComment(filter, namespace)
	return
}

func (s *Content) QueryComment(id int, namespace string) (ret *model.Comment, err error) {
	ret, err = s.contentDao.QueryComment(id, namespace)
	return
}

func (s *Content) CreateComment(ptr *common.CommentParam, creater int, namespace string) (ret *model.Comment, err error) {
	commentPtr := ptr.ToComment(nil)
	commentPtr.Creater = creater
	commentPtr.CreateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.CreateComment(commentPtr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentComment, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) UpdateComment(id int, ptr *common.CommentParam, updater int, namespace string) (ret *model.Comment, err error) {
	currentComment, currentErr := s.contentDao.QueryComment(id, namespace)
	if currentErr != nil {
		err = currentErr
		return
	}

	currentComment = ptr.ToComment(currentComment)
	currentComment.Creater = updater
	currentComment.CreateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.UpdateComment(currentComment, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentComment, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) DeleteComment(id int, namespace string) (ret *model.Comment, err error) {
	ret, err = s.contentDao.DeleteComment(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentComment, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
