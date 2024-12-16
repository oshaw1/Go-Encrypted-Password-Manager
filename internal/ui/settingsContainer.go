package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func NewSettingsContainer(pathToPasswordFile string, window fyne.Window) *fyne.Container {
	masterPassword := widget.NewPasswordEntry()
	masterPassword.PlaceHolder = "Please enter a new master password"

	settingsTitle := widget.NewLabel("Settings")
	settingsTitle.Alignment = fyne.TextAlignCenter

	MasterPasswordSetting := container.NewBorder(nil, nil, nil,
		widget.NewButton("Set Master Password", func() {
			confirmPassword := widget.NewPasswordEntry()
			confirmPassword.PlaceHolder = "Confirm new master password"

			content := container.NewVBox(
				confirmPassword,
				widget.NewLabel("Are you sure you want to set a new master password? This will ERASE all past passwords."),
			)

			dialog.NewCustomConfirm(
				"Confirm Setting Master Password",
				"Set",
				"Cancel",
				content,
				func(confirmed bool) {
					if confirmed {
						if masterPassword.Text != confirmPassword.Text {
							dialog.ShowError(fmt.Errorf("passwords do not match"), window)
							return
						}
						passwords.InitializePasswordManager(masterPassword.Text, pathToPasswordFile)
					}
				},
				window,
			).Show()
		}),
		masterPassword,
	)

	settingsContainer := container.NewVBox(
		settingsTitle,
		widget.NewSeparator(),
		MasterPasswordSetting,
	)
	return settingsContainer
}
