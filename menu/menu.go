package menu

type Menu struct {
	Items []Item
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) AddItem(itm Item) {
	m.Items = append(m.Items, itm)
}

func (m *Menu) SetActive(link string) bool {
	foundActive := false

	for _, item := range m.Items {
		item.Active = item.Link == link

		if !foundActive && item.Active {
			foundActive = true
		}

		for _, child := range item.Children {
			child.Active = child.Link == link
		}
	}

	return foundActive
}
