package xontrols

import "github.com/louisevanderlith/droxolite/context"

//APXController can only have POST Requests, They should be used to call other API's or workflows.
type APXController interface {
	Send(context.Contexer) error
}
