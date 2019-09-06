package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Deleteable handles controls that can delete content [DELETE]
type Deleteable interface {
	Nomad
	Delete(context.Requester) (int, interface{})
}
