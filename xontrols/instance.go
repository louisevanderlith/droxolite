package xontrols

import (
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
)

//InstanceCtrl is simply a base Controller that 'almost' implements a Controller.
//It should be inherited by the implenting Controller. API or UI, etc.
type InstanceCtrl struct {
	Ctx       context.Contexer
	AvoCookie *bodies.Cookies
}

//CreateInstance is used to setup Context on controllers.
func (ctrl *InstanceCtrl) CreateInstance(ctx context.Contexer) {
	ctrl.Ctx = ctx
}
