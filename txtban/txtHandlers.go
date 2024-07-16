package txtban

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thehxdev/txtban/tberr"
)

func (t *Txtban) readHandler(c *fiber.Ctx) error {
	txtid := c.Params("txtid")
	content, err := t.DB.GetTxtContentById(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}
	return c.SendString(content)
}

func (t *Txtban) teeHandler(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	authKey, err := getAuthKey(headers)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	txtName := c.Query("name", "")
	if len(txtName) == 0 {
		c.Status(fiber.StatusBadRequest)
		return sendError(c,
			tberr.New("txt name is empty", "choose a name for the txt you want to create"),
			fiber.StatusBadRequest)
	}

	user, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	body := c.Body()
	if len(body) == 0 {
		c.Status(fiber.StatusBadRequest)
		return sendError(c, tberr.New("txt content is empty"), fiber.StatusBadRequest)
	}

	txtid, err := t.DB.CreateTxt(user.ID, txtName, string(body))
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]string{
		"id":   txtid,
		"name": txtName,
	})
}

func (t *Txtban) rmHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return sendError(c, errEmptyTxtID, fiber.StatusBadRequest)
	}

	_, err = t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	err = t.DB.DeleteTxt(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return nil
}

func (t *Txtban) lsHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	user, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	txts, err := t.DB.GetAllTxts(user.ID)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return c.JSON(txts)
}

func (t *Txtban) chtxtHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return sendError(c, errEmptyTxtID, fiber.StatusBadRequest)
	}

	_, err = t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	err = t.DB.ChangeTxtContent(txtid, string(c.Body()))
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return nil
}

func (t *Txtban) mvHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return sendError(c, errEmptyTxtID, fiber.StatusBadRequest)
	}

	_, err = t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	newId, err := t.DB.ChangeTxtId(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]string{
		"id": newId,
	})
}

func (t *Txtban) renameHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errEmptyAuthorizationHeader, fiber.StatusBadRequest)
	}

	txtid := c.Query("txtid", "")
	if txtid == "" {
		c.Status(fiber.StatusNotFound)
		return sendError(c, errEmptyTxtID, fiber.StatusBadRequest)
	}

	var jdata JsonData

	if err := parseJsonBody(c.Body(), &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errBadJsonData, fiber.StatusBadRequest)
	}
    
    if len(jdata.Name) == 0 {
		return sendError(c, errEmptyTxtName, fiber.StatusBadRequest)
    }

	_, err = t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	err = t.DB.ChangeTxtName(txtid, jdata.Name)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return nil
}
