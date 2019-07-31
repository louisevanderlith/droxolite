package bodies

import (
	"strings"
)

type Menu struct {
	menu map[string][]MenuItem
}

/*
type MenuGroup struct {
	Label string
	Items []menuItem
}*/

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

/*
func (m *MenuGroup) AddItem(link, text string, children []menuItem) {
	id := fmt.Sprintf("m%v", m.Len())
	item := newItem(id, link, text, children)

	m.Items = append(m.Items, item)
}

func (m *MenuGroup) AddItemWithID(id, link, text string, children []menuItem) {
	item := newItem(id, link, text, children)

	m.Items = append(m.Items, item)
}
*/
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

type MenuItem struct {
	ID       string
	Text     string
	Enabled  bool
	Hidden   bool
	Link     string
	IsActive bool
	Children []MenuItem `json:",omitempty"`
}

func NewItem(id, link, text string, children []MenuItem) MenuItem {
	shortName := getUniqueName(text)
	result := MenuItem{
		ID:       id,
		Text:     text,
		Enabled:  true,
		Hidden:   false,
		Link:     link,
		IsActive: false,
	}

	if link == "#" {
		result.Link += shortName
	}

	if children != nil {
		result.Children = children
	}

	return result
}

func getUniqueName(raw string) string {
	if len(raw) == 0 {
		return "Home"
	}

	return strings.ToLower(strings.Replace(raw, " ", "", -1))
}
