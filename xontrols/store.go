package xontrols

import "github.com/louisevanderlith/droxolite/context"

//StoreController is used to access Storage API's
type Store interface {
	Nomad
	GetOne(context.Requester) (int, interface{})
	Create(context.Requester) (int, interface{})
	Update(context.Requester) (int, interface{})
	Delete(context.Requester) (int, interface{})
}
