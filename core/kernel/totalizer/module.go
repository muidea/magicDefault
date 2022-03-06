package totalizer

import (
	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicDefault/common"
	"github.com/muidea/magicDefault/core/kernel/totalizer/biz"
)

func init() {
	module.Register(New())
}

type Totalizer struct {
	batisClient client.Client

	biz *biz.Totalizer
}

func New() *Totalizer {
	return &Totalizer{}
}

func (s *Totalizer) ID() string {
	return common.TotalizerModule
}

func (s *Totalizer) BindBatisClient(clnt client.Client) {
	s.batisClient = clnt
}

func (s *Totalizer) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		s.batisClient,
		eventHub,
		backgroundRoutine,
	)
}

func (s *Totalizer) Teardown() {

}
