package common

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	"github.com/muidea/magicCas/model"
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

	// SaveEntity save entity
	SaveEntity = "/access/entity/save/"
	// DeleteEntity delete entity
	DeleteEntity = "/access/entity/delete/"
	// InitializeNamespace initialize namespace
	InitializeNamespace = "/initialize/namespace/"
	// WriteAccessLog write access log url
	WriteAccessLog = "/access/log/write/"

	// NotifyTimer notify timer
	NotifyTimer = "/base/timer/notify/"
)

const BaseModule = "/kernel/base"

type EntityView struct {
	ID    int    `json:"id"`
	EName string `json:"name"`
	EID   int    `json:"eID"`
	EType string `json:"eType"`
}

func (s *EntityView) FromEntity(ptr *model.Entity) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.EID = ptr.EID
	s.EName = ptr.EName
	s.EType = ptr.EType
}

func (s *EntityView) ToEntity() (ret *model.Entity) {
	ret = &model.Entity{
		ID:    s.ID,
		EName: s.EName,
		EID:   s.EID,
		EType: s.EType,
	}
	return
}

// Decode decode entityPtr
func (s *EntityView) Decode(req *http.Request) (err error) {
	str := req.URL.Query().Get("id")
	if str != "" {
		s.ID, err = strconv.Atoi(str)
		if err != nil {
			return
		}
	}
	str = req.URL.Query().Get("entityID")
	if str != "" {
		s.EID, err = strconv.Atoi(str)
		if err != nil {
			return
		}

	} else {
		err = fmt.Errorf("illegal entity info")
		return
	}

	s.EName = req.URL.Query().Get("entityName")
	s.EType = req.URL.Query().Get("entityType")
	return
}

// Encode encode entityPtr
func (s *EntityView) Encode(vals url.Values) url.Values {
	vals.Set("id", fmt.Sprintf("%d", s.ID))
	vals.Set("entityID", fmt.Sprintf("%d", s.EID))
	vals.Set("entityName", s.EName)
	vals.Set("entityType", s.EType)
	return vals
}

// IsValid check valid
func (s *EntityView) IsValid() bool {
	if s.EID <= 0 {
		return false
	}

	switch s.EType {
	case model.AccountType, model.EndpointType:
		return true
	}

	return false
}

func (s *EntityView) IsAdmin() bool {
	superUser := model.DefaultSuperAccount()

	return s.EType == model.AccountType && s.EName == superUser.Account
}

func (s *EntityView) IsSupper() bool {
	if s.EType == model.EndpointType {
		return true
	}

	return s.IsAdmin()
}

// UnknownEntity 未知账号
func UnknownEntity() *EntityView {
	return &EntityView{ID: -1, EID: -1, EName: "未知账号"}
}

const (
	// InvalidPermission 无效权限
	InvalidPermission = iota
	// UnSetPermission 未设权限
	UnSetPermission
	// ReadPermission 只读权限
	ReadPermission
	// WritePermission 可写权限
	WritePermission
	// DeletePermission 删除权限
	DeletePermission
	// AllPermission 全部权限
	AllPermission
)

// Permission info
type Permission struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

var permissionList = []*Permission{
	{Value: InvalidPermission, Name: "无效权限"},
	{Value: UnSetPermission, Name: "未设权限"},
	{Value: ReadPermission, Name: "只读权限"},
	{Value: WritePermission, Name: "可写权限"},
	{Value: DeletePermission, Name: "删除权限"},
	{Value: AllPermission, Name: "全部权限"},
}

// GetPermissionList get permission list
func GetPermissionList() []*Permission {
	return permissionList
}

// GetPermission get permission
func GetPermission(value int) (ret *Permission) {
	switch value {
	case UnSetPermission,
		ReadPermission,
		WritePermission,
		DeletePermission,
		AllPermission:
		ret = permissionList[value]
	default:
		ret = permissionList[0]
	}

	return
}

type LogView struct {
	ID         int         `json:"id"`
	Address    string      `json:"address"`
	Memo       string      `json:"memo"`
	Creater    *EntityView `json:"creater"`
	CreateTime int64       `json:"createTime"`
}

func (s *LogView) FromLog(ptr *model.Log, createrPtr *model.Entity) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Address = ptr.Address
	s.Memo = ptr.Memo
	s.Creater = &EntityView{}
	s.Creater.FromEntity(createrPtr)
	s.CreateTime = ptr.CreateTime
}

type RefreshResult struct {
	commonDef.Result
	Entity      *EntityView                `json:"entity"`
	SessionInfo *commonSession.SessionInfo `json:"sessionInfo"`
}

type LoginParam struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdatePasswordParam struct {
	Account     string `json:"account" validate:"required"`
	CurPassword string `json:"curPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type LoginResult RefreshResult

type LogoutResult struct {
	commonDef.Result
	SessionInfo *commonSession.SessionInfo `json:"sessionInfo"`
}

type VerifyEndpointParam struct {
	Endpoint   string `json:"endpoint" validate:"required"`
	IdentifyID string `json:"identifyID" validate:"required"`
	AuthToken  string `json:"authToken" validate:"required"`
}

type VerifyEndpointResult LoginResult

type EntityRoleResult struct {
	commonDef.Result
	Role *RoleView `json:"role"`
}

type QueryEntityResult struct {
	commonDef.Result
	Entity *EntityView `json:"entity"`
}

type AccessLogListResult struct {
	commonDef.Result
	Total     int64      `json:"total"`
	AccessLog []*LogView `json:"accessLog"`
}

// EnumPrivilegeItemResult enum privilege item result
type EnumPrivilegeItemResult struct {
	commonDef.Result
	Privilege []*Privilege `json:"privilege"`
}

type TimerNotify struct {
	PreTime time.Time
	CurTime time.Time
}
