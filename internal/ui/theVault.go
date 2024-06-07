package ui

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/model"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func NewPasswordVaultContainer(pathToPasswordFile string, masterPassword string, window fyne.Window) *fyne.Container {
	titleLabel := widget.NewLabel("Welcome To The Vault")
	vaultContentContainer := container.NewVBox()
	masterPasswordEntry := widget.NewPasswordEntry()
	addPasswordButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if masterPasswordEntry.Text == "enter master password" {
			dialog.NewInformation("Please enter masterpassword", "Dismiss", window)
			return
		} else {
			formItems := []*widget.FormItem{
				widget.NewFormItem("Title", widget.NewEntry()),
				widget.NewFormItem("Link", widget.NewEntry()),
				widget.NewFormItem("Username: ", widget.NewEntry()),
				widget.NewFormItem("Password", widget.NewPasswordEntry()),
			}
			newPasswordForm := dialog.NewForm("Add new password", "Add Password", "Dismiss", formItems, func(b bool) {
				if b {
					title := formItems[0].Widget.(*widget.Entry).Text
					link := formItems[1].Widget.(*widget.Entry).Text
					username := formItems[2].Widget.(*widget.Entry).Text
					password := formItems[3].Widget.(*widget.Entry).Text

					err := passwords.StorePassword(title, link, username, password, masterPassword, pathToPasswordFile)
					if err != nil {
						dialog.NewError(err, window)
						fmt.Println("Error storing password:", err)
					} else {
						refreshPasswordCards(pathToPasswordFile, masterPassword, window, vaultContentContainer)
					}
				}
			}, window)
			newPasswordForm.Show()
		}
	})
	headerContainer := container.NewHBox(titleLabel, layout.NewSpacer(), addPasswordButton)

	titleBorder := container.NewBorder(headerContainer, nil, nil, nil)
	masterPasswordEntry.SetPlaceHolder("enter master password")

	unlockButton := widget.NewButtonWithIcon("Unlock", theme.ListIcon(), func() {
		data, err := os.ReadFile(pathToPasswordFile)
		if err != nil {
			fmt.Printf("Failed to read password file: %v\n", err)
		}

		var passwordData model.PasswordData
		err = json.Unmarshal(data, &passwordData)
		if err != nil {
			fmt.Printf("Failed to parse JSON data: %v\n", err)
		}
		if encryption.HashMasterPassword(masterPasswordEntry.Text) == encryption.HashMasterPassword(masterPassword) {
			refreshPasswordCardsToDecrypedVersion(pathToPasswordFile, masterPassword, window, vaultContentContainer)
			go func() {
				time.Sleep(5 * time.Minute)
				refreshPasswordCards(pathToPasswordFile, masterPassword, window, vaultContentContainer)
			}()
		} else {
			fmt.Print(masterPasswordEntry.Text)
			masterPasswordEntry.SetText("Master password incorrect")
		}
	})

	masterPasswordContainer := container.NewBorder(nil, nil, nil, unlockButton, masterPasswordEntry)

	vaultContentContainer = container.NewVBox(masterPasswordContainer)
	passwordCards := createPasswordCard(pathToPasswordFile, masterPassword, window, vaultContentContainer)
	vaultContentContainer.Add(passwordCards)

	theVault := container.NewVBox(
		titleBorder,
		vaultContentContainer,
	)

	paddedVault := container.NewPadded(theVault)
	return paddedVault
}

func createPasswordCard(pathToPasswordFile string, masterPassword string, window fyne.Window, vaultContentContainer *fyne.Container) fyne.CanvasObject {
	data, err := os.ReadFile(pathToPasswordFile)
	if err != nil {
		fmt.Printf("Failed to read password file: %v\n", err)
		dialog.NewError(err, window)
		return nil
	}

	var passwordData model.PasswordData
	err = json.Unmarshal(data, &passwordData)
	if err != nil {
		fmt.Printf("Failed to parse JSON data: %v\n", err)
		dialog.NewError(err, window)
	}

	if passwordData.Passwords == nil {
		fmt.Print("No passwords found")
		return nil
	}

	var passwordCards []fyne.CanvasObject
	for _, password := range passwordData.Passwords {
		hyperLinkLabel := widget.NewLabel("Link:  " + randomStars())
		usernameLabel := widget.NewLabel("Username / Account:  " + randomStars())
		passwordLabel := widget.NewLabel("Password:  " + randomStars())
		labelContainer := container.NewVBox(hyperLinkLabel, usernameLabel, passwordLabel)
		passwordCard := widget.NewCard(password.Title, "", labelContainer)

		passwordButtonsContainer := container.NewVBox(
			widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {}),
			createDeleteButton(pathToPasswordFile, masterPassword, password.ID, window, vaultContentContainer),
		)

		passwordCardContent := container.NewBorder(nil, nil, nil, passwordButtonsContainer, passwordCard)
		passwordCards = append(passwordCards, passwordCardContent)
	}

	return container.New(layout.NewVBoxLayout(), passwordCards...)
}

func randomStars() string {
	minLength := 5
	maxLength := 15

	randomLength := rand.Intn(maxLength-minLength+1) + minLength
	asterisks := strings.Repeat("*", randomLength)
	return asterisks
}

func createDecryptedPasswordCard(pathToPasswordFile string, masterPassword string, window fyne.Window, vaultContentContainer *fyne.Container) fyne.CanvasObject {
	data, err := os.ReadFile(pathToPasswordFile)
	if err != nil {
		fmt.Printf("Failed to read password file: %v\n", err)
		dialog.NewError(err, window)
		return nil
	}

	var passwordData model.PasswordData
	err = json.Unmarshal(data, &passwordData)
	if err != nil {
		fmt.Printf("Failed to parse JSON data: %v\n", err)
		dialog.NewError(err, window)
	}

	if passwordData.Passwords == nil {
		fmt.Print("No passwords found")
		return nil
	}

	var passwordCards []fyne.CanvasObject

	for _, password := range passwordData.Passwords {
		passwordID := password.ID
		retrievedLink, retrievedUsername, retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Printf("Failed to retrieve password for ID %s: %v\n", passwordID, err)
			dialog.NewError(err, window)
		}
		hyperLinkLabel := widget.NewLabel("Link:  " + retrievedLink)
		usernameLabel := widget.NewLabel("Username / Account:  " + retrievedUsername)
		passwordLabel := widget.NewLabel("Password:  " + retrievedPassword)
		labelContainer := container.NewVBox(hyperLinkLabel, usernameLabel, passwordLabel)
		passwordCard := widget.NewCard(password.Title, "", labelContainer)

		passwordButtonsContainer := container.NewVBox(
			widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {}),
			createDeleteButton(pathToPasswordFile, masterPassword, password.ID, window, vaultContentContainer),
		)

		passwordCardContent := container.NewBorder(nil, nil, nil, passwordButtonsContainer, passwordCard)
		passwordCards = append(passwordCards, passwordCardContent)

		passwordTitle := password.Title
		fmt.Printf("Password for %s: %s\n", passwordTitle, retrievedPassword)
	}

	return container.New(layout.NewVBoxLayout(), passwordCards...)
}

func createDeleteButton(pathToPasswordFile string, masterPassword string, passwordID string, window fyne.Window, vaultContentContainer *fyne.Container) fyne.Widget {
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		confirmDialog := dialog.NewConfirm("Confirm Delete", "Are you sure you want to delete this password? It will be gone forever...", func(confirmed bool) {
			if confirmed {
				err := passwords.DeletePasswordByID(passwordID, masterPassword, pathToPasswordFile)
				if err != nil {
					dialog.NewError(err, window)
					fmt.Println("Error deleting password:", err)
				}
				refreshPasswordCards(pathToPasswordFile, masterPassword, window, vaultContentContainer)
			}
		}, window)
		confirmDialog.Show()

	})
	return deleteButton
}

func refreshPasswordCards(pathToPasswordFile string, masterPassword string, window fyne.Window, vaultContentContainer *fyne.Container) {
	passwordCards := createPasswordCard(pathToPasswordFile, masterPassword, window, vaultContentContainer)
	vaultContentContainer.Objects = []fyne.CanvasObject{
		vaultContentContainer.Objects[0],
		passwordCards,
	}
	vaultContentContainer.Refresh()
}

func refreshPasswordCardsToDecrypedVersion(pathToPasswordFile string, masterPassword string, window fyne.Window, vaultContentContainer *fyne.Container) {
	passwordCards := createDecryptedPasswordCard(pathToPasswordFile, masterPassword, window, vaultContentContainer)
	vaultContentContainer.Objects = []fyne.CanvasObject{
		vaultContentContainer.Objects[0],
		passwordCards,
	}
	vaultContentContainer.Refresh()
}
