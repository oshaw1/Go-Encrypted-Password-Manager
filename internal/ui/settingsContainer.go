package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewSettingsContainer() *fyne.Container {

	settingsContainer := container.NewVBox()
	return settingsContainer
}
