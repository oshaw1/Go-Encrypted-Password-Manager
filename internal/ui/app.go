package ui

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

var (
	AuthProgressBar *widget.ProgressBar
	AuthTimer       *time.Ticker
	RemainingTime   time.Duration
	authed          bool = false
)

func NewApp() {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	window := a.NewWindow("The Vault")
	pathToPasswordFile := "data/passwords.json"

	if !passwords.CheckPasswordFileExistsInDataDirectory(pathToPasswordFile) {
		getMasterPassword(window, pathToPasswordFile)
	}

	passwordVaultScrollContainer := container.NewVScroll(
		NewPasswordVaultContainer(pathToPasswordFile, "", window),
	)

	contentContainer := container.NewStack(
		passwordVaultScrollContainer,
	)

	vaultNavButton := widget.NewButton("The Vault", func() {
		passwordVaultScrollContainer := container.NewVScroll(
			NewPasswordVaultContainer(pathToPasswordFile, "", window),
		)
		contentContainer.Objects = []fyne.CanvasObject{passwordVaultScrollContainer}
		contentContainer.Refresh()
	})

	settingsNavButton := widget.NewButton("Settings", func() {
		settingsScrollContainer := container.NewVScroll(
			NewSettingsContainer(pathToPasswordFile, window),
		)
		contentContainer.Objects = []fyne.CanvasObject{settingsScrollContainer}
		contentContainer.Refresh()
	})

	AuthProgressBar = widget.NewProgressBar()
	AuthProgressBar.TextFormatter = func() string {
		minutes := int(RemainingTime.Minutes())
		seconds := int(RemainingTime.Seconds()) % 60
		return fmt.Sprintf("%d:%02d", minutes, seconds)
	}
	AuthTimer = time.NewTicker(time.Second)
	RemainingTime = 5 * time.Minute

	navContainer := container.NewVBox(
		vaultNavButton,
		settingsNavButton,
		layout.NewSpacer(),
		AuthProgressBar,
	)

	go func() {
		for range AuthTimer.C {
			if authed {
				RemainingTime -= time.Second
				progress := float64(RemainingTime) / float64(5*time.Minute)
				AuthProgressBar.SetValue(progress)

				if RemainingTime <= 0 {
					authed = false
					AuthProgressBar.SetValue(1.0)
					refreshPasswordCards(pathToPasswordFile, "", window, vaultContentContainer)
					return
				}
			}
		}
	}()

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

func getMasterPassword(window fyne.Window, pathToPasswordFile string) {
	formItems := []*widget.FormItem{
		widget.NewFormItem("Set Master Password: ", widget.NewPasswordEntry()),
		widget.NewFormItem("Confirm Password: ", widget.NewPasswordEntry()),
	}

	newPasswordForm := dialog.NewForm("Set Master Password", "Set", "Cancel", formItems, func(b bool) {
		if !b {
			fmt.Println("Master password not set. Exiting...")
			os.Exit(0)
		}

		password := formItems[0].Widget.(*widget.Entry).Text
		confirmPassword := formItems[1].Widget.(*widget.Entry).Text

		if password != confirmPassword {
			dialog.ShowError(fmt.Errorf("passwords do not match"), window)
			getMasterPassword(window, pathToPasswordFile)
			return
		}

		confirmDialog := dialog.NewConfirm("Confirm master password",
			"Are you sure you want this to be your master password? It is very important!",
			func(confirmed bool) {
				if confirmed {
					err := passwords.InitializePasswordManager(password, pathToPasswordFile)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to initialize password manager: %v", err), window)
						return
					}

					info := dialog.NewInformation("Setup Complete",
						"Master password has been set successfully. The app will now close. Please restart it.", window)
					info.Show()
					time.Sleep(time.Duration(5) * time.Second)
					window.Close()
				} else {
					getMasterPassword(window, pathToPasswordFile)
				}
			}, window)
		confirmDialog.Show()
	}, window)

	newPasswordForm.Resize(fyne.NewSize(400, 200))
	newPasswordForm.Show()
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(1000, 1000))
}
