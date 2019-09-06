package xontrols

//Queries signals that the path can accept querystrings
type Queries interface {
	Nomad
	AcceptsQuery() map[string]string
}
