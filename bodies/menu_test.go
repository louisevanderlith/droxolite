package bodies

import "testing"

func TestMenu_CorrectStructure(t *testing.T) {
	var m Menu

	groupa := NewMenuGroup("Home")
	groupa.AddItem("/a", "Going A", nil)
	bchildren := []menuItem{
		newItem("b1", "/b/1", "Going B1", nil),
		newItem("b2", "/b/2", "Going B2", nil),
	}
	groupa.AddItem("/b", "B has Children", bchildren)

	groupb := NewMenuGroup("Nowhere")
	groupb.AddItem("/no", "Going Nowhere", nil)

	m = append(m, groupa)
	m = append(m, groupb)

	for _, v := range m {
		t.Logf("Label %s", v.Label)
		for _, item := range v.Items {
			t.Logf("Item: %+v", item)

			for _, child := range item.Children {
				t.Logf("Child %v", child)
			}
		}
	}

	t.Fail()
}
