package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type EntradaConEvento struct {
	widget.Entry
	TeclaEvento func(key *fyne.KeyEvent)
}

func CrearEntradaConEvento() *EntradaConEvento {
	entry := &EntradaConEvento{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (ece *EntradaConEvento) TypedKey(key *fyne.KeyEvent) {
	ece.Entry.TypedKey(key)
	if ece.TeclaEvento != nil {
		ece.TeclaEvento(key)
	}
}
