package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Store handles controls that can Update [PUT]
type Updateable interface {
	Nomad
	Update(context.Requester) (int, interface{})
}
