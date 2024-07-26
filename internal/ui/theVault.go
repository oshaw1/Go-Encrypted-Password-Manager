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

var (
	vaultContentContainer = container.NewVBox()
)

func NewPasswordVaultContainer(pathToPasswordFile string, masterPassword string, window fyne.Window) *fyne.Container {
	titleLabel := widget.NewLabel("Welcome To The Vault")
	masterPasswordEntry := widget.NewPasswordEntry()
	addPasswordButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if masterPasswordEntry.Text == "enter master password" {
			dialog.NewInformation("Please enter masterpassword", "Dismiss", window).Show()
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
			authed = true
			RemainingTime = 5 * time.Minute
			AuthProgressBar.SetValue(1.0)
		} else {
			fmt.Print(masterPasswordEntry.Text)
			masterPasswordEntry.SetText("Master password incorrect")
		}
	})

	masterPasswordContainer := container.NewBorder(nil, nil, nil, unlockButton, masterPasswordEntry)

	vaultContentContainer = container.NewVBox(masterPasswordContainer)
	var passwordCards fyne.CanvasObject
	if !authed {
		passwordCards = createPasswordCard(pathToPasswordFile, masterPassword, window, vaultContentContainer)
	} else {
		passwordCards = createDecryptedPasswordCard(pathToPasswordFile, masterPassword, window, vaultContentContainer)
	}
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
	for index, password := range passwordData.Passwords {
		hyperLinkLabel := widget.NewLabel("Link:  " + randomStars())
		usernameLabel := widget.NewLabel("Username / Account:  " + randomStars())
		passwordLabel := widget.NewLabel("Password:  " + randomStars())
		labelContainer := container.NewVBox(hyperLinkLabel, usernameLabel, passwordLabel)
		passwordCard := widget.NewCard(password.Title, "", labelContainer)

		passwordButtonsContainer := container.NewVBox(
			createEditButton(pathToPasswordFile, masterPassword, index, window, vaultContentContainer),
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

	for index, password := range passwordData.Passwords {
		passwordID := password.ID
		retrievedLink, retrievedUsername, retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Printf("Failed to retrieve password for ID %s: %v\n", passwordID, err)
			dialog.NewError(err, window)
		}

		hyperLinkLabel := widget.NewHyperlink("Link:  "+retrievedLink, nil)
		hyperLinkLabel.OnTapped = func() {
			window.Clipboard().SetContent(retrievedLink)
			showPopupBelow(hyperLinkLabel, "Link copied to clipboard", window)
		}

		usernameLabel := widget.NewHyperlink("Username / Account:  "+retrievedUsername, nil)
		usernameLabel.OnTapped = func() {
			window.Clipboard().SetContent(retrievedUsername)
			showPopupBelow(usernameLabel, "Username copied to clipboard", window)
		}

		passwordLabel := widget.NewHyperlink("Password:  "+retrievedPassword, nil)
		passwordLabel.OnTapped = func() {
			window.Clipboard().SetContent(retrievedPassword)
			showPopupBelow(usernameLabel, "Password copied to clipboard", window)
		}

		labelContainer := container.NewVBox(hyperLinkLabel, usernameLabel, passwordLabel)
		passwordCard := widget.NewCard(password.Title, "", labelContainer)

		passwordButtonsContainer := container.NewVBox(
			createEditButton(pathToPasswordFile, masterPassword, index, window, vaultContentContainer),
			createDeleteButton(pathToPasswordFile, masterPassword, password.ID, window, vaultContentContainer),
		)

		passwordCardContent := container.NewBorder(nil, nil, nil, passwordButtonsContainer, passwordCard)
		passwordCards = append(passwordCards, passwordCardContent)
	}

	return container.New(layout.NewVBoxLayout(), passwordCards...)
}

func showPopupBelow(hyperlink fyne.Widget, message string, window fyne.Window) {
	popupLabel := widget.NewLabel(message)
	popupLabel.Alignment = fyne.TextAlignCenter
	popup := widget.NewPopUp(popupLabel, window.Canvas())

	widgetPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(hyperlink)
	widgetSize := hyperlink.Size()

	popupPos := fyne.NewPos(widgetPos.X, widgetPos.Y+widgetSize.Height)
	popup.ShowAtPosition(popupPos)

	go func() {
		time.Sleep(1 * time.Second)
		popup.Hide()
	}()
}

func createEditButton(pathToPasswordFile string, masterPassword string, passwordIndex int, window fyne.Window, vaultContentContainer *fyne.Container) fyne.Widget {
	editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		if !authed {
			dialog.NewInformation("Please Unlock the Vault", "Please unlock the vault before editing passwords...", window).Show()
			return
		}
		data, err := os.ReadFile(pathToPasswordFile)
		if err != nil {
			fmt.Printf("Failed to read password file: %v\n", err)
			dialog.NewError(err, window)
			return
		}

		var passwordData model.PasswordData
		err = json.Unmarshal(data, &passwordData)
		if err != nil {
			fmt.Printf("Failed to parse JSON data: %v\n", err)
			dialog.NewError(err, window)
			return
		}

		password := passwordData.Passwords[passwordIndex]

		formItems := []*widget.FormItem{
			widget.NewFormItem("Title", widget.NewEntry()),
			widget.NewFormItem("Link", widget.NewEntry()),
			widget.NewFormItem("Username", widget.NewEntry()),
			widget.NewFormItem("Password", widget.NewPasswordEntry()),
		}

		editPasswordForm := dialog.NewForm("Edit Password", "Save", "Cancel", formItems, func(ok bool) {
			if ok {
				newTitle := formItems[0].Widget.(*widget.Entry).Text
				newLink := formItems[1].Widget.(*widget.Entry).Text
				newUsername := formItems[2].Widget.(*widget.Entry).Text
				newPassword := formItems[3].Widget.(*widget.Entry).Text

				err := passwords.EditPassword(password.ID, newTitle, newLink, newUsername, newPassword, masterPassword, pathToPasswordFile)
				if err != nil {
					dialog.NewError(err, window)
					fmt.Println("Error editing password:", err)
				} else {
					if !authed {
						refreshPasswordCards(pathToPasswordFile, masterPassword, window, vaultContentContainer)
					} else {
						refreshPasswordCardsToDecrypedVersion(pathToPasswordFile, masterPassword, window, vaultContentContainer)
					}
				}
			}
		}, window)

		editPasswordForm.Resize(fyne.NewSize(600, 200))

		passwordID := password.ID
		retrievedLink, retrievedUsername, retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, pathToPasswordFile)
		if err != nil {
			dialog.NewError(err, window)
			fmt.Println("Error retrieving password:", err)
		}

		// Set the initial values of the form fields
		formItems[0].Widget.(*widget.Entry).SetText(password.Title)
		formItems[1].Widget.(*widget.Entry).SetText(retrievedLink)
		formItems[2].Widget.(*widget.Entry).SetText(retrievedUsername)
		formItems[3].Widget.(*widget.Entry).SetText(retrievedPassword)

		editPasswordForm.Show()
	})
	return editButton
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
