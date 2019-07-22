package droxolite

import (
	"github.com/louisevanderlith/droxolite/controllers"
)

//Epoxy puts everything together
type Epoxy struct {
	service     *Service
	controllers []controllers.Controller
}

func (e *Epoxy) Plak() {
	//avoc, err := bodies.GetAvoCookie(ctrl.GetMyToken(), ctrl.ctrlMap.GetPublicKeyPath())

}
