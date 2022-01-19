package common

import (
	"github.com/muidea/magicCas/model"
	commonDef "github.com/muidea/magicCommon/def"
)

const (
	FilterRole = "/role/query/"
	QueryRole  = "/role/query/:id"
	CreateRole = "/role/create/"
	UpdateRole = "/role/update/:id"
	DeleteRole = "/role/delete/:id"
)

const RoleModule = "/module/role"

const superID = 999999

// SuperRole get super role
func SuperRole() *model.Role {
	return &model.Role{
		ID:      superID,
		Name:    "superRole",
		Private: []*model.Private{{ID: 999999, Value: AllPrivate, Path: "*"}},
	}
}

type RoleView struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Private     []*PrivateItem `json:"private"`
}

func (s *RoleView) FromRole(ptr *model.Role) {
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	for _, val := range ptr.Private {
		item := &PrivateItem{}
		item.FromPrivate(val)
		s.Private = append(s.Private, item)
	}
}

func (s *RoleView) ToRole() (ret *model.Role) {
	ret = &model.Role{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
	}

	for _, val := range s.Private {
		ret.Private = append(ret.Private, val.ToPrivate())
	}

	return
}

func (s *RoleView) IsSuper() bool {
	return s.ID == superID
}

type RoleLite struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *RoleLite) FromRole(ptr *model.Role) {
	s.ID = ptr.ID
	s.Name = ptr.Name
}

func (s *RoleLite) ToRole() (ret *model.Role) {
	return &model.Role{ID: s.ID, Name: s.Name}
}

// PrivateItem 单条配置项
type PrivateItem struct {
	ID    int          `json:"id"`
	Path  string       `json:"path"`
	Value *PrivateInfo `json:"value"`
}

func (s *PrivateItem) FromPrivate(ptr *model.Private) {
	s.ID = ptr.ID
	s.Path = ptr.Path
	s.Value = GetPrivateInfo(ptr.Value)
}

func (s *PrivateItem) ToPrivate() *model.Private {
	return &model.Private{ID: s.ID, Value: s.Value.Value, Path: s.Path}
}

type RoleParam struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Private     []*PrivateItem `json:"private"`
}

func (s *RoleParam) ToRole() (ret *model.Role) {
	ret = &model.Role{
		Name:        s.Name,
		Description: s.Description,
	}

	for _, val := range s.Private {
		ret.Private = append(ret.Private, val.ToPrivate())
	}

	return
}

type RoleResult struct {
	commonDef.Result
	Role *RoleView `json:"role"`
}

type RoleLiteListResult struct {
	commonDef.Result
	Total int64       `json:"total"`
	Role  []*RoleLite `json:"role"`
}

type RoleListResult struct {
	commonDef.Result
	Total int64       `json:"total"`
	Role  []*RoleView `json:"role"`
}

type RoleStatisticResult struct {
	RoleListResult
}
