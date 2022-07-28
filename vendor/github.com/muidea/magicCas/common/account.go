package common

import (
	"github.com/muidea/magicCas/model"
	commonDef "github.com/muidea/magicCommon/def"
)

const (
	FilterAccount = "/account/query/"
	QueryAccount  = "/account/query/:id"
	CreateAccount = "/account/create/"
	UpdateAccount = "/account/update/:id"
	DeleteAccount = "/account/delete/:id"
	CheckAccount  = "/account/check/"
)

const AccountModule = "/module/account"

// AccountView account
type AccountView struct {
	//ID 唯一标示单元
	ID          int       `json:"id"`
	Account     string    `json:"account"`
	EMail       string    `json:"email"`
	Description string    `json:"description"`
	Status      *Status   `json:"status"`
	Role        *RoleLite `json:"role"`
	CreateTime  int64     `json:"createTime"`
}

// FromAccount from account
func (s *AccountView) FromAccount(ptr *model.Account) {
	if ptr == nil {
		return
	}

	s.ID = ptr.ID
	s.Account = ptr.Account
	s.EMail = ptr.EMail
	s.Description = ptr.Description
	s.Status = GetStatus(ptr.Status)
	if ptr.Role != nil {
		s.Role = &RoleLite{}
		s.Role.FromRole(ptr.Role)
	}

	s.CreateTime = ptr.CreateTime
}

func (s *AccountView) ToAccount() (ret *model.Account) {
	ret = &model.Account{
		ID:          s.ID,
		Account:     s.Account,
		EMail:       s.EMail,
		Description: s.Description,
		Status:      s.Status.ID,
		Role:        s.Role.ToRole(),
		CreateTime:  s.CreateTime,
	}

	return
}

// Entity get account entity
func (s *AccountView) Entity() *model.Entity {
	return &model.Entity{EID: s.ID, EType: model.AccountType, EName: s.Account}
}

type AccountLite struct {
	ID      int     `json:"id"`
	Account string  `json:"account"`
	Status  *Status `json:"status"`
}

func (s *AccountLite) FromAccount(ptr *model.Account) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Account = ptr.Account
	s.Status = GetStatus(ptr.Status)
}

func (s *AccountLite) FromAccountView(ptr *AccountView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Account = ptr.Account
	s.Status = ptr.Status
}

func (s *AccountLite) ToAccount() (ret *model.Account) {
	ret = &model.Account{
		ID:      s.ID,
		Account: s.Account,
		Status:  s.Status.ID,
	}

	return
}

type AccountParam struct {
	Account     string    `json:"account"`
	Password    string    `json:"password"`
	EMail       string    `json:"email"`
	Description string    `json:"description"`
	Status      *Status   `json:"status"`
	Role        *RoleLite `json:"role"`
}

func (s *AccountParam) ToAccount() (ret *model.Account) {
	ret = &model.Account{
		Account:     s.Account,
		Password:    s.Password,
		EMail:       s.EMail,
		Description: s.Description,
	}
	if s.Status != nil {
		ret.Status = s.Status.ID
	}
	if s.Role != nil {
		ret.Role = s.Role.ToRole()
	}
	return
}

type AccountResult struct {
	commonDef.Result
	Account *AccountView `json:"account"`
}

type AccountLiteListResult struct {
	commonDef.Result
	Total   int64          `json:"total"`
	Account []*AccountLite `json:"account"`
}

type AccountListResult struct {
	commonDef.Result
	Total   int64          `json:"total"`
	Account []*AccountView `json:"account"`
}

type AccountStatisticResult struct {
	AccountListResult
	Role []*RoleLite `json:"role"`
}
