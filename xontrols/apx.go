package xontrols

import "github.com/louisevanderlith/droxolite/context"

//APX can only have POST Requests, They should be used to call other API's or workflows.
//Default should start the workflow process
//Send must be used to progress forward
type APX interface {
	Nomad
	Send(context.Requester) (int, interface{})
}
