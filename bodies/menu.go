package bodies

type Menu struct {
	menu map[string][]MenuItem
}

func NewMenu() *Menu {
	return &Menu{make(map[string][]MenuItem)}
}

func (m *Menu) AddGroup(label string, items []MenuItem) {
	m.menu[label] = append(m.menu[label], items...)
}

func (m *Menu) Len() int {
	return len(m.menu)
}

func (m *Menu) Items() map[string][]MenuItem {
	return m.menu
}

func (m *Menu) SetActive(link string) bool {
	foundActive := false

	for _, items := range m.menu {
		for _, item := range items {
			item.IsActive = item.Link == link

			if !foundActive && item.IsActive {
				foundActive = true
			}

			for _, child := range item.Children {
				child.IsActive = child.Link == link
			}
		}
	}

	return foundActive
}
