package common

import (
	"github.com/muidea/magicCas/model"
	commonDef "github.com/muidea/magicCommon/def"
)

const (
	FilterNamespace = "/namespace/query/"
	QueryNamespace  = "/namespace/query/:id"
	CreateNamespace = "/namespace/create/"
	UpdateNamespace = "/namespace/update/:id"
	DeleteNamespace = "/namespace/delete/:id"
)

const NamespaceModule = "/module/namespace"

type ValidityView struct {
	ID        int   `json:"id"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Expired   int   `json:"expired"`
}

func (s *ValidityView) FromValidity(ptr *model.Validity) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.StartTime = ptr.StartTime
	s.EndTime = ptr.EndTime
	if ptr.EndTime > 0 {
		s.Expired = ptr.Expire()
		return
	}

	s.Expired = 36500
}

func (s *ValidityView) ToValidity() (ret *model.Validity) {
	ret = &model.Validity{
		ID:        s.ID,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
	}

	return
}

type NamespaceView struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Short       string        `json:"short"`
	Description string        `json:"description"`
	Status      *Status       `json:"status"`
	Validity    *ValidityView `json:"validity"`
	CreateTime  int64         `json:"createTime"`
}

func (s *NamespaceView) Disable() bool {
	return s.Status.ID == model.DisableStatus
}

func (s *NamespaceView) FromNamespace(ptr *model.Namespace) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Short = ptr.Short
	s.Description = ptr.Description
	s.Status = GetStatus(ptr.Status)
	s.Validity = &ValidityView{}
	s.Validity.FromValidity(&ptr.Validity)
	s.CreateTime = ptr.CreateTime
}

func (s *NamespaceView) ToNamespace() (ret *model.Namespace) {
	ret = &model.Namespace{
		ID:          s.ID,
		Name:        s.Name,
		Short:       s.Short,
		Description: s.Description,
		Status:      s.Status.ID,
		Validity:    *s.Validity.ToValidity(),
		CreateTime:  s.CreateTime,
	}

	return
}

type NamespaceLite struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

func (s *NamespaceLite) FromNamespace(ptr *model.Namespace) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Short = ptr.Short
}

func (s *NamespaceLite) ToNamespace() (ret *model.Namespace) {
	ret = &model.Namespace{
		ID:    s.ID,
		Name:  s.Name,
		Short: s.Short,
	}

	return
}

type NamespaceParam struct {
	Name        string  `json:"name" validate:"required"`
	Short       string  `json:"short"`
	Description string  `json:"description"`
	Status      *Status `json:"status"`
	Validity    int     `json:"validity" validate:"required"`
}

func (s *NamespaceParam) ToNamespace() (ret *model.Namespace) {
	ret = &model.Namespace{
		Name:        s.Name,
		Short:       s.Short,
		Description: s.Description,
		Validity:    model.NewValidity(s.Validity),
	}
	if s.Status != nil {
		ret.Status = s.Status.ID
	}
	return
}

type NamespaceResult struct {
	commonDef.Result
	Namespace *NamespaceView `json:"namespace"`
}

type NamespaceLiteListResult struct {
	commonDef.Result
	Total     int64            `json:"total"`
	Namespace []*NamespaceLite `json:"namespace"`
}

type NamespaceListResult struct {
	commonDef.Result
	Total     int64            `json:"total"`
	Namespace []*NamespaceView `json:"namespace"`
}

type NamespaceStatisticResult struct {
	NamespaceListResult
}
