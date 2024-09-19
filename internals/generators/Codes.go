package generators

import (
	"crypto/rand"
	"math/big"
)

// GenerateAlphaNumericCode generates an alpha numeric code by the given length
func GenerateAlphaNumericCode(length int) (string, error) {

	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

// Generate6DigitCode generates a numeric 6 digit code
func Generate6DigitCode() (int, error) {
	max := big.NewInt(1000000) // 10^6
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	code := int(n.Int64())
	// Ensure the code is 6 digits by re-rolling if it's less than 100000
	if code < 100000 {
		return Generate6DigitCode()
	}

	return code, nil
}
