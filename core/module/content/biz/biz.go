package biz

import (
	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/module/content/dao"
)

type Content struct {
	biz.Base

	contentDao dao.Content

	endpointName string
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *Content {

	ptr := &Content{
		Base:         biz.New(common.ContentModule, eventHub, backgroundRoutine),
		contentDao:   dao.New(batisClient),
		endpointName: endpointName,
	}

	return ptr
}

func (s *Content) Notify(event event.Event, result event.Result) {
	return
}
