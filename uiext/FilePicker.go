package uiext

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func FilePicker(w fyne.Window, callback func(string)) {
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			callback("")
			return
		}
		if f == nil {
			callback("")
			return
		}
		callback(f.URI().Path())
	}, w)
}
