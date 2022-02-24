package biz

import (
	"time"

	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/base/biz"
	"github.com/muidea/magicDefault/core/kernel/setting/dao"
)

type Setting struct {
	biz.Base

	settingDao dao.Setting
	startTime  time.Time
}

func New(
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Setting {
	ptr := &Setting{
		Base:       biz.New(common.SettingModule, eventHub, backgroundRoutine),
		settingDao: dao.New(batisClient),
		startTime:  time.Now(),
	}

	return ptr
}

func (s *Setting) Notify(event event.Event, result event.Result) {
	// TODO
}
