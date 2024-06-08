package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsContainer(window fyne.Window) *fyne.Container {
	masterPassword := widget.NewEntry()
	settingsTitle := widget.NewLabel("Settings")
	settingsTitle.Alignment = fyne.TextAlignCenter
	MasterPasswordSetting := container.NewBorder(nil, nil, nil, widget.NewButton("Set Master Password", func() {
		confirmDialog := dialog.NewConfirm("Confirm Setting Master Password", "Are you sure you want to set a new master password? It will break all past passwords...", func(confirmed bool) {
			if confirmed {
				err := os.Setenv("MASTER_PASSWORD", masterPassword.Text)
				if err != nil {
					fmt.Println("Failed to set MASTER_PASSWORD environment variable:", err)
				}
			}
		}, window)
		confirmDialog.Show()
	}), masterPassword)
	settingsContainer := container.NewVBox(settingsTitle, widget.NewSeparator(), MasterPasswordSetting)
	return settingsContainer
}
