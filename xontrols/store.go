package xontrols

//StoreController is used to access Storage API's
type StoreController interface {
	Create()
	Read()
	ReadOne()
	Update()
	Delete()
}
