package gear

import "golang.org/x/crypto/bcrypt"

func CryptPassword(text string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

func CheckPassword(hashed string, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(text))
	if err != nil {
		return false
	}
	return true
}
