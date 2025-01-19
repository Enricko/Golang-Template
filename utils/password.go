package utils

import "golang.org/x/crypto/bcrypt"

func ComparePasswords(hashedPassword string, plainPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}