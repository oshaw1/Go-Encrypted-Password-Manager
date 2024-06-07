package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/ui"
)

func main() {
	masterPassword := "master-password" // change to config later
	pathToPasswordFile := "data/passwords.json"

	//initialise the required files
	if !passwords.CheckPasswordFileExistsInDataDirectory(pathToPasswordFile) {
		err := passwords.InitializePasswordManager(masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Println("Failed to initialize password data:", err)
			return
		}
		fmt.Println("Password data initialized successfully")
	}

	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	window := a.NewWindow("The Vault")

	passwordVaultScrollContainer := container.NewVScroll(
		ui.NewPasswordVaultContainer(pathToPasswordFile, masterPassword, window),
	)

	contentContainer := container.NewStack(
		passwordVaultScrollContainer,
	)

	vaultNavButton := widget.NewButton("The Vault", func() {
		passwordVaultScrollContainer := container.NewVScroll(
			ui.NewPasswordVaultContainer(pathToPasswordFile, masterPassword, window),
		)
		contentContainer.Objects = []fyne.CanvasObject{passwordVaultScrollContainer}
		contentContainer.Refresh()
	})

	settingsNavButton := widget.NewButton("Settings", func() {
		settingsScrollContainer := container.NewVScroll(
			ui.NewSettingsContainer(),
		)
		contentContainer.Objects = []fyne.CanvasObject{settingsScrollContainer}
		contentContainer.Refresh()
	})

	navContainer := container.NewVBox(
		vaultNavButton,
		settingsNavButton,
	)

	splitContainer := container.NewHSplit(
		navContainer,
		contentContainer,
	)
	splitContainer.SetOffset(0.20)

	window.SetContent(splitContainer)
	window.Resize(fyne.NewSize(1000, 1000))
	window.CenterOnScreen()

	window.ShowAndRun()
}
