package mix

type Bag interface {
	SetValue(name string, val interface{})
	Values() map[string]interface{}
}

func NewBag() Bag {
	return make(bag)
}

func NewDataBag(data interface{}) Bag {
	b := NewBag()
	b.SetValue("Data", data)

	return b
}

type bag map[string]interface{}

func (b bag) SetValue(name string, val interface{}) {
	b[name] = val
}

func (b bag) Values() map[string]interface{} {
	return b
}
