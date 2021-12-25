package totalizer

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/assist/persistence"
	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/kernel/totalizer/biz"
)

func init() {
	module.Register(New())
}

type Totalizer struct {
	biz *biz.Totalizer
}

func New() module.Module {
	return &Totalizer{}
}

func (s *Totalizer) ID() string {
	return common.TotalizerModule
}

func (s *Totalizer) Setup(endpointName string, eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		persistence.GetBatisClient(),
		eventHub,
		backgroundRoutine,
	)
}

func (s *Totalizer) Teardown() {

}
