package bodies

import "strings"

type MenuItem struct {
	ID        string
	Text      string
	Enabled   bool
	Hidden    bool
	Link      string
	IsActive  bool
	CanCreate bool
	Children  []MenuItem `json:",omitempty"`
}

func NewItem(id, link, text string, creates bool, children []MenuItem) MenuItem {
	shortName := getUniqueName(text)
	result := MenuItem{
		ID:        id,
		Text:      text,
		Enabled:   true,
		Hidden:    false,
		Link:      link,
		IsActive:  false,
		CanCreate: creates,
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
