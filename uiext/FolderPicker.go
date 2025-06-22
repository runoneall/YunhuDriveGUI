package uiext

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func FolderPicker(w fyne.Window, callback func(string)) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			callback("")
			return
		}
		if uri == nil {
			callback("")
			return
		}
		selectedPath := uri.Path()
		callback(selectedPath)
	}, w)
}
