package txtban

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/models"
)

func (t *Txtban) useraddHandler(c *fiber.Ctx) error {
	var jdata map[string]string

	if err := c.BodyParser(&jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	pass := jdata["password"]
	minPassLen := viper.GetInt("limits.minPasswordLen")
	if len(pass) < minPassLen {
		c.Status(fiber.StatusBadRequest)
		return fmt.Errorf("password length is less than %d characters", minPassLen)
	}

	authKey := models.CreateAuthKey(uuid.String(), pass)
	err = t.Conn.CreateUser(uuid.String(), pass, authKey)
	if err != nil {
		return err
	}

	resp := map[string]string{
		"uuid":    uuid.String(),
		"authKey": authKey,
	}

	return c.JSON(&resp)
}

func (t *Txtban) userdelHandler(c *fiber.Ctx) error {
	authKey, err := getAuthKey(c.GetReqHeaders())
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	user, err := t.Conn.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	err = t.Conn.DeleteUser(user.ID)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (t *Txtban) whoamiHandler(c *fiber.Ctx) error {
	var jdata map[string]string

	if err := c.BodyParser(&jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	user, err := t.Conn.AuthenticateByPassword(jdata["uuid"], jdata["password"])
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	return c.JSON(map[string]string{
		"authKey": user.AuthKey,
	})
}

func (t *Txtban) passwdHandler(c *fiber.Ctx) error {
	var jdata map[string]string

	if err := c.BodyParser(&jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	user, err := t.Conn.AuthenticateByPassword(jdata["uuid"], jdata["old_password"])
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	newPass := jdata["new_password"]
	newAuthKey := models.CreateAuthKey(user.UUID, newPass)
	err = t.Conn.UpdateUserPassword(user.ID, newPass, newAuthKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		return err
	}

	return c.JSON(map[string]string{
		"authKey": newAuthKey,
	})
}
