package biz

import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/model"
)

func (s *Content) FilterCatalog(filter *bc.QueryFilter, namespace string) (ret []*model.Catalog, total int64, err error) {
	ret, total, err = s.contentDao.FilterCatalog(filter, namespace)
	return
}

func (s *Content) QueryCatalog(id int, namespace string) (ret *model.Catalog, err error) {
	ret, err = s.contentDao.QueryCatalog(id, namespace)
	return
}

func (s *Content) CreateCatalog(ptr *common.CatalogParam, creater int, namespace string) (ret *model.Catalog, err error) {
	catalogPtr := ptr.ToCatalog()
	catalogPtr.Creater = creater
	catalogPtr.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.contentDao.CreateCatalog(catalogPtr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentCatalog, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) UpdateCatalog(id int, ptr *common.CatalogParam, updater int, namespace string) (ret *model.Catalog, err error) {
	currentCatalog, currentErr := s.contentDao.QueryCatalog(id, namespace)
	if currentErr != nil {
		return
	}

	currentCatalog.Creater = updater
	currentCatalog.UpdateTime = time.Now().UTC().Unix()

	ret, err = s.contentDao.UpdateCatalog(currentCatalog, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentCatalog, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *Content) DeleteCatalog(id int, namespace string) (ret *model.Catalog, err error) {
	ret, err = s.contentDao.DeleteCatalog(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.NotifyContentCatalog, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
