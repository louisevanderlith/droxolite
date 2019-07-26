package xontrols

import (
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
)

//InstanceCtrl is simply a base Controller that 'almost' implements a Controller.
//It should be inherited by the implenting Controller. API or UI, etc.
type InstanceCtrl struct {
	ctx       context.Contexer
	AvoCookie *bodies.Cookies
}

//CreateInstance is used to setup Context on controllers.
func (ctrl *InstanceCtrl) CreateInstance(ctx context.Contexer) {
	ctrl.ctx = ctx
}

//FindParam returns path variables /{var}
func (ctrl *InstanceCtrl) FindParam(name string) string {
	return ctrl.ctx.FindParam(name)
}

//SetHeader sets a response header.
func (ctrl *InstanceCtrl) SetHeader(key string, val string) {
	ctrl.ctx.SetHeader(key, val)
}

//Body populates the container with the Request Body
func (ctrl *InstanceCtrl) Body(container interface{}) error {
	return ctrl.ctx.Body(container)
}

//Ctx returns the Context object.
func (ctrl *InstanceCtrl) Ctx() context.Contexer {
	return ctrl.ctx
}
