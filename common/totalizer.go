package common

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicDefault/model"
	"time"
)

const (
	QueryTotalizer  = "/kernel/totalizer/query/"
	CreateTotalizer = "/kernel/totalizer/create/"
	DeleteTotalizer = "/kernel/totalizer/delete/"
)

const (
	TotalizeRealtime = 1
	TotalizeWeek     = 2
	TotalizeMonth    = 3
	TotalizeDaily    = 4
)

const TotalizerModule = "/kernel/totalizer"

type Totalizer interface {
	Owner() string
	Trigger(event event.Event) (err error)
	Period(typeVal int, timeStamp time.Time) (err error)
	Get(typeVal int) (ret *model.Totalizer, err error)
	Same(owner, trigger string) bool
}

type Serializer interface {
	Load(owner string, typeVal, catalogVal int, namespace string) (ret *model.Totalizer, err error)
	Save(ptr *model.Totalizer) (ret *model.Totalizer, err error)
}

type TotalizeParam struct {
	Owner   string `json:"owner"`
	Trigger string `json:"trigger"`
	Period  []int  `json:"period"`
}

type TotalizerView struct {
	ID        int     `json:"id"`
	Owner     string  `json:"owner"`
	Type      int     `json:"type"`
	TimeStamp int64   `json:"timeStamp"`
	Value     float64 `json:"value"`
	Catalog   int     `json:"catalog"`
}

func (s *TotalizerView) FromTotalizer(ptr *model.Totalizer) {
	s.ID = ptr.ID
	s.Owner = ptr.Owner
	s.Type = ptr.Type
	s.TimeStamp = ptr.TimeStamp
	s.Value = ptr.Value
	s.Catalog = ptr.Catalog
}
