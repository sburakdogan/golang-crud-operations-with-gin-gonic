package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(plainPassword string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func VerifyPassword(plainPassword string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err
}
