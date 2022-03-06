package common

import "github.com/muidea/magicDefault/model"

const (
	QueryTotalizer  = "/base/totalizer/query/"
	CreateTotalizer = "/base/totalizer/create/"
	DeleteTotalizer = "/base/totalizer/delete/:id"
)

const (
	TotalizeRealtime = 1
	TotalizeWeek     = 2
	TotalizeMonth    = 3
)

const TotalizerModule = "/kernel/totalizer"

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
