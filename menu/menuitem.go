package menu

import "strings"

type Item struct {
	ID       string
	Text     string
	Link     string
	Active   bool
	Enabled  bool
	Hidden   bool
	Children []Item `json:",omitempty"`
}

func NewItem(id, link, text string, children []Item) Item {
	shortName := getUniqueName(text)
	result := Item{
		ID:      id,
		Text:    text,
		Enabled: true,
		Hidden:  false,
		Link:    link,
		Active:  false,
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
