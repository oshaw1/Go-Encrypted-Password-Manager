package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsContainer(window fyne.Window) *fyne.Container {
	masterPassword := widget.NewEntry()
	masterPassword.PlaceHolder = "Please enter a new master password"
	settingsTitle := widget.NewLabel("Settings")
	settingsTitle.Alignment = fyne.TextAlignCenter
	MasterPasswordSetting := container.NewBorder(nil, nil, nil, widget.NewButton("Set Master Password", func() {
		confirmDialog := dialog.NewConfirm("Confirm Setting Master Password", "Are you sure you want to set a new master password? It will break all past passwords... The app will require a restart", func(confirmed bool) {
			if confirmed {
				err := setMasterPasswordEnvVar(masterPassword.Text)
				if err != nil {
					fmt.Println("Failed to set MASTER_PASSWORD environment variable:", err)
				}
				fmt.Print("new master-password set restarting application")
				os.Exit(1)
			}
		}, window)
		confirmDialog.Show()
	}), masterPassword)
	settingsContainer := container.NewVBox(settingsTitle, widget.NewSeparator(), MasterPasswordSetting)
	return settingsContainer
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
