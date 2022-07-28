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
		ID:        superID,
		Name:      "superRole",
		Privilege: []*model.Privilege{{ID: 999999, Value: AllPermission, Path: "*"}},
	}
}

type RoleView struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Privilege   []*Privilege `json:"privilege"`
}

func (s *RoleView) FromRole(ptr *model.Role) {
	if ptr == nil {
		return
	}

	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	for _, val := range ptr.Privilege {
		item := &Privilege{}
		item.FromPrivilege(val)
		s.Privilege = append(s.Privilege, item)
	}
}

func (s *RoleView) ToRole() (ret *model.Role) {
	ret = &model.Role{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
	}

	for _, val := range s.Privilege {
		ret.Privilege = append(ret.Privilege, val.ToPrivilege())
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
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
}

func (s *RoleLite) ToRole() (ret *model.Role) {
	return &model.Role{ID: s.ID, Name: s.Name}
}

// Privilege 单条配置项
type Privilege struct {
	ID    int         `json:"id"`
	Path  string      `json:"path"`
	Value *Permission `json:"value"`
}

func (s *Privilege) FromPrivilege(ptr *model.Privilege) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Path = ptr.Path
	s.Value = GetPermission(ptr.Value)
}

func (s *Privilege) ToPrivilege() *model.Privilege {
	return &model.Privilege{ID: s.ID, Value: s.Value.Value, Path: s.Path}
}

type RoleParam struct {
	Name        string       `json:"name" validate:"required"`
	Description string       `json:"description"`
	Privilege   []*Privilege `json:"privilege"`
}

func (s *RoleParam) ToRole() (ret *model.Role) {
	ret = &model.Role{
		Name:        s.Name,
		Description: s.Description,
	}

	for _, val := range s.Privilege {
		ret.Privilege = append(ret.Privilege, val.ToPrivilege())
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
