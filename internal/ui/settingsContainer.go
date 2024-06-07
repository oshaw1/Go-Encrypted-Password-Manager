package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/oshaw1/Encrypted-Password-Manager/config"
)

func NewSettingsContainer(window fyne.Window) *fyne.Container {
	pathToMasterPassword := "data/masterPassword.json"

	masterPassword := widget.NewEntry()
	settingsTitle := widget.NewLabel("Settings")
	settingsTitle.Alignment = fyne.TextAlignCenter
	MasterPasswordSetting := container.NewBorder(nil, nil, nil, widget.NewButton("Set Master Password", func() {
		confirmDialog := dialog.NewConfirm("Confirm Setting Master Password", "Are you sure you want to set a new master password? It will break all past passwords...", func(confirmed bool) {
			if confirmed {
				err := config.StoreMasterPassword(masterPassword.Text, pathToMasterPassword)
				if err != nil {
					dialog.NewError(err, window)
					fmt.Println("Error deleting password:", err)
				}
			}
		}, window)
		confirmDialog.Show()
	}), masterPassword)
	settingsContainer := container.NewVBox(settingsTitle, widget.NewSeparator(), MasterPasswordSetting)
	return settingsContainer
}
