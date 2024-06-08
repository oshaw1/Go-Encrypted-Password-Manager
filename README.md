# Go-Encrypted-Password-Manager

this is a full stack encrypted password manager written entirely in Golang.

Fyne-io for the frontend

it uses AES encryption


## How to package
to package the application to a useful executable :

- navigate to cmd 
- `fyne package -os windows -icon Logo.png`
- move the executable to the root dir "Encrypted-Password-Manager/"
- right click, click create a shortcut, move shortcut to desktop
- encrypt away...

the first time you launch this app it will prompt you to set a master password, this uses environment variables and all encryption is tied to this password if you change it the passwords will not unencrypt until it is the same again.