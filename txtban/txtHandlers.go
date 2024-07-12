package txtban

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// func (t *Txtban) catHandler(c *fiber.Ctx) error {
// 	authKey, err := getAuthKey(c.GetReqHeaders())
// 	if err != nil {
// 		return err
// 	}
// 	var jdata map[string]string
//
// 	if err := c.BodyParser(&jdata); err != nil {
// 		t.ErrLogger.Println(err.Error())
// 		return err
// 	}
//
// 	user, err := t.Conn.AuthenticateByAuthKey(authKey)
// 	if err != nil {
// 		t.ErrLogger.Println(err.Error())
// 		return err
// 	}
//
// 	txt, err := t.Conn.GetTxtByName(user.ID, jdata["name"])
// 	if err != nil {
// 		t.ErrLogger.Println(err.Error())
// 		return err
// 	}
//
// 	return c.SendString(txt.Content)
// }

func (t *Txtban) readHandler(c *fiber.Ctx) error {
	txtid := c.Params("txtid")
	content, err := t.Conn.GetTxtContentById(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}
	return c.SendString(content)
}

func (t *Txtban) teeHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	txtName := c.Query("name", "")
	if len(txtName) == 0 {
		return errors.New("txt name could not be empty")
	}

	user, err := t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txtid, err := t.Conn.CreateTxt(user.ID, txtName, string(c.Body()))
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"id":   txtid,
		"name": txtName,
	})
}

func (t *Txtban) rmHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	_, err = t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return errors.New("txt id could not be empty")
	}

	err = t.Conn.DeleteTxt(txtid)
	if err != nil {
		return err
	}

	return nil
}

func (t *Txtban) lsHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	user, err := t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txts, err := t.Conn.GetAllTxts(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(txts)
}

func (t *Txtban) chtxtHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	_, err = t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return errors.New("txt id could not be empty")
	}

	err = t.Conn.ChangeTxtContent(txtid, string(c.Body()))
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	return nil
}

func (t *Txtban) mvHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	_, err = t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return errors.New("txt id could not be empty")
	}

	newId, err := t.Conn.ChangeTxtId(txtid)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"id":   newId,
	})
}

func (t *Txtban) renameHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		return err
	}

	var jdata map[string]string
	if err := c.BodyParser(&jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	_, err = t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return fiber.ErrUnauthorized
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return errors.New("txt id could not be empty")
	}

	err = t.Conn.ChangeTxtName(txtid, jdata["name"])
	if err != nil {
		return err
	}

	return nil
}
