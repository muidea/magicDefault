package model

// Private 权限表
// value: 0 未设权限, 1 读权限, 2 写权限, 3 删除权限, 4 全部权限
// namePath: 权限路径
type Private struct {
	ID    int    `json:"id" orm:"id key auto"`
	Value int    `json:"value" orm:"value" validate:"required"`
	Path  string `json:"path" orm:"path" validate:"required"`
}

// Role private role
type Role struct {
	ID          int        `json:"id" orm:"id key auto"`
	Name        string     `json:"name" orm:"name" validate:"required"`
	Description string     `json:"description" orm:"description"`
	Private     []*Private `json:"private" orm:"private"`
	Namespace   string     `json:"namespace" orm:"namespace"`
}

func (s *Role) IsSuper() bool {
	return s.ID == 999999
}
