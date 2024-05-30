package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) ([]byte, error) {
	const op = "Hash password"
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		fmt.Errorf("failed to generate passwordHash hash", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return passHash, nil
}
