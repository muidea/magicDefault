package biz

import (
	"fmt"
	"time"

	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/config"
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

func (s *Setting) QuerySetting(sessionInfo *session.SessionInfo, namespace string) (ret []*common.Content, err error) {
	ret = []*common.Content{}

	if namespace == config.SuperNamespace() {
		startTimeItem := &common.Content{Name: "系统启动时间", Content: s.startTime.Local().Format("2006-01-02 15:04:05")}
		ret = append(ret, startTimeItem)
		return
	}

	namespacePtr := s.queryNamespace(sessionInfo, namespace)
	if namespacePtr == nil {
		err = fmt.Errorf("指定租户不存在,%s", namespace)
		return
	}
	startTime := time.Unix(namespacePtr.Validity.StartTime, 0)
	startTimeItem := &common.Content{Name: "开通时间", Content: startTime.Local().Format("2006-01-02 15:04:05")}
	ret = append(ret, startTimeItem)

	endTime := time.Unix(namespacePtr.Validity.EndTime, 0)
	endTimeItem := &common.Content{Name: "截止时间", Content: endTime.Local().Format("2006-01-02 15:04:05")}
	ret = append(ret, endTimeItem)

	expireItem := &common.Content{Name: "有效期限", Content: fmt.Sprintf("剩余时间%d天", namespacePtr.Validity.Expired)}
	ret = append(ret, expireItem)
	return
}
