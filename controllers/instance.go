package controllers

import (
	"github.com/louisevanderlith/droxolite/context"
	secure "github.com/louisevanderlith/secure/core"
)

//InstanceCtrl is simply a base Controller that 'almost' implements a Controller.
//It should be inherited by the implenting Controller. API or UI, etc.
type InstanceCtrl struct {
	Ctx       context.Contexer
	AvoCookie *secure.Cookies
}

//CreateInstance is used to setup Context on controllers.
func (ctrl *InstanceCtrl) CreateInstance(ctx context.Contexer) {
	ctrl.Ctx = ctx
}