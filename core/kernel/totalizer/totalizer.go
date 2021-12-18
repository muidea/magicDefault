package totalizer

import (
	"github.com/muidea/magicCommon/application"
	"github.com/muidea/magicCommon/module"
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

func (s *Totalizer) Setup(endpointName string) {
	app := application.GetApp()
	s.biz = biz.New(
		app.EventHub(),
		app.BackgroundRoutine(),
	)
}

func (s *Totalizer) Teardown() {

}
