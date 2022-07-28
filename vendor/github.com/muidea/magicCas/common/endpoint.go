package common

import (
	"github.com/muidea/magicCas/model"
	commonDef "github.com/muidea/magicCommon/def"
)

const (
	FilterEndpoint = "/endpoint/query/"
	QueryEndpoint  = "/endpoint/query/:id"
	CreateEndpoint = "/endpoint/create/"
	UpdateEndpoint = "/endpoint/update/:id"
	DeleteEndpoint = "/endpoint/delete/:id"
)

const EndpointModule = "/module/endpoint"

// EndpointView endpoint view
type EndpointView struct {
	//ID 唯一标示单元
	ID          int       `json:"id"`
	Endpoint    string    `json:"endpoint"`
	Description string    `json:"description"`
	IdentifyID  string    `json:"identifyID"`
	AuthToken   string    `json:"authToken"`
	Status      *Status   `json:"status"`
	Role        *RoleLite `json:"role"`
	CreateTime  int64     `json:"createTime"`
}

// FromEndpoint from endpoint
func (s *EndpointView) FromEndpoint(ptr *model.Endpoint) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Endpoint = ptr.Endpoint
	s.Description = ptr.Description
	s.IdentifyID = ptr.IdentifyID
	s.AuthToken = ptr.AuthToken
	s.Status = GetStatus(ptr.Status)

	if ptr.Role != nil {
		s.Role = &RoleLite{}
		s.Role.FromRole(ptr.Role)
	}

	s.CreateTime = ptr.CreateTime
}

func (s *EndpointView) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		ID:          s.ID,
		Endpoint:    s.Endpoint,
		Description: s.Description,
		IdentifyID:  s.IdentifyID,
		AuthToken:   s.AuthToken,
		Status:      s.Status.ID,
		Role:        s.Role.ToRole(),
		CreateTime:  s.CreateTime,
	}

	return
}

// Entity get endpoint entity
func (s *EndpointView) Entity() *model.Entity {
	return &model.Entity{EID: s.ID, EType: model.EndpointType, EName: s.Endpoint}
}

type EndpointLite struct {
	ID         int     `json:"id"`
	Endpoint   string  `json:"endpoint"`
	IdentifyID string  `json:"identifyID"`
	Status     *Status `json:"status"`
}

func (s *EndpointLite) FromEndpoint(ptr *model.Endpoint) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Endpoint = ptr.Endpoint
	s.IdentifyID = ptr.IdentifyID
	s.Status = GetStatus(ptr.Status)
}

func (s *EndpointLite) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		ID:         s.ID,
		Endpoint:   s.Endpoint,
		IdentifyID: s.IdentifyID,
		Status:     s.Status.ID,
	}

	return
}

type EndpointParam struct {
	Endpoint    string    `json:"endpoint"`
	Description string    `json:"description"`
	IdentifyID  string    `json:"identifyID"`
	AuthToken   string    `json:"authToken"`
	Status      *Status   `json:"status"`
	Role        *RoleLite `json:"role"`
}

func (s *EndpointParam) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		Endpoint:    s.Endpoint,
		Description: s.Description,
		IdentifyID:  s.IdentifyID,
		AuthToken:   s.AuthToken,
	}
	if s.Status != nil {
		ret.Status = s.Status.ID
	}
	if s.Role != nil {
		ret.Role = s.Role.ToRole()
	}
	return
}

type EndpointResult struct {
	commonDef.Result
	Endpoint *EndpointView `json:"endpoint"`
}

type EndpointLiteListResult struct {
	commonDef.Result
	Total    int64           `json:"total"`
	Endpoint []*EndpointLite `json:"endpoint"`
}

type EndpointListResult struct {
	commonDef.Result
	Total    int64           `json:"total"`
	Endpoint []*EndpointView `json:"endpoint"`
}

type EndpointStatisticResult struct {
	EndpointListResult
	Role []*RoleLite `json:"role"`
}
