package tools

import "golang.org/x/crypto/bcrypt"

// HashPassword hash a given password
func HashPassword(pwd string) (string, error) {

	pwdHashed, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		return "", err
	}

	return string(pwdHashed), nil
}

// ComparePassword compare chiper password with entered password making sure are equal
func ComparePassword(hashedPWD, pwdEntered string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPWD), []byte(pwdEntered))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
