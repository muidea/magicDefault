package model

import (
	"time"
)

type Validity struct {
	ID        int   `json:"id" orm:"id key auto"`
	StartTime int64 `json:"startTime" orm:"startTime"`
	EndTime   int64 `json:"endTime" orm:"endTime"`
}

func (s *Validity) IsSame(ptr *Validity) bool {
	if ptr == nil {
		return false
	}

	return s.Expire() == ptr.Expire()
}

func (s *Validity) Validate() bool {
	if s.EndTime < 0 {
		return true
	}

	return time.Now().UTC().Unix() < s.EndTime
}

func (s *Validity) UpdateExpire(expire int) {
	if s.Validate() {
		s.EndTime = s.StartTime + int64(24*60*60*expire)
		return
	}

	s.ResetExpire(expire)
}

func (s *Validity) ResetExpire(expire int) {
	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour * time.Duration(24*expire))

	s.StartTime = startTime.Unix()
	s.EndTime = endTime.Unix()
}

func (s *Validity) Expire() int {
	expire := time.Unix(s.EndTime, 0).Sub(time.Now().UTC())
	if expire <= 0 {
		return 0
	}

	return int(expire.Hours()/24) + 1
}

func NewValidity(val int) Validity {
	if val <= 0 {
		val = 0
	}

	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour * time.Duration(24*val))
	return Validity{StartTime: startTime.Unix(), EndTime: endTime.Unix()}
}

type Namespace struct {
	ID          int      `json:"id" orm:"id key auto"`
	Name        string   `json:"name" orm:"name"`
	Short       string   `json:"short" orm:"short"`
	Description string   `json:"description" orm:"description"`
	Status      int      `json:"status" orm:"status"`
	Validity    Validity `json:"validity" orm:"validity"`
	CreateTime  int64    `json:"createTime" orm:"createTime"`
}

func (s *Namespace) Assign(ptr *Namespace) {
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Short = ptr.Short
	s.Description = ptr.Description
	s.Status = ptr.Status
	s.Validity = ptr.Validity
	s.CreateTime = ptr.CreateTime
}

func (s *Namespace) Disable() bool {
	return s.Status == DisableStatus
}

func (s *Namespace) Validate() bool {
	return s.Validity.Validate() && !s.Disable()
}
