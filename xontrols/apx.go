package xontrols

//APXController can only have POST Requests, They should be used to call other API's or workflows.
type APXController interface {
	Send()
}
