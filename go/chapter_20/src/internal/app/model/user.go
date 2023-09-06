package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/yukiHaga/web_server/src/pkg/henagin/auth"
	"github.com/yukiHaga/web_server/src/pkg/henagin/db"
)

type UserId int64

// PasswordとConfirmationはRailsの家蔵属性として入っていたから、一応入れておいた
type User struct {
	Id             UserId
	Name           string
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUser(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

// 既存のユーザーが存在するかの処理もかく
func (user *User) SignUp(password, passwordConfirmation string) error {
	if password != passwordConfirmation {
		return fmt.Errorf("not equal password: %s, %s", password, passwordConfirmation)
	}

	db, err := db.NewDB()
	if err != nil {
		return err
	}
	defer db.Close()

	registeredUser, err := FindUserByEmail(user.Email)
	if err != nil {
		return err
	}

	if registeredUser.Id > 0 {
		return errors.New("already register")
	}

	stmtIns, err := db.Prepare("INSERT INTO users (name, email, password_digest, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	passwordDigest, err := auth.EncryptPassword(password)
	if err != nil {
		return err
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	result, err := stmtIns.Exec(user.Name, user.Email, passwordDigest, createdAt, updatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = UserId(id)
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	user.PasswordDigest = passwordDigest

	return nil
}

func (user *User) Login(password string) error {
	if err := auth.CompareHashAndPassword(user.PasswordDigest, password); err != nil {
		return err
	}

	return nil
}

func FindUserByEmail(email string) (*User, error) {
	db, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// idでselectしているので、一人のユーザーしかヒットしない
	rows, err := db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := User{}
	// 結果がないなら、ここのループに入らない
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func FindUserById(userId string) (*User, error) {
	db, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// idでselectしているので、一人のユーザーしかヒットしない
	rows, err := db.Query("SELECT * FROM users WHERE id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
