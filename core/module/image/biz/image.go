package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/event"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/task"

	fClnt "github.com/muidea/magicFile/client"
	fc "github.com/muidea/magicFile/common"
	fm "github.com/muidea/magicFile/model"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
)

type Image struct {
	biz.Base

	endpointName string
	fileService  string
}

func New(
	endpointName string,
	fileService string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Image {
	ptr := &Image{
		Base:         biz.New(common.ImageModule, eventHub, backgroundRoutine),
		endpointName: endpointName,
		fileService:  fileService,
	}

	return ptr
}

func (s *Image) FilterImage(filter *fu.ContentFilter, namespace string) (ret []*fm.FileDetail, total int64, err error) {
	clnt := fClnt.NewClient(s.fileService)
	defer clnt.Release()

	clnt.AttachSource(fmt.Sprintf("%s_%s", s.endpointName, namespace))

	ret, total, err = clnt.FilterFile(filter)

	return
}

func (s *Image) UpdateImage(id int, param *fc.FileParam, namespace string) (ret *fm.FileDetail, err error) {
	clnt := fClnt.NewClient(s.fileService)
	defer clnt.Release()

	clnt.AttachSource(fmt.Sprintf("%s_%s", s.endpointName, namespace))

	ret, err = clnt.UpdateFile(id, param)
	return
}

func (s *Image) DeleteImage(id int, namespace string) (ret *fm.FileDetail, err error) {
	clnt := fClnt.NewClient(s.fileService)
	defer clnt.Release()

	clnt.AttachSource(fmt.Sprintf("%s_%s", s.endpointName, namespace))

	ret, err = clnt.DeleteFile(id)

	return
}

func (s *Image) QueryImage(id int, namespace string) (ret *fm.FileDetail, err error) {
	clnt := fClnt.NewClient(s.fileService)
	defer clnt.Release()

	clnt.AttachSource(fmt.Sprintf("%s_%s", s.endpointName, namespace))

	ret, err = clnt.QueryFile(id)

	return
}
