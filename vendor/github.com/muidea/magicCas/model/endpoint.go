package model

// Endpoint endpoint
type Endpoint struct {
	//ID 唯一标示单元
	ID          int    `json:"id" orm:"id key auto"`
	Endpoint    string `json:"endpoint" orm:"endpoint" validate:"required"`
	Description string `json:"description" orm:"description"`
	IdentifyID  string `json:"identifyID" orm:"identifyID"`
	AuthToken   string `json:"authToken" orm:"authToken"`
	Status      int    `json:"status" orm:"status"`
	Role        *Role  `json:"role" orm:"role"`
	CreateTime  int64  `json:"createTime" orm:"createTime"`
	Namespace   string `json:"namespace" orm:"namespace"`
}

// Enable enable status
func (s *Endpoint) Enable() bool {
	return s.Status == EnableStatus
}

// Disable disable status
func (s *Endpoint) Disable() bool {
	return s.Status == DisableStatus
}

// Entity get endpoint entity
func (s *Endpoint) Entity() *Entity {
	return &Entity{EName: s.Endpoint, EID: s.ID, EType: EndpointType}
}

//IsSame is same endpoint
func (s *Endpoint) IsSame(endpoint *Endpoint) bool {
	if s.ID != endpoint.ID {
		return false
	}

	if s.Endpoint != endpoint.Endpoint {
		return false
	}

	if s.Description != endpoint.Description {
		return false
	}

	if s.IdentifyID != endpoint.IdentifyID {
		return false
	}

	if s.AuthToken != endpoint.AuthToken {
		return false
	}

	if s.Status != endpoint.Status {
		return false
	}

	if s.Role.ID != endpoint.Role.ID {
		return false
	}

	return true
}
