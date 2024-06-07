package model

type PasswordData struct {
	MasterPasswordHash string `json:"master_password_hash"`
	Salt               string `json:"salt"`
	Passwords          []struct {
		ID                string `json:"id"`
		Title             string `json:"title"`
		Hyperlink         string `json:"hyperlink"`
		EncryptedPassword string `json:"encrypted_password"`
	} `json:"passwords"`
}
