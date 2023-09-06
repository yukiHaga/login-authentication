package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// 暗号(Hash化)
// Encryptは暗号化するという意味の動詞
func EncryptPassword(password string) (string, error) {
	// GenerateFromPassword は、与えられたコストでパスワードの bcrypt ハッシュを返します。
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
