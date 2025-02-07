package passwords

import "golang.org/x/crypto/bcrypt"

// HashPassword will hash the password and return the hashed password and an error.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsSamePassword returns true if the hash of password is the same as hashedPassword.
func IsSamePassword(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}
