package servicetype

//ServiceType Services must identify as an API, APX, APP
//API: Lowest logic layer. API can not call another API.
//APX: Workflow Executor
//APP: Presentation layer.
type Enum int

const (
	API Enum = iota
	APP
	ANY
	APX
)

var servicetypes = [...]string{
	"API",
	"APP",
	"ANY",
	"APX"}

func (s Enum) String() string {
	return servicetypes[s]
}
