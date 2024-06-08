package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/ui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	window := a.NewWindow("The Vault")
	var masterPassword string

	getMasterPassword(window, func(password string) {
		masterPassword = password

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
				ui.NewSettingsContainer(window),
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
	})

	window.ShowAndRun()
}

func getMasterPassword(window fyne.Window, callback func(string)) {
	// Check if the environment variable exists
	masterPassword, exists := os.LookupEnv("MASTER_PASSWORD")
	if exists {
		// Environment variable exists, use its value
		callback(masterPassword)
	} else {
		// Environment variable does not exist, prompt the user to set it using a form
		formItems := []*widget.FormItem{
			widget.NewFormItem("Master Password:", widget.NewPasswordEntry()),
		}
		newPasswordForm := dialog.NewForm("Set Master Password", "Set", "Cancel", formItems, func(b bool) {
			if b {
				masterPassword = formItems[0].Widget.(*widget.Entry).Text
				// Set the environment variable system-wide
				err := setMasterPasswordEnvVar(masterPassword)
				if err != nil {
					fmt.Println("Failed to set MASTER_PASSWORD environment variable:", err)
				}
				callback(masterPassword)
			} else {
				fmt.Println("Master password not set. Exiting...")
				os.Exit(0)
			}
		}, window)

		// Set the size and position of the dialog box
		newPasswordForm.Resize(fyne.NewSize(400, 200))
		window.CenterOnScreen()
		window.Resize(fyne.NewSize(400, 200))

		newPasswordForm.Show()
	}
}
func setMasterPasswordEnvVar(masterPassword string) error {
	// Set the environment variable system-wide
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("setx", "MASTER_PASSWORD", masterPassword)
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("echo 'export MASTER_PASSWORD=%s' >> ~/.bashrc", masterPassword))
	}
	return cmd.Run()
}
