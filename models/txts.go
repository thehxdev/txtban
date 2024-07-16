package models

import (
	"time"

	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/tberr"
	"github.com/thehxdev/txtban/tbrandom"
)

type Txt struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Content string     `json:"content,omitempty"`
	Created *time.Time `json:"created"`
	UID     int        `json:"uid,omitempty"`
}

func (c *Conn) CreateTxt(uid int, name, content string) (string, error) {
	idLen := tbrandom.GenRandNum(4, viper.GetInt("limits.maxTxtIdLen"))

	id := tbrandom.GenRandString(idLen)
	stmt := `INSERT INTO txts (id, name, content, uid) VALUES (?, ?, ?, ?)`

	res, err := c.DB.Exec(stmt, id, name, content, uid)
	if err != nil {
		return "", tberr.New(err.Error())
	}

	_, err = res.LastInsertId()
	if err != nil {
		return "", tberr.New(err.Error())
	}

	return id, nil
}

func (c *Conn) GetTxtByName(uid int, name string) (*Txt, error) {
	txt := new(Txt)
	stmt := `SELECT * FROM txts WHERE uid = ? AND name = ?`

	err := c.DB.QueryRow(stmt, uid, name).Scan(&txt.ID, &txt.Name, &txt.Content, &txt.Created, &txt.UID)
	if err != nil {
		return nil, tberr.New(err.Error())
	}

	return txt, nil
}

func (c *Conn) GetTxtById(txtid string) (*Txt, error) {
	txt := new(Txt)
	stmt := `SELECT * FROM txts WHERE id = ?`

	err := c.DB.QueryRow(stmt, txtid).Scan(&txt.ID, &txt.Name, &txt.Content, &txt.Created, &txt.UID)
	if err != nil {
		return nil, tberr.New(err.Error())
	}

	return txt, nil
}

func (c *Conn) GetTxtContentById(id string) (string, error) {
	var s string
	stmt := `SELECT content FROM txts WHERE id = ?`

	err := c.DB.QueryRow(stmt, id).Scan(&s)
	if err != nil {
		return "", tberr.New(err.Error())
	}

	return s, nil
}

func (c *Conn) DeleteTxt(id string) error {
	stmt := `DELETE FROM txts WHERE id = ?`

	_, err := c.DB.Exec(stmt, id)
	if err != nil {
		return tberr.New(err.Error())
	}

	return nil
}

func (c *Conn) GetAllTxts(uid int) ([]*Txt, error) {
	txts := []*Txt{}
	stmt := `SELECT id, name, created FROM txts WHERE uid = ? ORDER BY created DESC`

	rows, err := c.DB.Query(stmt, uid)
	if err != nil {
		return nil, tberr.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		txt := &Txt{}
		err := rows.Scan(&txt.ID, &txt.Name, &txt.Created)
		if err != nil {
			return nil, tberr.New(err.Error())
		}

		txts = append(txts, txt)
	}

	return txts, nil
}

func (c *Conn) ChangeTxtContent(txtid string, content string) error {
	stmt := `UPDATE txts SET content = ? WHERE id = ?`

	_, err := c.DB.Exec(stmt, content, txtid)
	if err != nil {
		return tberr.New(err.Error())
	}

	return nil
}

func (c *Conn) ChangeTxtId(txtid string) (string, error) {
	_, err := c.GetTxtById(txtid)
	if err != nil {
		return "", err
	}

	newId := tbrandom.GenRandString(tbrandom.GenRandNum(4, viper.GetInt("limits.maxTxtIdLen")))
	stmt := `UPDATE txts SET id = ? WHERE id = ?`

	_, err = c.DB.Exec(stmt, newId, txtid)
	if err != nil {
		return "", tberr.New(err.Error())
	}

	return newId, nil
}

func (c *Conn) ChangeTxtName(txtid, name string) error {
	stmt := `UPDATE txts SET name = ? WHERE id = ?`

	_, err := c.DB.Exec(stmt, name, txtid)
	if err != nil {
		return tberr.New(err.Error())
	}

	return nil
}
