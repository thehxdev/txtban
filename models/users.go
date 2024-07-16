package models

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/thehxdev/txtban/tberr"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID      int
	UUID    string
	PHash   string
	AuthKey string
}

func CreateAuthKey(uuid, password string) string {
	s := []byte(uuid + password)
	hash := sha256.New()
	hash.Write(s)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (c *Conn) CreateUser(uuid, password, authKey string) error {
	stmt := `INSERT INTO users (uuid, phash, authKey) VALUES (?, ?, ?)`

	phash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return tberr.New(err.Error())
	}

	_, err = c.DB.Exec(stmt, uuid, string(phash), authKey)
	if err != nil {
		return tberr.New(err.Error())
	}

	return nil
}

func (c *Conn) AuthenticateByPassword(uuid, pass string) (*User, error) {
	u := &User{}
	stmt := `SELECT id, uuid, phash, authKey FROM users WHERE uuid = ?`

	err := c.DB.QueryRow(stmt, uuid).Scan(&u.ID, &u.UUID, &u.PHash, &u.AuthKey)
	if err != nil {
		return nil, tberr.New(err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PHash), []byte(pass))
	if err != nil {
		return nil, tberr.New(err.Error(), "double check your password")
	}

	return u, nil
}

func (c *Conn) AuthenticateByAuthKey(authKey string) (*User, error) {
	u := new(User)
	stmt := `SELECT id, uuid, phash, authKey FROM users WHERE authKey = ?`

	err := c.DB.QueryRow(stmt, authKey).Scan(&u.ID, &u.UUID, &u.PHash, &u.AuthKey)
	if err != nil {
		return nil, tberr.New(err.Error())
	}

	return u, nil
}

func (c *Conn) DeleteUser(id int) error {
	stmt1 := `DELETE FROM users WHERE id = ?`
	stmt2 := `DELETE FROM txts WHERE uid = ?`

	_, err := c.DB.Exec(stmt1, id)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(stmt2, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) UpdateUserPassword(id int, password, authKey string) error {
	stmt := `UPDATE users SET phash = ?, authKey = ? WHERE id = ?`

	phash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return tberr.New(err.Error())
	}

	_, err = c.DB.Exec(stmt, phash, authKey, id)
	if err != nil {
		return tberr.New(err.Error())
	}

	return nil
}
