package ui

import (
	"fmt"
	"os"
	"os/exec"
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
	masterPassword := os.Getenv("MASTER_PASSWORD")

	pathToPasswordFile := "data/passwords.json"

	if !passwords.CheckPasswordFileExistsInDataDirectory(pathToPasswordFile) {
		err := passwords.InitializePasswordManager(masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Println("Failed to initialize password data:", err)
			return
		}
		fmt.Println("Password data initialized successfully")
	}

	passwordVaultScrollContainer := container.NewVScroll(
		NewPasswordVaultContainer(pathToPasswordFile, masterPassword, window),
	)

	contentContainer := container.NewStack(
		passwordVaultScrollContainer,
	)

	vaultNavButton := widget.NewButton("The Vault", func() {
		passwordVaultScrollContainer := container.NewVScroll(
			NewPasswordVaultContainer(pathToPasswordFile, masterPassword, window),
		)
		contentContainer.Objects = []fyne.CanvasObject{passwordVaultScrollContainer}
		contentContainer.Refresh()
	})

	settingsNavButton := widget.NewButton("Settings", func() {
		settingsScrollContainer := container.NewVScroll(
			NewSettingsContainer(window),
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
					refreshPasswordCards(pathToPasswordFile, masterPassword, window, vaultContentContainer)
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

	if os.Getenv("MASTER_PASSWORD") == "" {
		getMasterPassword(window, func(password string) {
			masterPassword = password
			window.Close()
		})
	}

	window.ShowAndRun()
}

func getMasterPassword(window fyne.Window, callback func(string)) {
	formItems := []*widget.FormItem{
		widget.NewFormItem("Set Master Password: ", widget.NewPasswordEntry()),
	}
	newPasswordForm := dialog.NewForm("Please set a master password.", "Set", "Cancel", formItems, func(b bool) {
		if b {
			masterPassword := formItems[0].Widget.(*widget.Entry).Text
			confirmDialog := dialog.NewConfirm("Confirm master password", "Are you sure you want this to be your master password? It is very important!", func(confirmed bool) {
				if confirmed {
					info := dialog.NewCustomWithoutButtons("Thankyou", widget.NewLabel("Thankyou, The app will now close. Please re-open it for the changes to take effect"), window)
					info.Show()
					time.Sleep(time.Duration(5) * time.Second)
					cmd := exec.Command("setx", "MASTER_PASSWORD", masterPassword)
					err := cmd.Run()
					if err != nil {
						fmt.Println("Failed to set MASTER_PASSWORD system environment variable:", err)
						os.Exit(1)
					}
					callback(masterPassword)
				} else {
					// user canceled the confirmation, show the password form again
					getMasterPassword(window, callback)
				}
			}, window)
			confirmDialog.Show()
		} else {
			fmt.Println("Master password not set. Exiting...")
			os.Exit(0)
		}
	}, window)

	// Set the size and position of the dialog box
	newPasswordForm.Resize(fyne.NewSize(400, 200))
	newPasswordForm.Show()
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(1000, 1000))
}
