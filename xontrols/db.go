package xontrols

//DBController is used to access Storage API's
type DBController interface {
	Create()
	Read()
	ReadOne()
	Update()
	Delete()
}
