package txtban

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/models"
	"github.com/thehxdev/txtban/tberr"
)

func (t *Txtban) useraddHandler(c *fiber.Ctx) error {
	var jdata JsonData

	if err := parseJsonBody(c.Body(), &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errBadJsonData, fiber.StatusBadRequest)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	pass := jdata.Password
	minPassLen := viper.GetInt("limits.minPasswordLen")
	if len(pass) < minPassLen {
		c.Status(fiber.StatusBadRequest)
		err := tberr.New(
			fmt.Sprintf("password length is less than %d characters", minPassLen),
			fmt.Sprintf("choose a strong password with more than %d characters", minPassLen))
		return sendError(c, err, fiber.StatusBadRequest)
	}

	authKey := models.CreateAuthKey(uuid.String(), pass)
	err = t.DB.CreateUser(uuid.String(), pass, authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]string{
		"uuid":    uuid.String(),
		"authKey": authKey,
	})
}

func (t *Txtban) userdelHandler(c *fiber.Ctx) error {
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

	err = t.DB.DeleteUser(user.ID)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return nil
}

func (t *Txtban) whoamiHandler(c *fiber.Ctx) error {
	var jdata JsonData

	if err := parseJsonBody(c.Body(), &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errBadJsonData, fiber.StatusBadRequest)
	}

	user, err := t.DB.AuthenticateByPassword(jdata.UserId, jdata.Password)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]string{
		"authKey": user.AuthKey,
	})
}

func (t *Txtban) passwdHandler(c *fiber.Ctx) error {
	var jdata JsonData

	if err := parseJsonBody(c.Body(), &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errBadJsonData, fiber.StatusBadRequest)
	}

	user, err := t.DB.AuthenticateByPassword(jdata.UserId, jdata.OldPassword)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errUnauthorized, fiber.StatusUnauthorized)
	}

	newPass := jdata.NewPassword
	newAuthKey := models.CreateAuthKey(user.UUID, newPass)
	err = t.DB.UpdateUserPassword(user.ID, newPass, newAuthKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return sendError(c, errInternalServerError, fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]string{
		"authKey": newAuthKey,
	})
}
