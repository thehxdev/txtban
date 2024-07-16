package txtban

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thehxdev/txtban/tberr"
)

var (
	errEmptyAuthorizationHeader = tberr.New(fiber.ErrBadRequest.Error(), "set 'Authorization' header")
	errUnauthorized             = tberr.New(fiber.ErrUnauthorized.Error())
	errInternalServerError      = tberr.New(fiber.ErrInternalServerError.Error())
	errEmptyTxtID               = tberr.New("txt id is empty", "txt id could not be empty")
	errEmptyTxtName             = tberr.New("txt name is empty", "txt name could not be empty")
	errBadJsonData              = tberr.New("failed to parse request json data")
)

func sendError(c *fiber.Ctx, err error, status int) error {
	c.Status(status)
	return err
}
