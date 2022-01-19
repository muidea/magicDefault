package model

import (
	"fmt"
)

const (
	// AccountType account
	AccountType = "account"
	// EndpointType endpoint
	EndpointType = "endpoint"
)

// Entity accessLog entity type
type Entity struct {
	ID        int    `json:"id" orm:"id key auto"`
	EName     string `json:"name" orm:"eName" validate:"required"`
	EID       int    `json:"eID" orm:"eID" validate:"required"`
	EType     string `json:"eType" orm:"eType" validate:"required"`
	Namespace string `json:"namespace" orm:"namespace" validate:"required"`
}

func (s *Entity) String() string {
	switch s.EType {
	case AccountType:
		return fmt.Sprintf("用户:%s", s.EName)
	case EndpointType:
		return fmt.Sprintf("终端:%s", s.EName)
	}

	return s.EName
}

func (s *Entity) IsAdmin() bool {
	superUser := DefaultSuperAccount()

	return s.EType == AccountType && s.EName == superUser.Account
}
