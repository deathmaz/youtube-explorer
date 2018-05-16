package ui

import "github.com/jroimartin/gocui"

// func customEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
// switch {
// case ch != 0 && mod == 0:
// v.EditWrite(ch)
// case ch == 'j':
// v.EditWrite('j')
// case key == gocui.KeySpace:
// v.EditWrite(' ')
// case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
// v.EditDelete(true)
// }
// }

func search(g *gocui.Gui, v *gocui.View) error {
	deleteGlobKeybindings(g)

	maxX, maxY := g.Size()
	if v, err := g.SetView(searchView, maxX/2-13, maxY/4, maxX/2+13, maxY/2+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// var DefaultEditor gocui.Editor = gocui.EditorFunc(customEditor)

		v.Editable = true
		v.Wrap = true
		// v.Editor = DefaultEditor

		if _, err := g.SetCurrentView(searchView); err != nil {
			return err
		}
	}

	return nil
}
