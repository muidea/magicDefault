package common

import (
	"time"

	cd "github.com/muidea/magicCommon/def"

	cc "github.com/muidea/magicCas/common"

	"github.com/muidea/magicDefault/model"
)

const (
	// LoginAccount login account url
	LoginAccount = "/access/account/login/"
	// LogoutAccount logout account url
	LogoutAccount = "/access/account/logout/"
	// UpdateAccountPassword update account password url
	UpdateAccountPassword = "/access/account/password/update/"
	// VerifyEndpoint verify endpoint url
	VerifyEndpoint = "/access/endpoint/verify/"
	// RefreshSession refresh session url
	RefreshSession = "/access/session/refresh/"
	// VerifyEntityRole verify entity role
	VerifyEntityRole = "/access/entity/role/verify/"
	// QueryEntity query entity
	QueryEntity = "/access/entity/query/:id"
	// QueryAccessLog query access log url
	QueryAccessLog = "/access/log/query/"
	// QueryOperateLog query operate log url
	QueryOperateLog = "/operate/log/query/"
	// WriteOperateLog write operate log url
	WriteOperateLog = "/operate/log/write/"

	// UploadFile upload file url
	UploadFile = "/static/file/upload/"
	// ViewFile view file url
	ViewFile = "/static/file/view/"

	// QueryBaseInfo query base info
	QueryBaseInfo = "/base/info/query/"
	// NotifyTimer notify timer
	NotifyTimer = "/base/timer/notify/"
)

const BaseModule = "/kernel/base"

type LogView struct {
	ID         int            `json:"id"`
	Address    string         `json:"address"`
	Memo       string         `json:"memo"`
	Creater    *cc.EntityView `json:"creater"`
	CreateTime int64          `json:"createTime"`
}

func (s *LogView) FromLog(ptr *model.Log, createrPtr *cc.EntityView) {
	s.ID = ptr.ID
	s.Address = ptr.Address
	s.Memo = ptr.Memo
	s.Creater = createrPtr
	s.CreateTime = ptr.CreateTime
}

// OperateLogListResult operate log list result
type OperateLogListResult struct {
	cd.Result
	Total      int64      `json:"total"`
	OperateLog []*LogView `json:"operateLog"`
}

type TimerNotify struct {
	PreTime time.Time
	CurTime time.Time
}
