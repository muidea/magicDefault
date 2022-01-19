package common

import "github.com/muidea/magicCas/model"

// Status status
type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IsInitStatus is init status
func (s *Status) IsInitStatus() bool {
	return s.ID == model.InitStatus
}

// GetStatus get status
func GetStatus(id int) *Status {
	switch id {
	case model.EnableStatus:
		return &Status{ID: model.EnableStatus, Name: "启用"}
	case model.DisableStatus:
		return &Status{ID: model.DisableStatus, Name: "停用"}
	}

	return &Status{ID: model.InitStatus, Name: "初始化"}
}
